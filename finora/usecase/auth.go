package usecase

import (
	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/inbound"
	"github.com/ayyoob-k-a/finora/model/response"
	"github.com/ayyoob-k-a/finora/repo"
	"github.com/ayyoob-k-a/finora/utils"
)

type Usecase struct {
	repo *repo.Repo

	// Add fields as needed for your repository
}

func NewUsecase(repo *repo.Repo) *Usecase {
	return &Usecase{
		repo: repo,
	}
}

func (u *Usecase) Signup(data domain.User) error {
	var err error

	data.ID, err = u.repo.Signup(data)
	if err != nil {
		return err
	}

	// Note: OTP functionality moved to new service layer
	return nil
}

func (u *Usecase) Login(data inbound.Login) (*response.AuthResponse, error) {
	_, err := u.repo.Login(data)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := utils.GenerateToken(1) // Placeholder ID
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		UserID:       1, // Placeholder ID since old system uses int
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (u *Usecase) VerifyUser(email string) error {
	return u.repo.VerifyUser(email)
}

func (u *Usecase) GetUserByEmail(email string) (*response.User, error) {
	return u.repo.GetUserByEmail(email)
}
