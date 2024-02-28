package server

import (
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var signingKey = os.Getenv("JWT_SIGNING_KEY")

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// POST
	e.POST("/users", s.registerUser)

	// GET
	e.GET("/login", s.getUserToken)
	e.GET("/users/:name", jwtMiddleWare(adminMiddleWare((s.getUserByUserName))))
	e.GET("/categories/:name", jwtMiddleWare(s.getCategorieByName))
	e.GET("/categories", jwtMiddleWare(s.getCategorieByName))
	// PUT
	// e.PUT("/users/:name", jwtMiddleWare(s.updateUser))
	e.PUT("/users/:name/profile-picture", jwtMiddleWare(s.postProfilePicture))
	// e.PUT("/users/:name/categories/", jwtMiddleWare(adminMiddleWare(s.setCategoryVisibilityOnUser)))

	// DELETE
	e.DELETE("/users/:name", jwtMiddleWare(s.deleteUser))
	e.DELETE("/users/:profile-pucture", jwtMiddleWare(s.deleteProfilePicture))
	// e.DELETE("/categories/:name", jwtMiddleWare(adminMiddleWare(s.deleteCategoryByName)))

	return e
}

func adminMiddleWare(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !isAdmin(c) {
			return echo.ErrUnauthorized
		}
		return next(c)
	}
}

// Middleware for JWT Authintication
var jwtMiddleWare = echojwt.WithConfig(echojwt.Config{
	SigningMethod: "HS512",
	SigningKey:    []byte(signingKey),
})
