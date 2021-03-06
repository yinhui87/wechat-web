package api

import (
	"bytes"
	"encoding/json"
	"github.com/ikuiki/wwdk/datastruct"
	"github.com/ikuiki/wwdk/tool"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

// GetContact 获取联系人
// 注：坑！此处获取到的居然不是完整的联系人，必须和init中获取到的合并后才是完整的联系人列表
// @return contact 联系人列表（需要与wxInit获得的列表合并才是完整联系人列表）
func (api *wechatwebAPI) GetContact() (contactList []datastruct.Contact, body []byte, err error) {
	params := url.Values{}
	params.Set("r", tool.GetWxTimeStamp())
	resp, err := api.client.Get("https://" + api.apiDomain + "/cgi-bin/mmwebwx-bin/webwxgetcontact?" + params.Encode())
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	respStruct := datastruct.GetContactRespond{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		err = errors.New("respond json Unmarshal to struct fail: " + err.Error())
		return
	}
	if respStruct.BaseResponse.Ret != 0 {
		err = errors.Errorf("respond ret error(%d): %s", respStruct.BaseResponse.Ret, string(body))
		return
	}
	contactList = respStruct.MemberList
	return
}

// BatchGetContact 获取群聊的成员
// 在需要获取群成员时，或者是通过群成员UserName与群ChatRoomID获取群成员信息的时候可以调用此方法
// @Param contactItemList 要获取的联系人的信息，获取好友与群成员要填UserName，获取群成员中非好友的联系人的信息需要带群的ChatRoomID
// @return contactList 要获取的联系人的信息
func (api *wechatwebAPI) BatchGetContact(contactItemList []datastruct.BatchGetContactRequestListItem) (contactList []datastruct.Contact, body []byte, err error) {
	dataStruct := datastruct.BatchGetContactRequest{
		BaseRequest: api.baseRequest(),
		List:        contactItemList,
		Count:       int64(len(contactItemList)),
	}
	if dataStruct.Count == 0 {
		return
	}
	reqBody, err := json.Marshal(dataStruct)
	if err != nil {
		err = errors.New("json.Marshal error: " + err.Error())
		return
	}
	params := url.Values{}
	params.Set("type", "ex")
	params.Set("r", tool.GetWxTimeStamp())
	resp, err := api.client.Post("https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxbatchgetcontact?"+params.Encode(),
		"application/json;charset=UTF-8",
		bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ = ioutil.ReadAll(resp.Body)
	respStruct := datastruct.BatchGetContactResponse{}
	err = json.Unmarshal(body, &respStruct)
	if err != nil {
		err = errors.New("respond json Unmarshal to struct fail: " + err.Error())
		return
	}
	if respStruct.BaseResponse.Ret != 0 {
		err = errors.Errorf("respond ret error(%d): %s", respStruct.BaseResponse.Ret, string(body))
		return
	}
	contactList = respStruct.ContactList
	return
}

// ModifyUserRemakName 修改联系人备注
// @param userName 要修改的联系人的UserName
// @param remarkName 新的备注
func (api *wechatwebAPI) ModifyUserRemakName(userName, remarkName string) (body []byte, err error) {
	murReq := datastruct.ModifyRemarkRequest{
		BaseRequest: api.baseRequest(),
		CmdID:       2,
		RemarkName:  remarkName,
		UserName:    userName,
	}
	reqBody, err := json.Marshal(murReq)
	if err != nil {
		err = errors.New("Marshal reqBody to json fail: " + err.Error())
		return
	}
	req, err := http.NewRequest("POST", "https://"+api.apiDomain+"/cgi-bin/mmwebwx-bin/webwxoplog", bytes.NewReader(reqBody))
	if err != nil {
		err = errors.New("create request error: " + err.Error())
		return
	}
	resp, err := api.request(req)
	if err != nil {
		err = errors.New("request error: " + err.Error())
		return
	}
	defer resp.Body.Close()
	var murResp datastruct.ModifyRemarkRespond
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		err = errors.New("read response body error: " + err.Error())
		return
	}
	err = json.Unmarshal(body, &murResp)
	if err != nil {
		err = errors.New("UnMarshal respond json fail: " + err.Error())
		return
	}
	if murResp.BaseResponse.Ret != 0 {
		err = errors.Errorf("Respond error ret(%d): %s", murResp.BaseResponse.Ret, murResp.BaseResponse.ErrMsg)
		return
	}
	return
}
