package util

import (
	"github.com/goccy/go-json"
)

type ServiceAccount struct {
	ProjectID string `json:"project_id"`
}

func GetProjectIdFromGCred(credJson string) (string, error) {
	var account ServiceAccount

	// Unmarshal the JSON data into the struct
	err := json.Unmarshal([]byte(credJson), &account)
	if err != nil {
		return "", err
	}
	return account.ProjectID, nil
}
