package github

type GetGithubUserInfoReply struct {
	Provider  string `json:"provider"`
	Login     string `json:"login"`
	Location  string `json:"location"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

type GetAuthURLResponse struct {
	Url string
}

type EmailInfo struct {
	Email      string `json:"email"`
	Primary    bool   `json:"primary"`
	Verified   bool   `json:"verified"`
	Visibility string `json:"visibility"`
}
