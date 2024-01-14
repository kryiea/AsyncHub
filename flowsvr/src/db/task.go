package db

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/niuniumart/gosdk/martlog"
)

var TaskNsp Task

// 任务表
type Task struct {
	Id               int       `gorm:"primary_key;AUTO_INCREMENT:15"`
	UserId           string    `gorm:"type:varchar(256);column:user_id;not null;default:'';index:idx_user_id"`
	TaskId           string    `gorm:"type:varchar(256);column:task_id;not null;default:'';unique_index:idx_task_id"`
	TaskType         string    `gorm:"type:varchar(128);column:task_type;not null"`
	TaskStage        string    `gorm:"type:varchar(128);column:task_stage;not null"`
	Status           int       `gorm:"type:tinyint(3);column:status;not null;default:1;index:idx_tasktype_status_modify_time"`
	Priority         int       `gorm:"type:int(11);column:priority;not null;default:0;comment:'优先级'"`
	CrtRetryNum      int       `gorm:"type:int(11);column:crt_retry_num;not null;default:0;comment:'已经重试几次了'"`
	MaxRetryNum      int       `gorm:"type:int(11);column:max_retry_num;not null;default:0;comment:'最大能重试几次'"`
	MaxRetryInterval int       `gorm:"type:int(11);column:max_retry_interval;not null;default:0;comment:'最大重试间隔'"`
	ScheduleLog      string    `gorm:"type:varchar(4096);column:schedule_log;not null;default:'';comment:'调度信息记录'"`
	TaskContext      string    `gorm:"type:varchar(8192);column:task_context;not null;default:'';comment:'任务上下文，用户自定义'"`
	OrderTime        int64     `gorm:"type:int(20);column:order_time;not null;default:0;comment:'调度时间，越小调度越优先';index:idx_tasktype_status_modify_time"`
	CreateTime       time.Time `gorm:"type:datetime;column:create_time;not null;autoCreateTime;"`
	ModifyTime       time.Time `gorm:"type:datetime;column:modify_time;not null;autoCreateTime;"`
}

// getTableName 返回一个字符串类型的表名，格式为"t_任务类型_当前表名_位置"
func (p *Task) getTableName(taskType string, pos string) string {
	return fmt.Sprintf("t_%s_%s_%s", taskType, p.TableName(), pos)
}

// TablName 返回表名
func (p *Task) TableName() string {
	return "task"
}

// GenTaskId 生成任务id
func GenTaskId(taskType string, pos string) string {
	// 替换下划线
	taskType = strings.Replace(taskType, "_", "-", -1)
	//  生成uuid
	return fmt.Sprintf("%+v_%s_%s", uuid.New(), taskType, pos)
}

// getTablePosFromTaskId 获取任务id中的位置, 返回两个字符串，第一个是任务类型，第二个是位置
func (p *Task) getTablePosFromTaskId(taskID string) (string, string) {
	s := strings.Split(taskID, "_")
	switch len(s) {
	case 3:
		s[1] = strings.Replace(s[1], "-", "_", -1)
		return s[1], s[2]
	default:
		martlog.Errorf("大错误,任务id格式不对,没有_匹配", taskID)
		return "", ""
	}
}

// BatchSetStatus 批量设置任务状态
func (p *Task) BatchSetStatus(db *gorm.DB, taskIdList []string, status TaskEnum) error {
	var dic = map[string]interface{}{
		"status":      status,
		"modify_time": time.Now().Format("2006-01-02 15:04:05"),
	}
	tmpTaskId := taskIdList[0]
	taskTyepe, pos := p.getTablePosFromTaskId(tmpTaskId)
	db = db.Table(p.getTableName(taskTyepe, pos)).Where("task_id IN (?)", taskIdList).Update(dic)
	err := db.Error
	if err != nil {
		return err
	}
	return nil
}

// CreateTable 创建任务信息表，需要传入任务类型和位置
func (p *Task) CreateTable(db *gorm.DB, taskType, pos string) error {
	// 构建表名
	newTableName := p.getTableName(taskType, pos)
	return db.Table(newTableName).AutoMigrate(&Task{}).Error
}

// Find 查找任务, 返回任务信息，如果找不到，则返回nil
func (p *Task) Find(db *gorm.DB, taskId string) (*Task, error) {
	var data = &Task{}
	taskType, pos := p.getTablePosFromTaskId(taskId)
	err := db.Table(p.getTableName(taskType, pos)).Where("task_id = ?", taskId).First(data).Error
	return data, err
}

// Create 创建任务
func (p *Task) Create(db *gorm.DB, taskType, pos string, task *Task) error {
	err := db.Table(p.getTableName(taskType, pos)).Create(task).Error
	return err
}

// Save 保存任务记录
func (p *Task) Save(db *gorm.DB, task *Task) error {
	taskType, pos := p.getTablePosFromTaskId(task.TaskId)
	err := db.Table(p.getTableName(taskType, pos)).Save(task).Error
	return err
}

// GetTaskList 获取任务列表
func (p *Task) GetTaskList(db *gorm.DB, taskType string, pos string, status TaskEnum, limit int) ([]*Task, error) {
	var taskList = make([]*Task, 0)
	err := db.
		Table(p.getTableName(taskType, pos)).
		Where("status = ?", status).
		Order("order_time").
		Limit(limit).
		Find(&taskList).Error
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// GetAliveTaskList 获取处于激活状态的任务列表
func (p *Task) GetAliveTaskList(db *gorm.DB, taskType string, pos string, limit int) ([]*Task, error) {
	// 构建任务列表
	var taskList = make([]*Task, 0)
	// 设置状态集合
	var statusSet = []TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING}
	// 在数据库中查找符合条件的任务列表，并返回结果和错误信息
	err := db.
		Table(p.getTableName(taskType, pos)).
		Where("status in (?)", statusSet).
		Order("modify_time").
		Limit(limit).
		Find(&taskList).Error
	// 返回任务列表和错误信息
	if err != nil {
		return nil, err
	}
	return taskList, nil
}

// GetAliveTaskCount 获取处于激活状态的任务数
func (p *Task) GetAliveTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	// 调用getTaskCount方法，传入数据库连接、任务类型、位置和任务状态数组
	return p.getTaskCount(db, taskType, pos,
		[]TaskEnum{TASK_STATUS_PENDING, TASK_STATUS_PROCESSING})
}

// GetTaskCount 获取任务数
func (p *Task) getTaskCount(db *gorm.DB, taskType, pos string, statusSet []TaskEnum) (int, error) {
	// 定义计数变量
	var count int
	// 调用数据库操作方法
	err := db.
		Table(p.getTableName(taskType, pos)).
		// 添加查询条件，根据状态集合查询
		Where("status in (?)", statusSet).
		// 统计符合条件的记录数
		Count(&count).Error
	if err != nil {
		// 如果发生错误，返回错误和计数
		return count, err
	}
	// 返回计数和nil错误
	return count, nil
}

// GetAllTaskCount 获取所有任务数
func (p *Task) GetAllTaskCount(db *gorm.DB, taskType, pos string) (int, error) {
	return p.GetTableCount(db, taskType, pos)
}

// 获取表记录总数
func (p *Task) GetTableCount(db *gorm.DB, taskType, pos string) (int, error) {
	// 定义计数变量
	var count int
	// 调用数据库的Table方法，获取指定表名的表对象
	err := db.Table(p.getTableName(taskType, pos)).Count(&count).Error
	// 如果出错
	if err != nil {
		// 返回计数和错误
		return count, err
	}
	// 返回计数和nil
	return count, nil
}

// GetTaskCountByStatus 根据状态获取任务数
func (p *Task) GetTaskCountByStatus(db *gorm.DB, taskType, pos string, status int) (int, error) {
	var count int
	err := db.Table(p.getTableName(taskType, pos)).Where("status = ?", status).Count(&count).Error
	if err != nil {
		return count, err
	}
	return count, nil
}

// GetFinishTaskCount 获取完成的任务数
func (p *Task) GetFinishTaskCount(db *gorm.DB, taskType string, pos string, statusSet []TaskEnum) (int, error) {
	return p.getTaskCount(db, taskType, pos, []TaskEnum{TASK_STATUS_FAILED, TASK_STATUS_SUCCESS})
}

// SetStatusPending 更新任务状态为等待处理
func (p *Task) SetStatusPending(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_PENDING)
}

// SetStatusSucc 设置任务为成功状态
func (p *Task) SetStatusSucc(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_SUCCESS)
}

// SetStatusFailed 设置任务为失败状态
func (p *Task) SetStatusFailed(db *gorm.DB, taskId string) error {
	return p.SetStatus(db, taskId, TASK_STATUS_FAILED)
}

TODO 1.14
// setStatus 更新任务状态
func (p *Task) SetStatus(db *gorm.DB, taskId string, status TaskEnum) error {
	// 创建一个包含状态字段的字典
	var dic = map[string]interface{}{
		"status": status,
	}
	// 根据任务ID获取任务类型和位置信息
	taskType, pos := p.getTablePosFromTaskId(taskId)
	// 根据任务类型和位置信息更新任务状态
	err := db.
		Table(p.getTableName(taskType, pos)).
		Where("task_id = ?", taskId).
		Updates(dic).Error
	if err != nil {
		return err
	}
	return nil
}
