package define

import "sync"

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
var ReqToken sync.Map
