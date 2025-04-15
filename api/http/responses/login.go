package responses

type LoginResponse struct {
	LoginSessionToken string `json:"login_session_token"`
	LoginExpireIn     int64  `json:"login_session_token_expiry"`
}
