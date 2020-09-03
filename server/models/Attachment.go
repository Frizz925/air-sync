package models

type BaseAttachment struct {
	Type string `json:"type"`
	Mime string `json:"mime"`
	Name string `json:"name"`
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

func NewCreateAttachment(name string, typ string, mime string) CreateAttachment {
	return CreateAttachment{
		BaseAttachment: BaseAttachment{
			Name: name,
			Type: typ,
			Mime: mime,
		},
	}
}
