package wxweb

import (
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/tool"
	"net/http"
)

type wechatCookie struct {
	Skey       string
	Wxsid      string
	Wxuin      string
	Uvid       string
	DataTicket string
	AuthTicket string
	PassTicket string
}

type WechatWeb struct {
	cookie      *wechatCookie
	userAgent   string
	deviceId    string
	contactList []*datastruct.Contact
	user        *datastruct.User
	syncKey     *datastruct.SyncKey
	sKey        string
	messageHook map[datastruct.MessageType][]interface{}
}

func NewWechatWeb() (wxweb WechatWeb) {
	return WechatWeb{
		userAgent:   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36",
		deviceId:    "e" + tool.GetRandomStringFromNum(15),
		messageHook: make(map[datastruct.MessageType][]interface{}),
	}
}

func setWechatCookie(request *httplib.BeegoHTTPRequest, cookie *wechatCookie) {
	request.SetCookie(&http.Cookie{Name: "wxsid", Value: cookie.Wxsid})
	request.SetCookie(&http.Cookie{Name: "webwx_data_ticket", Value: cookie.DataTicket})
	request.SetCookie(&http.Cookie{Name: "webwxuvid", Value: cookie.Uvid})
	request.SetCookie(&http.Cookie{Name: "webwx_auth_ticket", Value: cookie.AuthTicket})
	request.SetCookie(&http.Cookie{Name: "wxuin", Value: cookie.Wxuin})
}

func getBaseRequest(cookie *wechatCookie, deviceId string) (baseRequest *datastruct.BaseRequest) {
	return &datastruct.BaseRequest{
		Uin:      cookie.Wxuin,
		Sid:      cookie.Wxsid,
		Skey:     cookie.Skey,
		DeviceID: deviceId,
	}
}

func (this *WechatWeb) GetContact(username string) (contact *datastruct.Contact, err error) {
	for _, v := range this.contactList {
		if v.UserName == username {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

func (this *WechatWeb) GetContactByAlias(alias string) (contact *datastruct.Contact, err error) {
	for _, v := range this.contactList {
		if v.Alias == alias {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

func (this *WechatWeb) GetContactByNickname(nickname string) (contact *datastruct.Contact, err error) {
	for _, v := range this.contactList {
		if v.NickName == nickname {
			return v, nil
		}
	}
	return nil, errors.New("User not found")
}

func (this *WechatWeb) GetContactList() (contacts []*datastruct.Contact) {
	return this.contactList
}
