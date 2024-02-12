package model

type Response struct {
	Data   interface{}    `json:"data"`
	Status ResponseStatus `json:"status"`
	Msg    string         `json:"msg"`
}

type ResponseStatus string

const (
	ResponseSuccess ResponseStatus = "success"
	ResponseError   ResponseStatus = "error"
	ResponseNop     ResponseStatus = "nop"
)
