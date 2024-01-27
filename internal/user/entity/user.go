package entity

import (
	"errors"
	"math/rand"
)

type UserRepository interface {
	Create(user *User) error
	FindAll() ([]*User, error)
	Update(user *User) error
	DeleteByID(id uint) error
	GetByID(id uint) (*User, error)
}

type User struct {
	ID    uint
	Name  string
	Email string
}

func NewUser(name, email string) *User {
	return &User{
		ID:    uint(rand.Uint32()),
		Name:  name,
		Email: email,
	}
}

func (d *User) Update(name, email string) {
	d.Name = name
	d.Email = email
}

type InMemoryUserRepository struct {
	users map[string]*User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*User),
	}
}

func (r *InMemoryUserRepository) DeleteByID(id string) error {
	if _, exists := r.users[id]; !exists {
		return errors.New("user not found")
	}

	delete(r.users, id)
	return nil
}

func (r *InMemoryUserRepository) FindAll() ([]*User, error) {
	var allUsers []*User
	for _, user := range r.users {
		allUsers = append(allUsers, user)
	}
	return allUsers, nil
}
