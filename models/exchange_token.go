package models

type ExchangeToken struct {
	Code        string `json:"code"`
	RedirectURL string `json:"redirect_url"`
}
