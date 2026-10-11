package user

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(input RegisterInput) (User, error) {
	input.Nickname = strings.TrimSpace(input.Nickname)
	input.Name = strings.TrimSpace(input.Name)

	if input.Nickname == "" || utf8.RuneCountInString(input.Nickname) > 30 {
		return User{}, fmt.Errorf("Nickname must contain 1 to 30 characters: %w", ErrInvalidInput)
	}
	if input.Name == "" || utf8.RuneCountInString(input.Name) > 50 {
		return User{}, fmt.Errorf("Name must contain 1 to 50 characters: %w", ErrInvalidInput)
	}
	if utf8.RuneCountInString(input.Password) < 8 {
		return User{}, fmt.Errorf("Password must contain at least 8 characters: %w", ErrInvalidInput)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("Hash password: %w", err)
	}

	u := User{
		Nickname:     input.Nickname,
		Name:         input.Name,
		PasswordHash: string(hash),
	}

	created, err := s.repo.Create(u)
	if err != nil {
		return User{}, err
	}

	return created, nil
}

func (s *Service) Get(id int64) (User, error) {
	if id <= 0 {
		return User{}, fmt.Errorf("Id must be positive: %w", ErrInvalidInput)
	}

	return s.repo.GetByID(id)
}

func (s *Service) Delete(id int64) error {
	if id <= 0 {
		return fmt.Errorf("Id must be positive: %w", ErrInvalidInput)
	}

	return s.repo.Delete(id)
}
