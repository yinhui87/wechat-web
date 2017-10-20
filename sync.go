package wxweb

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/astaxie/beego/httplib"
	"github.com/yinhui87/wechat-web/datastruct"
	"github.com/yinhui87/wechat-web/datastruct/appmsg"
	"github.com/yinhui87/wechat-web/tool"
	"io/ioutil"
	"log"
	"net/url"
	"strconv"
	"strings"
	"time"
)

// assembleSyncKey 组装synckey
// 将同步需要的synckey组装为请求字符串
func assembleSyncKey(syncKey *datastruct.SyncKey) string {
	keys := make([]string, 0)
	for _, v := range syncKey.List {
		keys = append(keys, strconv.FormatInt(v.Key, 10)+"_"+strconv.FormatInt(v.Val, 10))
	}
	ret := strings.Join(keys, "|")
	return url.QueryEscape(ret)
}

// analysisSyncResp 解析同步状态返回值
// 同步状态返回的接口
func analysisSyncResp(syncResp string) (result datastruct.SyncCheckRespond) {
	syncResp = strings.TrimPrefix(syncResp, "{")
	syncResp = strings.TrimSuffix(syncResp, "}")
	arr := strings.Split(syncResp, ",")
	result = datastruct.SyncCheckRespond{}
	for _, v := range arr {
		if strings.HasPrefix(v, "retcode") {
			result.Retcode = strings.TrimPrefix(strings.TrimSuffix(v, `"`), `retcode:"`)
		}
		if strings.HasPrefix(v, "selector") {
			result.Selector = strings.TrimPrefix(strings.TrimSuffix(v, `"`), `selector:"`)
		}
	}
	return result
}

// syncCheck 同步状态
// 轮询微信服务器，如果有新的状态，会通过此接口返回需要同步的信息
func (wxwb *WechatWeb) syncCheck() (selector string, err error) {
	// req := httplib.Get("https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck")
	// req.Param("r", tool.GetWxTimeStamp())
	// req.Param("skey", wxwb.sKey)
	// req.Param("sid", wxwb.cookie.Wxsid)
	// req.Param("uin", wxwb.cookie.Wxuin)
	// req.Param("deviceid", wxwb.deviceID)
	// req.Param("synckey", assembleSyncKey(wxwb.syncKey))
	// req.Param("_", tool.GetWxTimeStamp())
	// setWechatCookie(req, wxwb.cookie)
	// resp, err := req.String()
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	params.Set("skey", wxwb.sKey)
	params.Set("sid", wxwb.cookie.Wxsid)
	params.Set("uin", wxwb.cookie.Wxuin)
	params.Set("deviceid", wxwb.deviceID)
	params.Set("synckey", assembleSyncKey(wxwb.syncKey))
	params.Set("_", tool.GetWxTimeStamp())
	resp, err := wxwb.client.Get("https://webpush.wx2.qq.com/cgi-bin/mmwebwx-bin/synccheck?" + params.Encode())
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	retArr := tool.AnalysisWxWindowRespond(string(body))

	ret := analysisSyncResp(retArr["window.synccheck"])
	if ret.Retcode != "0" {
		if ret.Retcode == "1101" {
			return "Logout", nil
		}
		return "", errors.New("respond Retcode " + ret.Retcode)
	}
	return ret.Selector, nil
}

// getMessage 同步消息
// 如果同步状态接口返回有新消息需要同步，通过此接口从服务器中获取新消息
func (wxwb *WechatWeb) getMessage() (gmResp datastruct.GetMessageRespond, err error) {
	// req := httplib.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync")
	// req.Param("sid", wxwb.cookie.Wxsid)
	// req.Param("skey", wxwb.sKey)
	// req.Param("pass_ticket", wxwb.PassTicket)
	// setWechatCookie(req, wxwb.cookie)
	// req.Body(body)
	// resp, err := req.Bytes()
	gmResp = datastruct.GetMessageRespond{}
	data, err := json.Marshal(datastruct.GetMessageRequest{
		BaseRequest: wxwb.baseRequest,
		SyncKey:     wxwb.syncKey,
		Rr:          ^time.Now().Unix() + 1,
	})
	if err != nil {
		return gmResp, errors.New("Marshal request body to json fail: " + err.Error())
	}
	params := url.Values{}
	params.Set("sid", wxwb.cookie.Wxsid)
	params.Set("skey", wxwb.sKey)
	params.Set("pass_ticket", wxwb.PassTicket)
	resp, err := wxwb.client.Post("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxsync?"+params.Encode(),
		"application/json;charset=UTF-8",
		bytes.NewReader(data))
	if err != nil {
		return gmResp, errors.New("request error: " + err.Error())
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &gmResp)
	if err != nil {
		return gmResp, errors.New("Unmarshal respond json fail: " + err.Error())
	}
	if gmResp.BaseResponse.Ret != 0 {
		return gmResp, errors.New("respond error ret: " + strconv.FormatInt(gmResp.BaseResponse.Ret, 10))
	}
	// if gmResp.AddMsgCount > 0 {
	// 	fmt.Println(string(resp))
	// 	panic(nil)
	// }
	return gmResp, nil
}

// SaveMessageImage 保存消息图片到指定位置
func (wxwb *WechatWeb) SaveMessageImage(msg datastruct.Message) (filename string, err error) {
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetmsgimg")
	req.Param("MsgID", msg.MsgID)
	req.Param("skey", wxwb.sKey)
	// req.Param("type", "slave")
	setWechatCookie(req, wxwb.cookie)
	filename = msg.MsgID + ".jpg"
	err = req.ToFile(filename)
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	return filename, nil
}

// SaveMessageVoice 保存消息声音到指定位置
func (wxwb *WechatWeb) SaveMessageVoice(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.VoiceMsg {
		return "", errors.New("Message type wrong")
	}
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvoice")
	req.Param("msgid", msg.MsgID)
	req.Param("skey", wxwb.sKey)
	setWechatCookie(req, wxwb.cookie)
	filename = msg.MsgID + ".mp3"
	err = req.ToFile(filename)
	if err != nil {

		return "", errors.New("request error: " + err.Error())
	}
	return filename, nil
}

// SaveMessageVideo 保存消息视频到指定位置
func (wxwb *WechatWeb) SaveMessageVideo(msg datastruct.Message) (filename string, err error) {
	if msg.MsgType != datastruct.LittleVideoMsg {
		return "", errors.New("Message type wrong")
	}
	var videoContent appmsg.VideoMsgContent
	err = xml.Unmarshal([]byte(msg.Content), &videoContent)
	if err != nil {
		return "", errors.New("Unmarshal message content to struct: " + err.Error())
	}
	req := httplib.Get("https://wx2.qq.com/cgi-bin/mmwebwx-bin/webwxgetvideo")
	req.Param("msgid", msg.MsgID)
	req.Param("skey", wxwb.sKey)
	setWechatCookie(req, wxwb.cookie)
	req.Header("Range", "bytes=0-")
	filename = msg.MsgID + ".mp4"
	// err = req.ToFile(filename)
	resp, err := req.Response()
	if err != nil {
		return "", errors.New("request error: " + err.Error())
	}
	length, err := strconv.ParseInt(videoContent.VideoMsg.Length, 10, 64)
	if err != nil {
		return "", errors.New("Parse Video Content Length error: " + err.Error())
	}
	if resp.ContentLength != length {
		return "", errors.New("Respond content length wrong")
	}
	n, err := tool.WriteToFile(filename, resp.Body)
	if err != nil {
		return "", errors.New("WriteToFile error: " + err.Error())
	}
	if int64(n) != length {
		return filename, errors.New("File size wrong")
	}
	return filename, nil
}

// StartServe 启动消息同步服务
func (wxwb *WechatWeb) StartServe() {
Serve:
	for true {
		selector, err := wxwb.syncCheck()
		if err != nil {
			log.Printf("SyncCheck error: %s\n", err.Error())
			continue
		}
		// log.Println("selector: ", selector)
		switch selector {
		case "7":
			// log.Println("SyncCheck 7")
			gmResp, err := wxwb.getMessage()
			if err != nil {
				log.Printf("GetMessage error: %s\n", err.Error())
				continue
			}
			wxwb.syncKey = gmResp.SyncKey
			for _, msg := range gmResp.AddMsgList {
				err = wxwb.messageProcesser(msg)
				if err != nil {
					log.Printf("MessageProcesser error: %s\n", err.Error())
					continue
				}
			}
		case "Logout":
			log.Println("User has logout web wechat, exit...")
			break Serve
		default:
			log.Printf("SyncCheck Unknow selector: %s\n", selector)
		}
		time.Sleep(100 * time.Microsecond)
	}
}
