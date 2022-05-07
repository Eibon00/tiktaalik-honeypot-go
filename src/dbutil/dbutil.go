package dbutil

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	bolt "go.etcd.io/bbolt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	conf "tiktaalik-honeypot-go/src/configurator"
	"tiktaalik-honeypot-go/src/dbutil/dbstructs"
	"tiktaalik-honeypot-go/src/dbutil/dbwork"
	"tiktaalik-honeypot-go/src/write-line/colors"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

// GetDbPath 获取数据库位置
func GetDbPath() string {
	config := conf.ParseConfigFile()
	return fmt.Sprintf("%s/%s", config.Data.DBPath, config.Data.DBName)
}

// CreateDatabase 创建数据库
func CreateDatabase() {
	//DbPath := GetDbPath()
	//db, err := bolt.Open(DbPath+"/tiktaalik.db", 0600, nil)
	db, err := bolt.Open(GetDbPath(), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}
	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			panic(0)
		}
	}(db)
}

// CreateTable 创建表
func CreateTable(TableName string) {
	db, err := bolt.Open(GetDbPath(), 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(TableName))
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)
}

// WriteConfigFile 将配置文件读取到的信息写入数据库,每次启动都会自动写入这个值
func WriteConfigFile(config conf.Config) {
	if conf.FileExists(GetDbPath()) != true {
		log.Printf("[+] Database %s%s%s not exists, create...", colors.Blue, config.Data.DBName, colors.Reset)
	}
	db, err := bolt.Open(fmt.Sprintf("%s/%s", config.Data.DBPath, config.Data.DBName), 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		log.Fatal(err)
	}

	//获取上次启动时的配置,不一致则弹出提示
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("config"))
		if bucket == nil {
			bucket, err = tx.CreateBucketIfNotExists([]byte("config"))
			log.Printf("[+] Bucket %sconfig%s not exists, create...", colors.Blue, colors.Reset)
			if err != nil {
				return err
			}
		}
		var oldConfig conf.Config
		val := bucket.Get([]byte("config_json"))
		err = json.Unmarshal(val, &oldConfig)
		if oldConfig != config {
			log.Printf("[+] Config file changed, use new configuration..")
			newConfig, _ := json.Marshal(config)
			err = bucket.Put([]byte("config_json"), newConfig)
		}
		return nil
	})

	//默认命令允许列表存入数据库
	var cmdFilePath string
	if os.Getenv("HONEYPOT_CONFIG") == "" {
		cmdFilePath = fmt.Sprintf("%s/config/cmds.txt", conf.GetRootPath())
	} else {
		cmdFilePath = fmt.Sprintf("%s/cmds.txt", os.Getenv("HONEYPOT_CONFIG"))
	}
	bytes, _ := ioutil.ReadFile(cmdFilePath)
	commandsList := strings.Split(string(bytes), "\r\n")
	err = db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("commands"))
		if bucket == nil {
			bucket, err = tx.CreateBucketIfNotExists([]byte("commands"))
			log.Printf("[+] Bucket %scommands%s not exists, create...", colors.Purple, colors.Reset)
			if err != nil {
				return err
			}
		}

		//遍历命令列表
		for _, command := range commandsList {
			if bucket.Get([]byte(command)) == nil {
				defaultCommands := dbstructs.Commands{
					IsAllowed: true,
					IsUnknown: false,
				}

				log.Printf("[+] add command %s%s%s to allowed list", colors.Purple, command, colors.Reset)
				jsonByte := dbwork.DumpStruct(defaultCommands)
				err = bucket.Put([]byte(command), []byte(jsonByte))
			}
		}
		return nil
	})

	defer func(db *bolt.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(db)

	log.Printf("[+] Databases configuration complete !")
}
