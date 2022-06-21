package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.design/x/clipboard"
)

type Sync struct {
	Clipboard Clipboard
	API       ClipboardAPI
	mu        sync.Mutex
}

func NewSync(address string) *Sync {
	//	初始化剪贴板
	if err := clipboard.Init(); err != nil {
		log.Fatal(err)
	}

	log.Infof("服务器地址: %s", address)

	return &Sync{
		API: ClipboardAPI{address},
	}
}

func (s *Sync) TextToCloud() (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//	将本地剪贴板的文本同步到云剪贴板
	content := string(clipboard.Read(clipboard.FmtText))
	if c, err := s.API.Set(content, "text/plain"); err != nil {
		return content, err
	} else {
		// 更新剪贴板的创建时间
		s.Clipboard.CreateAt = c.CreateAt
		log.Debugf("Sync: [客户端 => 服务器] %v", c)
		return content, nil
	}
}

func (s *Sync) ImageToCloud() (PNG, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//	将本地剪贴板的图片同步到云剪贴板
	png := PNG(clipboard.Read(clipboard.FmtImage))
	if c, err := s.API.Set(png.ToBase64(), "image/png"); err != nil {
		return png, err
	} else {
		// 更新剪贴板的创建时间
		s.Clipboard.CreateAt = c.CreateAt
		log.Debugf("Sync: [客户端 => 服务器] %v", c)
		return png, nil
	}
}

func (s *Sync) FromCloud() ([]byte, string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	//	从云剪贴板获取内容
	if remote, err := s.API.Get(); err != nil {
		return nil, "", err
	} else {
		s.Clipboard = remote
		log.Debugf("Sync: [客户端 <= 服务器]: %v", s.Clipboard)
		switch s.Clipboard.Mime {
		case "text/plain":
			clipboard.Write(clipboard.FmtText, []byte(s.Clipboard.Data))
			return []byte(s.Clipboard.Data), "text/plain", nil
		case "image/png":
			png, err := NewPNGFromBase64(s.Clipboard.Data)
			if err != nil {
				return png, s.Clipboard.Mime, err
			} else {
				clipboard.Write(clipboard.FmtImage, []byte(png))
				return png, s.Clipboard.Mime, nil
			}
		default:
			return nil, s.Clipboard.Mime, fmt.Errorf("不支持的MIME类型: %s", s.Clipboard.Mime)
		}
	}
}

type OnTextChanged func(string) error
type OnImageChanged func(PNG) error

func (s *Sync) WatchLocal(ctx context.Context, handleText OnTextChanged, handleImage OnImageChanged) {
	//	监听本地剪贴板变化
	chText := clipboard.Watch(ctx, clipboard.FmtText)
	chImage := clipboard.Watch(ctx, clipboard.FmtImage)

	for {
		select {
		case <-chText:
			//	将本地剪贴板的文本同步到云剪贴板
			content, err := s.TextToCloud()
			if err != nil {
				log.Error(err)
			}
			handleText(content)
		case <-chImage:
			//	将本地剪贴板的图像进行Base64编码后同步到云剪贴板
			content, err := s.ImageToCloud()
			if err != nil {
				log.Error(err)
			}
			handleImage(content)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Sync) handleFromCloud(handleText OnTextChanged, handleImage OnImageChanged) {
	//	从云剪切板获取内容
	data, mime, err := s.FromCloud()
	if err != nil {
		log.Error(err)
		return
	}
	//	通知剪贴板变化
	switch mime {
	case "text/plain":
		if err := handleText(string(data)); err != nil {
			log.Error(err)
		}
	case "image/png":
		if err := handleImage(PNG(data)); err != nil {
			log.Error(err)
		}
	}
}

func (s *Sync) WatchCloud(ctx context.Context, handleText OnTextChanged, handleImage OnImageChanged) {
	//	从云剪切板获取内容
	s.handleFromCloud(handleText, handleImage)
	//	监听云剪贴板变化
	for {
		select {
		case <-time.Tick(3 * time.Second):
			//	首先从服务器上检查最新的剪切板内容
			latest, err := s.API.Info()
			if err != nil {
				log.Error(err)
			}
			if latest.CreateAt != s.Clipboard.CreateAt {
				//	云剪贴板时间有所变动，需更新本地剪切板内容
				s.handleFromCloud(handleText, handleImage)
			}
		case <-ctx.Done():
			return
		}
	}
}
