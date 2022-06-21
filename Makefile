run:
	go run .

build:
	go mod tidy
	go build -o ./bin/clipboard-client .

fonts:
	mkdir -p ./assets/fonts
	wget -nc -P ./assets/fonts/ https://github.com/anthonyfok/fonts-wqy-microhei/raw/master/wqy-microhei.ttc

icons:
	mkdir -p ./assets/icons
	wget -nc -P ./assets/icons/ https://gitee.com/zinface/qt.qgo-clipboard-client/raw/master/systray.png

clean:
	rm -rf ./clipboard-client

clean-all: clean
	rm -rf ./assets/*
