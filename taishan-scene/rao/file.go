package rao

type FileInfo struct {
	ID          int32    `json:"id"`
	PlanID      int32    `json:"plan_id"`
	Name        string   `json:"name"`
	Size        int32    `json:"size"`
	Rows        int32    `json:"rows"`
	Column      []Column `json:"column"`
	Status      bool     `json:"status"`
	CreatedTime string   `json:"created_time"`
	UpdatedTime string   `json:"updated_time"`
}

type Column struct {
	Col          string `json:"col"`
	Alias        string `json:"alias"`
	FileName     string `json:"file_name"`
	FirstLineVal string `json:"first_line_val"`
	ColIndex     int    `json:"col_index"`
	ReadType     int    `json:"read_type"`
}

type ColumnUpdateReq struct {
	PlanID int32 `json:"plan_id" binding:"required,gt=0"`
	Column
}

type FileDownload struct {
	Path     string `json:"path"`
	FileName string `json:"file_name"`
}
