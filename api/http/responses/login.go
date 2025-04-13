package responses

type LoginResponse struct {
	LoginSessionToken string `json:"loginSessionToken"`
	LoginExpireIn     int64  `json:"login_expireIn"`
}
