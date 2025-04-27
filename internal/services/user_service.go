package services

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interface"
	"context"
	"database/sql"
)

type UserService struct {
	userRepository domain_interface.UserRepositoryInterface
}

func NewUserService(userRepository domain_interface.UserRepositoryInterface) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) Create(payload *entity.RegisterPayload, tx *sql.Tx) (*entity.RegisterResponse, error) {
	if err := s.userRepository.Create(context.Background(), tx, &entity.User{
		Email:    payload.Email,
		Name:     payload.Name,
		Password: payload.Password,
	}); err != nil {
		return nil, err
	}

	return &entity.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *UserService) GetByID(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error) {
	return s.userRepository.GetByID(ctx, tx, id)
}

func (s *UserService) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error) {
	return s.userRepository.GetByEmail(ctx, tx, email)
}

func (s *UserService) Update(ctx context.Context, tx *sql.Tx, payload *entity.User) error {
	user, err := s.userRepository.GetByID(ctx, tx, payload.ID)
	if err != nil {
		return err
	}

	return s.userRepository.Update(ctx, tx, &entity.User{
		ID:       user.ID,
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	})
}

func (s *UserService) Delete(ctx context.Context, tx *sql.Tx, id string) error {
	return s.userRepository.Delete(ctx, tx, id)
}
