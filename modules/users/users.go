package users

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int     `db:"id" json:"id"`
	Email    string  `db:"email" json:"email"`
	Username string  `db:"username" json:"username"`
	Image    *string `db:"image" json:"image"`
	Bio      *string `db:"bio" json:"bio"`
}

type Profile struct {
	Username  string  `json:"username"`
	Image     *string `json:"image"`
	Bio       *string `json:"bio"`
	Following bool    `json:"following"`
}

type UserProfile struct {
	Profile *Profile `json:"profile"`
}

type UserRegisterReq struct {
	Username string `db:"username" json:"username"`
	Email    string `db:"email" json:"email"`
	Password string `db:"password" json:"password"`
}

type RegisterReq struct {
	User *UserRegisterReq `json:"user"`
}

type ResponsePassport struct {
	User UserPassport `json:"user"`
}

type UserPassport struct {
	// Id       int     `db:"id" json:"id"`
	Email    string  `db:"email" json:"email"`
	Username string  `db:"username" json:"username"`
	Image    *string `db:"image" json:"image"`
	Bio      *string `db:"bio" json:"bio"`
	Token    string  `db:"access_token" json:"token"`
}

type UserToken struct {
	Id          string `db:"id" json:"id"`
	User_Id     int    `db:"user_id" json:"user_id"`
	AccessToken string `db:"access_token" json:"access_token"`
}

type UserCredential struct {
	Email    string `db:"email" json:"email" form:"email"`
	Password string `db:"password" json:"password" form:"password"`
}

type UserSignin struct {
	User UserCredential `json:"user"`
}

type UserCredentialInput struct {
	User UserCredentialCheck `json:"user"`
}

type UserCredentialCheck struct {
	Id          int     `json:"id" db:"id"`
	Email       string  `json:"email" db:"email"`
	Password    string  `json:"password" db:"password"`
	Username    string  `json:"username" db:"username"`
	Image       *string `json:"image" db:"image"`
	Bio         *string `json:"bio" db:"bio"`
	AccessToken string  `json:"access_token"`
}

type UserClaims struct {
	Id int `db:"id" json:"id"`
}

type OauthToken struct {
	AccessToken string `json:"access_token" form:"access_token"`
}

type Oauth struct {
	Id     string `db:"id" json:"id"`
	UserId int    `db:"user_id" json:"user_id"`
}

// type UserRemoveCredential struct {
// 	OauthId string `json:"oauth_id" form:"oauth_id"`
// }

func (obj *UserRegisterReq) BcryptHashing() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)

	}
	obj.Password = string(hashedPassword)
	return nil
}

func (obj *UserCredentialCheck) BcryptHashingUpdate() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(obj.Password), 10)
	if err != nil {
		return fmt.Errorf("hash password failed: %v", err)

	}
	obj.Password = string(hashedPassword)
	return nil
}

func (obj *UserRegisterReq) IsEmail() bool {
	match, err := regexp.MatchString(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`, obj.Email)
	if err != nil {
		return false
	}
	return match
}
