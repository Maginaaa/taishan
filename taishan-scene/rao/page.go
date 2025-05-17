package rao

type PageResponse[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}
