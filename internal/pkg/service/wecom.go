package service

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/walkerdu/wecom-backend/pkg/chatbot"
	"github.com/walkerdu/wecom-backend/pkg/wecom"
	"github.com/walkerdu/wecom-read-it-later/configs"
	"github.com/walkerdu/wecom-read-it-later/internal/pkg/handler"
)

type WeComServer struct {
	httpSvr *http.Server
	wc      *wecom.WeCom
}

func NewWeComServer(config *configs.WeComConfig) (*WeComServer, error) {
	log.Printf("[INFO] NewWeComServer")

	svr := &WeComServer{}

	// 初始化企业微信API
	svr.wc = wecom.NewWeCom(&config.AgentConfig)

	mux := http.NewServeMux()
	mux.Handle("/wecom", svr.wc)

	svr.httpSvr = &http.Server{
		Addr:    config.Addr,
		Handler: mux,
	}

	svr.initHandler()

	// 注册聊天消息的异步推送回调
	chatbot.MustChatbot().RegsiterMessagePublish(svr.wc.PushTextMessage)

	return svr, nil
}

// 注册企业微信消息处理的业务逻辑Handler
func (svr *WeComServer) initHandler() error {
	for msgType, handler := range handler.HandlerInst().GetLogicHandlerMap() {
		svr.wc.RegisterLogicMsgHandler(msgType, handler.HandleMessage)
	}

	return nil
}

func (svr *WeComServer) Serve() error {
	log.Printf("[INFO] Server()")

	if err := svr.httpSvr.ListenAndServe(); nil != err {
		log.Fatalf("httpSvr ListenAndServe() failed, err=%s", err)
		return err
	}

	return nil
}

func (svr *WeComServer) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := svr.httpSvr.Shutdown(ctx); err != nil {
		log.Printf("httpSvr ListenAndServe() failed, err=%s", err)
		return err
	}

	log.Println("[INFO]close httpSvr success")
	return nil
}
