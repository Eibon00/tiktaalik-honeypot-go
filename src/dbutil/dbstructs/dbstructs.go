package dbstructs

//将key-value存储形式封装为id-json的形式

// Dashboard 后台页面表结构
type Dashboard struct {
	Admin     string   `json:"admin"`
	Password  string   `json:"password"`
	Session   string   `json:"session"`
	TrustedIP []string `json:"trusted_ip"`
}

// Commands 命令过滤器表结构
type Commands struct {
	Args      []string `json:"args"`
	IsAllowed bool     `json:"is_allowed"`
	IsUnknown bool     `json:"is_unknown"`
}

// Logger 记录入侵者单次访问的使用记录
type Logger struct {
	IP        string   `json:"ip_address"`
	LastLogin string   `json:"last_login"`
	Record    []string `json:"usage_record"`
}

// BlackList 黑名单表结构
type BlackList struct {
	LastLogin string `json:"last_login"`
	Banned    bool   `json:"banned"`
}

// MyFiles 记录启动时整个文件系统的结构
type MyFiles struct {
	Filename   string `json:"filename"`
	FatherPath string `json:"father_path"`
	IsFolder   bool   `json:"is_folder"`
}
