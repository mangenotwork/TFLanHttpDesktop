package assets

import (
	_ "embed"
)

//go:embed download.html
var DownloadPg string

//go:embed memo.html
var MemoPg string

//go:embed tailwindcss.js
var TailwindcssData string

//go:embed upload.html
var UploadPg string
