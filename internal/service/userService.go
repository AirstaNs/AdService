package service

import (
	"golang.org/x/net/context"
	"homework10/internal/adapters/repository/userrepo"
	"homework10/internal/entities"
)

type usersService struct {
	userRepository userrepo.UserRepository
}

//go:generate go run github.com/vektra/mockery/v2@v2.25.0 --name=UserService --filename=mockUserService.go --output ../mocks/servicemocks
type UserService interface {
	CreateUser(ctx context.Context, nickname string, email string) (*entities.User, error)
	UpdateUser(ctx context.Context, UserID int64, Nickname string, Email string) (*entities.User, error)
	GetUserByID(ctx context.Context, userID int64) (*entities.User, error)
	RemoveUser(ctx context.Context, userID int64) error
}

func NewUserService(userRepository userrepo.UserRepository) UserService {
	return &usersService{userRepository: userRepository}
}

func (a *usersService) CreateUser(ctx context.Context, nickname string, email string) (*entities.User, error) {
	user := entities.User{
		Nickname: nickname,
		Email:    email,
	}
	id, err := a.userRepository.AddUser(user)
	user.ID = id

	if err != nil {
		return &user, err
	}

	return &user, nil
}

func (a *usersService) UpdateUser(ctx context.Context, UserID int64, Nickname string, Email string) (*entities.User, error) {
	userByID, err := a.userRepository.GetUserByID(UserID)
	if err != nil {
		return userByID, err
	}
	setUser := entities.User{}

	setUser.ID = userByID.ID

	if Nickname != "" {
		setUser.Nickname = Nickname
	}
	if Email != "" {
		setUser.Email = Email
	}

	return a.userRepository.EditUser(setUser)
}

func (a *usersService) GetUserByID(ctx context.Context, userID int64) (*entities.User, error) {
	return a.userRepository.GetUserByID(userID)
}

func (a *usersService) RemoveUser(ctx context.Context, userID int64) error {
	_, err := a.userRepository.GetUserByID(userID)
	if err != nil {
		return err
	}
	return a.userRepository.DeleteUser(userID)
}
