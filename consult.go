package consult

import (
	"errors"
	"reflect"

	"github.com/xxjwxc/consult/consulkv"
)

// AutoLoadConfig 自动加载config配置
func AutoLoadConfig(conf *consulkv.Config, obj interface{}) error {
	dataStruct := reflect.Indirect(reflect.ValueOf(obj))
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	elm := &consulElement{conf}

	dataStructType := dataStruct.Type()
	for i := 0; i < dataStructType.NumField(); i++ { // 第一轮
		fieldv := dataStruct.Field(i)
		field := dataStructType.Field(i)

		err := elm.scanObject(fieldv, field, "")
		if err != nil {
			return err
		}
	}
	return nil
}

// AutoSetConfig 自动设置config配置
func AutoSetConfig(conf *consulkv.Config, obj interface{}, isUpdate bool) error {
	if !isUpdate { // 不用更新
		list, err := conf.List()
		if err != nil {
			return err
		}
		if len(list) > 0 { // 不用更新
			return nil
		}
	}

	dataStruct := reflect.Indirect(reflect.ValueOf(obj))
	if dataStruct.Kind() != reflect.Struct {
		return errors.New("expected a pointer to a struct")
	}

	elm := &consulElement{conf}

	dataStructType := dataStruct.Type()
	for i := 0; i < dataStructType.NumField(); i++ { // 第一轮
		fieldv := dataStruct.Field(i)
		field := dataStructType.Field(i)

		err := elm.setObject(fieldv, field, "")
		if err != nil {
			return err
		}
	}
	return nil
}
