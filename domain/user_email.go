package domain

import "context"

type UserOTPEmail struct {
	Name string
	OTP  string
}

type UserEmailService interface {
	SendConfirmationCode(ctx context.Context, user User) error
}
