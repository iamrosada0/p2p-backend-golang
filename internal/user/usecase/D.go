package usecase

import "p2p/internal/user/entity"

type DeleteUserInputDto struct {
	ID uint `json:"id"`
}

type DeleteUserOutputDto struct {
	ID uint `json:"id"`
}

type DeleteUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewDeleteUserUseCase(UserRepository entity.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{UserRepository: UserRepository}
}

func (u *DeleteUserUseCase) Execute(input DeleteUserInputDto) (*DeleteUserOutputDto, error) {
	err := u.UserRepository.DeleteByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &DeleteUserOutputDto{ID: input.ID}, nil
}
