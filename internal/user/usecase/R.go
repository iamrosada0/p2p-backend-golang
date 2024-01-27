package usecase

import "p2p/internal/user/entity"

type GetAllUsersOutputDto struct {
	Users []*entity.User `json:"users"`
}

type GetAllUsersUseCase struct {
	UserRepository entity.UserRepository
}

func NewGetAllUsersUseCase(UserRepository entity.UserRepository) *GetAllUsersUseCase {
	return &GetAllUsersUseCase{UserRepository: UserRepository}
}

func (u *GetAllUsersUseCase) Execute() (*GetAllUsersOutputDto, error) {
	Users, err := u.UserRepository.FindAll()
	if err != nil {
		return nil, err
	}

	return &GetAllUsersOutputDto{Users: Users}, nil
}
