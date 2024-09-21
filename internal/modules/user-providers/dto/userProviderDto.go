package dto

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/medium-messenger/messenger-backend/utils/enums"
)

type UserProviderDto struct {
	Name        string         `json:"name" validate:"required,gt=0"`
	Type        enums.Provider `json:"type" validate:"required,oneof=twilio plivo"`
	Credentials interface{}    `json:"credentials" validate:"required"`
}

type TwilioCredDto struct {
	TwilioAccountSid          string `json:"twilio_account_sid" validate:"required,gt=0"`
	TwilioAuthToken           string `json:"twilio_auth_token" validate:"required,gt=0"`
	TwilioMessagingServiceSid string `json:"twilio_messaging_service_sid" validate:"required,gt=0"`
	TwilioFromPhoneNumber     string `json:"twilio_from_phone_number" validate:"required,e164"`
}

func GetCredFromDto[T any](upDto UserProviderDto) (*T, error) {
	jsonByte, err := json.Marshal(upDto.Credentials)
	if err != nil {
		return nil, fmt.Errorf("error on provider credentials: %s", err.Error())
	}
	var detail T
	if err := json.Unmarshal(jsonByte, &detail); err != nil {
		return nil, fmt.Errorf("error on provider credentials: %s", err.Error())
	}
	return &detail, nil
}

func GetCredFromBytes[T any](dataByte []byte) (*T, error) {
	var detail T
	if err := json.Unmarshal(dataByte, &detail); err != nil {
		return nil, fmt.Errorf("error on provider credentials: %s", err.Error())
	}
	return &detail, nil
}
