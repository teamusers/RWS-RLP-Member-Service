package requests

import "rlp-member-service/model"

type RegisterVerification struct {
	Email string `json:"email"`
}

type Register struct {
	User model.User `json:"user"`
}
