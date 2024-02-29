package server

import (
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var signingKey = os.Getenv("JWT_SIGNING_KEY")

func (s *Server) RegisterRoutes() *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Static(os.Getenv("PICTURE_DIR")))

	// POST
	e.POST("api/users", s.registerUser)
	e.POST("api/categories", jwtMiddleWare(adminMiddleWare(s.postCategory)))
	e.POST("api/login", s.login)
	e.POST("api/user-categories", jwtMiddleWare(adminMiddleWare(s.setCategoryVisibilityOnUser)))

	// GET
	e.GET("api/users/:name", jwtMiddleWare((s.getUserByUserName)))
	e.GET("api/users/:name/profile-picture", jwtMiddleWare(s.getProfilePicture))
	e.GET("api/categories", jwtMiddleWare(s.getAllCategories))

	// PUT
	e.PUT("api/users/:id", jwtMiddleWare(s.updateUser))
	e.PUT("api/users/:name/profile-picture", jwtMiddleWare(s.updateProfilePicture))

	// DELETE
	e.DELETE("api/users/:name", jwtMiddleWare(s.deleteUser))
	e.DELETE("api/users/:name/profile-picture", jwtMiddleWare(s.deleteProfilePicture))
	e.DELETE("api/categories/:name", jwtMiddleWare(adminMiddleWare(s.deleteCategory)))

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
