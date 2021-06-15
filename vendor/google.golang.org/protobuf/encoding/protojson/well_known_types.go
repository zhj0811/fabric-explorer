// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package protojson

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"time"

	"google.golang.org/protobuf/internal/detectknown"
	"google.golang.org/protobuf/internal/encoding/json"
	"google.golang.org/protobuf/internal/errors"
	"google.golang.org/protobuf/internal/fieldnum"
	"google.golang.org/protobuf/internal/strs"
	"google.golang.org/protobuf/proto"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// isCustomType returns true if type name has special JSON conversion rules.
// The list of custom types here has to match the ones in marshalCustomType and
// unmarshalCustomType.
func isCustomType(name pref.FullName) bool {
	switch detectknown.Which(name) {
	case detectknown.AnyProto:
	case detectknown.TimestampProto:
	case detectknown.DurationProto:
	case detectknown.WrappersProto:
	case detectknown.StructProto:
	case detectknown.FieldMaskProto:
	case detectknown.EmptyProto:
	default:
		return false
	}
	return true
}

// marshalCustomType marshals given well-known type message that have special
// JSON conversion rules. It needs to be a message type where isCustomType
// returns true, else it will panic.
func (e encoder) marshalCustomType(m pref.Message) error {
	name := m.Descriptor().FullName()
	switch detectknown.Which(name) {
	case detectknown.AnyProto:
		return e.marshalAny(m)
	case detectknown.TimestampProto:
		return e.marshalTimestamp(m)
	case detectknown.DurationProto:
		return e.marshalDuration(m)
	case detectknown.WrappersProto:
		return e.marshalWrapperType(m)
	case detectknown.StructProto:
		return e.marshalStructType(m)
	case detectknown.FieldMaskProto:
		return e.marshalFieldMask(m)
	case detectknown.EmptyProto:
		return e.marshalEmpty(m)
	default:
		panic(fmt.Sprintf("%s does not have a custom marshaler", name))
	}
}

// unmarshalCustomType unmarshals given well-known type message that have
// special JSON conversion rules. It needs to be a message type where
// isCustomType returns true, else it will panic.
func (d decoder) unmarshalCustomType(m pref.Message) error {
	name := m.Descriptor().FullName()
	switch detectknown.Which(name) {
	case detectknown.AnyProto:
		return d.unmarshalAny(m)
	case detectknown.TimestampProto:
		return d.unmarshalTimestamp(m)
	case detectknown.DurationProto:
		return d.unmarshalDuration(m)
	case detectknown.WrappersProto:
		return d.unmarshalWrapperType(m)
	case detectknown.StructProto:
		return d.unmarshalStructType(m)
	case detectknown.FieldMaskProto:
		return d.unmarshalFieldMask(m)
	case detectknown.EmptyProto:
		return d.unmarshalEmpty(m)
	default:
		panic(fmt.Sprintf("%s does not have a custom unmarshaler", name))
	}
}

// The JSON representation of an Any message uses the regular representation of
// the deserialized, embedded message, with an additional field `@type` which
// contains the type URL. If the embedded message type is well-known and has a
// custom JSON representation, that representation will be embedded adding a
// field `value` which holds the custom JSON in addition to the `@type` field.

func (e encoder) marshalAny(m pref.Message) error {
	fds := m.Descriptor().Fields()
	fdType := fds.ByNumber(fieldnum.Any_TypeUrl)
	fdValue := fds.ByNumber(fieldnum.Any_Value)

	// Start writing the JSON object.
	e.StartObject()
	defer e.EndObject()

	if !m.Has(fdType) {
		if !m.Has(fdValue) {
			// If message is empty, marshal out empty JSON object.
			return nil
		} else {
			// Return error if type_url field is not set, but value is set.
			return errors.New("%s: type_url is not set", m.Descriptor().FullName())
		}
	}

	typeVal := m.Get(fdType)
	valueVal := m.Get(fdValue)

	// Marshal out @type field.
	typeURL := typeVal.String()
	e.WriteName("@type")
	if err := e.WriteString(typeURL); err != nil {
		return err
	}

	// Resolve the type in order to unmarshal value field.
	emt, err := e.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil {
		return errors.New("%s: unable to resolve %q: %v", m.Descriptor().FullName(), typeURL, err)
	}

	em := emt.New()
	err = proto.UnmarshalOptions{
		AllowPartial: true, // never check required fields inside an Any
		Resolver:     e.opts.Resolver,
	}.Unmarshal(valueVal.Bytes(), em.Interface())
	if err != nil {
		return errors.New("%s: unable to unmarshal %q: %v", m.Descriptor().FullName(), typeURL, err)
	}

	// If type of value has custom JSON encoding, marshal out a field "value"
	// with corresponding custom JSON encoding of the embedded message as a
	// field.
	if isCustomType(emt.Descriptor().FullName()) {
		e.WriteName("value")
		return e.marshalCustomType(em)
	}

	// Else, marshal out the embedded message's fields in this Any object.
	if err := e.marshalFields(em); err != nil {
		return err
	}

	return nil
}

func (d decoder) unmarshalAny(m pref.Message) error {
	// Peek to check for json.ObjectOpen to avoid advancing a read.
	start, err := d.Peek()
	if err != nil {
		return err
	}
	if start.Kind() != json.ObjectOpen {
		return d.unexpectedTokenError(start)
	}

	// Use another decoder to parse the unread bytes for @type field. This
	// avoids advancing a read from current decoder because the current JSON
	// object may contain the fields of the embedded type.
	dec := decoder{d.Clone(), UnmarshalOptions{}}
	tok, err := findTypeURL(dec)
	switch err {
	case errEmptyObject:
		// An empty JSON object translates to an empty Any message.
		d.Read() // Read json.ObjectOpen.
		d.Read() // Read json.ObjectClose.
		return nil

	case errMissingType:
		if d.opts.DiscardUnknown {
			// Treat all fields as unknowns, similar to an empty object.
			return d.skipJSONValue()
		}
		// Use start.Pos() for line position.
		return d.newError(start.Pos(), err.Error())

	default:
		if err != nil {
			return err
		}
	}

	typeURL := tok.ParsedString()
	emt, err := d.opts.Resolver.FindMessageByURL(typeURL)
	if err != nil {
		return d.newError(tok.Pos(), "unable to resolve %v: %q", tok.RawString(), err)
	}

	// Create new message for the embedded message type and unmarshal into it.
	em := emt.New()
	if isCustomType(emt.Descriptor().FullName()) {
		// If embedded message is a custom type,
		// unmarshal the JSON "value" field into it.
		if err := d.unmarshalAnyValue(em); err != nil {
			return err
		}
	} else {
		// Else unmarshal the current JSON object into it.
		if err := d.unmarshalMessage(em, true); err != nil {
			return err
		}
	}
	// Serialize the embedded message and assign the resulting bytes to the
	// proto value field.
	b, err := proto.MarshalOptions{
		AllowPartial:  true, // No need to check required fields inside an Any.
		Deterministic: true,
	}.Marshal(em.Interface())
	if err != nil {
		return d.newError(start.Pos(), "error in marshaling Any.value field: %v", err)
	}

	fds := m.Descriptor().Fields()
	fdType := fds.ByNumber(fieldnum.Any_TypeUrl)
	fdValue := fds.ByNumber(fieldnum.Any_Value)

	m.Set(fdType, pref.ValueOfString(typeURL))
	m.Set(fdValue, pref.ValueOfBytes(b))
	return nil
}

var errEmptyObject = fmt.Errorf(`empty object`)
var errMissingType = fmt.Errorf(`missing "@type" field`)

// findTypeURL returns the token for the "@type" field value from the given
// JSON bytes. It is expected that the given bytes start with json.ObjectOpen.
// It returns errEmptyObject if the JSON object is empty or errMissingType if
// @type field does not exist. It returns other error if the @type field is not
// valid or other decoding issues.
func findTypeURL(d decoder) (json.Token, error) {
	var typeURL string
	var typeTok json.Token
	numFields := 0
	// Skip start object.
	d.Read()

Loop:
	for {
		tok, err := d.Read()
		if err != nil {
			return json.Token{}, err
		}

		switch tok.Kind() {
		case json.ObjectClose:
			if typeURL == "" {
				// Did not find @type field.
				if numFields > 0 {
					return json.Token{}, errMissingType
				}
				return json.Token{}, errEmptyObject
			}
			break Loop

		case json.Name:
			numFields++
			if tok.Name() != "@type" {
				// Skip value.
				if err := d.skipJSONValue(); err != nil {
					return json.Token{}, err
				}
				continue
			}

			// Return error if this was previously set already.
			if typeURL != "" {
				return json.Token{}, d.newError(tok.Pos(), `duplicate "@type" field`)
			}
			// Read field value.
			tok, err := d.Read()
			if err != nil {
				return json.Token{}, err
			}
			if tok.Kind() != json.String {
				return json.Token{}, d.newError(tok.Pos(), `@type field value is not a string: %v`, tok.RawString())
			}
			typeURL = tok.ParsedString()
			if typeURL == "" {
				return json.Token{}, d.newError(tok.Pos(), `@type field contains empty value`)
			}
			typeTok = tok
		}
	}

	return typeTok, nil
}

// skipJSONValue parses a JSON value (null, boolean, string, number, object and
// array) in order to advance the read to the next JSON value. It relies on
// the decoder returning an error if the types are not in valid sequence.
func (d decoder) skipJSONValue() error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	// Only need to continue reading for objects and arrays.
	switch tok.Kind() {
	case json.ObjectOpen:
		for {
			tok, err := d.Read()
			if err != nil {
				return err
			}
			switch tok.Kind() {
			case json.ObjectClose:
				return nil
			case json.Name:
				// Skip object field value.
				if err := d.skipJSONValue(); err != nil {
					return err
				}
			}
		}

	case json.ArrayOpen:
		for {
			tok, err := d.Peek()
			if err != nil {
				return err
			}
			switch tok.Kind() {
			case json.ArrayClose:
				d.Read()
				return nil
			default:
				// Skip array item.
				if err := d.skipJSONValue(); err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// unmarshalAnyValue unmarshals the given custom-type message from the JSON
// object's "value" field.
func (d decoder) unmarshalAnyValue(m pref.Message) error {
	// Skip ObjectOpen, and start reading the fields.
	d.Read()

	var found bool // Used for detecting duplicate "value".
	for {
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		case json.ObjectClose:
			if !found {
				return d.newError(tok.Pos(), `missing "value" field`)
			}
			return nil

		case json.Name:
			switch tok.Name() {
			case "@type":
				// Skip the value as this was previously parsed already.
				d.Read()

			case "value":
				if found {
					return d.newError(tok.Pos(), `duplicate "value" field`)
				}
				// Unmarshal the field value into the given message.
				if err := d.unmarshalCustomType(m); err != nil {
					return err
				}
				found = true

			default:
				if d.opts.DiscardUnknown {
					if err := d.skipJSONValue(); err != nil {
						return err
					}
					continue
				}
				return d.newError(tok.Pos(), "unknown field %v", tok.RawString())
			}
		}
	}
}

// Wrapper types are encoded as JSON primitives like string, number or boolean.

// The "value" field has the same field number for all wrapper types.
const wrapperFieldNumber = fieldnum.BoolValue_Value

func (e encoder) marshalWrapperType(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(wrapperFieldNumber)
	val := m.Get(fd)
	return e.marshalSingular(val, fd)
}

func (d decoder) unmarshalWrapperType(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(wrapperFieldNumber)
	val, err := d.unmarshalScalar(fd)
	if err != nil {
		return err
	}
	m.Set(fd, val)
	return nil
}

// The JSON representation for Empty is an empty JSON object.

func (e encoder) marshalEmpty(pref.Message) error {
	e.StartObject()
	e.EndObject()
	return nil
}

func (d decoder) unmarshalEmpty(pref.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.ObjectOpen {
		return d.unexpectedTokenError(tok)
	}

	for {
		tok, err := d.Read()
		if err != nil {
			return err
		}
		switch tok.Kind() {
		case json.ObjectClose:
			return nil

		case json.Name:
			if d.opts.DiscardUnknown {
				if err := d.skipJSONValue(); err != nil {
					return err
				}
				continue
			}
			return d.newError(tok.Pos(), "unknown field %v", tok.RawString())

		default:
			return d.unexpectedTokenError(tok)
		}
	}
}

func (e encoder) marshalStructType(m pref.Message) error {
	switch m.Descriptor().Name() {
	case "Struct":
		return e.marshalStruct(m)
	case "ListValue":
		return e.marshalListValue(m)
	case "Value":
		return e.marshalKnownValue(m)
	default:
		panic(fmt.Sprintf("invalid struct type: %v", m.Descriptor().FullName()))
	}
}

func (d decoder) unmarshalStructType(m pref.Message) error {
	switch m.Descriptor().Name() {
	case "Struct":
		return d.unmarshalStruct(m)
	case "ListValue":
		return d.unmarshalListValue(m)
	case "Value":
		return d.unmarshalKnownValue(m)
	default:
		panic(fmt.Sprintf("invalid struct type: %v", m.Descriptor().FullName()))
	}
}

// The JSON representation for Struct is a JSON object that contains the encoded
// Struct.fields map and follows the serialization rules for a map.

func (e encoder) marshalStruct(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(fieldnum.Struct_Fields)
	return e.marshalMap(m.Get(fd).Map(), fd)
}

func (d decoder) unmarshalStruct(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(fieldnum.Struct_Fields)
	return d.unmarshalMap(m.Mutable(fd).Map(), fd)
}

// The JSON representation for ListValue is JSON array that contains the encoded
// ListValue.values repeated field and follows the serialization rules for a
// repeated field.

func (e encoder) marshalListValue(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(fieldnum.ListValue_Values)
	return e.marshalList(m.Get(fd).List(), fd)
}

func (d decoder) unmarshalListValue(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(fieldnum.ListValue_Values)
	return d.unmarshalList(m.Mutable(fd).List(), fd)
}

// The JSON representation for a Value is dependent on the oneof field that is
// set. Each of the field in the oneof has its own custom serialization rule. A
// Value message needs to be a oneof field set, else it is an error.

func (e encoder) marshalKnownValue(m pref.Message) error {
	od := m.Descriptor().Oneofs().ByName("kind")
	fd := m.WhichOneof(od)
	if fd == nil {
		return errors.New("%s: none of the oneof fields is set", m.Descriptor().FullName())
	}
	return e.marshalSingular(m.Get(fd), fd)
}

func (d decoder) unmarshalKnownValue(m pref.Message) error {
	tok, err := d.Peek()
	if err != nil {
		return err
	}

	var fd pref.FieldDescriptor
	var val pref.Value
	switch tok.Kind() {
	case json.Null:
		d.Read()
		fd = m.Descriptor().Fields().ByNumber(fieldnum.Value_NullValue)
		val = pref.ValueOfEnum(0)

	case json.Bool:
		tok, err := d.Read()
		if err != nil {
			return err
		}
		fd = m.Descriptor().Fields().ByNumber(fieldnum.Value_BoolValue)
		val = pref.ValueOfBool(tok.Bool())

	case json.Number:
		tok, err := d.Read()
		if err != nil {
			return err
		}
		fd = m.Descriptor().Fields().ByNumber(fieldnum.Value_NumberValue)
		var ok bool
		val, ok = unmarshalFloat(tok, 64)
		if !ok {
			return d.newError(tok.Pos(), "invalid google.protobuf.Value: %v", tok.RawString())
		}

	case json.String:
		// A JSON string may have been encoded from the number_value field,
		// e.g. "NaN", "Infinity", etc. Parsing a proto double type also allows
		// for it to be in JSON string form. Given this custom encoding spec,
		// however, there is no way to identify that and hence a JSON string is
		// always assigned to the string_value field, which means that certain
		// encoding cannot be parsed back to the same field.
		tok, err := d.Read()
		if err != nil {
			return err
		}
		fd = m.Descriptor().Fields().ByNumber(fieldnum.Value_StringValue)
		val = pref.ValueOfString(tok.ParsedString())

	case json.ObjectOpen:
		fd = m.Descriptor().Fields().ByNumber(fieldnum.Value_StructValue)
		val = m.NewField(fd)
		if err := d.unmarshalStruct(val.Message()); err != nil {
			return err
		}

	case json.ArrayOpen:
		fd = m.Descriptor().Fields().ByNumber(fieldnum.Value_ListValue)
		val = m.NewField(fd)
		if err := d.unmarshalListValue(val.Message()); err != nil {
			return err
		}

	default:
		return d.newError(tok.Pos(), "invalid google.protobuf.Value: %v", tok.RawString())
	}

	m.Set(fd, val)
	return nil
}

// The JSON representation for a Duration is a JSON string that ends in the
// suffix "s" (indicating seconds) and is preceded by the number of seconds,
// with nanoseconds expressed as fractional seconds.
//
// Durations less than one second are represented with a 0 seconds field and a
// positive or negative nanos field. For durations of one second or more, a
// non-zero value for the nanos field must be of the same sign as the seconds
// field.
//
// Duration.seconds must be from -315,576,000,000 to +315,576,000,000 inclusive.
// Duration.nanos must be from -999,999,999 to +999,999,999 inclusive.

const (
	secondsInNanos       = 999999999
	maxSecondsInDuration = 315576000000
)

func (e encoder) marshalDuration(m pref.Message) error {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(fieldnum.Duration_Seconds)
	fdNanos := fds.ByNumber(fieldnum.Duration_Nanos)

	secsVal := m.Get(fdSeconds)
	nanosVal := m.Get(fdNanos)
	secs := secsVal.Int()
	nanos := nanosVal.Int()
	if secs < -maxSecondsInDuration || secs > maxSecondsInDuration {
		return errors.New("%s: seconds out of range %v", m.Descriptor().FullName(), secs)
	}
	if nanos < -secondsInNanos || nanos > secondsInNanos {
		return errors.New("%s: nanos out of range %v", m.Descriptor().FullName(), nanos)
	}
	if (secs > 0 && nanos < 0) || (secs < 0 && nanos > 0) {
		return errors.New("%s: signs of seconds and nanos do not match", m.Descriptor().FullName())
	}
	// Generated output always contains 0, 3, 6, or 9 fractional digits,
	// depending on required precision, followed by the suffix "s".
	f := "%d.%09d"
	if nanos < 0 {
		nanos = -nanos
		if secs == 0 {
			f = "-%d.%09d"
		}
	}
	x := fmt.Sprintf(f, secs, nanos)
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, ".000")
	e.WriteString(x + "s")
	return nil
}

func (d decoder) unmarshalDuration(m pref.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.String {
		return d.unexpectedTokenError(tok)
	}

	secs, nanos, ok := parseDuration(tok.ParsedString())
	if !ok {
		return d.newError(tok.Pos(), "invalid google.protobuf.Duration value %v", tok.RawString())
	}
	// Validate seconds. No need to validate nanos because parseDuration would
	// have covered that already.
	if secs < -maxSecondsInDuration || secs > maxSecondsInDuration {
		return d.newError(tok.Pos(), "google.protobuf.Duration value out of range: %v", tok.RawString())
	}

	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(fieldnum.Duration_Seconds)
	fdNanos := fds.ByNumber(fieldnum.Duration_Nanos)

	m.Set(fdSeconds, pref.ValueOfInt64(secs))
	m.Set(fdNanos, pref.ValueOfInt32(nanos))
	return nil
}

// parseDuration parses the given input string for seconds and nanoseconds value
// for the Duration JSON format. The format is a decimal number with a suffix
// 's'. It can have optional plus/minus sign. There needs to be at least an
// integer or fractional part. Fractional part is limited to 9 digits only for
// nanoseconds precision, regardless of whether there are trailing zero digits.
// Example values are 1s, 0.1s, 1.s, .1s, +1s, -1s, -.1s.
func parseDuration(input string) (int64, int32, bool) {
	b := []byte(input)
	size := len(b)
	if size < 2 {
		return 0, 0, false
	}
	if b[size-1] != 's' {
		return 0, 0, false
	}
	b = b[:size-1]

	// Read optional plus/minus symbol.
	var neg bool
	switch b[0] {
	case '-':
		neg = true
		b = b[1:]
	case '+':
		b = b[1:]
	}
	if len(b) == 0 {
		return 0, 0, false
	}

	// Read the integer part.
	var intp []byte
	switch {
	case b[0] == '0':
		b = b[1:]

	case '1' <= b[0] && b[0] <= '9':
		intp = b[0:]
		b = b[1:]
		n := 1
		for len(b) > 0 && '0' <= b[0] && b[0] <= '9' {
			n++
			b = b[1:]
		}
		intp = intp[:n]

	case b[0] == '.':
		// Continue below.

	default:
		return 0, 0, false
	}

	hasFrac := false
	var frac [9]byte
	if len(b) > 0 {
		if b[0] != '.' {
			return 0, 0, false
		}
		// Read the fractional part.
		b = b[1:]
		n := 0
		for len(b) > 0 && n < 9 && '0' <= b[0] && b[0] <= '9' {
			frac[n] = b[0]
			n++
			b = b[1:]
		}
		// It is not valid if there are more bytes left.
		if len(b) > 0 {
			return 0, 0, false
		}
		// Pad fractional part with 0s.
		for i := n; i < 9; i++ {
			frac[i] = '0'
		}
		hasFrac = true
	}

	var secs int64
	if len(intp) > 0 {
		var err error
		secs, err = strconv.ParseInt(string(intp), 10, 64)
		if err != nil {
			return 0, 0, false
		}
	}

	var nanos int64
	if hasFrac {
		nanob := bytes.TrimLeft(frac[:], "0")
		if len(nanob) > 0 {
			var err error
			nanos, err = strconv.ParseInt(string(nanob), 10, 32)
			if err != nil {
				return 0, 0, false
			}
		}
	}

	if neg {
		if secs > 0 {
			secs = -secs
		}
		if nanos > 0 {
			nanos = -nanos
		}
	}
	return secs, int32(nanos), true
}

// The JSON representation for a Timestamp is a JSON string in the RFC 3339
// format, i.e. "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z" where
// {year} is always expressed using four digits while {month}, {day}, {hour},
// {min}, and {sec} are zero-padded to two digits each. The fractional seconds,
// which can go up to 9 digits, up to 1 nanosecond resolution, is optional. The
// "Z" suffix indicates the timezone ("UTC"); the timezone is required. Encoding
// should always use UTC (as indicated by "Z") and a decoder should be able to
// accept both UTC and other timezones (as indicated by an offset).
//
// Timestamp.seconds must be from 0001-01-01T00:00:00Z to 9999-12-31T23:59:59Z
// inclusive.
// Timestamp.nanos must be from 0 to 999,999,999 inclusive.

const (
	maxTimestampSeconds = 253402300799
	minTimestampSeconds = -62135596800
)

func (e encoder) marshalTimestamp(m pref.Message) error {
	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(fieldnum.Timestamp_Seconds)
	fdNanos := fds.ByNumber(fieldnum.Timestamp_Nanos)

	secsVal := m.Get(fdSeconds)
	nanosVal := m.Get(fdNanos)
	secs := secsVal.Int()
	nanos := nanosVal.Int()
	if secs < minTimestampSeconds || secs > maxTimestampSeconds {
		return errors.New("%s: seconds out of range %v", m.Descriptor().FullName(), secs)
	}
	if nanos < 0 || nanos > secondsInNanos {
		return errors.New("%s: nanos out of range %v", m.Descriptor().FullName(), nanos)
	}
	// Uses RFC 3339, where generated output will be Z-normalized and uses 0, 3,
	// 6 or 9 fractional digits.
	t := time.Unix(secs, nanos).UTC()
	x := t.Format("2006-01-02T15:04:05.000000000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, "000")
	x = strings.TrimSuffix(x, ".000")
	e.WriteString(x + "Z")
	return nil
}

func (d decoder) unmarshalTimestamp(m pref.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.String {
		return d.unexpectedTokenError(tok)
	}

	t, err := time.Parse(time.RFC3339Nano, tok.ParsedString())
	if err != nil {
		return d.newError(tok.Pos(), "invalid google.protobuf.Timestamp value %v", tok.RawString())
	}
	// Validate seconds. No need to validate nanos because time.Parse would have
	// covered that already.
	secs := t.Unix()
	if secs < minTimestampSeconds || secs > maxTimestampSeconds {
		return d.newError(tok.Pos(), "google.protobuf.Timestamp value out of range: %v", tok.RawString())
	}

	fds := m.Descriptor().Fields()
	fdSeconds := fds.ByNumber(fieldnum.Timestamp_Seconds)
	fdNanos := fds.ByNumber(fieldnum.Timestamp_Nanos)

	m.Set(fdSeconds, pref.ValueOfInt64(secs))
	m.Set(fdNanos, pref.ValueOfInt32(int32(t.Nanosecond())))
	return nil
}

// The JSON representation for a FieldMask is a JSON string where paths are
// separated by a comma. Fields name in each path are converted to/from
// lower-camel naming conventions. Encoding should fail if the path name would
// end up differently after a round-trip.

func (e encoder) marshalFieldMask(m pref.Message) error {
	fd := m.Descriptor().Fields().ByNumber(fieldnum.FieldMask_Paths)
	list := m.Get(fd).List()
	paths := make([]string, 0, list.Len())

	for i := 0; i < list.Len(); i++ {
		s := list.Get(i).String()
		// Return error if conversion to camelCase is not reversible.
		cc := strs.JSONCamelCase(s)
		if s != strs.JSONSnakeCase(cc) {
			return errors.New("%s.paths contains irreversible value %q", m.Descriptor().FullName(), s)
		}
		paths = append(paths, cc)
	}

	e.WriteString(strings.Join(paths, ","))
	return nil
}

func (d decoder) unmarshalFieldMask(m pref.Message) error {
	tok, err := d.Read()
	if err != nil {
		return err
	}
	if tok.Kind() != json.String {
		return d.unexpectedTokenError(tok)
	}
	str := strings.TrimSpace(tok.ParsedString())
	if str == "" {
		return nil
	}
	paths := strings.Split(str, ",")

	fd := m.Descriptor().Fields().ByNumber(fieldnum.FieldMask_Paths)
	list := m.Mutable(fd).List()

	for _, s := range paths {
		s = strings.TrimSpace(s)
		// Convert to snake_case. Unlike encoding, no validation is done because
		// it is not possible to know the original path names.
		list.Append(pref.ValueOfString(strs.JSONSnakeCase(s)))
	}
	return nil
}
