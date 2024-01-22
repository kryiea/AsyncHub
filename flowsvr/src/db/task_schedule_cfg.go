package db

import (
	"time"

	"github.com/jinzhu/gorm"
)

/*
 * @desc: task schedule config
 * @TaskType：任务类型
 * @ScheduleLimit：每次取多少个任务来执行
 * @MaxProcessingTime：单个任务一次最大执行时间，单位秒
 * @MaxRetryNum:  短任务特有属性，最大重试次数，表示最大重试多少次
 * @MaxRetryInterval: 短任务特有属性，表示短任务重试间隔
 * @
 */

// TaskScheduleCfg cfg
type TaskScheduleCfg struct {
	TaskType          string
	ScheduleLimit     int
	ScheduleInterval  int
	MaxProcessingTime int64
	MaxRetryNum       int
	RetryInterval     int // 初始重试间隔
	MaxRetryInterval  int
	CreateTime        *time.Time
	ModifyTime        *time.Time
}

// TableName 表名
func (p *TaskScheduleCfg) TableName() string {
	return "t_schedule_cfg"
}

// Create 创建记录
func (p *TaskScheduleCfg) Create(db *gorm.DB, task *TaskScheduleCfg) error {
	err := db.
		Table(p.TableName()).
		Create(task).Error
	return err
}

// Save 保存记录
func (p *TaskScheduleCfg) Save(db *gorm.DB, task *TaskScheduleCfg) error {
	// 使用 gorm 数据库操作库进行数据库操作
	err := db.
		// 设置要操作的表名为 TaskScheduleCfg 的表
		Table(p.TableName()).
		// 将 task 保存到数据库中
		Save(task).Error
	return err
}

// GetTaskTypeCfg 获取记录
func (p *TaskScheduleCfg) GetTaskTypeCfg(db *gorm.DB, taskType string) (*TaskScheduleCfg, error) {
	// 创建一个新的TaskScheduleCfg对象
	var cfg = new(TaskScheduleCfg)
	// 查询数据库，获取符合条件的TaskScheduleCfg记录
	err := db.
		// 设置查询的表名为TaskScheduleCfg对应的表名
		Table(p.TableName()).
		// 设置查询条件为task_type等于给定的taskType
		Where("task_type = ?", taskType).
		// 查询第一条符合条件的记录，并保存到cfg中
		First(&cfg).Error
	if err != nil {
		// 如果查询出错，返回nil和错误信息
		return nil, err
	}
	// 返回查询到的TaskScheduleCfg对象和nil错误信息
	return cfg, nil
}

// GetTaskTypeCfgList 获取记录列表
func (p *TaskScheduleCfg) GetTaskTypeCfgList(db *gorm.DB) ([]*TaskScheduleCfg, error) {
	// 创建一个 TaskScheduleCfg 类型的切片，用于存储查询结果
	var taskTypeCfgList = make([]*TaskScheduleCfg, 0)
	// 使用当前对象的 TableName 方法获取表名，并设置到 db 中
	db = db.Table(p.TableName())
	// 使用 db.Find 方法查询 taskTypeCfgList，并将结果存储到 taskTypeCfgList 中
	err := db.Find(&taskTypeCfgList).Error
	// 如果查询过程中出现错误，则返回 nil 和错误信息
	if err != nil {
		return nil, err
	}
	// 返回查询结果和 nil 错误信息
	return taskTypeCfgList, nil
}
