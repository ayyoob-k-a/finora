package repo

import (
	"errors"

	"github.com/ayyoob-k-a/finora/domain"
	"github.com/ayyoob-k-a/finora/model/inbound"
	"github.com/ayyoob-k-a/finora/model/response"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Repo struct {
	db *gorm.DB

	// Add fields as needed for your repository
}

func NewRepo(db *gorm.DB) *Repo {
	return &Repo{
		db: db,
	}
}
func (r *Repo) Signup(data domain.User) (int, error) {
	var existing domain.User
	err := r.db.Where("email = ?", data.Email).First(&existing).Error
	if err == nil {
		return 0, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, err
	}

	data.Password, err = hashPassword(data.Password)
	if err != nil {
		return 0, err
	}

	if err := r.db.Create(&data).Error; err != nil {
		return 0, err
	}

	return data.ID, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func (r *Repo) Login(data inbound.Login) (*domain.User, error) {
	var user domain.User

	err := r.db.
		Where("email = ? OR phone = ?", data.Identifier, data.Identifier).
		First(&user).Error

	if err != nil {
		return nil, errors.New("invalid email/phone or password")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		return nil, errors.New("invalid email/phone or password")
	}

	return &user, nil
}

// VerifyUser marks the user as verified in the DB
func (r *Repo) VerifyUser(email string) error {
	return r.db.Model(&domain.User{}).
		Where("email = ?", email).
		Update("is_verified", true).Error
}

// GetUserByEmail fetches user by email
func (r *Repo) GetUserByEmail(email string) (*response.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &response.User{
		ID:       int(user.ID),
		Email:    user.Email,
		Username: user.Username,
	}, nil
}

func (r *Repo) InsertOtp(data domain.User) error {

	err := r.db.Save(&data).Error
	if err != nil {
		return err
	}

	return nil

}
