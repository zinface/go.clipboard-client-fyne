package main

import (
	"context"
	"flag"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	log "github.com/sirupsen/logrus"
)

var (
	address = flag.String("addr", "http://localhost:9090", "服务器地址")
	verbose = flag.Bool("v", false, "打印详细日志")
)

func main() {
	//	初始化命令行参数
	flag.Parse()

	//	初始化日志
	if *verbose {
		log.SetLevel(log.TraceLevel)
	} else {
		log.SetLevel(log.DebugLevel)
	}

	//	应用
	a := app.New()
	a.Settings().SetTheme(&UnicodeTheme{})
	a.SetIcon(resourceIconSystray)

	//	窗口
	w := a.NewWindow("云剪贴板")
	//	窗口内容
	switcher := NewImageTextSwitcher(&w)
	switcher.Reset()
	//	设置窗口内容
	w.SetContent(switcher.Container)

	//	系统托盘图标
	if a, ok := a.(desktop.App); ok {
		menu := fyne.NewMenu("",
			fyne.NewMenuItem("显示窗口", w.Show))
		a.SetSystemTrayMenu(menu)
	}

	//	启动本地及云之剪切板监控
	sync := NewSync(*address)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go sync.WatchLocal(ctx, switcher.ShowText, switcher.ShowImage)
	go sync.WatchCloud(ctx, switcher.ShowText, switcher.ShowImage)

	//	运行显示
	w.ShowAndRun()
}
