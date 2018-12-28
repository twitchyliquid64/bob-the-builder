package auth

import (
	"errors"
	"github.com/hoisie/web"
)

var ErrNotAuthenticated = errors.New("Not authenticated")
var ErrUserDoesntExist = errors.New("User not found")

type UserSource interface {
	GetUserByUsername(username string) (User, error)
}

type User interface {
	Name() string
	CheckPassword(pass string) (bool, error)
}

type OTPUser interface {
	User
	VerifyOTP(code string) bool
	IsOTPEnrolled() bool
}

type AuthInfo struct {
	User User
}

type Auther interface {
	AuthInfo(*web.Context) (*AuthInfo, error)
	DoLogin(*web.Context) (*AuthInfo, error)
}
