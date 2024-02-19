package usersusecases

import (
	"fmt"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/users"
	usersRepositories "github.com/NattpkJsw/real-world-api-go/modules/users/usersRepositories"
	"github.com/NattpkJsw/real-world-api-go/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type IUsersUsecase interface {
	InsertCustomer(req *users.UserRegisterReq) (*users.ResponsePassport, error)
	GetPassport(req *users.UserCredential) (*users.ResponsePassport, error)
	DeleteOauth(accessToken string) error
	GetUser(token string) (*users.ResponsePassport, error)
	UpdateUser(user *users.UserCredentialCheck) (*users.ResponsePassport, error)
}

type usersUsecase struct {
	cfg             config.IConfig
	usersRepository usersRepositories.IUsersRepository
}

func UsersUsecase(cfg config.IConfig, usersRepository usersRepositories.IUsersRepository) IUsersUsecase {
	return &usersUsecase{
		cfg:             cfg,
		usersRepository: usersRepository,
	}
}

func (u *usersUsecase) InsertCustomer(req *users.UserRegisterReq) (*users.ResponsePassport, error) {
	password := req.Password
	if err := req.BcryptHashing(); err != nil {
		return nil, err
	}

	// Insert user
	_, err := u.usersRepository.InsertUser(req)
	if err != nil {
		return nil, err
	}

	loginUser := &users.UserCredential{
		Email:    req.Email,
		Password: password,
	}

	return u.GetPassport(loginUser)

}

func (u *usersUsecase) GetPassport(req *users.UserCredential) (*users.ResponsePassport, error) {
	//Find user
	user, err := u.usersRepository.FindOneUserByEmail(req.Email)
	if err != nil {
		return nil, err
	}
	//Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, fmt.Errorf("password is invalid")
	}
	// sign token
	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: user.Id,
	})
	if err != nil {
		return nil, err
	}

	// set user token
	userToken := &users.UserToken{
		User_Id:     user.Id,
		AccessToken: accessToken.SignToken(),
	}
	if err := u.usersRepository.InsertOauth(userToken); err != nil {
		return nil, err
	}

	//Set passport
	passport := &users.UserPassport{
		// Id:       user.Id,
		Email:    user.Email,
		Username: user.Username,
		Image:    user.Image,
		Bio:      user.Bio,
		Token:    userToken.AccessToken,
	}
	passportOutput := &users.ResponsePassport{
		User: *passport,
	}
	return passportOutput, nil
}

func (u *usersUsecase) DeleteOauth(accessToken string) error {
	if err := u.usersRepository.DeleteOauth(accessToken); err != nil {
		return err
	}
	return nil
}

func (u *usersUsecase) GetUser(token string) (*users.ResponsePassport, error) {

	oauth, err := u.usersRepository.FindOneOath(token)
	if err != nil {
		return nil, err
	}

	profile, err := u.usersRepository.GetProfile(oauth.UserId)
	if err != nil {
		return nil, err
	}

	passport, err := u.refreshToken(token)
	if err != nil {
		return nil, err
	}

	userPassport := &users.UserPassport{
		Email:    profile.Email,
		Username: profile.Username,
		Image:    profile.Image,
		Bio:      profile.Bio,
		Token:    passport.AccessToken,
	}

	resPassport := &users.ResponsePassport{
		User: *userPassport,
	}

	return resPassport, nil
}

func (u *usersUsecase) UpdateUser(user *users.UserCredentialCheck) (*users.ResponsePassport, error) {
	updatedUser, err := u.usersRepository.UpdateUser(user)
	if err != nil {
		return nil, err
	}

	passport, err := u.refreshToken(user.AccessToken)
	if err != nil {
		return nil, err
	}

	userOut := &users.UserPassport{
		Email:    updatedUser.Email,
		Username: updatedUser.Username,
		Image:    updatedUser.Image,
		Bio:      updatedUser.Bio,
		Token:    passport.AccessToken,
	}
	userRes := &users.ResponsePassport{
		User: *userOut,
	}
	return userRes, nil
}

func (u *usersUsecase) refreshToken(accessTokenIn string) (*users.UserToken, error) {
	oauthID, err := u.usersRepository.FindOneOath(accessTokenIn)
	if err != nil {
		return nil, err
	}

	accessToken, err := auth.NewAuth(auth.Access, u.cfg.Jwt(), &users.UserClaims{
		Id: oauthID.UserId,
	})
	if err != nil {
		return nil, err
	}
	passport := &users.UserToken{
		Id:          oauthID.Id,
		User_Id:     oauthID.UserId,
		AccessToken: accessToken.SignToken(),
	}

	if err := u.usersRepository.UpdateOauth(passport); err != nil {
		return nil, err
	}

	return passport, nil
}
