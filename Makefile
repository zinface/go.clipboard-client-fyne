appID = com.gitee.zinface.go.clipboard-client
icon = assets/icons/icon.png

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
	wget -nc -O ./assets/icons/icon.png https://gitee.com/zinface/qt.qgo-clipboard-client/raw/master/systray.png

clean:
	rm -rf ./clipboard-client
	rm -rf ./bin
	rm -rf ./fyne-cross
	rm -rf ./*.app

clean-all: clean
	rm -rf ./assets/*

tools:
	go install fyne.io/fyne/v2/cmd/fyne@latest
	go install github.com/fyne-io/fyne-cross@latest

package-arm64:
	# fyne package -os darwin --appID=$(appID) --icon=$(icon)
	fyne-cross darwin -arch=* -app-id $(appID) -icon $(icon)

package-amd64:
	fyne-cross linux -arch=amd64 -app-id $(appID) -icon $(icon)
	fyne-cross windows -arch=amd64 -app-id $(appID) -icon $(icon)
	fyne-cross android -app-id $(appID) -icon $(icon)
