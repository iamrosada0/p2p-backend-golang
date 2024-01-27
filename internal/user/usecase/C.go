package usecase

import "p2p/internal/user/entity"

type CreateUserInputDto struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserOutputDto struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type CreateUserUseCase struct {
	UserRepository entity.UserRepository
}

func NewCreateUserUseCase(UserRepository entity.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{UserRepository: UserRepository}
}

func (u *CreateUserUseCase) Execute(input CreateUserInputDto) (*CreateUserOutputDto, error) {
	User := entity.NewUser(
		input.Name,
		input.Email,
	)

	err := u.UserRepository.Create(User)
	if err != nil {
		return nil, err
	}

	return &CreateUserOutputDto{
		ID:    User.ID,
		Name:  User.Name,
		Email: User.Email,
	}, nil
}
