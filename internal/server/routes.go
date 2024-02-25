package server

import (
	"Sadeem-RestAPI/internal/models"
	"Sadeem-RestAPI/internal/validation"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var Validator = &validation.CustomValidator{V: validator.New()}

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.helloWorldHandler)
	e.POST("/auth", s.registerUser)
	e.GET("/user/:id", s.getUserByID)
	e.GET("/user/:name", s.getUserByUserName)
	e.GET("/user/:name/test", s.getUserByUserNameT)

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
	user := new(models.User)

	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if msgs, err := Validator.Validate(user, lang); err != nil {
		return c.JSON(http.StatusBadRequest, msgs)
	}

	if err := models.Models.User.Insert(user); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("User %s Registered Successfully", user.UserName))
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
