package helper

import (
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

type Response struct {
	Meta Meta        `json:"meta"`
	Data interface{} `json:"data"`
}

type Meta struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
	Status  string `json:"status"`
}

func APIResponse(message string, code int, status string, data interface{}) Response {
	meta := Meta{
		Message: message,
		Code:    code,
		Status:  status,
	}

	jsonResponse := Response{
		Meta: meta,
		Data: data,
	}

	return jsonResponse
}

func FormatValidationError(err error, trans ut.Translator) []string {
	var errors []string
	//validate := validator.New()
	//english := en.New()
	//uni := ut.New(english, english)
	//trans, _ = uni.GetTranslator("en")
	//_ = enTranslations.RegisterDefaultTranslations(validate, trans)
	//if err == nil {
	//	return nil
	//}

	validatorErrs := err.(validator.ValidationErrors)
	for _, e := range validatorErrs {
		//translatedErr := fmt.Errorf(e.Translate(trans))
		errors = append(errors, e.Translate(trans))
	}
	return errors
}
