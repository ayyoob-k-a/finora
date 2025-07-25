package usecase

import (
	"github.com/ayyoob-k-a/finora/configs"
	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/inbound"
	"github.com/ayyoob-k-a/finora/model/response"
	"github.com/ayyoob-k-a/finora/repo"
	"github.com/ayyoob-k-a/finora/utils"
)

type Usecase struct {
	repo *repo.Repo
	mail configs.Mail

	// Add fields as needed for your repository
}

func NewUsecase(repo *repo.Repo, Mail configs.Mail) *Usecase {
	return &Usecase{
		repo: repo,
		mail: Mail,
	}
}

func (u *Usecase) Signup(data domain.User) error {
	var err error

	data.ID, err = u.repo.Signup(data)
	if err != nil {
		return err
	}

	data.Otp, err = utils.GenerateOTP(6)
	if err != nil {
		return err
	}

	err = u.repo.InsertOtp(data)
	if err != nil {
		return err
	}

	go utils.SendVerificationEmail(data.Email, data.Otp, u.mail.SecretKey)
	return nil

}

func (u *Usecase) Login(data inbound.Login) (*response.AuthResponse, error) {
	user, err := u.repo.Login(data)
	if err != nil {
		return nil, err
	}

	accessToken, refreshToken, err := utils.GenerateToken(int(user.ID))
	if err != nil {
		return nil, err
	}

	return &response.AuthResponse{
		UserID:       int(user.ID),
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
