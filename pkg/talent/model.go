package talent

import "net/http"

type TalentObject struct {
	ServiceEndpoint string
	UserId          string
	Cookie          *http.Cookie
}

type Information struct {
	Success bool `json:"success"`
	User    struct {
		Name        string `json:"name"`
		PhoneNumber string `json:"phoneNumber"`
		ID          string `json:"id"`
	} `json:"user"`
}
