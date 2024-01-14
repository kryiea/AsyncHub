package db

import (
	"asynchub/flowsvr/src/config"
	"fmt"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/niuniumart/gosdk/gormcli"
)

// DB 数据库实例
var DB *gorm.DB

// InitDB 初始化数据库
func InitDB() error {
	var err error

	// 设置gormcli 的factory配置
	gormcli.Factory = gormcli.GormFactory{
		MaxIdleConn: config.Conf.MySQL.MaxIdleConn, // 最大空闲连接数
		MaxConn:     config.Conf.MySQL.MaxConn,     // 最大连接数
		IdleTimeout: config.Conf.MySQL.IdleTimeout, // 空闲连接超时时间
	}

	// 创建gorm实例
	DB, err = gormcli.Factory.CreateGorm(
		config.Conf.MySQL.User,
		config.Conf.MySQL.Pwd,
		config.Conf.MySQL.Url,
		config.Conf.MySQL.Dbname,
	)
	if err != nil {
		return err
	}

	// 尝试发送 ping 请求
	err = DB.DB().Ping()
	if err != nil {
		fmt.Println("Database connection is not available:", err)
		return err
	}

	fmt.Println("Database connection is available")

	return nil
}

// 数据库错误码 GORM_DUPLACATE_ERR_KEY
const (
	GORM_DUPLACATE_ERR_KEY = "Duplicated entry"
)

// IsDupErr 判断是否是重复错误
func IsDupErr(err error) bool {
	return strings.Contains(err.Error(), GORM_DUPLACATE_ERR_KEY)
}

// GetTaskTableName 获取任务表名, 返回格式为 t_taskType_task
func GetTaskTableName(taskType string) string {
	taskTableName := fmt.Sprintf("t_%s_task", taskType)
	return taskTableName
}
