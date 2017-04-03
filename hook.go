package wxweb

import (
	"errors"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
)

type TextMessageHook func(*Context, datastruct.Message)

type ImageMessageHook func(*Context, datastruct.Message, appmsg.ImageMsgContent)

type EmotionMessageHook func(*Context, datastruct.Message, appmsg.EmotionMsgContent)

type RevokeMessageHook func(*Context, datastruct.Message, appmsg.RevokeMsgContent)

type VideoMessageHook func(*Context, datastruct.Message, appmsg.VideoMsgContent)

func (this *WechatWeb) RegisterMessageHook(hook interface{}) error {
	switch hook.(type) {
	case TextMessageHook:
		this.messageHook[datastruct.TEXT_MSG] = append(this.messageHook[datastruct.TEXT_MSG], hook)
	case ImageMessageHook:
		this.messageHook[datastruct.IMAGE_MSG] = append(this.messageHook[datastruct.IMAGE_MSG], hook)
	case EmotionMessageHook:
		this.messageHook[datastruct.ANIMATION_EMOTIONS_MSG] = append(this.messageHook[datastruct.ANIMATION_EMOTIONS_MSG], hook)
	case RevokeMessageHook:
		this.messageHook[datastruct.REVOKE_MSG] = append(this.messageHook[datastruct.REVOKE_MSG], hook)
	case VideoMessageHook:
		this.messageHook[datastruct.LITTLE_VIDEO_MSG] = append(this.messageHook[datastruct.LITTLE_VIDEO_MSG], hook)
	default:
		return errors.New("Unknown hook function")
	}
	return nil
}
