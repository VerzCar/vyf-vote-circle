package model

type Response struct {
	Status ResponseStatus `json:"status"`
	Msg    string         `json:"msg"`
	Data   interface{}    `json:"data"`
}

type ResponseStatus string

const (
	ResponseSuccess ResponseStatus = "success"
	ResponseError   ResponseStatus = "error"
	ResponseNop     ResponseStatus = "nop"
)
