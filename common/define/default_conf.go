package define

const (
	DBFileDirName  = "TFLanHttpDesktop"
	DBFileFileName = "data.db"
)

var (
	LanIP    = ""
	HttpPort int
	DoMain   = ""
)

var DownloadMem = make(map[string]string)
