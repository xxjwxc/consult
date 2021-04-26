[![Build Status](https://travis-ci.org/xxjwxc/consult.svg?branch=master)](https://travis-ci.org/xxjwxc/consult)
[![Go Report Card](https://goreportcard.com/badge/github.com/xxjwxc/consult)](https://goreportcard.com/report/github.com/xxjwxc/consult)
[![GoDoc](https://godoc.org/github.com/xxjwxc/consult?status.svg)](https://godoc.org/github.com/xxjwxc/consult)

## [中文文档](README_zh.md)
# consult
A consul key/value tool for golang

## Usage

### install
```
go get -u github.com/xxjwxc/consult@master
```

### New Config
```golang
conf := consulkv.NewConfig()
```

### With Options
```golang
conf := consulkv.NewConfig(
    consulkv.WithPrefix(prefix),             // consul kv prefix
    consulkv.WithAddress(address),           // consul address
    consulkv.WithAuth(username, password),   // cosul auth
    consulkv.WithToken(token),               // cousl token
    consulkv.WithLoger(loger),               // loger
)

```

### Init
```golang
if err := conf.Init();err !=nil {
    return err
}
```

### Put
```golang
if err := conf.Put(key, value);err !=nil {
    return err
}
```

### Delete
```golang
if err := conf.Delete(key);err !=nil {
    return err
}
```

### Get
```golang
// scan
if err := conf.Get(key).Scan(x);err !=nil {
    return err
}

// get float
float := conf.Get(key).Float()

// get float with default
float := conf.Get(key).Float(defaultFloat)

// get int
i := conf.Get(key).Int()

// get int with default
i := conf.Get(key).Int(defaultInt)

// get uint
uInt := conf.Get(key).Uint()

// get uint with default
uInt := conf.Get(key).Uint(defaultUint)

// get bool
b := conf.Get(key).Bool()

// get bool with default
b := conf.Get(key).Bool(defaultBool)

// get []byte
bytes := conf.Get(key).Bytes()

// get uint with default
bytes := conf.Get(key).bytes(defaultBytes)

// get string
str := conf.Get(key).String()

// get string with default
str := conf.Get(key).String(defaultStr)

// get time
t := conf.Get(key).Time()

// get time with default
t := conf.Get(key).Time(defaultTime)

// get nested key values
conf.Get(key).Get(nextKey1).Get(nextKey2).String()
```

### Watch
```golang
conf.Watch(path, func(r *Result){
    r.Scan(x)
})

```

### Stop Watch
```golang
// stop single watcher
conf.StopWatch(path)

// stop multiple watcher
conf.StopWatch(path1, path2)

// stop all watcher
conf.StopWatch()
```


### Automatic acquisition / update through go tag

- When defining a variable, add the label `consult:""` to register and obtain automatically
```go
import (
	"github.com/xxjwxc/consult"
)

type Info struct {
    Port  string  `yaml:"port" consul:"port"` // port
}

var info Info
consult.AutoLoadConfig(conf, &info) //  AutoLoad

consult.AutoSetConfig(conf, &info, false) // AutoUpdate

```

## example

```go 

import (
	"fmt"
	"testing"

	"github.com/xxjwxc/consult/consulkv"
    "github.com/xxjwxc/consult"
)

type Config struct {
	MySQLInfo    MysqlDbInfo `yaml:"mysql_info" consul:"mysql_info"`
	Port         string      `yaml:"port" consul:"port"`                   // 端口号
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

func main() {
	conf := consulkv.NewConfig(
		consulkv.WithPrefix("service/servername"),      // consul kv prefix
		consulkv.WithAddress("192.155.1.150:8500"), // consul address
	)
	if err := conf.Init(); err != nil {
		mylog.Error(err)
		return
	}

	var config Config
	consult.AutoLoadConfig(conf, &config) //  自动加载
	fmt.Println(config)

	consult.AutoSetConfig(conf, &config, false) // 执行一次更新
	fmt.Println(config)
}

```