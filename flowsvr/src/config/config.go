package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

var Conf *TomlConfig

// 全局配置
type TomlConfig struct {
	Common commonConfig
	MySQL  mysqlConfig
	Redis  redisConfig
	Task   TaskConfig
}

// 普通配置
type commonConfig struct {
	Port    int  `toml:"port"`
	OpenTLS bool `toml:"open_tls"`
}

// mysql数据库配置
type mysqlConfig struct {
	Url         string `toml:"url"`
	User        string `toml:"user"`
	Pwd         string `toml:"pwd"`
	Dbname      string `toml:"db_name"`
	MaxIdleConn int    `toml:"max_idle"`
	MaxConn     int    `toml:"max_conn"`
	IdleTimeout int    `toml:"idle_timeout"`
}

// redis配置
type redisConfig struct {
	Url             string `toml:"url"`
	Auth            string `toml:"auth"`
	MaxIdle         int    `toml:"max_idle"`
	MaxActive      int    `toml:"max_active"`
	IdleTimeout     int    `toml:"idle_timeout"`
	CacheTimeoutDay int    `toml:"cache_timeout_day"`
}

// 任务配置
type TaskConfig struct {
	TableMaxRows        int `toml:"table_max_rows"`        // 每个表最大行数
	AliveThreshold      int `toml:"alive_threshold"`       // 存活时间阈值
	SplitInteval        int `toml:"split_inteval"`         // 分表间隔
	LongProcessInterval int `toml:"long_process_interval"` // 长任务间隔
	MoveInterval        int `toml:"move_interval"`         // 更新begin下标的时间间隔
	MaxProcessTime      int `toml:"max_process_time"`      // 最大处理时间
}

// 包init函数
func init() {
	inferRootDir()
}

// 初始化配置
func Init() {
	initConf()
}

// 初始化配置文件
func initConf() {
	Conf = new(TomlConfig)
	Conf.LoadConfig()
}

func (c *TomlConfig) LoadConfig() {
	// 判断配置文件是否存在
	if _, err := os.Stat(GetConfigPath()); err != nil {
		panic(err)
	}
	// 读取配置文件
	if _, err := toml.DecodeFile(GetConfigPath(), &c); err != nil {
		panic(err)
	}

}

// 项目主目录
var rootDir string

// 推断root目录
func inferRootDir() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	// 推断
	var infer func(string) string
	infer = func(dir string) string {
		if exists(dir + "/main.go") {
			return dir
		}
		parent := filepath.Dir(dir)
		return infer(parent)
	}

	rootDir = infer(pwd)
}

// exists 函数判断给定的目录是否存在
// 参数dir是要判断的目录路径
// 返回值是bool类型，true表示目录存在，false表示目录不存在
func exists(dir string) bool {
	// 查找主机是不是存在dir目录
	_, err := os.Stat(dir)
	// 如果目录存在，返回true；否则返回false
	return err == nil || os.IsExist(err)
}

func GetConfigPath() string {
	return rootDir + "/config/config.toml"
}
