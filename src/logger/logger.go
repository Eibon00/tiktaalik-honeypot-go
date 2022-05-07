package logger

import (
	"tiktaalik-honeypot-go/src/dbutil"
	"tiktaalik-honeypot-go/src/dbutil/dbstructs"
)

func ReceiveResult(record []string) dbstructs.Logger {
	var log dbstructs.Logger
	log.Record = record
	return log
}

func LogClient(ip, loginTime string, Record []string) bool {
	var log dbstructs.Logger
	log.IP = ip
	log.LastLogin = loginTime
	log.Record = Record
	//return dbutil.RecordAll("logger", strconv.FormatInt(time.Now().UnixNano(), 10), log)
	return dbutil.RecordAll("logger", ip, log)
}

func LogBlacklist(ip, loginTime string) bool {
	var blacklist dbstructs.BlackList
	blacklist.LastLogin = loginTime
	blacklist.Banned = false
	return dbutil.RecordAll("blacklist", ip, blacklist)
}
