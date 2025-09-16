package user

import (
	"profile-service/pkg/entities"
	"regexp"

	"github.com/google/uuid"
)

type Service interface {
	FindByID(id string) (*entities.User, error)
	GetByID(id string) (*entities.User, error)
	IsFileExist(id string) (bool, error)
	UpdateEmail(userIDString string, email string) (*entities.User, error)
	UpdatePhone(userIDString string, phone string) (*entities.User, error)
	UpdateProfile(userIDString string, user entities.UpdateUserRequest) (*entities.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) FindByID(id string) (*entities.User, error) {
	return s.repo.FindByID(id)
}

func (s *service) IsFileExist(id string) (bool, error) {
	// Anda bisa memanggil fungsi IsFileExist pada repo disini
	return s.repo.IsFileExist(id)
}

func (s *service) UpdateProfile(userIDString string, user entities.UpdateUserRequest) (*entities.User, error) {
	// cek fileID
	fileID, err := uuid.Parse(user.FileId)
	if err != nil {
		return nil, entities.ErrInvalidFileID
	}

	// cek file exist
	cekFile, err := s.repo.IsFileExist(user.FileId)
	if err != nil {
		return nil, err
	}
	if !cekFile {
		return nil, entities.ErrFileNotFound
	}

	// cek userID
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, entities.ErrInvalidUserID
	}

	newData := &entities.User{
		ID:                userID,
		BankAccountName:   user.BankAccountName,
		BankAccountNumber: user.BankAccountNumber,
		BankAccountHolder: user.BankAccountHolder,
		FileId:            &fileID,
	}

	return s.repo.UpdateProfile(newData)
}

func (s *service) UpdateEmail(userIDString string, email string) (*entities.User, error) {
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, entities.ErrInvalidUserID
	}

	new_data := &entities.User{
		ID:    userID,
		Email: email,
	}
	return s.repo.UpdateEmail(new_data)
}

func (s *service) UpdatePhone(userIDString string, phone string) (*entities.User, error) {
	userID, err := uuid.Parse(userIDString)
	if err != nil {
		return nil, entities.ErrInvalidUserID
	}

	//cek panjang phone
	if len(phone) < 8 || len(phone) > 16 {
		return nil, entities.ErrInvalidUserID
	}

	//cek format phone
	// re := regexp.MustCompile(`^\+[1-9]\d{1,14}$`)
	re := regexp.MustCompile(`^(?:\+62)8[1-9][0-9]{6,10}$`)
	if !re.MatchString(phone) {
		return nil, entities.ErrInvalidPhoneNumber
	}

	new_data := &entities.User{
		ID:    userID,
		Phone: phone,
	}
	return s.repo.UpdatePhone(new_data)

}

func (s *service) GetByID(id string) (*entities.User, error) {
	return s.repo.GetByID(id)
}
