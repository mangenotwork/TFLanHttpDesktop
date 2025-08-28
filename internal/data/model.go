package data

import "time"

// DownloadNow 当前下载的文件
type DownloadNow struct {
	Path       string // 文件路径
	IsPassword bool   // 是否设置密码
	Password   string // 密码
}

// DownloadLog 下载的日志 key是时间戳
type DownloadLog struct {
	Time      string // 请求的时间 时间戳
	IP        string // 请求端ip
	UserAgent string // 请求端的user-agent
	Path      string // 下载文件的路径
	Size      string // 下载的文件的信息
}

// UploadNow 当前上传的路径
type UploadNow struct {
	Path       string // 接收上传文件的路径
	IsPassword bool   // 是否设置密码
	Password   string // 密码
}

// UploadLog 上传的日志 key是时间戳
type UploadLog struct {
	Time      string // 请求的时间 时间戳
	IP        string // 请求端ip
	UserAgent string // 请求端的user-agent
	Path      string // 上传文件的路径
	Size      string // 上传的文件的信息
}

// Memo 备忘录  key是备忘录id
type Memo struct {
	Id         string
	Name       string
	CreateTime time.Time
	LastTime   time.Time
	Authority  int    // 客户端权限 0无权限 1只读 2可读写
	IsPassword bool   // 是否设置密码
	Password   string // 密码
}

// MemoContent 备忘录具体内容 key是备忘录id
type MemoContent string

// OperationLog 操作日志  key是时间戳
type OperationLog struct {
	Time  int64  // 操作时间 时间戳
	Event string // 操作事件
}
