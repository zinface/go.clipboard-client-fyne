package main

import (
	"flag"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var address = flag.String("addr", "http://localhost:9090", "服务器地址")

func main() {
	flag.Parse()

	//	应用
	a := app.New()
	a.Settings().SetTheme(&UnicodeTheme{})
	a.SetIcon(resourceIconSystray)

	//	窗口
	w := a.NewWindow("云剪贴板")
	w.Resize(fyne.NewSize(400, 300))

	//	窗口内容
	label := widget.NewLabel("云剪贴板!")
	img := canvas.NewImageFromResource(resourceIconSystray)
	img.FillMode = canvas.ImageFillOriginal
	h := container.New(layout.NewCenterLayout(), img, label)
	w.SetContent(h)

	//	系统托盘图标
	menu := fyne.NewMenu("",
		fyne.NewMenuItem("复制", func() {
			log.Println("TODO: 复制 ...")
		}),
		fyne.NewMenuItem("复制为Base64", func() {
			log.Println("TODO: 复制为Base64...")
		}),
		fyne.NewMenuItem("显示窗口", func() {
			log.Println("显示窗口")
			w.Show()
		}),
	)
	if a, ok := a.(desktop.App); ok {
		a.SetSystemTrayMenu(menu)
	}

	//	启动剪切板刷新
	api := ClipboardApi{Address: *address}
	log.Printf("服务器地址: %s\n", api.Address)

	tick := time.NewTicker(time.Second * 5)
	defer tick.Stop()
	done := make(chan struct{})
	defer close(done)

	var clipboard ClipBoard
	var err error

	if clipboard, err = api.Info(); err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-tick.C:
				//	首先从服务器上检查最新的剪切板内容
				var latest ClipBoard
				if latest, err = api.Info(); err != nil {
					log.Println(err)
				}
				if latest.CreateAt != clipboard.CreateAt {
					//	更新剪切板内容
					if clipboard, err = api.Get(); err != nil {
						log.Println(err)
					} else {
						log.Println("剪贴板：[客户端 <= 服务器]:", clipboard.Mime, clipboard.CreateAt)
						switch clipboard.Mime {
						case "text/plain":
							w.Clipboard().SetContent(clipboard.Data)
						case "image/png":
							//	TODO: 图片
						}
					}
				}
				//	如果本地剪切板内容有所更新
				content := w.Clipboard().Content()
				if content != clipboard.Data {
					//	更新服务器上的剪切板内容
					//	TODO: 图片
					if c, err := api.Set(content, "text/plain"); err != nil {
						log.Println(err)
					} else {

						log.Println("剪贴板：[客户端 => 服务器]:", c.Mime, c.CreateAt)
					}
				}
			case <-done:
				return
			}
		}
	}()

	//	运行显示
	w.ShowAndRun()
}
