package server

import (
	"Sadeem-RestAPI/internal/models"
	"Sadeem-RestAPI/internal/translation"
	"Sadeem-RestAPI/internal/validation"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var Validator = &validation.CustomValidator{V: validator.New()}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.helloWorldHandler)
	e.POST("/users", s.registerUser)
	e.PUT("/users/:id/profile-picture", s.postProfilePicture)
	e.GET("/users/:id", s.getUserByID)
	e.GET("/users/:name", s.getUserByUserName)
	e.GET("/users/:name/test", s.getUserByUserNameT)

	return e
}

func (s *Server) helloWorldHandler(c echo.Context) error {
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) registerUser(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)
	user := new(models.User)

	errorMessage := make(map[string]any)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if msgs, err := Validator.Validate(user, lang); err != nil {
		errorMessage["error"] = msgs
		return c.JSON(http.StatusBadRequest, errorMessage)
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

func (s *Server) postProfilePicture(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	id := c.Param("id")

	pictureDir := os.Getenv("PICTURE_DIR")

	Message := make(map[string]string)
	defer c.Request().Body.Close()
	byte, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return c.JSON(http.StatusBadRequest, "Could not read body")
	}

	mimeType := http.DetectContentType(byte)

	if mimeType != "image/jpeg" && mimeType != "image/png" {
		Message["error"] = localizer.MustLocalize(
			&i18n.LocalizeConfig{
				DefaultMessage: &i18n.Message{
					ID:    "NotPngOrJpeg",
					Other: "Profile Picture must be a PNG or a JPG",
				},
			})

		return c.JSON(http.StatusBadRequest, Message)
	}

	Message["message"] = localizer.MustLocalize(
		&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "SuccessUpdateProfilePicture",
				Other: "Profile Picture Updated Successfully",
			},
		},
	)

	_, fileType, _ := strings.Cut(mimeType, "/")
	fileName := fmt.Sprintf("user_%s-profile_picture.%s", id, fileType)
	filePath := path.Join(pictureDir, fileName)
	Message["path"] = filePath

	_, err = os.Create(filePath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err = os.WriteFile(filePath, byte, os.ModeAppend)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user_id, err := strconv.Atoi(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	user := &models.User{
		ID:          user_id,
		PicturePath: filePath,
	}

	err = models.Models.User.UpdatePicture(user)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, Message)
}

func (s *Server) getUserByID(c echo.Context) error {
	resp := map[string]string{
		"id":   "1",
		"user": "CollCaz",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getUserByUserName(c echo.Context) error {
	resp := map[string]string{
		"id":   "1",
		"user": "CollCaz",
	}

	return c.JSON(http.StatusOK, resp)
}

func (s *Server) getUserByUserNameT(c echo.Context) error {
	resp := map[string]string{
		"id":   "1",
		"user": "Test",
	}

	return c.JSON(http.StatusOK, resp)
}
