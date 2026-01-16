package auth

import "errors"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(email, password string) (*User, error) {
	hash, err := HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &User{
		Email:        email,
		PasswordHash: hash,
		Role:         "user", // 👈 default role
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(email, password string) (*User, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return nil, err
	}

	if !CheckPassword(password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}
	return user, nil
}
