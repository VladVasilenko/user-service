package service

import "errors"

// UserRepository описывает контракт для работы с базой данных
type UserRepository interface {
	GetByID(id int) (string, bool, error)
}

type UserService struct {
	repo UserRepository
}

// NewUserService — конструктор сервиса
func NewUserService(r UserRepository) *UserService {
	return &UserService{
		repo: r,
	}
}

// GetUserProfile — метод, где будет жить чистая бизнес-логика
func (s *UserService) GetUserProfile(id int) (string, error) {
	if id <= 0 {
		return "", errors.New("невалидный ID пользователя")
	}

	name, exists, err := s.repo.GetByID(id)
	if err != nil {
		return "", err
	}
	if !exists {
		return "", errors.New("пользователь не найден")
	}

	// Здесь в будущем можно добавить проверки: забанен ли юзер,
	// залогинен ли, отформатировать его имя и т.д.
	return name, nil
}
