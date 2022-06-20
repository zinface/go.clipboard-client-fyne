package main

import (
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	log "github.com/sirupsen/logrus"
)

var (
	MINIMUM_WINDOW_SIZE = fyne.NewSize(400, 300)
	MAXIMUM_WINDOW_SIZE = fyne.NewSize(1000, 800)
	MINIMUM_IMAGE_SIZE  = fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
)

func MaxSize(s1, s2 fyne.Size) fyne.Size {
	return fyne.NewSize(
		fyne.Max(s1.Width, s2.Width),
		fyne.Max(s1.Height, s2.Height))
}

func MinSize(s1, s2 fyne.Size) fyne.Size {
	return fyne.NewSize(
		fyne.Min(s1.Width, s2.Width),
		fyne.Min(s1.Height, s2.Height))
}

func ClipSize(s, minimum, maximum fyne.Size) fyne.Size {
	size := MaxSize(s, minimum)
	size = MinSize(size, maximum)
	return size
}

type ImageTextSwitcher struct {
	Window    *fyne.Window    //	记录所属窗口，用以刷新窗口大小
	Container *fyne.Container //	包含图片和文本的上层容器，用以获得最小尺寸
	Image     *canvas.Image
	Text      *widget.Label
	mu        sync.RWMutex
}

func NewImageTextSwitcher(w *fyne.Window) *ImageTextSwitcher {
	s := ImageTextSwitcher{
		Window: w, //	记录所属窗口
		Text:   widget.NewLabel(""),
	}

	//	生成图像控件
	var err error
	s.Image, err = PNG(icon_systray).ToCanvasImage()
	if err != nil {
		log.Error(err)
		return nil
	}

	//	生成包含容器
	// s.Container = container.NewCenter(
	// 	container.NewHBox(s.Image, s.Text))
	s.Container = container.NewCenter(s.Image, s.Text)
	return &s
}

func (c *ImageTextSwitcher) Reset() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	//	默认图片
	c.Image.Image, err = PNG(icon_systray).ToImage()
	if err != nil {
		return err
	}
	iconSize := fyne.NewSize(theme.IconInlineSize(), theme.IconInlineSize())
	c.Image.SetMinSize(iconSize)
	//	默认文字
	c.Text.SetText("云剪切板")
	//	同时显示图片和文本
	c.Image.Hide()
	c.Text.Show()
	//	刷新显示
	c.Refresh()

	return nil
}

func (c *ImageTextSwitcher) ShowText(content string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	//	设置文本标签
	c.Text.SetText(content)
	//	交换图像/文本显隐开关
	c.Image.Hide()
	c.Text.Show()
	//	刷新显示
	c.Refresh()

	return nil
}

func (c *ImageTextSwitcher) ShowImage(png PNG) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	//	设置图片显示给入图像
	var err error
	c.Image.Image, err = png.ToImage()
	if err != nil {
		return err
	}
	//	设置图片最小尺寸
	imgSize := png.Size()
	size := ClipSize(imgSize, MINIMUM_IMAGE_SIZE, MAXIMUM_WINDOW_SIZE)
	c.Image.SetMinSize(size)
	//	交换图像/文本显隐开关
	c.Image.Show()
	c.Text.Hide()
	//	刷新显示
	c.Refresh()

	return nil
}

func (c *ImageTextSwitcher) Refresh() {
	if c.Window != nil {
		//	计算最小尺寸，并确保大于设定最小尺寸且小于设定最大尺寸
		size := ClipSize(c.Container.MinSize(), MINIMUM_WINDOW_SIZE, MAXIMUM_WINDOW_SIZE)
		//	设置窗口尺寸（否则窗口只会扩增，而不会缩小）
		(*c.Window).Resize(size)
		//	刷新各级容器以显示内容
		(*c.Window).Content().Refresh()
	}
}
