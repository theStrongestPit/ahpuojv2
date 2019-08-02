package utils

import (
	"mime/multipart"
	"os"
	"testing"
)

func Test_ImportFps(t *testing.T) {
	var file multipart.File
	file, _ = os.Open("fps.xml")
	data, _ := ImportFps(file)
	t.Logf("%+v", data.Item[0])
}
