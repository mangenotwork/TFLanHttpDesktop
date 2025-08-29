package define

const (
	DBFileDirName  = "TFLanHttpDesktop"
	DBFileFileName = "data.db"
	CsrfAuthKey    = "https://github.com/mangenotwork/TFLanHttpDesktop"
	CsrfName       = "TFLanHttpDesktop"
)

var (
	LanIP    = ""
	HttpPort int
	DoMain   = ""
)

var DownloadMem = make(map[string]string)
var UploadMem = make(map[string]string)

// 签名配置（实际应用中应从环境变量或配置文件读取）
const (
	SignSecretKey  = "TFLanHttpDesktop"
	SignSaltLength = 16             // 随机盐值长度（字节）
	SignExpiresIn  = 3600 * 24 * 30 // 签名有效期（秒），0表示永久有效
)
