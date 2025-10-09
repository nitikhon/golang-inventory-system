package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	mock_port "github.com/nitikhon/golang-inventory-system/internal/core/port/mock"
	mock_hash "github.com/nitikhon/golang-inventory-system/internal/util/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewUserService(t *testing.T) {
	// arrage
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)

	// act
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	// assert
	assert.NotNil(t, mockUserService, "expect userService, got nil")
	assert.NotNil(t, mockUserService.repo, "expect userRepo, got nil")
	assert.NotNil(t, mockUserService.crypto, "expect Crypto, got nil")
}

func TestCreateUser_Success(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "1234",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}

	mockUserRepo.EXPECT().GetUserByUsername(userInput.Username).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(nil, nil)
	mockHashApp.EXPECT().HashPassword(gomock.Any()).Return("hashedpassword", nil)
	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(&userInput, nil)

	// act
	user, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.Equal(t, &userInput, user)
	assert.Nil(t, err, fmt.Sprintf("expect nil, got %v", err))
}

func TestCreateUser_Failed_Validation_MissingField(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	tests := []struct {
		name      string
		userInput entity.User
		mockErr   error
	}{
		{
			name: "missing username",
			userInput: entity.User{
				Username:  "",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockErr: errors.New("username is required"),
		},
		{
			name: "missing email",
			userInput: entity.User{
				Username:  "test",
				Email:     "",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockErr: errors.New("email is required"),
		},
		{
			name: "missing password",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockErr: errors.New("password is required"),
		},
		{
			name: "missing phone",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockErr: errors.New("phone is required"),
		},
		{
			name: "missing first name",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "",
				LastName:  "Corner",
			},
			mockErr: errors.New("first name is required"),
		},
		{
			name: "missing last name",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "",
			},
			mockErr: errors.New("last name is required"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// act
			_, err := mockUserService.CreateUser(&tt.userInput)

			// assert
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestCreateUser_Failed_Validation_Email(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	userInput := entity.User{
		Username:  "test",
		Email:     "test.com",
		Password:  "1234",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}
	mockErr := errors.New("invalid email format")

	// act
	_, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.Equal(t, mockErr, err)
}

// validate both success and fail, because there are more than 1 success case
func TestCreateUser_Validation_Phone(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	tests := []struct {
		name       string
		userInput  entity.User
		mockReturn string // mock return only phone field
		mockErr    error
	}{
		{
			name: "success: phone number not contains hyphen(-)",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "0987654321",
			mockErr:    nil,
		},
		{
			name: "success: phone number contains hyphen(-)",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "098-765-4321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "0987654321",
			mockErr:    nil,
		},
		{
			name: "success: phone number less than 10 digits",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "78787",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("phone number must be at 10 digits"),
		},
		{
			name: "success: phone number more than 10 digits",
			userInput: entity.User{
				Username:  "test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "9999999999999999999999999999999999",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("phone number must be at 10 digits"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			if tt.mockErr == nil {
				mockUserRepo.EXPECT().GetUserByUsername(tt.userInput.Username).Return(nil, nil)
				mockUserRepo.EXPECT().GetUserByEmail(tt.userInput.Email).Return(nil, nil)
				mockUserRepo.EXPECT().GetUserByPhone(tt.mockReturn).Return(nil, nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(&tt.userInput, nil)
			}

			// act
			user, err := mockUserService.CreateUser(&tt.userInput)

			// assert
			if tt.mockErr == nil {
				assert.Equal(t, tt.mockReturn, user.Phone)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestCreateUser_Failed_Validation_Username(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	tests := []struct {
		name       string
		userInput  entity.User
		mockReturn string // mock return only username field
		mockErr    error
	}{
		{
			name: "contains @",
			userInput: entity.User{
				Username:  "H@l31y",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username can only contain alphabets, numbers, dots, and underscores"),
		},
		{
			name: "contains non-english",
			userInput: entity.User{
				Username:  "ทดสอบ",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username can only contain alphabets, numbers, dots, and underscores"),
		},
		{
			name: "username legnth is less than 3",
			userInput: entity.User{
				Username:  "a",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username must be 3-20 characters"),
		},
		{
			name: "username legnth is greater than 20",
			userInput: entity.User{
				Username:  "abcdefghijklmnopqrstuvwxyz",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username must be 3-20 characters"),
		},
		{
			name: "username not starts with a letter",
			userInput: entity.User{
				Username:  "999test",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username must start with a letter"),
		},
		{
			name: "username ends with '.'",
			userInput: entity.User{
				Username:  "test.",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username cannot start or end with '.' or '_'"),
		},
		{
			name: "username ends with '_'",
			userInput: entity.User{
				Username:  "test_",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username cannot start or end with '.' or '_'"),
		},
		{
			name: "username contains '__'",
			userInput: entity.User{
				Username:  "te__st",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username cannot contain '..', '__', '._', or '_.'"),
		},
		{
			name: "username contains '..'",
			userInput: entity.User{
				Username:  "te..st",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username cannot contain '..', '__', '._', or '_.'"),
		},
		{
			name: "username contains '._'",
			userInput: entity.User{
				Username:  "te._st",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username cannot contain '..', '__', '._', or '_.'"),
		},
		{
			name: "username contains '_.'",
			userInput: entity.User{
				Username:  "te_.st",
				Email:     "test@gmail.com",
				Password:  "1234",
				Phone:     "0987654321",
				FirstName: "John",
				LastName:  "Corner",
			},
			mockReturn: "",
			mockErr:    errors.New("username cannot contain '..', '__', '._', or '_.'"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			if tt.mockErr == nil {
				mockUserRepo.EXPECT().GetUserByUsername(tt.userInput.Username).Return(nil, nil)
				mockUserRepo.EXPECT().GetUserByEmail(tt.userInput.Email).Return(nil, nil)
				mockUserRepo.EXPECT().GetUserByPhone(tt.mockReturn).Return(nil, nil)
				mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(&tt.userInput, nil)
			}

			// act
			user, err := mockUserService.CreateUser(&tt.userInput)

			// assert
			if tt.mockErr == nil {
				assert.Equal(t, tt.mockReturn, user.Phone)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestCreateUser_UsernameExists(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "1234",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}
	mockErr := errors.New("a user with the provided credentials already exists")

	mockUserRepo.EXPECT().GetUserByUsername(gomock.Any()).Return(&userInput, nil)

	// act
	user, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.NotNil(t, user, fmt.Sprintf("expect existed user, got %v", user))
	assert.Equal(t, mockErr, err)
}

func TestCreateUser_EmailExists(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "1234",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}
	mockErr := errors.New("a user with the provided credentials already exists")

	mockUserRepo.EXPECT().GetUserByUsername(gomock.Any()).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(&userInput, nil)

	// act
	user, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.NotNil(t, user, fmt.Sprintf("expect existed user, got %v", user))
	assert.Equal(t, mockErr, err)
}

func TestCreateUser_PhoneExists(t *testing.T) {
	// arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_hash.NewMockCryptoUtil(ctrl)
	mockUserService := NewUserService(mockUserRepo, mockHashApp)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "1234",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}
	mockErr := errors.New("a user with the provided credentials already exists")

	mockUserRepo.EXPECT().GetUserByUsername(gomock.Any()).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByEmail(gomock.Any()).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByPhone(gomock.Any()).Return(&userInput, nil)

	// act
	user, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.NotNil(t, user, fmt.Sprintf("expect existed user, got %v", user))
	assert.Equal(t, mockErr, err)
}
