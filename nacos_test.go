package consult

import (
	"fmt"
	"testing"

	nacos "github.com/xxjwxc/consult/nacoskv"
	"github.com/xxjwxc/public/mylog"
)

type Config1 struct {
	MysqlInfo MysqlInfo `yaml:"mysql_info" nacos:"mysql_info"`
	EtcdInfo  EtcdInfo1 `yaml:"etcd_info" nacos:"etcd_info"`
	Demo      string    `yaml:"demo" nacos:"demo"`
	AliTTSUri string    `yaml:"ali_tts_uri" nacos:"ali_tts_uri"`
}

type MysqlInfo struct {
	Host     string `validate:"required" nacos:"host"`     // Host. 地址
	Port     int    `validate:"required" nacos:"port"`     // Port 端口号
	Username string `validate:"required" nacos:"username"` // Username 用户名
	Password string `nacos:"password"`                     // Password 密码
	Database string `validate:"required" nacos:"database"` // Database 数据库名
	Type     int    // 数据库类型: 0:mysql , 1:sqlite , 2:mssql
}

// EtcdInfo etcd config info
type EtcdInfo1 struct {
	Addrs   []string `yaml:"addrs" nacos:"addrs"`     // Host. 地址
	Timeout int      `yaml:"timeout" nacos:"timeout"` // 超时时间(秒)
}

func TestInit(t *testing.T) {
	kVer := nacos.NewConfig("192.155.1.150", 8848,
		nacos.WithPrefix("nacos"),
		nacos.WithNamespaceId("08066f57-5fbd-4d4b-8188-45893cd8b9a2"),
		nacos.WithLogDir("18Nacos\\nacos\\log"),
		nacos.WithCacheDir("18Nacos\\nacos\\cache"),
	)
	if err := kVer.Init(); err != nil {
		mylog.Error(err)
		return
	}

	var config Config
	AutoLoadConfig(kVer, &config)
	fmt.Println("Config: ", config)

	//config.Demo = "RichieCheng"
	//config.EtcdInfo.Timeout = 3
	//config.EtcdInfo.Addrs = []string{"192.155.1.150:2379", "127.0.0.1:3306"}
	//config.AliTTSUri = "https://nls-gateway.cn-shanghai.aliyuncs.com/stream/v1/tts"

	AutoSetConfig(kVer, &config, true) // 执行一次更新
	fmt.Println(config)
}
