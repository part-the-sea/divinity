package main

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUser_ValidUser_ReturnsNoError(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := validateUser(user)

	assert.NoError(t, err)
}

func TestValidateUser_ReturnsErrorForInvalidPassword(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "",
	}

	err := validateUser(user)

	assert.Error(t, err)
	assert.Equal(t, "password is required", err.Error())
}

func TestValidateUser_ReturnsErrorForInvalidFirstName(t *testing.T) {
	user := &User{
		FirstName: "",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := validateUser(user)

	assert.Error(t, err)
	assert.Equal(t, "first name is required", err.Error())
}

func TestValidateUser_ReturnsErrorForInvalidLastName(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := validateUser(user)

	assert.Error(t, err)
	assert.Equal(t, "last name is required", err.Error())
}

func TestValidateUser_ReturnsErrorForInvalidEmail(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "",
		Password:  "password",
	}

	err := validateUser(user)

	assert.Error(t, err)
	assert.Equal(t, "email is required", err.Error())
}

func TestValidateUser_ReturnsErrorForInvalidEmailFormat(t *testing.T) {
	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe",
		Password:  "password",
	}

	err := validateUser(user)

	assert.Error(t, err)
	assert.Equal(t, "invalid email format", err.Error())
}

type MockUserStore struct {
	CreateFunc     func(ctx context.Context, user *User) error
	GetByIDFunc    func(ctx context.Context, id string) (*User, error)
	GetByEmailFunc func(ctx context.Context, email string) (*User, error)
	UpdateFunc     func(ctx context.Context, user *User) error
	DeleteFunc     func(ctx context.Context, id string) error
}

func (m *MockUserStore) Create(ctx context.Context, user *User) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, user)
	}

	return nil
}

func (m *MockUserStore) GetByID(ctx context.Context, id string) (*User, error) {
	if m.GetByIDFunc != nil {
		return m.GetByIDFunc(ctx, id)
	}

	return nil, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	if m.GetByEmailFunc != nil {
		return m.GetByEmailFunc(ctx, email)
	}

	return nil, nil
}

func (m *MockUserStore) Update(ctx context.Context, user *User) error {
	if m.UpdateFunc != nil {
		return m.UpdateFunc(ctx, user)
	}

	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, id string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, id)
	}

	return nil
}

func TestUserService_Create_ReturnsErrorForFailingToCheckForExistingUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := userService.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Create_ReturnsErrorForExistingUserEmail(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{ID: "1", Email: "john.doe@example.com"}, nil
		},
	})

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := userService.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())
}

func TestUserService_Create_ReturnsErrorForFailingToCreateUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		CreateFunc: func(ctx context.Context, user *User) error {
			return errors.New("random error")
		},
	})

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := userService.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Create_ReturnsErrorForFailingToHashPassword(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		CreateFunc: func(ctx context.Context, user *User) error {
			return errors.New("random error")
		},
	})

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := userService.Create(context.Background(), user)

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Create_ReturnsNoErrorForValidUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		CreateFunc: func(ctx context.Context, user *User) error {
			return nil
		},
	})

	user := &User{
		FirstName: "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		Password:  "password",
	}

	err := userService.Create(context.Background(), user)

	assert.NoError(t, err)
}

func TestUserService_GetByID_ReturnsErrorForFailingToGetUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	user, err := userService.GetByID(context.Background(), "1")

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
	assert.Nil(t, user)
}

func TestUserService_GetByID_ReturnsErrorForUserNotFound(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, nil
		},
	})

	user, err := userService.GetByID(context.Background(), "1")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
}

func TestUserService_GetByID_ReturnsUserForValidID(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	user, err := userService.GetByID(context.Background(), "1")

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john.doe@example.com", user.Email)
}

func TestUserService_GetByEmail_ReturnsErrorForFailingToGetUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	user, err := userService.GetByEmail(context.Background(), "john.doe@example.com")

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
	assert.Nil(t, user)
}

func TestUserService_GetByEmail_ReturnsErrorForUserNotFound(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return nil, nil
		},
	})

	user, err := userService.GetByEmail(context.Background(), "john.doe@example.com")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
	assert.Nil(t, user)
}

func TestUserService_GetByEmail_ReturnsUserForValidEmail(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	user, err := userService.GetByEmail(context.Background(), "john.doe@example.com")

	assert.NoError(t, err)
	assert.Equal(t, "1", user.ID)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, "john.doe@example.com", user.Email)
}

func TestUserService_Update_ReturnsErrorForFailingToCheckForExistingUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	err := userService.Update(context.Background(), "1", &UpdateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Update_ReturnsErrorForFailingToUpdateUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		UpdateFunc: func(ctx context.Context, user *User) error {
			return errors.New("random error")
		},
	})

	err := userService.Update(context.Background(), "1", &UpdateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Update_ReturnsValidUserForValidRequest(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		UpdateFunc: func(ctx context.Context, user *User) error {
			return nil
		},
	})

	err := userService.Update(context.Background(), "1", &UpdateUserRequest{
		FirstName: "Jane",
		LastName:  "Doe",
	})

	assert.NoError(t, err)
}

func TestUserService_UpdatePassword_ReturnsErrorForFailingToGetUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	err := userService.UpdatePassword(context.Background(), "1", &UpdatePasswordRequest{
		Password: "password",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_UpdatePassword_ReturnsErrorForUserNotFound(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, nil
		},
	})

	err := userService.UpdatePassword(context.Background(), "1", &UpdatePasswordRequest{
		Password: "password",
	})

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestUserService_UpdatePassword_ReturnsErrorForEmptyPassword(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	err := userService.UpdatePassword(context.Background(), "1", &UpdatePasswordRequest{
		Password: "",
	})

	assert.Error(t, err)
	assert.Equal(t, "password is required", err.Error())
}

func TestUserService_UpdatePassword_ReturnsNoErrorForValidRequest(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		UpdateFunc: func(ctx context.Context, user *User) error {
			return nil
		},
	})

	err := userService.UpdatePassword(context.Background(), "1", &UpdatePasswordRequest{
		Password: "password",
	})

	assert.NoError(t, err)
}

func TestUserService_UpdatePassword_ReturnsErrorForFailingToUpdateUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		UpdateFunc: func(ctx context.Context, user *User) error {
			return errors.New("random error")
		},
	})

	err := userService.UpdatePassword(context.Background(), "1", &UpdatePasswordRequest{
		Password: "password",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_UpdateEmail_ReturnsErrorForFailingToGetUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "john.doe@example.com",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_UpdateEmail_ReturnsErrorForUserNotFound(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, nil
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "john.doe@example.com",
	})

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestUserService_UpdateEmail_ReturnsErrorForEmptyEmail(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "",
	})

	assert.Error(t, err)
	assert.Equal(t, "email is required", err.Error())
}

func TestUserService_UpdateEmail_ReturnsErrorForFailingToCheckForExistingUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "john.doe@example.com",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_UpdateEmail_ReturnsErrorForExistingEmail(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		GetByEmailFunc: func(ctx context.Context, email string) (*User, error) {
			return &User{ID: "2", FirstName: "Jane", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "john.doe@example.com",
	})

	assert.Error(t, err)
	assert.Equal(t, "user with this email already exists", err.Error())
}

func TestUserService_UpdateEmail_ReturnsErrorForFailingToUpdateUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		UpdateFunc: func(ctx context.Context, user *User) error {
			return errors.New("random error")
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "john.doe@example.com",
	})

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_UpdateEmail_ReturnsNoErrorForValidRequest(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	err := userService.UpdateEmail(context.Background(), "1", &UpdateEmailRequest{
		Email: "john.doe@example.com",
	})

	assert.NoError(t, err)
}

func TestUserService_Delete_ReturnsErrorForFailingToGetUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, errors.New("random error")
		},
	})

	err := userService.Delete(context.Background(), "1")

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Delete_ReturnsErrorForUserNotFound(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return nil, nil
		},
	})

	err := userService.Delete(context.Background(), "1")

	assert.Error(t, err)
	assert.Equal(t, "user not found", err.Error())
}

func TestUserService_Delete_ReturnsErrorForFailingToDeleteUser(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
		DeleteFunc: func(ctx context.Context, id string) error {
			return errors.New("random error")
		},
	})

	err := userService.Delete(context.Background(), "1")

	assert.Error(t, err)
	assert.Equal(t, "an internal error occurred", err.Error())
}

func TestUserService_Delete_ReturnsNoErrorForValidRequest(t *testing.T) {
	userService := NewUserService(&MockUserStore{
		GetByIDFunc: func(ctx context.Context, id string) (*User, error) {
			return &User{ID: "1", FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"}, nil
		},
	})

	err := userService.Delete(context.Background(), "1")

	assert.NoError(t, err)
}
