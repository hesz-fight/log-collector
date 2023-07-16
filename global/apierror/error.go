package apierror

import (
	"errors"
	"fmt"
)

// ApiError 全局错误对象
type ApiError struct {
	no  int
	msg string
	err error
}

// NewApiError 返回错误对象指针
func NewApiError(no int, msg string) *ApiError {
	return &ApiError{
		no:  no,
		msg: msg,
		err: errors.New(msg),
	}
}

// GetNo 获取标号
func (e *ApiError) GetNo() int {
	return e.no
}

// ToError 获取 error 类型
func (e *ApiError) ToError() error {
	return e.err
}

// Error 获取错误信息
func (e *ApiError) Error() string {
	return e.err.Error()
}

// String 格式化输出
func (e *ApiError) String() string {
	return fmt.Sprintf("No:%v:Msg:%v", e.no, e.msg)
}
