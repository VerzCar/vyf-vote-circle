package utils

var imageMimeTypeMap = map[string]struct{}{
	"image/apng":    {},
	"image/avif":    {},
	"image/gif":     {},
	"image/jpeg":    {},
	"image/png":     {},
	"image/svg+xml": {},
	"image/webp":    {},
}

func IsImageMimeType(mimeType string) bool {
	_, ok := imageMimeTypeMap[mimeType]
	return ok
}
