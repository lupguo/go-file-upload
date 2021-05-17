package app

var extMIME = map[string]string{
	"image/jpeg":               ".jpg",
	"image/gif":                ".gif",
	"image/png":                ".png",
	"image/bmp":                ".bmp",
	"image/webp":               ".webp",
	"image/x-icon":             ".ico",
	"image/vnd.microsoft.icon": ".ico",
}

// ExtByMIME 通过mime获取图片的后缀名，返回以.为前缀的扩展，如.jpg
func ExtByMIME(mime string) string {
	if ext, ok := extMIME[mime]; ok {
		return ext
	}
	return ".omg"
}

// IsAllowedUpload 是否允许上传的MIMEType
func IsAllowedUpload(mime string) bool {
	if _, ok := extMIME[mime]; ok {
		return true
	}
	return false
}
