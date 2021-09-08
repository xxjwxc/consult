package nacos

import (
	"fmt"
	"strings"

	"github.com/xxjwxc/consult/consulkv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

// Option ...
type Option func(opt *Config)

// NewConfig ...
func NewConfig(ipAddr string, port uint64, opts ...Option) *Config {
	c := &Config{
		serverConfig: constant.NewServerConfig(ipAddr, port),
		clientConfig: constant.NewClientConfig(),
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

func WithPrefix(prefix string) Option {
	return func(c *Config) {
		c.prefix = prefix
	}
}

func WithNamespaceId(NamespaceId string) Option {
	return func(c *Config) {
		c.clientConfig.NamespaceId = NamespaceId
	}
}

func WithLogDir(LogDir string) Option {
	return func(c *Config) {
		c.clientConfig.LogDir = LogDir
	}
}

func WithCacheDir(CacheDir string) Option {
	return func(c *Config) {
		c.clientConfig.CacheDir = CacheDir
	}
}

type Config struct {
	prefix string

	client       config_client.IConfigClient
	serverConfig *constant.ServerConfig
	clientConfig *constant.ClientConfig
}

func (c *Config) Init() error {
	client, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig: c.clientConfig,
			ServerConfigs: []constant.ServerConfig{
				*c.serverConfig,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("init failed: %w", err)
	}

	c.client = client
	return nil
}

func (c *Config) SearchConfig(path string) (map[string]string, error) {
	var dataId string
	var group string
	strList := strings.Split(path, "/")
	if len(strList) < 4 {
		dataId = strList[1]
		group = "DEFAULT_GROUP"
	} else {
		dataId = strList[2]
		group = strList[1]
	}
	searchPage, err := c.client.SearchConfig(vo.SearchConfigParam{
		Search:   "accurate", // blur 模糊匹配
		DataId:   dataId,
		Group:    group,
		PageNo:   1,
		PageSize: 10,
	})
	if err != nil {
		return nil, err
	}

	mp := make(map[string]string)
	for _, item := range searchPage.PageItems {
		mp[path] = item.Content
	}

	return mp, nil
}

func (c *Config) List() ([]string, error) {
	searchPage, err := c.client.SearchConfig(vo.SearchConfigParam{
		Search:   "blur", // blur 模糊匹配
		DataId:   "",
		Group:    "",
		PageNo:   1,
		PageSize: 10,
	})
	if err != nil {
		return nil, err
	}

	var list []string
	for _, item := range searchPage.PageItems {
		str := item.Group + "/" + item.DataId + "/" + item.Content
		list = append(list, str)
	}
	return list, nil
}

func (c *Config) Put(path string, value interface{}) error {
	data := []byte(fmt.Sprintf("%v", value))
	//fmt.Println("data: ", data)
	//fmt.Println("path: ", path)

	var dataId string
	var group string
	strList := strings.Split(path, "/")
	if len(strList) > 1 {
		dataId = strList[1]
		group = strList[0]
	} else {
		dataId = strList[0]
		group = "DEFAULT_GROUP"
	}
	_, err := c.client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: string(data),
	})
	if err != nil {
		fmt.Printf("PublishConfig err:%+v \n", err)
	}
	//fmt.Println("b: ", b)
	return nil
}

func (c *Config) Get(keys ...string) (ret *consulkv.Result) {
	var (
		path = c.absPath(keys...) + "/"
		//fields []string
	)

	//fmt.Println("Path: ", path)
	//fmt.Println("Fields: ", fields)

	ret = &consulkv.Result{}
	list, err := c.SearchConfig(path)
	if err != nil {
		ret.SetErr(fmt.Errorf("Get SearchConfig failed: %w", err))
		return
	}

	for k, v := range list {
		ret.SetG([]byte(v))
		ret.SetK(k)
	}

	return
}

func (c *Config) absPath(keys ...string) string {
	if len(keys) == 0 {
		return c.prefix
	}

	if len(keys[0]) == 0 {
		return c.prefix
	}

	if len(c.prefix) == 0 {
		return strings.Join(keys, "/")
	}

	return c.prefix + "/" + strings.Join(keys, "/")
}
