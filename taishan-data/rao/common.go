package rao

type CommonResponse[T any] struct {
	Code int64  `json:"code"`
	Em   string `json:"em"`
	Et   string `json:"et"`
	Data T      `json:"data"`
}
