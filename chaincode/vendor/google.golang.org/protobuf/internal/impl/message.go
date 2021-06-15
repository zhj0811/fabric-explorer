// Copyright 2018 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package impl

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"google.golang.org/protobuf/internal/genname"
	"google.golang.org/protobuf/reflect/protoreflect"
	pref "google.golang.org/protobuf/reflect/protoreflect"
)

// MessageInfo provides protobuf related functionality for a given Go type
// that represents a message. A given instance of MessageInfo is tied to
// exactly one Go type, which must be a pointer to a struct type.
//
// The exported fields must be populated before any methods are called
// and cannot be mutated after set.
type MessageInfo struct {
	// GoReflectType is the underlying message Go type and must be populated.
	GoReflectType reflect.Type // pointer to struct

	// Desc is the underlying message descriptor type and must be populated.
	Desc pref.MessageDescriptor

	// Exporter must be provided in a purego environment in order to provide
	// access to unexported fields.
	Exporter exporter

	// OneofWrappers is list of pointers to oneof wrapper struct types.
	OneofWrappers []interface{}

	initMu   sync.Mutex // protects all unexported fields
	initDone uint32

	reflectMessageInfo // for reflection implementation
	coderMessageInfo   // for fast-path method implementations
}

// exporter is a function that returns a reference to the ith field of v,
// where v is a pointer to a struct. It returns nil if it does not support
// exporting the requested field (e.g., already exported).
type exporter func(v interface{}, i int) interface{}

// getMessageInfo returns the MessageInfo for any message type that
// is generated by our implementation of protoc-gen-go (for v2 and on).
// If it is unable to obtain a MessageInfo, it returns nil.
func getMessageInfo(mt reflect.Type) *MessageInfo {
	m, ok := reflect.Zero(mt).Interface().(pref.ProtoMessage)
	if !ok {
		return nil
	}
	mr, ok := m.ProtoReflect().(interface{ ProtoMessageInfo() *MessageInfo })
	if !ok {
		return nil
	}
	return mr.ProtoMessageInfo()
}

func (mi *MessageInfo) init() {
	// This function is called in the hot path. Inline the sync.Once logic,
	// since allocating a closure for Once.Do is expensive.
	// Keep init small to ensure that it can be inlined.
	if atomic.LoadUint32(&mi.initDone) == 0 {
		mi.initOnce()
	}
}

func (mi *MessageInfo) initOnce() {
	mi.initMu.Lock()
	defer mi.initMu.Unlock()
	if mi.initDone == 1 {
		return
	}

	t := mi.GoReflectType
	if t.Kind() != reflect.Ptr && t.Elem().Kind() != reflect.Struct {
		panic(fmt.Sprintf("got %v, want *struct kind", t))
	}
	t = t.Elem()

	si := mi.makeStructInfo(t)
	mi.makeReflectFuncs(t, si)
	mi.makeCoderMethods(t, si)

	atomic.StoreUint32(&mi.initDone, 1)
}

// getPointer returns the pointer for a message, which should be of
// the type of the MessageInfo. If the message is of a different type,
// it returns ok==false.
func (mi *MessageInfo) getPointer(m pref.Message) (p pointer, ok bool) {
	switch m := m.(type) {
	case *messageState:
		return m.pointer(), m.messageInfo() == mi
	case *messageReflectWrapper:
		return m.pointer(), m.messageInfo() == mi
	}
	return pointer{}, false
}

type (
	SizeCache       = int32
	WeakFields      = map[int32]protoreflect.ProtoMessage
	UnknownFields   = []byte
	ExtensionFields = map[int32]ExtensionField
)

var (
	sizecacheType       = reflect.TypeOf(SizeCache(0))
	weakFieldsType      = reflect.TypeOf(WeakFields(nil))
	unknownFieldsType   = reflect.TypeOf(UnknownFields(nil))
	extensionFieldsType = reflect.TypeOf(ExtensionFields(nil))
)

type structInfo struct {
	sizecacheOffset offset
	weakOffset      offset
	unknownOffset   offset
	extensionOffset offset

	fieldsByNumber        map[pref.FieldNumber]reflect.StructField
	oneofsByName          map[pref.Name]reflect.StructField
	oneofWrappersByType   map[reflect.Type]pref.FieldNumber
	oneofWrappersByNumber map[pref.FieldNumber]reflect.Type
}

func (mi *MessageInfo) makeStructInfo(t reflect.Type) structInfo {
	si := structInfo{
		sizecacheOffset: invalidOffset,
		weakOffset:      invalidOffset,
		unknownOffset:   invalidOffset,
		extensionOffset: invalidOffset,

		fieldsByNumber:        map[pref.FieldNumber]reflect.StructField{},
		oneofsByName:          map[pref.Name]reflect.StructField{},
		oneofWrappersByType:   map[reflect.Type]pref.FieldNumber{},
		oneofWrappersByNumber: map[pref.FieldNumber]reflect.Type{},
	}

fieldLoop:
	for i := 0; i < t.NumField(); i++ {
		switch f := t.Field(i); f.Name {
		case genname.SizeCache, genname.SizeCacheA:
			if f.Type == sizecacheType {
				si.sizecacheOffset = offsetOf(f, mi.Exporter)
			}
		case genname.WeakFields, genname.WeakFieldsA:
			if f.Type == weakFieldsType {
				si.weakOffset = offsetOf(f, mi.Exporter)
			}
		case genname.UnknownFields, genname.UnknownFieldsA:
			if f.Type == unknownFieldsType {
				si.unknownOffset = offsetOf(f, mi.Exporter)
			}
		case genname.ExtensionFields, genname.ExtensionFieldsA, genname.ExtensionFieldsB:
			if f.Type == extensionFieldsType {
				si.extensionOffset = offsetOf(f, mi.Exporter)
			}
		default:
			for _, s := range strings.Split(f.Tag.Get("protobuf"), ",") {
				if len(s) > 0 && strings.Trim(s, "0123456789") == "" {
					n, _ := strconv.ParseUint(s, 10, 64)
					si.fieldsByNumber[pref.FieldNumber(n)] = f
					continue fieldLoop
				}
			}
			if s := f.Tag.Get("protobuf_oneof"); len(s) > 0 {
				si.oneofsByName[pref.Name(s)] = f
				continue fieldLoop
			}
		}
	}

	// Derive a mapping of oneof wrappers to fields.
	oneofWrappers := mi.OneofWrappers
	for _, method := range []string{"XXX_OneofFuncs", "XXX_OneofWrappers"} {
		if fn, ok := reflect.PtrTo(t).MethodByName(method); ok {
			for _, v := range fn.Func.Call([]reflect.Value{reflect.Zero(fn.Type.In(0))}) {
				if vs, ok := v.Interface().([]interface{}); ok {
					oneofWrappers = vs
				}
			}
		}
	}
	for _, v := range oneofWrappers {
		tf := reflect.TypeOf(v).Elem()
		f := tf.Field(0)
		for _, s := range strings.Split(f.Tag.Get("protobuf"), ",") {
			if len(s) > 0 && strings.Trim(s, "0123456789") == "" {
				n, _ := strconv.ParseUint(s, 10, 64)
				si.oneofWrappersByType[tf] = pref.FieldNumber(n)
				si.oneofWrappersByNumber[pref.FieldNumber(n)] = tf
				break
			}
		}
	}

	return si
}

func (mi *MessageInfo) New() protoreflect.Message {
	return mi.MessageOf(reflect.New(mi.GoReflectType.Elem()).Interface())
}
func (mi *MessageInfo) Zero() protoreflect.Message {
	return mi.MessageOf(reflect.Zero(mi.GoReflectType).Interface())
}
func (mi *MessageInfo) Descriptor() protoreflect.MessageDescriptor { return mi.Desc }
