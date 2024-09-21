package dto

type MessageDetailDto struct {
	PhoneNumber       string
	FromPhoneNumber   string
	ServiceId         string
	TemplateId        string
	TemplateVariables interface{}
}
