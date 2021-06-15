package common

import "errors"

const (
	Success = iota //成功
	RequestFormatErr
	OpenFileErr
	ExcelFormatErr

	GetDBErr
	InsertDBErr
	UpdateDBErr

	GetPolicyErr
	GetServiceErr

	MarshalJSONErr
	UnmarshalJSONErr
	InvokeErr
	QueryErr

	TokenNilErr
	TokenInvalidErr

	UserInvalidErr
	UserNameOrPasswdErr
)

var (
	ErrExcelFormat = errors.New("excel format error")
)
