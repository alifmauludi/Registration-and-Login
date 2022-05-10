package handler

import (
	"fmt"
	"net/http"
	"reglog/auth"
	"reglog/helper"
	"reglog/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		err := validate.Struct(&input)
		errors := helper.FormatValidationError(err, trans)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Register account failed!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Register account failed!", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	token, err := h.authService.GenerateToken(newUser.ID)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Register account failed!", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)
	response := helper.APIResponse("Register account success!", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err, nil)
		errorMessage := gin.H{"errors": errors}
		response := helper.APIResponse("Login failed!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	loggedinUser, err := h.userService.Login(input)
	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}
		response := helper.APIResponse("Login failed!", http.StatusBadRequest, "error", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loggedinUser, token)
	response := helper.APIResponse("Successfuly loggedin", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("Failed to upload avatar image", http.StatusBadRequest, "error", data)

		c.JSON(http.StatusBadRequest, response)
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("Avatar successfuly uploaded", http.StatusOK, "success", data)

	c.JSON(http.StatusOK, response)
}
