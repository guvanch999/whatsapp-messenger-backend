package enums

type Provider string

const (
	Twilio Provider = "twilio"
	Plivo  Provider = "plivo"
)

type Platform string

const (
	WhatsApp Platform = "WhatsApp"
	Sms      Platform = "sms"
	Email    Platform = "email"
)
