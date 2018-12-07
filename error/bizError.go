package error

import "fmt"

type BizError struct {
	code int
	message string
}

func (p *BizError) Error() string {
	return fmt.Sprintf("BizError code: %v, message %s", p.code, p.message)
}

func newBizError(code int, message string) *BizError {
	return &BizError{
		code: code,
		message: message,
	}
}
