package validation

import (
	"Sadeem-RestAPI/internal/translation"
	"errors"

	"github.com/go-playground/validator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// Email regex taken from the let's go book
// Makes sure email is valid
// And no i don't know how it works
type ApiError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

type CustomValidator struct {
	V *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}, errLang string) ([]ApiError, error) {
	if err := cv.V.Struct(i); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			out := make([]ApiError, len(ve))
			for i, fe := range ve {
				out[i] = ApiError{msgForField(fe.Field()), msgForTag(fe.Tag(), errLang)}
			}
			return out, err
		}
		return nil, err
	}

	return nil, nil
}

func msgForField(field string) string {
	switch field {

	case "UnhashedPassword":
		return "password"
	case "Email":
		return "email"
	case "Name":
		return "name"
	case "UserName":
		return "userName"
	default:
		return field
	}
}

func msgForTag(tag, lang string) string {
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)
	var msg string
	switch tag {
	case "required":
		msg = localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Required",
				Other: "This field is required",
			},
		})
	case "email":
		msg = localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "Email",
				Other: "Invalid Email address",
			},
		})
	default:
		msg = tag
	}

	return msg
}
