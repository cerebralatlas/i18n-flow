package response

type LoginResp struct {
	AccessToken  string       `json:"access_token"`
	RefreshToken string       `json:"refresh_token"`
	User         UserInfoResp `json:"user"`
}

type UserInfoResp struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}
