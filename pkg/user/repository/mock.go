package repository

import "github.com/ivanterekh/go-skeleton/pkg/user"

type mock []user.User

// NewMock returns new mock repository.
func NewMock() UserRepository {
	return mock{
		user.User{
			Email:    "user1@gmail.com",
			Name:     "User Friendly",
			Role:     "user",
			ID:       42,
			Password: "user1",
		},
		user.User{
			Email:    "user2@gmail.com",
			Name:     "John Walker",
			Role:     "user",
			ID:       45,
			Password: "12345qwerty",
		},
	}
}

// GetByCreds returns user with provided credentials
// or returns ErrNoSuchUser if it does not exist.
func (repo mock) GetByCreds(email, password string) (*user.User, error) {
	for _, u := range repo {
		if u.Email == email && u.Password == password {
			return &u, nil
		}
	}

	return nil, ErrNoSuchUser
}

// GetByID returns user with provided id or
// returns ErrNoSuchUser if it does not exist.
func (repo mock) GetByID(id int) (*user.User, error) {
	for _, u := range repo {
		if u.ID == id {
			return &u, nil
		}
	}

	return nil, ErrNoSuchUser
}