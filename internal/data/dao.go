package data

// GetDownloadData 获取当前下载文件数据
func GetDownloadData() (*DownloadNow, error) {
	result := &DownloadNow{}
	err := DB.Get(DownloadNowTable, DownloadNowTableKey, &result)
	return result, err
}

// GetUploadData 获取当前上传文件数据
func GetUploadData() (*UploadNow, error) {
	result := &UploadNow{}
	err := DB.Get(UploadNowTable, UploadNowTableKey, &result)
	return result, err
}

// todo... 创建当前下载文件数据，不存在创建，存在更新
func SetDownloadData(value *DownloadNow) error {
	return DB.Set(DownloadNowTable, DownloadNowTableKey, &value)
}

// todo... 创建当前上传文件路径，不存在创建，存在更新

// todo... 删除当前下载文件

// todo... 删除当前上传文件

// todo... 记录下载日志

// todo... 记录上传日志

// todo... 查看下载日志

// todo... 查看上传日志

// todo... 记录操作日志

// todo... 查看操作日志

// todo... 创建备忘录

// todo... 修改备忘录内容

// todo... 查看备忘录信息

// todo... 查看备忘录内容

// todo... 删除备忘录
