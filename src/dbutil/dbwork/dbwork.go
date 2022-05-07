package dbwork

import (
	jsoniter "github.com/json-iterator/go"
	. "tiktaalik-honeypot-go/src/configurator"
	. "tiktaalik-honeypot-go/src/dbutil/dbstructs"
)

//对数据库操作提供封装的接口，调用格式为:表名，列名，字段，数据，全部为空则删除
var json = jsoniter.ConfigCompatibleWithStandardLibrary

//序列化
type myStruct interface {
	Dashboard | Commands | Logger | BlackList | MyFiles | Config
}

func DumpStruct[T myStruct](dump T) []byte {
	dumped, _ := json.Marshal(dump)
	return dumped
}

//反序列化
func loadDashboard(jsonByte []byte) Dashboard {
	var result Dashboard
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func loadCommands(jsonByte []byte) Commands {
	var result Commands
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func loadLogger(jsonByte []byte) Logger {
	var result Logger
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func loadBlackList(jsonByte []byte) BlackList {
	var result BlackList
	_ = json.Unmarshal(jsonByte, &result)
	return result
}

func loadMyFiles(jsonByte []byte) MyFiles {
	var result MyFiles
	_ = json.Unmarshal(jsonByte, &result)
	return result
}
