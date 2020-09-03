package models

type BaseAttachment struct {
	Filename string `json:"filename"`
	Type     string `json:"type"`
	Mime     string `json:"mime"`
}

type Attachment struct {
	BaseAttachment
	ID        string `json:"id"`
	CreatedAt int64  `json:"created_at"`
}

type CreateAttachment struct {
	BaseAttachment
}

var EmptyAttachment = Attachment{}

func NewCreateAttachment(filename string, atype string, mime string) CreateAttachment {
	return CreateAttachment{
		BaseAttachment: BaseAttachment{
			Filename: filename,
			Type:     atype,
			Mime:     mime,
		},
	}
}
