package dbwork

import (
	jsoniter "github.com/json-iterator/go"
	"reflect"
	. "tiktaalik-honeypot-go/src/configurator"
	. "tiktaalik-honeypot-go/src/dbutil/dbstructs"
)

//对数据库操作提供封装的接口，调用格式为:表名，列名，字段，数据，全部为空则删除
var json = jsoniter.ConfigCompatibleWithStandardLibrary

// MyStruct 序列化
type MyStruct interface {
	Dashboard | Commands | Logger | BlackList | MyFiles | Config
}

func DumpStruct[T MyStruct](dump T) []byte {
	dumped, _ := json.Marshal(dump)
	return dumped
}

func IsEmpty[T MyStruct](s T) bool {
	switch (interface{})(s).(type) {
	case Commands:
		return reflect.DeepEqual(s, Commands{})
	case Dashboard:
		return reflect.DeepEqual(s, Dashboard{})
	default:
		return reflect.DeepEqual(s, Config{})
	}
}

//反序列化

func LoadDashboard(jsonByte []byte) Dashboard {
	var result Dashboard
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func LoadCommands(jsonByte []byte) Commands {
	var result Commands
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func LoadLogger(jsonByte []byte) Logger {
	var result Logger
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func LoadBlackList(jsonByte []byte) BlackList {
	var result BlackList
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func LoadMyFiles(jsonByte []byte) MyFiles {
	var result MyFiles
	_ = json.Unmarshal(jsonByte, &result)
	return result
}
