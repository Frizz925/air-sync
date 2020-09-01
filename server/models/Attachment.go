package models

type BaseAttachment struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Mime     string `json:"mime"`
}

type Attachment struct {
	BaseAttachment
	Id        string `json:"id"`
	CreatedAt int64  `gorm:"autoCreateTime"`
}

type CreateAttachment struct {
	BaseAttachment
}

var EmptyAttachment = Attachment{}
