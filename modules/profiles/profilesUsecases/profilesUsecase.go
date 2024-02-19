package profilesusecases

import (
	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/NattpkJsw/real-world-api-go/modules/profiles"
	profilesrepositories "github.com/NattpkJsw/real-world-api-go/modules/profiles/profilesRepositories"
)

type IProfilesUsecase interface {
	GetProfile(username string, curUserId int) (*profiles.JsonProfile, error)
	FollowUser(username string, curUserId int) (*profiles.JsonProfile, error)
	UnfollowUser(username string, curUserId int) (*profiles.JsonProfile, error)
}

type profilesUsecase struct {
	cfg                config.IConfig
	profilesRepository profilesrepositories.IProfilesRepository
}

func ProfilesUsecase(cfg config.IConfig, profilesRepository profilesrepositories.IProfilesRepository) IProfilesUsecase {
	return &profilesUsecase{
		cfg:                cfg,
		profilesRepository: profilesRepository,
	}
}

func (u *profilesUsecase) GetProfile(username string, curUserId int) (*profiles.JsonProfile, error) {
	profile, err := u.profilesRepository.FindOneUserProfileByUsername(username, curUserId)
	if err != nil {
		return nil, err
	}
	jsonProfile := &profiles.JsonProfile{
		Profile: *profile,
	}
	return jsonProfile, nil
}

func (u *profilesUsecase) FollowUser(username string, curUserId int) (*profiles.JsonProfile, error) {
	profile, err := u.profilesRepository.FollowUser(username, curUserId)
	if err != nil {
		return nil, err
	}
	jsonProfile := &profiles.JsonProfile{
		Profile: *profile,
	}
	return jsonProfile, nil
}

func (u *profilesUsecase) UnfollowUser(username string, curUserId int) (*profiles.JsonProfile, error) {
	profile, err := u.profilesRepository.UnfollowUser(username, curUserId)
	if err != nil {
		return nil, err
	}
	jsonProfile := &profiles.JsonProfile{
		Profile: *profile,
	}
	return jsonProfile, nil
}
