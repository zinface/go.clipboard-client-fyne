package main

import (
	"bytes"
	"encoding/base64"
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	log "github.com/sirupsen/logrus"
)

type PNG []byte

func NewPNGFromBase64(s string) (PNG, error) {
	if b, err := base64.StdEncoding.DecodeString(s); err != nil {
		return nil, err
	} else {
		return b, nil
	}
}

func (p PNG) ToBase64() string {
	return base64.StdEncoding.EncodeToString(p)
}

func (p PNG) ToImage() (image.Image, error) {
	if i, _, err := image.Decode(bytes.NewReader(p)); err != nil {
		return nil, err
	} else {
		return i, nil
	}
}

func (p PNG) ToCanvasImage() (*canvas.Image, error) {
	if i, err := p.ToImage(); err != nil {
		return nil, err
	} else {
		img := canvas.NewImageFromImage(i)
		img.FillMode = canvas.ImageFillContain
		//	由于图片默认大小为0x0，所以这里手动计算最小尺寸
		// TODO: 在 macOS 中如果调整显示器缩放，将导致这里图片尺寸成比例变化，应该寻找自动的方法。
		img.SetMinSize(getImageSize(i))

		return img, nil
	}
}

func getImageSize(img image.Image) fyne.Size {
	intSize := img.Bounds().Size()
	imgSize := fyne.NewSize(float32(intSize.X), float32(intSize.Y))
	return imgSize
}

func (p PNG) Size() fyne.Size {
	if image, err := p.ToImage(); err != nil {
		log.Error(err)
		return fyne.Size{}
	} else {
		return getImageSize(image)
	}
}
