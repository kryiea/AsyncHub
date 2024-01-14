package db

type TaskEnum int

const (
	TASK_STATUS_PENDING    TaskEnum = 1 // 待处理
	TASK_STATUS_PROCESSING TaskEnum = 2 // 正在处理
	TASK_STATUS_SUCCESS    TaskEnum = 3 // 成功
	TASK_STATUS_FAILED     TaskEnum = 4 // 失败
)

const MAX_PRIORITY = 3600 * 24 * 30 * 12 // 优先1年

func IsValidStatus(status TaskEnum) bool {
	switch status {
	case TASK_STATUS_PENDING, TASK_STATUS_PROCESSING, TASK_STATUS_SUCCESS, TASK_STATUS_FAILED:
		return true
	default:
		return false
	}
}
