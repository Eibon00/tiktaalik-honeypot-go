package configurator

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Config struct {
	Data Data `json:"data"`
	Auth Auth `json:"auth"`
}

type Data struct {
	DBPath string `json:"db_path"`
	DBName string `json:"db_name"`
}

type Auth struct {
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
}

// ParseConfigFile 从docker的预设变量中读取配置文件位置，若没有则使用自身运行目录下的配置文件目录
func ParseConfigFile() Config {
	var config Config
	var configPath string

	if os.Getenv("HONEYPOT_CONFIG") == "" {
		configPath = fmt.Sprintf("%s/config/config.json", GetRootPath())
	} else {
		configPath = fmt.Sprintf("%s/config.json", os.Getenv("HONEYPOT_CONFIG"))
	}

	configBytes, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(0)
	}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic(0)
	}
	return config
}

// GetRootPath 检测当前运行目录
func GetRootPath() string {
	//_, b, _, _ := runtime.Caller(0)
	//return filepath.Join(filepath.Dir(b), "../../")
	dir := getCurrentAbPathByExecutable()
	tmpDir, _ := filepath.EvalSymlinks(os.TempDir())
	if strings.Contains(dir, tmpDir) {
		return filepath.Join(filepath.Dir(getCurrentAbPathByCaller()), "../")
	}
	return filepath.Join(filepath.Dir(dir), "../")
}

func getCurrentAbPathByExecutable() string {
	exePath, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	res, _ := filepath.EvalSymlinks(filepath.Dir(exePath))
	return res
}

func getCurrentAbPathByCaller() string {
	var abPath string
	_, filename, _, ok := runtime.Caller(0)
	if ok {
		abPath = path.Dir(filename)
	}
	return abPath
}

func FileExists(filepath string) bool {
	stat, err := os.Stat(filepath)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return !stat.IsDir()
}

func DirExists(dirpath string) bool {
	_, err := os.Stat(dirpath)
	if errors.Is(err, fs.ErrNotExist) {
		return false
	}
	return true
}

func Config2Byte(config interface{}, order binary.ByteOrder) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	err := binary.Write(buf, order, config)
	if err != nil {
		panic(err)
	}
	return buf.Bytes(), nil
}

// ConfigCheck 检测config合理性，待实现
func ConfigCheck(conf Config) error {
	return nil
}
