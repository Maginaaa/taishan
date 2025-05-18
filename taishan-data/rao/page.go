package rao

type PageResponse[T any] struct {
	Total int64 `json:"total"`
	List  []T   `json:"list"`
}

type CommonPageSearch struct {
	StartTime int64 `json:"start_time"`
	EndTime   int64 `json:"end_time"`
	Page      int   `json:"page" binding:"required,gt=0"`      // 页码
	PageSize  int   `json:"page_size" binding:"required,gt=0"` // 每页条数
}
