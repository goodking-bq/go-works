// Code generated by go-bindata. DO NOT EDIT.
// source: ui/dist/css/chunk-vendors.6dc57fe0.css
package bindata


import (
	"os"
	"time"
)

var _bindataUiDistCssChunkvendors6dc57fe0css = []byte(
	"\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xd2\x4b\x4d\xce\x48\x2c\x2a\x29\xae\x2e\xcf\x4c\x29\xc9\xb0\x32\x33\x30" +
	"\x28\xa8\xb0\xce\x48\xcd\x4c\xcf\x28\xb1\x32\x01\x71\x6a\x01\x01\x00\x00\xff\xff\x67\x7a\x4a\x64\x22\x00\x00\x00" +
	"")

func bindataUiDistCssChunkvendors6dc57fe0cssBytes() ([]byte, error) {
	return bindataRead(
		_bindataUiDistCssChunkvendors6dc57fe0css,
		"ui/dist/css/chunk-vendors.6dc57fe0.css",
	)
}



func bindataUiDistCssChunkvendors6dc57fe0css() (*asset, error) {
	bytes, err := bindataUiDistCssChunkvendors6dc57fe0cssBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{
		name: "ui/dist/css/chunk-vendors.6dc57fe0.css",
		size: 34,
		md5checksum: "",
		mode: os.FileMode(420),
		modTime: time.Unix(1588942021, 0),
	}

	a := &asset{bytes: bytes, info: info}

	return a, nil
}

