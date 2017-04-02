package main

import (
	"fmt"
	"github.com/yinhui87/wechat-web"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
	"log"
	"time"
)

func main() {
	wx := wxweb.NewWechatWeb()
	t := testServ{}
	err := wx.RegisterMessageHook(wxweb.TextMessageHook(t.ProcessTextMessage))
	if err != nil {
		panic("RegisterMessageHook TextMessageHook: " + err.Error())
	}
	err = wx.RegisterMessageHook(wxweb.ImageMessageHook(t.ProcessImageMessage))
	if err != nil {
		panic("RegisterMessageHook ImageMessageHook: " + err.Error())
	}
	err = wx.Login()
	if err != nil {
		panic("WxWeb Login error: " + err.Error())
	}
	wx.StartServe()
}

type testServ struct {
}

func (this *testServ) ProcessTextMessage(ctx *wxweb.Context, msg datastruct.Message) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived a text msg from %s: %s", from.NickName, msg.Content)
	// reply the same message
	smResp, err := ctx.App.SendTextMessage(msg.FromUserName, msg.Content)
	if err != nil {
		log.Println("sendTextMessage error: " + err.Error())
	}
	log.Println("messageSent, msgId: " + smResp.MsgID + ", Local ID: " + smResp.LocalID)
	// Set message to readed at phone
	err = ctx.App.StatusNotify(msg.ToUserName, msg.FromUserName)
	if err != nil {
		log.Println("StatusNotify error: " + err.Error())
	}
	go func() {
		time.Sleep(10 * time.Second)
		ctx.App.SendRevokeMessage(smResp.MsgID, smResp.LocalID, msg.FromUserName)
	}()
}

func (this *testServ) ProcessImageMessage(ctx *wxweb.Context, msg datastruct.Message, imgContent appmsg.ImageMsgContent) {
	from, err := ctx.App.GetContact(msg.FromUserName)
	if err != nil {
		log.Println("getContact error: " + err.Error())
	}
	log.Printf("Recived a image msg from %s", from.NickName)
	fmt.Println("aeskey: ", imgContent.Img.AesKey)
}
