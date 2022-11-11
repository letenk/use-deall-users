package tests

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/letenk/use_deal_user/helper"
	"github.com/letenk/use_deal_user/models/web"
	"github.com/letenk/use_deal_user/repository"
	"github.com/letenk/use_deal_user/service"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUserHandler(t *testing.T) {
	// Create new user
	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	dataUser := web.UserCreateRequest{
		Fullname: fullname,
		Username: strings.ToLower(helper.RandomPerson()),
		Password: "password",
		Role:     "admin",
	}

	repository := repository.NewUserRepository(ConnTest)
	service := service.NewServiceUser(repository)
	service.Create(dataUser)
	// End Create new user

	// Test Cases
	testCases := []struct {
		name string
		req  web.UserLoginRequest
	}{
		{
			name: "success_login",
			req: web.UserLoginRequest{
				Username: dataUser.Username,
				Password: dataUser.Password,
			},
		},
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
			name: "failed_login_validation_error",
			req: web.UserLoginRequest{
				Username: "",
				Password: "",
			},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Data body
			dataBody := fmt.Sprintf(`{"username": "%s", "password": "%s"}`, tc.req.Username, tc.req.Password)
			// New reader
			requestBody := strings.NewReader(dataBody)
			// Create new request
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/login", requestBody)
			// Added header content type
			request.Header.Add("Content-Type", "application/json")
			// Create new recorder
			recorder := httptest.NewRecorder()

			// Run http test
			RouteTest.ServeHTTP(recorder, request)

			// Get response
			response := recorder.Result()

			// Read all response
			body, _ := io.ReadAll(response.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			// Check pass test
			if tc.name == "success_login" {
				assert.Equal(t, 200, response.StatusCode)
				assert.Equal(t, 200, int(responseBody["code"].(float64)))
				assert.Equal(t, "success", responseBody["status"])
				assert.Equal(t, "login success", responseBody["message"])
				dataToken := responseBody["data"].(map[string]interface{})["token"]
				assert.NotEmpty(t, dataToken)
			} else if tc.name == "failed_login_wrong_username" || tc.name == "failed_login_wrong_password" {
				assert.Equal(t, 400, response.StatusCode)
				assert.Equal(t, 400, int(responseBody["code"].(float64)))
				assert.Equal(t, "error", responseBody["status"])
				assert.Equal(t, "login failed", responseBody["message"])
				assert.Equal(t, "username or password incorrect", responseBody["data"].(map[string]interface{})["errors"])
			} else {
				assert.Equal(t, 400, response.StatusCode)
				assert.Equal(t, 400, int(responseBody["code"].(float64)))
				assert.Equal(t, "error", responseBody["status"])
				assert.Equal(t, "login failed", responseBody["message"])
				assert.NotEqual(t, 0, len((responseBody["data"].(map[string]interface{})["errors"].([]interface{}))))
			}
		})
	}
}

func TestCreateUserHandler(t *testing.T) {
	// Todo Create success with role empty (user)
	// Todo Create success with role not empty (admin)
	// todo Username already exist
	// todo Validation error
	// todo unathorized
	// todo forbidden

	fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())
	username := strings.ToLower(helper.RandomPerson())

	// Test cases
	testCases := []struct {
		name string
		req  web.UserCreateRequest
	}{
		{
			name: "success_create_user",
			req: web.UserCreateRequest{
				Fullname: fullname,
				Username: username,
				Password: "password",
				Role:     "admin",
			},
		},
		{
			name: "success_create_user_without_role",
			req: web.UserCreateRequest{
				Fullname: fullname,
				Username: username,
				Password: "password",
			},
		},
		{
			name: "failed_create_user_username_is_exists",
			req: web.UserCreateRequest{
				Fullname: fullname,
				Username: "same",
				Password: "password",
			},
		},
		{
			name: "failed_create_user_validation_error",
			req:  web.UserCreateRequest{},
		},
		{
			name: "failed_create_user_unauthorized",
			req: web.UserCreateRequest{
				Fullname: fullname,
				Username: username,
				Password: "password",
			},
		},
		{
			name: "failed_create_user_forbidden",
			req: web.UserCreateRequest{
				Fullname: fullname,
				Username: username,
				Password: "password",
			},
		},
	}

	// Test
	for i := range testCases {
		tc := testCases[i]

		// Create new user and login for get token
		fullname := fmt.Sprintf("%s %s", helper.RandomPerson(), helper.RandomPerson())

		dataUser := web.UserCreateRequest{}
		if tc.name != "failed_create_user_forbidden" {
			// Create data user with field role
			dataUser.Fullname = fullname
			dataUser.Username = strings.ToLower(helper.RandomPerson())
			dataUser.Password = "password"
			dataUser.Role = "admin"
		} else {
			// Create data user withouth field role (default role: `user`)
			dataUser.Fullname = fullname
			dataUser.Username = strings.ToLower(helper.RandomPerson())
			dataUser.Password = "password"
		}

		repository := repository.NewUserRepository(ConnTest)
		service := service.NewServiceUser(repository)

		// Create user
		_, err := service.Create(dataUser)
		helper.ErrLogPanic(err)

		// Login
		dataLogin := web.UserLoginRequest{
			Username: dataUser.Username,
			Password: dataUser.Password,
		}
		token, err := service.Login(dataLogin)
		helper.ErrLogPanic(err)
		strToken := fmt.Sprintf("Bearer %s", token)
		// End Create new user and login for get token

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			var dataBody string
			// Data body
			if tc.name == "success_create_user" {
				// Body with field role
				dataBody = fmt.Sprintf(`{"fullname": "%s", "username": "%s", "password": "%s", "role": "%s"}`, tc.req.Fullname, tc.req.Username, tc.req.Password, tc.req.Role)
			} else if tc.name == "failed_create_user_username_is_exists" {
				// Value field username same with username in database
				dataBody = fmt.Sprintf(`{"fullname": "%s", "username": "%s", "password": "%s", "role": "%s"}`, tc.req.Fullname, dataUser.Username, tc.req.Password, tc.req.Role)
			} else {
				// Body without field role
				dataBody = fmt.Sprintf(`{"fullname": "%s", "username": "%s", "password": "%s"}`, tc.req.Fullname, tc.req.Username, tc.req.Password)
			}

			// New reader
			requestBody := strings.NewReader(dataBody)
			// Create new request
			request := httptest.NewRequest(http.MethodPost, "http://localhost:8080/api/v1/users", requestBody)
			// Added header `content-type`
			request.Header.Add("Content-Type", "application/json")

			// Added header `Authorization` if test case not "failed_create_user_unauthorized"
			if tc.name != "failed_create_user_unauthorized" {
				request.Header.Add("Authorization", strToken)
			}

			// Create new recorder
			recorder := httptest.NewRecorder()

			// Run http test
			RouteTest.ServeHTTP(recorder, request)

			// Get response
			response := recorder.Result()

			// Read all response
			body, _ := io.ReadAll(response.Body)
			var responseBody map[string]interface{}
			json.Unmarshal(body, &responseBody)

			if tc.name == "success_create_user" {
				assert.Equal(t, 201, response.StatusCode)
				assert.Equal(t, 201, int(responseBody["code"].(float64)))
				assert.Equal(t, "success", responseBody["status"])
				assert.Equal(t, "User has been created", responseBody["message"])

				// Get one by username
				user, err := service.GetOne(tc.req.Username)
				helper.ErrLogPanic(err)

				assert.NotEmpty(t, user.ID)
				assert.NotEmpty(t, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)

				assert.Equal(t, dataUser.Username, user.Username)
				assert.Equal(t, dataUser.Fullname, user.Fullname)
				assert.Equal(t, "admin", user.Role)

				// Compare password must be same
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dataUser.Password))
				assert.NoError(t, err)

			} else if tc.name == "success_create_user_without_role" {

				assert.Equal(t, 201, response.StatusCode)
				assert.Equal(t, 201, int(responseBody["code"].(float64)))
				assert.Equal(t, "success", responseBody["status"])
				assert.Equal(t, "User has been created", responseBody["message"])

				// Get one by username
				user, err := service.GetOne(tc.req.Username)
				helper.ErrLogPanic(err)

				assert.NotEmpty(t, user.ID)
				assert.NotEmpty(t, user.CreatedAt)
				assert.NotEmpty(t, user.UpdatedAt)

				assert.Equal(t, dataUser.Username, user.Username)
				assert.Equal(t, dataUser.Fullname, user.Fullname)
				assert.Equal(t, "user", user.Role)

				// Compare password must be same
				err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dataUser.Password))
				assert.NoError(t, err)

			} else if tc.name == "failed_create_user_username_is_exists" {

				assert.Equal(t, 400, response.StatusCode)
				assert.Equal(t, 400, int(responseBody["code"].(float64)))
				assert.Equal(t, "error", responseBody["status"])
				assert.Equal(t, "username already exist", responseBody["message"])

			} else if tc.name == "failed_create_user_validation_error" {

				assert.Equal(t, 400, response.StatusCode)
				assert.Equal(t, 400, int(responseBody["code"].(float64)))
				assert.Equal(t, "error", responseBody["status"])
				assert.Equal(t, "create user failed", responseBody["message"])
				assert.NotEqual(t, 0, len((responseBody["data"].(map[string]interface{})["errors"].([]interface{}))))

			} else if tc.name == "failed_create_user_unauthorized" {

				assert.Equal(t, 401, response.StatusCode)
				assert.Equal(t, 401, int(responseBody["code"].(float64)))
				assert.Equal(t, "error", responseBody["status"])
				assert.Equal(t, "Unauthorized", responseBody["message"])

			} else {

				assert.Equal(t, 403, response.StatusCode)
				assert.Equal(t, 403, int(responseBody["code"].(float64)))
				assert.Equal(t, "error", responseBody["status"])
				assert.Equal(t, "forbidden", responseBody["message"])

			}
		})
	}
}
