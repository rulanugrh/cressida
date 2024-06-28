package web

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"msg"`
	Data    interface{} `json:"data"`
}

func (r *Response) Error() string {
	return r.Message
}

func Success(message string, data any) error {
	return &Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

func Created(message string, data any) error {
	return &Response{
		Code:    201,
		Message: message,
		Data:    data,
	}
}

func NotFound(message string) error {
	return &Response{
		Code:    404,
		Message: message,
	}
}

func InternalServerError(message string) error {
	return &Response{
		Code:    500,
		Message: message,
	}
}

func BadRequest(message string) error {
	return &Response{
		Code:    400,
		Message: message,
	}
}

func Forbidden(message string) error {
	return &Response{
		Code:    403,
		Message: message,
	}
}

func Unauthorized(message string) error {
	return &Response{
		Code:    401,
		Message: message,
	}
}
