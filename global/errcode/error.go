package errcode

import (
	"errors"
	"fmt"
)

// WrapError 全局错误对象
type WrapError struct {
	no  int
	msg string
	err error
}

// NewWrapError 返回错误对象指针
func NewWrapError(no int, msg string) *WrapError {
	return &WrapError{
		no:  no,
		msg: msg,
		err: errors.New(msg),
	}
}

// GetNo 获取标号
func (e *WrapError) GetNo() int {
	return e.no
}

// ToError 获取 error 类型
func (e *WrapError) ToError() error {
	return e.err
}

// Error 获取错误信息
func (e *WrapError) Error() string {
	return e.err.Error()
}

// String 格式化输出
func (e *WrapError) String() string {
	return fmt.Sprintf("No:%v:Msg:%v", e.no, e.msg)
}

func (e *WrapError) WithDetail(text string) *WrapError {
	errMsg := fmt.Sprintf("%v. %v", e.err.Error(), text)
	e.err = errors.New(errMsg)
	return e
}
