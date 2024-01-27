package usecase

import "p2p/internal/user/entity"

type UpdateUserInputDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserOutputDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewUpdateUserUseCase(UserRepository entity.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{UserRepository: UserRepository}
}

func (u *UpdateUserUseCase) Execute(input UpdateUserInputDto) (*UpdateUserOutputDto, error) {
	existingUser, err := u.UserRepository.GetByID(input.ID)
	if err != nil {
		return nil, err
	}

	existingUser.Update(
		input.Name,
		input.Email,
	)

	err = u.UserRepository.Update(existingUser)
	if err != nil {
		return nil, err
	}

	return &UpdateUserOutputDto{
		ID:    existingUser.ID,
		Name:  existingUser.Name,
		Email: existingUser.Email,
	}, nil
}
