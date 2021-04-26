package consult

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/xxjwxc/consult/consulkv"
	"github.com/xxjwxc/public/mylog"
	"github.com/xxjwxc/public/tools"
)

type consulElement struct {
	conf *consulkv.Config
}

func (s *consulElement) scanObject(fieldv reflect.Value, field reflect.StructField, prefix string) error {
	consulTag := field.Tag.Get("consul")
	if consulTag == "" || consulTag == "-" { // 空,或者未设置表示忽略
		return nil
	}

	key := prefix + consulTag

	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		result := s.conf.Get(key)
		if result.Exists() {
			fieldv.Set(reflect.ValueOf(int(result.Int(0))))
		}
	case reflect.Bool:
		result := s.conf.Get(key)
		if result.Exists() {
			fieldv.Set(reflect.ValueOf(result.Bool()))
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := s.conf.Get(key)
		if result.Exists() {
			fieldv.Set(reflect.ValueOf(result.Uint()))
		}
	case reflect.Float32, reflect.Float64:
		result := s.conf.Get(key)
		if result.Exists() {
			fieldv.Set(reflect.ValueOf(result.Float()))
		}
	case reflect.String:
		result := s.conf.Get(key)
		if result.Exists() {
			fieldv.Set(reflect.ValueOf(result.String()))
		}
	case reflect.Slice:
		result := s.conf.Get(key)
		if result.Exists() {
			result.Scan(fieldv.Interface())
		}
		// if fieldv.Type().String() == "[]uint8" {
		// 	x := []byte(data.(string))
		// 	v = x
		// } else if fieldv.Type().String() == "[]string" {
		// 	mp := data.([]interface{})
		// 	var ss []string
		// 	for _, v := range mp {
		// 		ss = append(ss, v.(string))
		// 	}
		// 	v = ss
		// } else if fieldv.Type().String() == "[]int" {
		// 	mp := data.([]interface{})
		// 	var ss []int
		// 	for _, v := range mp {
		// 		ss = append(ss, int(v.(float64)))
		// 	}
		// 	v = ss

		// } else {
		// 	v = data
		// }
	case reflect.Struct:
		if fieldv.Type().String() == "time.Time" {
			result := s.conf.Get(key)
			if result.Exists() {
				fieldv.Set(reflect.ValueOf(result.Time()))
			}
		} else { // 其它类型
			dataStructType := fieldv.Type()
			for i := 0; i < dataStructType.NumField(); i++ { // 第一轮
				fieldv1 := fieldv.Field(i)
				field1 := dataStructType.Field(i)

				err := s.scanObject(fieldv1, field1, key+"/")
				if err != nil {
					return err
				}
			}
		}
	default:
		mylog.Errorf("%v not support", reflect.Struct)
	}

	return nil
}

func (s *consulElement) setObject(fieldv reflect.Value, field reflect.StructField, prefix string) error {
	consulTag := field.Tag.Get("consul")
	if consulTag == "" || consulTag == "-" { // 空,或者未设置表示忽略
		return nil
	}

	key := prefix + consulTag

	switch field.Type.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		value := fmt.Sprintf("%v", fieldv.Interface())
		err := s.conf.Put(key, value)
		if err != nil {
			return err
		}
	case reflect.String:
		err := s.conf.Put(key, fieldv.Interface().(string))
		if err != nil {
			return err
		}
	case reflect.Bool:
		value := "false"
		b := fieldv.Interface().(bool)
		if b {
			value = "true"
		}
		err := s.conf.Put(key, value)
		if err != nil {
			return err
		}
	case reflect.Slice, reflect.Array:
		var values []string
		for i := 0; i < fieldv.Len(); i++ {
			values = append(values, fmt.Sprintf("%v", fieldv.Index(i).Interface()))
		}

		err := s.conf.Put(key, strings.Join(values, ","))
		if err != nil {
			return err
		}
	case reflect.Struct:
		if fieldv.Type().String() == "time.Time" {
			err := s.conf.Put(key, tools.GetTimeStr(fieldv.Interface().(time.Time)))
			if err != nil {
				return err
			}
		} else { // 其它类型
			dataStructType := fieldv.Type()
			for i := 0; i < dataStructType.NumField(); i++ { // 第一轮
				fieldv1 := fieldv.Field(i)
				field1 := dataStructType.Field(i)
				err := s.setObject(fieldv1, field1, key+"/")
				if err != nil {
					return err
				}
			}
		}
	default:
		mylog.Errorf("%v not support", reflect.Struct)
	}

	return nil
}
