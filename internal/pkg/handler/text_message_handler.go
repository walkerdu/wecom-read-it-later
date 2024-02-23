package handler

import (
	"context"
	"log"
	"strings"
	"time"

	//"github.com/walkerdu/wecom-backend/pkg/chatbot"
	"github.com/walkerdu/wecom-backend/pkg/wecom"
)

const WeChatTimeOutSecs = 5

func init() {
	handler := &TextMessageHandler{}

	HandlerInst().RegisterLogicHandler(wecom.MessageTypeText, handler)
}

type TextMessageHandler struct {
}

var commandsMap = map[string]struct{}{
	"/day":       struct{}{},
	"/yesterday": struct{}{},
	"/week":      struct{}{},
	"/month":     struct{}{},
}

func (t *TextMessageHandler) GetHandlerType() wecom.MessageType {
	return wecom.MessageTypeText
}

func (t *TextMessageHandler) HandleMessage(msg wecom.MessageIF) (wecom.MessageIF, error) {
	textMsg := msg.(*wecom.TextMessageReq)

	var chatRsp string

	if !strings.HasPrefix(textMsg.Content, "/") {
		_, err := t.DBSet(textMsg)
		if err != nil {
			chatRsp = err.Error()
		} else {
			chatRsp = "success"
		}
	} else {
	}

	//chatRsp, err := chatbot.MustChatbot().GetResponse(textMsg.FromUserName, textMsg.Content)
	//if err != nil {
	//	log.Printf("[ERROR][HandleMessage] chatbot.GetResponse failed, err=%s", err)
	//	chatRsp = "chatbot something wrong, errMsg:" + err.Error()
	//}

	textMsgRsp := wecom.TextMessageRsp{
		Content: chatRsp,
	}

	return &textMsgRsp, nil
}

func (t *TextMessageHandler) DBSet(msg *wecom.TextMessageReq) (*wecom.TextMessageRsp, error) {
	ctx := context.Background()
	date := time.Now().Format("20060102")
	key := msg.FromUserName + "_" + date
	result, err := HandlerInst().redisClient.LPush(ctx, key, msg.Content).Result()
	if err != nil {
		log.Printf("[ERROR][DBSet] redis LPush failed, err=%s", err)
		return nil, err
	}

	log.Printf("[DEBUG][DBSet] redis LPush success, key:%v, value:%v, result=%v", key, msg.Content, result)
	return nil, nil
}
