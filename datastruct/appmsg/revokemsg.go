package appmsg

import (
	"encoding/xml"
)

// RevokeMsgContentRevoke 撤回消息的content的Revoke节点
type RevokeMsgContentRevoke struct {
	RevokeMsg  xml.Name `xml:"revokemsg"`
	Session    string   `xml:"session"`
	OldMsgID   string   `xml:"oldmsgid"`
	MsgID      string   `xml:"msgid"`
	ReplaceMsg string   `xml:"replacemsg"`
}

// RevokeMsgContent 撤回消息的content
type RevokeMsgContent struct {
	SysMsg    xml.Name                `xml:"sysmsg"`
	RevokeMsg *RevokeMsgContentRevoke `xml:"revokemsg"`
}
