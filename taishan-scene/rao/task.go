package rao

const (
	TypePlanTask = iota + 1
)

type TaskDetail[T any] struct {
	TaskID   int32  `json:"task_id"`
	Cron     string `json:"cron"`
	TaskData T      `json:"task_data"`
	Enable   bool   `json:"enable"`
}

type PlanTask struct {
	PlanID int32 `json:"plan_id"`
}

type TaskParam struct {
	PlanID   int32  `json:"plan_id"`
	ID       int32  `json:"id"`
	Type     int32  `json:"type"`
	Cron     string `json:"cron"`
	TaskInfo any
	Enable   bool
	UserID   int32
}
