package handler

import (
	"fmt"
	"net/http"
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
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	validate := validator.New()
	english := en.New()
	uni := ut.New(english, english)
	trans, _ := uni.GetTranslator("en")
	_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	var inputBind user.RegisterUserInputBind

	err := c.ShouldBindJSON(&inputBind)
	if err != nil {
		fmt.Println(inputBind.Name)
		fmt.Println(inputBind.Email)
		fmt.Println(inputBind.Password)

		validation := user.RegisterUserInputVal{
			Name:     inputBind.Name,
			Email:    inputBind.Email,
			Password: inputBind.Password,
		}

		err := validate.Struct(validation)

		errors := helper.FormatValidationError(err, trans)

		fmt.Println(errors)

		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Register account failed!!", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	newUser, err := h.userService.RegisterUser(inputBind)

	if err != nil {
		response := helper.APIResponse("Register account failed!", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "token")

	response := helper.APIResponse("Register account success!", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
