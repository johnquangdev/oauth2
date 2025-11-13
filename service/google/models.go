package google

type AuthURLResp struct {
	AuthURL     string `json:"auth_url"`
	State       string `json:"state"`
	RedirectURL string `json:"redirect_url"`
	ClientID    string `json:"client_id"`
}

type UserInfoResp struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`
	Id       string `json:"id"`
	Provider string `json:"provider"`
}
