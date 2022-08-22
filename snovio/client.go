package snovio

import (
	"strconv"

	"github.com/GeoinovaDev/lower-resultys/time/interval"
	"github.com/GeoinovaDev/lower-resultys/convert"
	"github.com/GeoinovaDev/lower-resultys/convert/decode"
	"github.com/GeoinovaDev/lower-resultys/exec/try"
	"github.com/GeoinovaDev/lower-resultys/net/request"
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
func (client *Client) FindEmails(dominio string, limit int) (_emails []Email, _err error) {
	url := "https://api.snov.io/v2/domain-emails-with-info?lastId=0&type=all&limit=" + strconv.Itoa(limit) + "&domain=" + dominio + "&access_token=" + client.AccessToken
	try.New().SetTentativas(3).Run(func() {
		response, err := request.New(url).SetTimeout(15).Get()
		if err != nil {
			panic(err)
		}

		protocol := Protocol{}
		decode.JSON(response, &protocol)

		_emails = protocol.Emails
		_err = nil
	}).Catch(func(msg string) {
		_emails = nil
		_err = newError(msg)
	})

	return
}

// AddEmailVerification ...
func (client *Client) AddEmailVerification(email string) (_result bool, _err error) {
	form := make(map[string]string)
	form["access_token"] = client.AccessToken

	url := "https://api.snov.io/v1/add-emails-to-verification?emails[]=" + email
	try.New().SetTentativas(3).Run(func() {
		response, err := request.New(url).SetTimeout(3).AddHeader("access_token", client.AccessToken).Post(form)
		if err != nil {
			panic(err)
		}

		protocol := make(map[string]map[string]bool)
		decode.JSON(response, &protocol)

		if protocol[email]["sent"] {
			_result = true
			_err = nil
		}
	}).Catch(func(msg string) {
		_result = false
		_err = newError(msg)
	})

	return
}

// CheckEmailStatus ...
func (client *Client) CheckEmailStatus(email string) (_result ProtocolStatus, _err error) {
	url := "https://api.snov.io/v1/get-emails-verification-status?emails[]=" + email
	form := make(map[string]string)
	form["access_token"] = client.AccessToken

	try.New().SetTentativas(3).Run(func() {
		response, err := request.New(url).SetTimeout(3).AddHeader("access_token", client.AccessToken).Post(form)
		if err != nil {
			panic(err)
		}

		protocol := map[string]ProtocolStatus{}
		decode.JSON(response, &protocol)

		if _, ok := protocol[email]; !ok {
			panic("Email " + email + " não retornado pelo https://api.snov.io/v1/get-emails-verification-status")
		}

		_result = protocol[email]
		_err = nil
	}).Catch(func(msg string) {
		_result = ProtocolStatus{}
		_err = newError(msg)
	})

	return
}

// CheckEmailValid ...
func (client *Client) CheckEmailValid(email string) (_result string, _err error) {
	check, err := client.CheckEmailStatus(email)
	if err != nil {
		_result = "unknown"
		_err = err
		return
	}

	if check.Status.Identifier == "not_verified" {
		result, err := client.AddEmailVerification(email)
		if err != nil {
			_result = "unknown"
			_err = err
			return
		}

		if result {
			waiting := true

			interval.New().Wait(10, func() {
				waiting = false
			})

			for waiting {
				check, err = client.CheckEmailStatus(email)
				if err != nil {
					_result = "unknown"
					_err = newError("Não foi possível checar o status do email " + email)
					return
				}

				if check.Status.Identifier == "in_progress" {
					continue
				}

				break
			}
		}
	}

	if check.Status.Identifier == "complete" {
		if check.Data.SMTPStatus == "valid" {
			_result = "verified"
		} else if check.Data.SMTPStatus == "not_valid" {
			_result = "notVerified"
		} else {
			_result = "unknown"
		}
		_err = nil
		return
	}

	return
}

// UpdateToken solicita token de autorização
func (client *Client) UpdateToken() error {
	form := make(map[string]string)
	form["grant_type"] = "client_credentials"
	form["client_id"] = client.ID
	form["client_secret"] = client.Secret

	response, err := request.New("https://api.snov.io/v1/oauth/access_token").Post(form)
	if err != nil {
		return err
	}

	json := make(map[string]string)
	convert.StringToJSON(response, &json)

	client.AccessToken = json["access_token"]

	return nil
}
