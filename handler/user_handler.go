package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/letenk/use_deal_user/models/domain"
	"github.com/letenk/use_deal_user/models/web"
	"github.com/letenk/use_deal_user/service"
)

type userHandler struct {
	service service.UserService
}

func NewHandlerUser(service service.UserService) *userHandler {
	return &userHandler{service}
}

func (h *userHandler) Login(c *gin.Context) {
	var req web.UserLoginRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		errors := web.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"login failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Login
	token, err := h.service.Login(req)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"login failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	fieldToken := gin.H{"token": token}
	// Create format response
	response := web.JSONResponseWithData(
		http.StatusOK,
		"success",
		"login success",
		fieldToken,
	)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Create(c *gin.Context) {
	// Check Authorization
	// Get current user login
	// currentUser := c.MustGet("currentUser").(domain.User)
	// if currentUser.Role != "admin" {
	// 	response := web.JSONResponseWithoutData(
	// 		http.StatusForbidden,
	// 		"error",
	// 		"forbidden",
	// 	)
	// 	c.JSON(http.StatusForbidden, response)
	// 	return
	// }

	// Get payload body
	var req web.UserCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		errors := web.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create user failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get by username
	user, _ := h.service.GetOneByUsername(req.Username)

	// If user.id not empty (user is available in db)
	if user.ID != "" {
		errorMessage := gin.H{"errors": "username already exist"}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create user failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create new user
	_, err = h.service.Create(req)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create user failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Create format response
	response := web.JSONResponseWithoutData(
		http.StatusCreated,
		"success",
		"User has been created",
	)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) GetAll(c *gin.Context) {
	// Get all
	users, err := h.service.GetAll()
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"create user failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"Success",
		"List of users",
		web.FormatUsersResponse(users),
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *userHandler) GetOne(c *gin.Context) {
	// Check Authorization
	// Get current user login
	currentUser := c.MustGet("currentUser").(domain.User)
	if currentUser.Role != "admin" {
		response := web.JSONResponseWithoutData(
			http.StatusForbidden,
			"error",
			"forbidden",
		)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// Get id from uri
	var userID web.UserGetIDUri
	err := c.ShouldBindUri(&userID)
	if err != nil {
		resp := gin.H{"errors": err}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"Data of users",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Get one
	user, err := h.service.GetOne(userID.ID)
	if err != nil {
		resp := gin.H{"errors": err.Error()}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"Data of users",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"Success",
		"Data of users",
		web.FormatUserResponse(user),
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *userHandler) Update(c *gin.Context) {
	// Check Authorization
	// Get current user login
	currentUser := c.MustGet("currentUser").(domain.User)
	if currentUser.Role != "admin" {
		response := web.JSONResponseWithoutData(
			http.StatusForbidden,
			"error",
			"forbidden",
		)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// Get id from uri
	var userID web.UserGetIDUri
	err := c.ShouldBindUri(&userID)
	if err != nil {
		resp := gin.H{"errors": err}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Get payload body
	var req web.UserUpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		errorMessage := gin.H{"errors": err}
		response := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update user failed",
			errorMessage,
		)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	// Get one
	user, err := h.service.GetOne(userID.ID)
	if err != nil {
		resp := gin.H{"errors": err.Error()}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Update
	_, err = h.service.Update(user.ID, req)
	if err != nil {
		resp := gin.H{"errors": err}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Get one user updated
	userUpdated, err := h.service.GetOne(userID.ID)
	if err != nil {
		resp := gin.H{"errors": err.Error()}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"update user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponseWithData(
		http.StatusOK,
		"Success",
		"update user success",
		web.FormatUserResponse(userUpdated),
	)
	c.JSON(http.StatusOK, jsonResponse)
}

func (h *userHandler) Delete(c *gin.Context) {
	// Check Authorization
	// Get current user login
	currentUser := c.MustGet("currentUser").(domain.User)
	if currentUser.Role != "admin" {
		response := web.JSONResponseWithoutData(
			http.StatusForbidden,
			"error",
			"forbidden",
		)
		c.JSON(http.StatusForbidden, response)
		return
	}

	// Get id from uri
	var userID web.UserGetIDUri
	err := c.ShouldBindUri(&userID)
	if err != nil {
		resp := gin.H{"errors": err}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"delete user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Get one
	_, err = h.service.GetOne(userID.ID)
	if err != nil {
		resp := gin.H{"errors": err.Error()}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"delete user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	// Delete
	_, err = h.service.Delete(userID.ID)
	if err != nil {
		resp := gin.H{"errors": err}
		jsonResponse := web.JSONResponseWithData(
			http.StatusBadRequest,
			"error",
			"delete user failed",
			resp,
		)
		c.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	jsonResponse := web.JSONResponseWithoutData(
		http.StatusOK,
		"Success",
		"User has been deleted",
	)
	c.JSON(http.StatusOK, jsonResponse)
}
