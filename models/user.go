package models

type User struct {
	Sub           string `json:"sub"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	Locale        Locale `json:"locale"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Email         string `json:"email"`
	Picture       string `json:"picture"`
}

type Locale struct {
	Country  string `json:"country"`
	Language string `json:"language"`
}
