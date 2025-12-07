package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	mock_port "github.com/nitikhon/golang-inventory-system/internal/core/port/mock"
	mock_util "github.com/nitikhon/golang-inventory-system/internal/util/mock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupUserServiceMock(t *testing.T) (*UserService, *mock_port.MockUserRepository, *mock_util.MockCryptoUtil, *mock_util.MockJWTUtil) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockHashApp := mock_util.NewMockCryptoUtil(ctrl)
	mockJWTApp := mock_util.NewMockJWTUtil(ctrl)

	service := NewUserService(mockUserRepo, mockHashApp, mockJWTApp)
	return service, mockUserRepo, mockHashApp, mockJWTApp
}

func TestNewUserService(t *testing.T) {
	// arrange
	mockUserService, _, _, _ := setupUserServiceMock(t)

	// assert
	assert.NotNil(t, mockUserService, "expect userService, got nil")
}

func TestCreateUser_Success(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, mockHashApp, _ := setupUserServiceMock(t)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "P@ssw0rd",
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

func TestCreateUser_Failed_HashedFail(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, mockHashApp, _ := setupUserServiceMock(t)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "P@ssw0rd",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}
	mockErr := errors.New("hash error")

	mockUserRepo.EXPECT().GetUserByUsername(userInput.Username).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(nil, nil)
	mockHashApp.EXPECT().HashPassword(gomock.Any()).Return("", mockErr)

	// act
	user, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.Equal(t, &entity.User{}, user)
	assert.EqualError(t, err, mockErr.Error())
}

func TestCreateUser_UsernameExists(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "P@ssw0rd",
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
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "P@ssw0rd",
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
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "P@ssw0rd",
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

func TestCreateUser_RepoError(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, mockHashApp, _ := setupUserServiceMock(t)

	userInput := entity.User{
		Username:  "test",
		Email:     "test@gmail.com",
		Password:  "P@ssw0rd",
		Phone:     "0987654321",
		FirstName: "John",
		LastName:  "Corner",
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByUsername(userInput.Username).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(nil, nil)
	mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(nil, nil)
	mockHashApp.EXPECT().HashPassword(gomock.Any()).Return("hashedpassword", nil)
	mockUserRepo.EXPECT().CreateUser(gomock.Any()).Return(&entity.User{}, mockErr)

	// act
	user, err := mockUserService.CreateUser(&userInput)

	// assert
	assert.Equal(t, &entity.User{}, user)
	assert.Equal(t, mockErr, err)
}

func TestUpdateUser(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		// arrange
		mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

		userInput := entity.User{
			Model:     gorm.Model{ID: 1},
			Email:     "test@gmail.com",
			Password:  "P@ssw0rd",
			Phone:     "0987654321",
			FirstName: "John",
			LastName:  "Corner",
		}

		mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(&userInput, nil)

		// act
		user, err := mockUserService.UpdateUser(&userInput)

		// assert
		assert.NotEqual(t, &entity.User{}, user)
		assert.Nil(t, err)
	})

	t.Run("success case - same user's email and phone", func(t *testing.T) {
		// arrange
		mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

		userInput := entity.User{
			Model:     gorm.Model{ID: 1},
			Email:     "test@gmail.com",
			Password:  "P@ssw0rd",
			Phone:     "0987654321",
			FirstName: "John",
			LastName:  "Corner",
		}

		// same user found (same ID) - should not block update
		mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(&userInput, nil)
		mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(&userInput, nil)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(&userInput, nil)

		// act
		user, err := mockUserService.UpdateUser(&userInput)

		// assert
		assert.NotEqual(t, &entity.User{}, user)
		assert.Nil(t, err)
	})

	t.Run("duplicate email error", func(t *testing.T) {
		// arrange
		mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

		userInput := entity.User{
			Model:     gorm.Model{ID: 1},
			Email:     "existing@gmail.com",
			Password:  "P@ssw0rd",
			Phone:     "0987654321",
			FirstName: "John",
			LastName:  "Corner",
		}

		existingUser := entity.User{
			Model: gorm.Model{ID: 2}, // different user
			Email: "existing@gmail.com",
		}

		mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(&existingUser, nil)

		// act
		user, err := mockUserService.UpdateUser(&userInput)

		// assert
		assert.Equal(t, &entity.User{}, user)
		assert.EqualError(t, err, "email already taken")
	})

	t.Run("duplicate phone error", func(t *testing.T) {
		// arrange
		mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

		userInput := entity.User{
			Model:     gorm.Model{ID: 1},
			Email:     "test@gmail.com",
			Password:  "P@ssw0rd",
			Phone:     "0987654321",
			FirstName: "John",
			LastName:  "Corner",
		}

		existingUser := entity.User{
			Model: gorm.Model{ID: 2}, // different user
			Phone: "0987654321",
		}

		mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(&existingUser, nil)

		// act
		user, err := mockUserService.UpdateUser(&userInput)

		// assert
		assert.Equal(t, &entity.User{}, user)
		assert.EqualError(t, err, "phone already taken")
	})

	t.Run("repo error", func(t *testing.T) {
		// arrange
		mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

		userInput := entity.User{
			Model:     gorm.Model{ID: 1},
			Email:     "test@gmail.com",
			Password:  "P@ssw0rd",
			Phone:     "0987654321",
			FirstName: "John",
			LastName:  "Corner",
		}

		mockUserRepo.EXPECT().GetUserByEmail(userInput.Email).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.EXPECT().GetUserByPhone(userInput.Phone).Return(nil, gorm.ErrRecordNotFound)
		mockUserRepo.EXPECT().UpdateUser(gomock.Any()).Return(&entity.User{}, errors.New("database error"))

		// act
		user, err := mockUserService.UpdateUser(&userInput)

		// assert
		assert.Equal(t, &entity.User{}, user)
		assert.EqualError(t, err, "database error")
	})
}

func TestDeleteUser(t *testing.T) {
	// arrange
	tests := []struct {
		name    string
		userID  uint
		mockErr error
	}{
		{
			name:   "success case",
			userID: 1,
		},
		{
			name:    "fail case (error from repo)",
			userID:  999,
			mockErr: gorm.ErrRecordNotFound,
		},
	}

	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.EXPECT().DeleteUser(tt.userID).Return(tt.mockErr)

			// act
			err := mockUserService.DeleteUser(tt.userID)

			// arrange
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestGetAllUsers(t *testing.T) {
	// arrange
	tests := []struct {
		name    string
		mockErr error
	}{
		{
			name: "success case",
		},
		{
			name:    "fail case (error from repo)",
			mockErr: schema.ErrUnsupportedDataType,
		},
	}

	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.EXPECT().GetAllUsers().Return([]*entity.User{}, tt.mockErr)

			// act
			users, err := mockUserService.GetAllUsers()

			// assert
			if tt.mockErr == nil {
				assert.Equal(t, []*entity.User{}, users)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestGetUserByID(t *testing.T) {
	// arrange
	tests := []struct {
		name    string
		userID  uint
		mockErr error
	}{
		{
			name:   "success case",
			userID: 1,
		},
		{
			name:    "fail case (error from repo)",
			userID:  999,
			mockErr: schema.ErrUnsupportedDataType,
		},
	}

	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.EXPECT().
				GetUserByID(tt.userID).
				Return(&entity.User{Model: gorm.Model{ID: tt.userID}}, tt.mockErr)

			// act
			user, err := mockUserService.GetUserByID(tt.userID)

			// assert
			if tt.mockErr == nil {
				assert.NotEqual(t, &entity.User{}, user)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestGetUserByUsername(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	tests := []struct {
		name        string
		username    string
		expectCalls func()
		mockErr     error
	}{
		{
			name:     "success case",
			username: "testuser",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("testuser").
					Return(&entity.User{Username: "testuser"}, nil)
			},
		},
		{
			name:     "repo error",
			username: "testuser",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("testuser").
					Return(nil, gorm.ErrRecordNotFound)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls()

			// act
			user, err := mockUserService.GetUserByUsername(tt.username)

			// assert
			if tt.mockErr != nil {
				assert.Nil(t, user)
			} else {
				assert.Equal(t, tt.username, user.Username)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	tests := []struct {
		name        string
		email       string
		expectCalls func()
		mockErr     error
	}{
		{
			name:  "success case",
			email: "test@gmail.com",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByEmail("test@gmail.com").
					Return(&entity.User{Email: "test@gmail.com"}, nil)
			},
		},
		{
			name:  "repo error",
			email: "test@gmail.com",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByEmail("test@gmail.com").
					Return(nil, gorm.ErrRecordNotFound)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls()

			// act
			user, err := mockUserService.GetUserByEmail(tt.email)

			// assert
			if tt.mockErr != nil {
				assert.Nil(t, user)
			} else {
				assert.Equal(t, tt.email, user.Email)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestGetUserByPhone(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	tests := []struct {
		name        string
		phone       string
		expectCalls func()
		mockErr     error
	}{
		{
			name:  "success case",
			phone: "0987654321",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByPhone("0987654321").
					Return(&entity.User{Phone: "0987654321"}, nil)
			},
		},
		{
			name:  "repo error",
			phone: "0987654321",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByPhone("0987654321").
					Return(nil, gorm.ErrRecordNotFound)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls()

			// act
			user, err := mockUserService.GetUserByPhone(tt.phone)

			// assert
			if tt.mockErr != nil {
				assert.Nil(t, user)
			} else {
				assert.Equal(t, tt.phone, user.Phone)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestLogin(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, mockHashApp, mockJWTApp := setupUserServiceMock(t)

	t.Setenv("ACCESS_TOKEN_SECRET", "somesecreta")
	t.Setenv("REFRESH_TOKEN_SECRET", "somesecretb")

	mockUser := entity.User{Model: gorm.Model{ID: 1}}

	tests := []struct {
		name        string
		username    string
		password    string
		expectCalls func()
		mockErr     error
	}{
		{
			name:     "success case",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByUsername("tester").
						Return(&mockUser, nil),
					mockHashApp.EXPECT().
						CheckPasswordHash(gomock.Any(), "P@ssw0rd").
						Return(nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("accessToken", nil),
					mockJWTApp.EXPECT().
						GenerateRefreshToken(mockUser).
						Return("refreshToken", nil),
					mockUserRepo.EXPECT().
						UpdateRefreshToken(uint(1), gomock.Any()).
						Return(nil),
				)
			},
		},
		{
			name:     "error from database",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("tester").
					Return(&entity.User{}, errors.New("some db error"))
			},
			mockErr: errors.New("some db error"),
		},
		{
			name:     "username not found",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				mockUserRepo.EXPECT().
					GetUserByUsername("tester").
					Return(nil, nil)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:     "incorrect password",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByUsername("tester").
						Return(&mockUser, nil),
					mockHashApp.EXPECT().
						CheckPasswordHash(gomock.Any(), "P@ssw0rd").
						Return(bcrypt.ErrMismatchedHashAndPassword),
				)
			},
			mockErr: errors.New("invalid credentials"),
		},
		{
			name:     "errors at GenerateAccessToken",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByUsername("tester").
						Return(&mockUser, nil),
					mockHashApp.EXPECT().
						CheckPasswordHash(gomock.Any(), "P@ssw0rd").
						Return(nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("", errors.New("GenerateAccessToken error")),
				)
			},
			mockErr: errors.New("GenerateAccessToken error"),
		},
		{
			name:     "errors at GenerateRefreshToken",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByUsername("tester").
						Return(&mockUser, nil),
					mockHashApp.EXPECT().
						CheckPasswordHash(gomock.Any(), "P@ssw0rd").
						Return(nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("accessToken", nil),
					mockJWTApp.EXPECT().
						GenerateRefreshToken(mockUser).
						Return("", errors.New("GenerateRefreshToken error")),
				)
			},
			mockErr: errors.New("GenerateRefreshToken error"),
		},
		{
			name:     "errors at UpdateRefreshToken",
			username: "tester",
			password: "P@ssw0rd",
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByUsername("tester").
						Return(&mockUser, nil),
					mockHashApp.EXPECT().
						CheckPasswordHash(gomock.Any(), "P@ssw0rd").
						Return(nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("accessToken", nil),
					mockJWTApp.EXPECT().
						GenerateRefreshToken(mockUser).
						Return("refreshToken", nil),
					mockUserRepo.EXPECT().
						UpdateRefreshToken(uint(1), "refreshToken").
						Return(gorm.ErrRecordNotFound),
				)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls()

			// act
			accessToken, refreshToken, err := mockUserService.Login(tt.username, tt.password)

			// assert
			if tt.mockErr == nil {
				assert.NotEmpty(t, accessToken)
				assert.NotEmpty(t, refreshToken)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestRefreshToken(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, _, mockJWTApp := setupUserServiceMock(t)

	t.Setenv("ACCESS_TOKEN_SECRET", "somesecreta")
	t.Setenv("REFRESH_TOKEN_SECRET", "somesecretb")

	mockRefreshToken := "refreshToken"
	mockUser := entity.User{Model: gorm.Model{ID: 1}, RefreshToken: mockRefreshToken}

	tests := []struct {
		name             string
		userID           uint
		mockRefreshToken string
		expectCalls      func()
		mockErr          error
	}{
		{
			name:             "success case",
			mockRefreshToken: mockRefreshToken,
			userID:           1,
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken(mockRefreshToken).
						Return(mockUser.ID, nil),
					mockUserRepo.EXPECT().
						GetUserByID(mockUser.ID).
						Return(&mockUser, nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("newAccessToken", nil),
					mockJWTApp.EXPECT().
						GenerateRefreshToken(mockUser).
						Return("newRefreshToken", nil),
					mockUserRepo.EXPECT().
						UpdateRefreshToken(mockUser.ID, "newRefreshToken").
						Return(nil),
				)
			},
		},
		{
			name:             "validate refresh token failed",
			userID:           1,
			mockRefreshToken: mockRefreshToken,
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken(mockRefreshToken).
						Return(uint(0), errors.New("invalid token")),
				)
			},
			mockErr: errors.New("invalid token"),
		},
		{
			name:             "error at GetUserByID",
			userID:           1,
			mockRefreshToken: mockRefreshToken,
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken(mockRefreshToken).
						Return(mockUser.ID, nil),
					mockUserRepo.EXPECT().
						GetUserByID(mockUser.ID).
						Return(&entity.User{}, errors.New("error from GetUserByID")),
				)
			},
			mockErr: errors.New("error from GetUserByID"),
		},
		{
			name:             "invalid refresh token",
			userID:           1,
			mockRefreshToken: "invalidToken",
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken("invalidToken").
						Return(mockUser.ID, nil),
					mockUserRepo.EXPECT().
						GetUserByID(mockUser.ID).
						Return(&mockUser, nil),
				)
			},
			mockErr: errors.New("invalid refresh token"),
		},
		{
			name:             "error at GenerateAccessToken",
			mockRefreshToken: mockRefreshToken,
			userID:           1,
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken(mockRefreshToken).
						Return(mockUser.ID, nil),
					mockUserRepo.EXPECT().
						GetUserByID(mockUser.ID).
						Return(&mockUser, nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("", errors.New("GenerateAccessToken error")),
				)
			},
			mockErr: errors.New("GenerateAccessToken error"),
		},
		{
			name:             "error at GenerateRefreshToken",
			mockRefreshToken: mockRefreshToken,
			userID:           1,
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken(mockRefreshToken).
						Return(mockUser.ID, nil),
					mockUserRepo.EXPECT().
						GetUserByID(mockUser.ID).
						Return(&mockUser, nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("newAccessToken", nil),
					mockJWTApp.EXPECT().
						GenerateRefreshToken(mockUser).
						Return("", errors.New("GenerateRefreshToken error")),
				)
			},
			mockErr: errors.New("GenerateRefreshToken error"),
		},
		{
			name:             "error at UpdateRefreshToken",
			mockRefreshToken: mockRefreshToken,
			userID:           1,
			expectCalls: func() {
				gomock.InOrder(
					mockJWTApp.EXPECT().
						ValidateRefreshToken(mockRefreshToken).
						Return(mockUser.ID, nil),
					mockUserRepo.EXPECT().
						GetUserByID(mockUser.ID).
						Return(&mockUser, nil),
					mockJWTApp.EXPECT().
						GenerateAccessToken(mockUser).
						Return("newAccessToken", nil),
					mockJWTApp.EXPECT().
						GenerateRefreshToken(mockUser).
						Return("newRefreshToken", nil),
					mockUserRepo.EXPECT().
						UpdateRefreshToken(mockUser.ID, "newRefreshToken").
						Return(gorm.ErrRecordNotFound),
				)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls()

			// act
			tokens, err := mockUserService.RefreshToken(tt.mockRefreshToken)

			// assert
			if tt.mockErr == nil {
				assert.NotEmpty(t, tokens)
				assert.NotEmpty(t, tokens.AccessToken)
				assert.NotEmpty(t, tokens.RefreshToken)
			}
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestLogout(t *testing.T) {
	// arrange
	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	mockUserWithRefreshToken := &entity.User{
		Model:        gorm.Model{ID: 1},
		RefreshToken: "refreshToken",
	}

	tests := []struct {
		name        string
		userID      uint
		expectCalls func()
		mockErr     error
	}{
		{
			name:   "success case",
			userID: 1,
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByID(uint(1)).
						Return(mockUserWithRefreshToken, nil),
					mockUserRepo.EXPECT().
						UpdateRefreshToken(uint(1), "").
						Return(nil),
				)
			},
		},
		{
			name:   "error at GetUserById",
			userID: 1,
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByID(uint(1)).
						Return(&entity.User{}, errors.New("GetUserByID error")),
				)
			},
			mockErr: errors.New("GetUserByID error"),
		},
		{
			name:   "user not found from GetUserById",
			userID: 1,
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByID(uint(1)).
						Return(nil, nil),
				)
			},
			mockErr: gorm.ErrRecordNotFound,
		},
		{
			name:   "error at UpdateUser",
			userID: 1,
			expectCalls: func() {
				gomock.InOrder(
					mockUserRepo.EXPECT().
						GetUserByID(uint(1)).
						Return(mockUserWithRefreshToken, nil),
					mockUserRepo.EXPECT().
						UpdateRefreshToken(uint(1), "").
						Return(errors.New("UpdateRefreshToken Error")),
				)
			},
			mockErr: errors.New("UpdateRefreshToken Error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls()

			// act
			err := mockUserService.Logout(tt.userID)

			// assert
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestUpdateUserProfile(t *testing.T) {
	// arrange
	tests := []struct {
		name        string
		userInput   entity.User
		expectCalls func(mockUserRepo *mock_port.MockUserRepository)
		mockErr     error
	}{
		{
			name: "success case",
			userInput: entity.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "0987654321",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				updatedFields := map[string]any{
					"first_name": "John",
					"last_name":  "Doe",
					"phone":      "0987654321",
				}
				mockUserRepo.EXPECT().
					UpdateUserProfile(uint(1), updatedFields).
					Return(&entity.User{Model: gorm.Model{ID: 1}}, nil)
				mockUserRepo.EXPECT().
					GetUserByID(uint(1)).
					Return(&entity.User{
						Model:     gorm.Model{ID: 1},
						FirstName: "John",
						LastName:  "Doe",
						Phone:     "0987654321",
					}, nil)
			},
		},
		{
			name: "repo error - UpdateUserProfile fails",
			userInput: entity.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "0987654321",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				updatedFields := map[string]any{
					"first_name": "John",
					"last_name":  "Doe",
					"phone":      "0987654321",
				}
				mockUserRepo.EXPECT().
					UpdateUserProfile(uint(1), updatedFields).
					Return(&entity.User{}, errors.New("database error"))
			},
			mockErr: errors.New("database error"),
		},
		{
			name: "repo error - GetUserByID fails",
			userInput: entity.User{
				Model:     gorm.Model{ID: 1},
				FirstName: "John",
				LastName:  "Doe",
				Phone:     "0987654321",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				updatedFields := map[string]any{
					"first_name": "John",
					"last_name":  "Doe",
					"phone":      "0987654321",
				}
				mockUserRepo.EXPECT().
					UpdateUserProfile(uint(1), updatedFields).
					Return(&entity.User{Model: gorm.Model{ID: 1}}, nil)
				mockUserRepo.EXPECT().
					GetUserByID(uint(1)).
					Return(&entity.User{}, errors.New("user not found"))
			},
			mockErr: errors.New("error while trying to get updated user"),
		},
	}

	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls(mockUserRepo)

			// act
			user, err := mockUserService.UpdateUserProfile(&tt.userInput)

			// assert
			if tt.mockErr == nil {
				assert.NotEqual(t, &entity.User{}, user)
				assert.Nil(t, err)
			} else {
				assert.Equal(t, &entity.User{}, user)
				assert.Equal(t, tt.mockErr, err)
			}
		})
	}
}

func TestUpdateUserPassword(t *testing.T) {
	// arrange
	tests := []struct {
		name        string
		userInput   entity.User
		expectCalls func(mockUserRepo *mock_port.MockUserRepository, mockHashApp *mock_util.MockCryptoUtil)
		mockErr     error
	}{
		{
			name: "success case",
			userInput: entity.User{
				Model:    gorm.Model{ID: 1},
				Password: "P@ssw0rd",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository, mockHashApp *mock_util.MockCryptoUtil) {
				mockHashApp.EXPECT().
					HashPassword("P@ssw0rd").
					Return("hashedpassword", nil)
				mockUserRepo.EXPECT().
					UpdateUserPassword(uint(1), "hashedpassword").
					Return(nil)
			},
		},
		{
			name: "hash error",
			userInput: entity.User{
				Model:    gorm.Model{ID: 1},
				Password: "P@ssw0rd",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository, mockHashApp *mock_util.MockCryptoUtil) {
				mockHashApp.EXPECT().
					HashPassword("P@ssw0rd").
					Return("", errors.New("hash error"))
			},
			mockErr: errors.New("hash error"),
		},
		{
			name: "repo error",
			userInput: entity.User{
				Model:    gorm.Model{ID: 1},
				Password: "P@ssw0rd",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository, mockHashApp *mock_util.MockCryptoUtil) {
				mockHashApp.EXPECT().
					HashPassword("P@ssw0rd").
					Return("hashedpassword", nil)
				mockUserRepo.EXPECT().
					UpdateUserPassword(uint(1), "hashedpassword").
					Return(errors.New("database error"))
			},
			mockErr: errors.New("database error"),
		},
	}

	mockUserService, mockUserRepo, mockHashApp, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls(mockUserRepo, mockHashApp)

			// act
			err := mockUserService.UpdateUserPassword(&tt.userInput)

			// assert
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestUpdateUserEmail(t *testing.T) {
	// arrange
	tests := []struct {
		name        string
		userInput   entity.User
		expectCalls func(mockUserRepo *mock_port.MockUserRepository)
		mockErr     error
	}{
		{
			name: "success case",
			userInput: entity.User{
				Model: gorm.Model{ID: 1},
				Email: "test@example.com",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				mockUserRepo.EXPECT().
					UpdateUserEmail(uint(1), "test@example.com").
					Return(nil)
			},
		},
		{
			name: "repo error",
			userInput: entity.User{
				Model: gorm.Model{ID: 1},
				Email: "test@example.com",
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				mockUserRepo.EXPECT().
					UpdateUserEmail(uint(1), "test@example.com").
					Return(errors.New("database error"))
			},
			mockErr: errors.New("database error"),
		},
	}

	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls(mockUserRepo)

			// act
			err := mockUserService.UpdateUserEmail(&tt.userInput)

			// assert
			assert.Equal(t, tt.mockErr, err)
		})
	}
}

func TestUpdateUserAdminStatus(t *testing.T) {
	// arrange
	tests := []struct {
		name        string
		userInput   entity.User
		expectCalls func(mockUserRepo *mock_port.MockUserRepository)
		mockErr     error
	}{
		{
			name: "success case - set admin to true",
			userInput: entity.User{
				Model:   gorm.Model{ID: 1},
				IsAdmin: true,
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				mockUserRepo.EXPECT().
					UpdateUserAdminStatus(uint(1), true).
					Return(nil)
			},
		},
		{
			name: "success case - set admin to false",
			userInput: entity.User{
				Model:   gorm.Model{ID: 1},
				IsAdmin: false,
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				mockUserRepo.EXPECT().
					UpdateUserAdminStatus(uint(1), false).
					Return(nil)
			},
		},
		{
			name: "repo error",
			userInput: entity.User{
				Model:   gorm.Model{ID: 1},
				IsAdmin: true,
			},
			expectCalls: func(mockUserRepo *mock_port.MockUserRepository) {
				mockUserRepo.EXPECT().
					UpdateUserAdminStatus(uint(1), true).
					Return(errors.New("database error"))
			},
			mockErr: errors.New("database error"),
		},
	}

	mockUserService, mockUserRepo, _, _ := setupUserServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			tt.expectCalls(mockUserRepo)

			// act
			err := mockUserService.UpdateUserAdminStatus(&tt.userInput)

			// assert
			assert.Equal(t, tt.mockErr, err)
		})
	}
}
