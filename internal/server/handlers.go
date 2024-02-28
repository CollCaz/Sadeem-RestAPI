package server

import (
	"Sadeem-RestAPI/internal/auth"
	"Sadeem-RestAPI/internal/models"
	"Sadeem-RestAPI/internal/translation"
	"Sadeem-RestAPI/internal/validation"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var Validator = &validation.CustomValidator{V: validator.New()}

func (s *Server) registerUser(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if msgs, err := Validator.Validate(user, lang); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": msgs})
	}

	if err := models.Models.User.Insert(user); err != nil {
		msg := err.Error()
		if pgerr, ok := err.(*pgconn.PgError); ok {
			if pgerr.Code == "23505" {
				msg = localizer.MustLocalize(&i18n.LocalizeConfig{
					DefaultMessage: &i18n.Message{
						ID:    "ErrorDuplicateEmailOrUsername",
						One:   "Email or Username Already Exists",
						Other: "Email or Username Already Exists",
					},
				})
			} else {
				msg = err.Error()
			}
		}
		return c.JSON(http.StatusBadRequest, msg)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("User %s Registered Successfully", user.UserName))
}

func (s *Server) getUserToken(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	err := models.Models.User.ValidateLogin(user)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrorFailedLogin",
				Other: "Username or Password incorrect",
			},
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	models.Models.User.SetUserRole(user)

	token, err := auth.CreateJwtToken(user)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": message})
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token, "user": user})
}

func (s *Server) deleteProfilePicture(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	if !ValidTokenForParam(c) {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUnAuthorized",
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	userName := c.Param("name")
	err := models.Models.User.ResetPicture(userName)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		return c.JSON(http.StatusOK, echo.Map{"message": message})
	}

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SuccessUserUpdate",
			Other: "User info update successfully",
		},
	})

	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (s *Server) deleteUser(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	if !ValidTokenForParam(c) {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUnAuthorized",
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	userName := c.Param("name")
	err := models.Models.User.DeleteUserByName(userName)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		return c.JSON(http.StatusOK, echo.Map{"message": message, "user": userName})
	}

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SuccessUserDelete",
			Other: "User deleted successfully",
		},
	})
	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (s *Server) postProfilePicture(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	fmt.Println(c.ParamNames())
	fmt.Println(c.Param("name"))
	if !ValidTokenForParam(c) {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUnAuthorized",
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	userName := c.Param("name")

	pictureDir := os.Getenv("PICTURE_DIR")

	var message string
	defer c.Request().Body.Close()
	byte, err := io.ReadAll(c.Request().Body)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "CouldNotReadImage",
				Other: "Could not proccess your image, plasea try again with a new image",
			},
		})
		return c.JSON(http.StatusBadRequest, echo.Map{"error": message})
	}

	mimeType := http.DetectContentType(byte)

	if mimeType != "image/jpeg" && mimeType != "image/png" {
		message = localizer.MustLocalize(
			&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "NotPngOrJpeg",
					Other: "Profile Picture must be a PNG or a JPG",
				},
			})

		return c.JSON(http.StatusBadRequest, echo.Map{"message": message})
	}

	_, fileType, _ := strings.Cut(mimeType, "/")
	fileName := fmt.Sprintf("user_%s-profile_picture.%s", userName, fileType)
	filePath := path.Join(pictureDir, fileName)

	message = localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "ErrorGenericInternal",
			Other: "We encountred an error proccessing you're request, please try again later",
		},
	})

	_, err = os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": message})
	}

	err = os.WriteFile(filePath, byte, os.ModeAppend)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": message})
	}

	user := &models.User{
		UserName:    userName,
		PicturePath: filePath,
	}

	err = models.Models.User.UpdatePicture(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	message = localizer.MustLocalize(
		&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "SuccessUpdateProfilePicture",
				Other: "Profile Picture Updated Successfully",
			},
		},
	)

	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (s *Server) getUserByUserName(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	if !isAdmin(c) {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUnAuthorized",
		})
		return c.JSON(http.StatusUnauthorized, message)
	}

	userName := c.Param("name")
	user, err := models.Models.User.GetUserByName(userName)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrorUserNotExists",
				Other: "No user with that name has been found",
			},
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": user})
}

func (s *Server) getCategorieByName(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	name := c.Param("name")

	cat, err := models.Models.Catagory.GetByName(name)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": echo.Map{"category": cat}})
}

func (s *Server) getAllCategories(c echo.Context) error {
	input := &models.Filters{}
	err := c.Bind(input)
	if err != nil {
		return err
	}

	input.SortSafeList = []string{"name", "-name"}

	var cats []*models.Catagory
	var metadata models.Metadata
	if isAdmin(c) {
		cats, metadata, err = models.Models.Catagory.GetAll(*input)
		if err != nil {
			return err
		}
	} else {
		token := c.Get("user").(*jwt.Token)
		claims := token.Claims.(jwt.MapClaims)
		userID := claims["id"].(int)
		cats, metadata, err = models.Models.Catagory.GetAllActive(userID, *input)
		if err != nil {
			return err
		}
	}
	return c.JSON(http.StatusOK, echo.Map{"message": echo.Map{"categories": cats, "metadata": metadata}})
}

// Returns ture if the JWT user is the same
// as the user in the url params OR if the jwt
// user is an admin
func ValidTokenForParam(c echo.Context) bool {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	tokenUser := claims["name"].(string)
	isAdmin := claims["admin"].(bool)

	paramUser := c.Param("name")

	return paramUser == tokenUser || isAdmin
}

func isAdmin(c echo.Context) bool {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	return claims["admin"].(bool)
}
