package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"im/serializer"
)

// ErrResponse 返回错误信息
func ErrResponse(err error) serializer.Response {
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		return serializer.Response{
			Status: 400,
			Msg:    "参数错误",
			Error:  fmt.Sprint(err),
		}
	}

	var unmarshalTypeError *json.UnmarshalTypeError
	if errors.As(err, &unmarshalTypeError) {
		return serializer.Response{
			Status: 400,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}

	return serializer.Response{
		Status: 400,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
