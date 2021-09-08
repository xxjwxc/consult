package consult

import (
	"github.com/xxjwxc/consult/consulkv"
)

type KVer interface {
	Init() error
	Put(path string, value interface{}) error
	Get(keys ...string) *consulkv.Result
	List() ([]string, error)
}
