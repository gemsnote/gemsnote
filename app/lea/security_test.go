package lea

import (
	"bytes"
	"image"
	"image/png"
	"testing"
	"time"
)

func TestRasterImageContentValidation(t *testing.T) {
	var pngData bytes.Buffer
	if err := png.Encode(&pngData, image.NewRGBA(image.Rect(0, 0, 2, 2))); err != nil {
		t.Fatal(err)
	}
	if ext, ok := RasterImageExtension(pngData.Bytes()); !ok || ext != ".png" {
		t.Fatal("valid PNG rejected")
	}
	for _, data := range []string{"<html><script>alert(1)</script>", `<svg xmlns="http://www.w3.org/2000/svg" onload="alert(1)"/>`, "not an image", ""} {
		if _, ok := RasterImageExtension([]byte(data)); ok {
			t.Fatal("active/non-image content accepted")
		}
	}
}

func TestPDFTokenIsScopedAndExpires(t *testing.T) {
	now := time.Unix(10000, 0)
	token := PDFToken("secret", "note-a", now)
	if !ValidPDFToken("secret", "note-a", token, now) {
		t.Fatal("valid token rejected")
	}
	if ValidPDFToken("secret", "note-b", token, now) || ValidPDFToken("other", "note-a", token, now) || ValidPDFToken("secret", "note-a", "secret", now) || ValidPDFToken("secret", "note-a", token, now.Add(5*time.Minute)) || ValidPDFToken("", "note-a", token, now) {
		t.Fatal("unscoped or expired PDF capability accepted")
	}
}
