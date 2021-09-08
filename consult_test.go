package consult

import (
	"fmt"
	"testing"

	"github.com/xxjwxc/consult/consulkv"
	"github.com/xxjwxc/public/mylog"
)

type Config struct {
	MySQLInfo    MysqlDbInfo `yaml:"mysql_info" consul:"mysql_info"`
	RedisDbInfo  RedisDbInfo `yaml:"redis_info" consul:"redis_info"`
	EtcdInfo     EtcdInfo    `yaml:"etcd_info" consul:"etcd_info"`
	Oauth2Url    string      `yaml:"oauth2_url"`
	Port         string      `yaml:"port" consul:"port"`                   // 端口号
	TtvNum       int         `yaml:"ttv_num" consul:"ttv_num"`             // ttv数量
	RegexHost    string      `yaml:"regex_host" consul:"regex_host"`       // 正则入库
	AliTtsURI    string      `yaml:"ali_tts_uri"`                          // 阿里云tts地址
	RegexMaxLen  int         `yaml:"regex_max_len" consul:"regex_max_len"` // 一次最大文字解析数
	HaibuildHost []string    `yaml:"haibuild_host" consul:"haibuild_host"` // motion grpc 地址列表
	ConsulAddr   string      `yaml:"consul_addr" consul:"consul_addr" `    // consul 地址
}

// MysqlDbInfo mysql database information. mysql 数据库信息
type MysqlDbInfo struct {
	Host     string `validate:"required" consul:"host"`     // Host. 地址
	Port     int    `validate:"required" consul:"port"`     // Port 端口号
	Username string `validate:"required" consul:"username"` // Username 用户名
	Password string `consul:"password"`                     // Password 密码
	Database string `validate:"required" consul:"database"` // Database 数据库名
	Type     int    // 数据库类型: 0:mysql , 1:sqlite , 2:mssql
}

// EtcdInfo etcd config info
type EtcdInfo struct {
	Addrs   []string `yaml:"addrs" consul:"addrs"`     // Host. 地址
	Timeout int      `yaml:"timeout" consul:"timeout"` // 超时时间(秒)
}

// RedisDbInfo redis database information. redis 数据库信息
type RedisDbInfo struct {
	Addrs     []string `yaml:"addrs" consul:"addrs"` // Host. 地址
	Password  string   // Password 密码
	GroupName string   `yaml:"group_name" consul:"group_name"` // 分组名字
	DB        int      `yaml:"db" consul:"db"`                 // 数据库序号
}

func TestMain(t *testing.T) {
	conf := consulkv.NewConfig(
		consulkv.WithPrefix("service/servername"),  // consul kv prefix
		consulkv.WithAddress("192.155.1.150:8500"), // consul address
	)
	if err := conf.Init(); err != nil {
		mylog.Error(err)
		return
	}

	var config Config
	AutoLoadConfig(conf, &config) //  自动加载

	fmt.Println(config)
	config.EtcdInfo.Addrs = append(config.EtcdInfo.Addrs, "192.155.1.150", "192.155.1.151")
	AutoSetConfig(conf, &config, false) // 执行一次更新

	fmt.Println(config)
}
