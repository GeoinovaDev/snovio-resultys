package snovio

import (
	"git.resultys.com.br/lib/lower/convert"
	"git.resultys.com.br/lib/lower/net/request"
)

// Client struct
type Client struct {
	ID          string
	Secret      string
	AccessToken string
}

// New cria um client
func New(id string, secret string) *Client {
	client := &Client{
		ID:     id,
		Secret: secret,
	}

	return client
}

// FindEmails pesqisa email por dominio
// Return array, error
func (client *Client) FindEmails(dominio string) ([]string, error) {
	form := make(map[string]string)

	url := "https://app.snov.io/restapi/get-domain-emails-with-info?type=all&limit=100&domain=" + dominio
	response, err := request.New(url).AddHeader("Authorization", "Bearer "+client.AccessToken).Post(form)
	if err != nil {
		return nil, err
	}

	emails := []string{}
	json := make(map[string]interface{})
	convert.StringToJSON(response, &json)
	jsonEmails := json["emails"].([]interface{})
	for _, value := range jsonEmails {
		result := value.(map[string]interface{})
		emails = append(emails, result["email"].(string))
	}

	return emails, nil
}

// GetToken solicita token de autorização
func (client *Client) UpdateToken() error {
	form := make(map[string]string)
	form["grant_type"] = "client_credentials"
	form["client_id"] = client.ID
	form["client_secret"] = client.Secret

	response, err := request.New("https://app.snov.io/oauth/access_token").Post(form)
	if err != nil {
		return err
	}

	json := make(map[string]string)
	convert.StringToJSON(response, &json)

	client.AccessToken = json["access_token"]

	return nil
}
