package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/letenk/use_deal_user/helper"
	"github.com/letenk/use_deal_user/models/domain"
	"github.com/letenk/use_deal_user/models/web"
	"github.com/letenk/use_deal_user/repository"
	"github.com/letenk/use_deal_user/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	dataUser := web.UserCreateRequest{
		Fullname: fullname,
		Username: strings.ToLower(helper.RandomPerson()),
		Password: "password",
		Role:     "admin",
	}

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)

	// Create new user
	_, err := service.Create(dataUser)
	helper.ErrLogPanic(err)

	// Test case
	testCases := []struct {
		name string
		req  web.UserLoginRequest
	}{
		{
			name: "failed_login_wrong_username",
			req: web.UserLoginRequest{
				Username: "wrong",
				Password: dataUser.Password,
			},
		},
		{
			name: "failed_login_wrong_password",
			req: web.UserLoginRequest{
				Username: dataUser.Username,
				Password: "wrong",
			},
		},
		{
			name: "success_login",
			req: web.UserLoginRequest{
				Username: dataUser.Username,
				Password: dataUser.Password,
			},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Login
			token, err := service.Login(tc.req)

			if tc.name == "failed_login_wrong_username" || tc.name == "failed_login_wrong_password" {
				assert.Empty(t, token)
				assert.Error(t, err)
				assert.Equal(t, "email or password incorrect", err.Error())
			} else {
				assert.NotEmpty(t, token)
				assert.NoError(t, err)
			}
		})

	}

}

func TestCreateUserService(t *testing.T) {
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())

	// Test case
	testCases := []struct {
		name string
		user web.UserCreateRequest
	}{
		{
			name: "success_with_field_role_admin",
			user: web.UserCreateRequest{
				Fullname: fullname,
				Username: strings.ToLower(helper.RandomPerson()),
				Password: "password1",
				Role:     "admin",
			},
		},
		{
			name: "success_with_field_role_user",
			user: web.UserCreateRequest{
				Fullname: fullname,
				Username: strings.ToLower(helper.RandomPerson()),
				Password: "password2",
				Role:     "user",
			},
		},
		{
			name: "success_without_field_role",
			user: web.UserCreateRequest{
				Fullname: fullname,
				Username: strings.ToLower(helper.RandomPerson()),
				Password: "password3",
			},
		},
	}

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)
	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userID, err := service.Create(tc.user)
			helper.ErrLogPanic(err)

			// Get one use repository
			dataUser, err := repository.GetOne(userID)

			if tc.name == "success_with_field_role_admin" {
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)
				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)

				assert.Equal(t, tc.user.Fullname, dataUser.Fullname)
				assert.Equal(t, tc.user.Username, dataUser.Username)
				assert.Equal(t, "admin", dataUser.Role)

				// Compare password
				err := bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(tc.user.Password))
				assert.NoError(t, err)
			} else if tc.name == "success_with_field_role_user" {
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)
				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)

				assert.Equal(t, tc.user.Fullname, dataUser.Fullname)
				assert.Equal(t, tc.user.Username, dataUser.Username)
				assert.Equal(t, "user", dataUser.Role)

				// Compare password
				err := bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(tc.user.Password))
				assert.NoError(t, err)
			} else {
				// If role is empty, set default role "user"
				assert.NoError(t, err)
				assert.NotEmpty(t, userID)
				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)

				assert.Equal(t, tc.user.Fullname, dataUser.Fullname)
				assert.Equal(t, tc.user.Username, dataUser.Username)
				assert.Equal(t, "user", dataUser.Role)

				// Compare password
				err := bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(tc.user.Password))
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetAllUserService(t *testing.T) {
	t.Parallel()

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)

	// Get all
	users, err := service.GetAll()
	helper.ErrLogPanic(err)

	// Pass
	assert.NoError(t, err)
	assert.NotZero(t, len(users))

	for _, data := range users {
		assert.NotEmpty(t, data.ID)
		assert.NotEmpty(t, data.Fullname)
		assert.NotEmpty(t, data.Username)
		assert.NotEmpty(t, data.Password)
		assert.NotEmpty(t, data.Role)
		assert.NotEmpty(t, data.CreatedAt)
		assert.NotEmpty(t, data.UpdatedAt)
	}
}

func TestGetOneUserService(t *testing.T) {
	t.Parallel()

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)

	// Create new sample object user
	newUser := domain.User{}
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	newUser.Fullname = fullname
	newUser.Username = strings.ToLower(helper.RandomPerson())
	newUser.Role = "admin"
	// Hash password
	password := "password"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	helper.ErrLogPanic(err)
	newUser.Password = string(passwordHash)
	// Create one user use repository
	id, err := repository.Insert(newUser)
	helper.ErrLogPanic(err)
	// End reate new sample object user

	testCases := []struct {
		name   string
		userID string
	}{
		{
			name:   "success_get_one_by_username",
			userID: id,
		},
		{
			name:   "failed_get_one_not_found",
			userID: "4112d578-6163-11ed-9b6a-0242ac120002",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		// Get one
		user, err := service.GetOne(tc.userID)

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.name == "success_get_one_by_username" {
				// Pass
				assert.NoError(t, err)

				assert.NotEmpty(t, user.ID)
				assert.Equal(t, newUser.Fullname, user.Fullname)
				assert.Equal(t, newUser.Username, user.Username)
				assert.Equal(t, newUser.Role, user.Role)

				assert.NotEmpty(t, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)

				// Compare password
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
				respErrorMessage := fmt.Sprintf("user with ID %s Not Found", tc.userID)
				assert.Equal(t, respErrorMessage, err.Error())
			}
		})

	}

}

func TestGetOneByUsernameUserService(t *testing.T) {
	t.Parallel()

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)

	// Create new sample object user
	newUser := domain.User{}
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	newUser.Fullname = fullname
	newUser.Username = strings.ToLower(helper.RandomPerson())
	newUser.Role = "admin"
	// Hash password
	password := "password"
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	helper.ErrLogPanic(err)
	newUser.Password = string(passwordHash)
	// Create one user use repository
	id, err := repository.Insert(newUser)
	helper.ErrLogPanic(err)
	// End reate new sample object user

	// Get one
	user, err := service.GetOne(id)
	helper.ErrLogPanic(err)

	testCases := []struct {
		name     string
		username string
	}{
		{
			name:     "success_get_one_by_username",
			username: user.Username,
		},
		{
			name:     "failed_get_one_by_username_not_found",
			username: "wrong",
		},
	}

	for i := range testCases {
		tc := testCases[i]

		// Get one by username
		userByUsername, err := service.GetOneByUsername(tc.username)

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.name == "success_get_one_by_username" {
				// Pass
				assert.NoError(t, err)

				assert.Equal(t, userByUsername.ID, id)
				assert.Equal(t, userByUsername.Fullname, user.Fullname)
				assert.Equal(t, userByUsername.Username, user.Username)
				assert.Equal(t, userByUsername.Role, user.Role)

				assert.NotEmpty(t, userByUsername.CreatedAt)
				assert.NotEmpty(t, userByUsername.UpdatedAt)

				// Compare password
				err = bcrypt.CompareHashAndPassword([]byte(userByUsername.Password), []byte(password))
				assert.NoError(t, err)
			} else {
				// Pass
				assert.Error(t, err)
				respErrorMessage := fmt.Sprintf("user with username %s Not Found", tc.username)
				assert.Equal(t, respErrorMessage, err.Error())
			}
		})
	}
}

func TestUpdateUserService(t *testing.T) {
	// Test Case
	testCases := []struct {
		name string
		req  web.UserUpdateRequest
	}{
		{
			name: "success_updated_field_fullname_only",
			req: web.UserUpdateRequest{
				Fullname: helper.RandomPerson(),
			},
		},
		{
			name: "success_updated_field_password_only",
			req: web.UserUpdateRequest{
				Password: "updated",
			},
		},
		{
			name: "success_updated_field_role_only",
			req: web.UserUpdateRequest{
				Role: "user",
			},
		},
		{
			name: "success_updated_all_field_update_request",
			req: web.UserUpdateRequest{
				Fullname: helper.RandomPerson(),
				Password: helper.RandomString(10),
				Role:     "user",
			},
		},
		{
			name: "failed_update_when_id_not_found",
			req: web.UserUpdateRequest{
				Fullname: helper.RandomPerson(),
				Password: helper.RandomString(10),
				Role:     "user",
			},
		},
	}

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)

	// Test
	for i := range testCases {
		tc := testCases[i]
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			// Create sample newUser
			fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
			newUser := web.UserCreateRequest{
				Fullname: fullname,
				Username: strings.ToLower(helper.RandomPerson()),
				Password: "password",
			}
			userID, err := service.Create(newUser)
			helper.ErrLogPanic(err)

			if tc.name == "success_updated_field_fullname_only" {
				// Test update
				ok, err := service.Update(userID, tc.req)

				assert.NoError(t, err)
				assert.True(t, ok)

				// Check each field must match
				// Get One
				dataUser, err := repository.GetOne(userID)
				helper.ErrLogPanic(err)

				// Fullname must be same
				assert.Equal(t, tc.req.Fullname, dataUser.Fullname)

				assert.NotEmpty(t, dataUser.Role)
				assert.NotEmpty(t, dataUser.Password)
				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)

			} else if tc.name == "success_updated_field_password_only" {
				// Test update
				ok, err := service.Update(userID, tc.req)

				assert.NoError(t, err)
				assert.True(t, ok)

				// Check each field must match
				// Get One
				dataUser, err := repository.GetOne(userID)
				helper.ErrLogPanic(err)

				// Compare password must be not same
				err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte("updated"))
				assert.NoError(t, err)

				assert.NotEmpty(t, dataUser.Fullname)
				assert.NotEmpty(t, dataUser.Role)
				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)

			} else if tc.name == "success_updated_field_role_only" {
				// Test update
				ok, err := service.Update(userID, tc.req)

				assert.NoError(t, err)
				assert.True(t, ok)

				// Check each field must match
				// Get One
				dataUser, err := repository.GetOne(userID)
				helper.ErrLogPanic(err)

				// Role must be same
				assert.Equal(t, tc.req.Role, dataUser.Role)

				assert.NotEmpty(t, dataUser.Fullname)
				assert.NotEmpty(t, dataUser.Password)
				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)

			} else if tc.name == "failed_update_when_id_not_found" {
				// Test update
				wrongID := "4112d578-6163-11ed-9b6a-0242ac120002"
				ok, err := service.Update(wrongID, tc.req)

				assert.Error(t, err)
				respErrMessage := fmt.Sprintf("user with ID %s Not Found", wrongID)
				assert.Equal(t, respErrMessage, err.Error())
				assert.False(t, ok)

			} else {
				// Test update
				ok, err := service.Update(userID, tc.req)
				assert.NoError(t, err)
				assert.True(t, ok)

				// Check each field must match
				// Get One
				dataUser, err := repository.GetOne(userID)
				helper.ErrLogPanic(err)

				// Fullname, Role and password must be same
				assert.Equal(t, tc.req.Fullname, dataUser.Fullname)
				assert.Equal(t, tc.req.Role, dataUser.Role)
				// Compare password must same
				err = bcrypt.CompareHashAndPassword([]byte(dataUser.Password), []byte(tc.req.Password))
				assert.NoError(t, err)

				assert.NotEmpty(t, dataUser.CreatedAt)
				assert.NotEmpty(t, dataUser.UpdatedAt)
			}
		})
	}
}

func TestDeleteUserService(t *testing.T) {
	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)

	// Create sample newUser
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	newUser := web.UserCreateRequest{
		Fullname: fullname,
		Username: strings.ToLower(helper.RandomPerson()),
		Password: "password",
	}
	userID, err := service.Create(newUser)
	helper.ErrLogPanic(err)

	// Test case
	testCase := []struct {
		name   string
		userID string
	}{
		{
			name:   "success_deleted",
			userID: userID,
		},
		{
			name:   "failed_delete_user_not_found",
			userID: "4112d578-6163-11ed-9b6a-0242ac120002",
		},
	}

	// Test
	for i := range testCase {
		tc := testCase[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ok, err := service.Delete(tc.userID)

			if tc.name == "success_deleted" {
				assert.True(t, ok)
				assert.NoError(t, err)
			} else {
				assert.False(t, ok)
				assert.Error(t, err)
				respErrMessage := fmt.Sprintf("user with ID %s Not Found", tc.userID)
				assert.Equal(t, respErrMessage, err.Error())
			}
		})
	}
}
