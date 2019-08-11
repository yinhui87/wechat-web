package appmsg

import (
	"encoding/xml"
)

// ImageMsgContentImg 图片消息的content中的img节点
type ImageMsgContentImg struct {
	ImgName        xml.Name `xml:"img"`
	AesKey         string   `xml:"aeskey,attr"`
	EncryVer       string   `xml:"encryver,attr"`
	CdnThumbAesKey string   `xml:"cdnthumbaeskey,attr"`
	CdnThumbURL    string   `xml:"cdnthumburl,attr"`
	CdnThumbLength string   `xml:"cdnthumblength,attr"`
	CdnThumbHeight string   `xml:"cdnthumbheight,attr"`
	CdnThumbWidth  string   `xml:"cdnthumbwidth,attr"`
	CdnMidHeight   string   `xml:"cdnmidheight,attr"`
	CdnMidWidth    string   `xml:"cdnmidwidth,attr"`
	CdnHdHeight    string   `xml:"cdnhdheight,attr"`
	CdnHdWidth     string   `xml:"cdnhdwidth,attr"`
	CdnMidImgURL   string   `xml:"cdnmidimgurl,attr"`
	Length         string   `xml:"length,attr"`
	Md5            string   `xml:"md5,attr"`
}

// ImageMsgContent 图片消息的content
type ImageMsgContent struct {
	MsgName xml.Name            `xml:"msg"`
	Img     *ImageMsgContentImg `xml:"img"`
}
