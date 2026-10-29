package lea

import (
	"bytes"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
)

// Never trust a user-supplied image suffix or MIME type. SVG/HTML must not be
// served as an image from the application's authenticated origin.
func RasterImageExtension(data []byte) (string, bool) {
	// BMP/WebP are served as opaque raster files; their headers are recognized
	// by the standard library without adding a server-side image decoder.
	switch http.DetectContentType(data) {
	case "image/bmp":
		return ".bmp", true
	case "image/webp":
		return ".webp", true
	}
	config, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil || config.Width <= 0 || config.Height <= 0 || int64(config.Width)*int64(config.Height) > 100_000_000 {
		return "", false
	}
	ext := map[string]string{"gif": ".gif", "jpeg": ".jpg", "png": ".png", "bmp": ".bmp", "webp": ".webp"}[format]
	return ext, ext != ""
}
