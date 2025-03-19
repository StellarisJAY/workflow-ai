package common

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Total   int    `json:"total"`
}

func NewSuccessResponse(data any) *Response {
	return &Response{
		Code:    200,
		Message: "ok",
		Data:    data,
	}
}

func NewErrorResponse(message string) *Response {
	return &Response{
		Code:    500,
		Message: message,
	}
}

func NewSuccessResponseWithTotal(data any, total int) *Response {
	return &Response{
		Code:    200,
		Message: "ok",
		Data:    data,
		Total:   total,
	}
}

func NewErrorResponseWithData(message string, data interface{}) *Response {
	return &Response{
		Code:    500,
		Message: message,
		Data:    data,
	}
}
