package consulkv

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/xxjwxc/public/tools"
)

// Result ...
type Result struct {
	g   []byte
	k   string
	err error
}

// Err ...
func (r *Result) Err() error {
	return r.err
}

// Get ...
func (r *Result) Get() []byte {
	return r.g
}

// Scan ...
func (r *Result) Scan(x interface{}) error {
	bt := fmt.Sprintf("[%v]", string(r.g))
	if reflect.ValueOf(x).Type().String() == "[]string" { // 字符串
		tmp := strings.Split(string(r.g), ",")
		if len(tmp) > 0 {
			bt = fmt.Sprintf(`["%v"]`, strings.Join(tmp, `","`))
		}
	}
	return json.Unmarshal([]byte(bt), x)
}

// Exists ...
func (r *Result) Exists() bool {
	return len(r.g) > 0
}

// Float ...
func (r *Result) Float(defaultValue ...float64) float64 {
	var df float64
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}

	data, _ := strconv.ParseFloat(string(r.g), 64)
	return data
}

// Int ...
func (r *Result) Int(defaultValue ...int64) int64 {
	var df int64
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}
	data, _ := strconv.ParseInt(string(r.g), 10, 64)

	return int64(data)
	//return int64(binary.BigEndian.Uint32(r.g))
}

// Uint ...
func (r *Result) Uint(defaultValue ...uint64) uint64 {
	var df uint64
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}

	data, _ := strconv.ParseInt(string(r.g), 10, 64)
	return uint64(data)
}

// Bool ...
func (r *Result) Bool(defaultValue ...bool) bool {
	var df bool
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}

	b, _ := strconv.ParseBool(string(r.g))

	return b
}

// Bytes ...
func (r *Result) Bytes(defaultValue ...[]byte) []byte {
	var df []byte
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}

	return r.g
}

// String
func (r *Result) String(defaultValue ...string) string {
	var df string
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}

	return string(r.g)
}

// Time ...
func (r *Result) Time(defaultValue ...time.Time) time.Time {
	var df time.Time
	if len(defaultValue) != 0 {
		df = defaultValue[0]
	}

	if !r.Exists() {
		return df
	}

	return tools.StrToTime(string(r.g), "", nil)
}

// Key ...
func (r *Result) Key() string {
	return r.k
}
