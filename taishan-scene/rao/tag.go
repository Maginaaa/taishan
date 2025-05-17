package rao

type TagSearchReq struct {
	Type int32 `json:"type"`
}

type Tag struct {
	ID    int32  `json:"id"`
	Label string `json:"label"`
}
