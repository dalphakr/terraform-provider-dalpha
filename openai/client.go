package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var (
	ADMIN_KEY = os.Getenv("OPENAI_ADMIN_KEY")
)

// CreateServiceAccount creates a service account in OpenAI
// and returns the API key
func CreateServiceAccount(projectId string, entity string) (string, string, error) {
	payload, err := json.Marshal(map[string]string{"name": entity})
	if err != nil {
		return "", "", err
	}

	cli := &http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf("https://api.openai.com/v1/organization/projects/%s/service_accounts", projectId),
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return "", "", err
	}
	req.Header.Add("Authorization", "Bearer "+ADMIN_KEY)
	req.Header.Add("Content-Type", "application/json")
	resp, err := cli.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != 201 {
		return "", "", fmt.Errorf("failed to create service account: %s", resp.Status)
	}
	// read response
	var result map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", "", err
	}

	return result["id"].(string), result["api_key"].(map[string]interface{})["value"].(string), nil
}

type ServiceAccountResponse struct {
	Object string `json:"object"`
	Data   []*ServiceAccount
}

type ServiceAccount struct {
	Object    string `json:"object"`
	Id        string `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	CreatedAt int    `json:"created_at"`
}

func FindServiceAccount(projectId string, entity string) (*ServiceAccount, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api.openai.com/v1/organization/projects/%s/service_accounts", projectId),
		nil,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+ADMIN_KEY)
	req.Header.Add("Content-Type", "application/json")

	cli := &http.Client{}
	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Failed to read service account: %s", resp.Status)
	}

	var accs ServiceAccountResponse
	if err = json.NewDecoder(resp.Body).Decode(&accs); err != nil {
		return nil, err
	}

	for _, acc := range accs.Data {
		if acc.Name == entity {
			return acc, nil
		}
	}

	return nil, fmt.Errorf("service account not found")
}
