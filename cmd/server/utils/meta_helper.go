package utils

import (
	"github.com/febriandani/ecommerce-be-system-service/protogen/golang/proto/system"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

// NewMetaSuccess creates a standard successful meta
func NewMetaSuccess(message, traceId string) *system.Meta {
	return &system.Meta{
		Code:    200,
		Status:  "success",
		Message: message,
		TraceId: traceId,
	}
}

func NewMetaSuccessWithData(message string, data interface{}, traceId string) *system.Meta {
	var anyData *anypb.Any
	if msg, ok := data.(proto.Message); ok && data != nil {
		anyData, _ = anypb.New(msg)
	}
	return &system.Meta{
		Code:    200,
		Status:  "success",
		Message: message,
		Data:    anyData,
		TraceId: traceId,
	}
}

func NewMetaSuccessWithAny(message string, data *anypb.Any, traceId string) *system.Meta {
	return &system.Meta{
		Code:    200,
		Status:  "success",
		Message: message,
		Data:    data,
		TraceId: traceId,
	}
}

// NewMetaError creates a standard error meta
func NewMetaError(code int32, message, internalMessage, traceId string) *system.Meta {
	return &system.Meta{
		Code:            code,
		Status:          "error",
		Message:         message,
		InternalMessage: internalMessage,
		TraceId:         traceId,
	}
}

// NewMetaFail is for failed but not error (e.g. validation)
func NewMetaFail(message, traceId string) *system.Meta {
	return &system.Meta{
		Code:    400,
		Status:  "fail",
		Message: message,
		TraceId: traceId,
	}
}
