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
	"strconv"
	"strings"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/labstack/echo/v4"
	"github.com/nicksnyder/go-i18n/v2/i18n"
)

var Validator = &validation.CustomValidator{V: validator.New()}

// pictureDir = os.Getenv("PICTURE_DIR")

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

	return c.JSON(http.StatusCreated, fmt.Sprintf("User %s Registered Successfully", user.UserName))
}

func (s *Server) updateUser(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	if !ValidTokenForParam(c) {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUnAuthorized",
		})

		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	type inputStruct struct {
		Name     string `json:"userName,omitempty"`
		Email    string `json:"email,omitempty" validate:"email"`
		Password string `json:"password" validate:"required"`
	}

	input := &inputStruct{}

	err := c.Bind(input)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err})
	}

	if msgs, err := Validator.Validate(input, lang); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": msgs})
	}

	id, err := getIDFromParam(c)
	if err != nil {
		c.Logger().Error(err)
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGeniricBadRequest",
		})

		return c.JSON(http.StatusBadRequest, echo.Map{"error": message})
	}

	user := &models.User{
		ID:       id,
		UserName: input.Name,
		Email:    input.Email,
	}

	err = models.Models.User.UpdateUser(user)
	if err != nil {
		c.Logger().Error(err)
		return err
	}

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "UserUpdateSuccess",
			Other: "User info updated successfully",
		},
	})

	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (s *Server) login(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	type input struct {
		Email    string `json:"email" validate:"required"`
		Password string `json:"password" validate:"required"`
	}

	i := &input{}

	if err := c.Bind(i); err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if msgs, err := Validator.Validate(i, lang); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": msgs})
	}

	user := &models.User{
		Email:            i.Email,
		UnhashedPassword: i.Password,
	}

	err := models.Models.User.ValidateLogin(user)
	if err != nil {
		c.Logger().Error(err, user)
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrorFailedLogin",
				Other: "Username or Password incorrect",
			},
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	models.Models.User.SetUserRole(user)
	models.Models.User.SetID(user)
	c.Logger().Error(user)

	token, err := auth.CreateJwtToken(user)
	if err != nil {
		c.Logger().Error(err, user)
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": message})
	}

	return c.JSON(http.StatusOK, token)
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

	user := c.Param("name")

	err := models.Models.User.ResetPicture(user)
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

	name := c.Param("name")

	err := models.Models.User.DeleteUser(name)
	if err != nil {
		c.Logger().Error(err)
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		return c.JSON(http.StatusOK, echo.Map{"message": message})
	}

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "SuccessUserDelete",
			Other: "User deleted successfully",
		},
	})
	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (s *Server) setCategoryVisibilityOnUser(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	type inputStruct struct {
		UserName   string   `json:"userName" validate:"required"`
		Categories []string `json:"categories" validate:"required"`
		Activate   bool     `json:"activate" validate:"required"`
	}

	input := &inputStruct{}
	err := c.Bind(input)
	if err != nil {
		c.Logger().Error(err)
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrorGenericBadRequest",
				Other: "Your request doe not match the specified format, please fix and try again",
			},
		})
		return c.JSON(http.StatusBadRequest, message)
	}

	err = models.Models.Catagory.EditOnUser(input.UserName, input.Categories, input.Activate)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) updateProfilePicture(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

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
		PicturePath: fileName,
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

	userName := c.Param("name")
	user, err := models.Models.User.GetUserByName(userName)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUserNotExists",
		})
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": message})
	}

	if !ValidTokenForParam(c) {
		return c.JSON(http.StatusUnauthorized, echo.Map{"user": echo.Map{"userName": user.UserName, "email": user.Email}})
	}

	return c.JSON(http.StatusOK, echo.Map{"message": user})
}

func (s *Server) getAllCategories(c echo.Context) error {
	input := &models.Filters{}

	var err error
	input.Page, err = strconv.Atoi(c.QueryParam("page"))
	if err != nil {
		return err
	}
	input.PageSize, err = strconv.Atoi(c.QueryParam("pageSize"))
	if err != nil {
		return err
	}

	var cats []*models.Catagory
	var metadata models.Metadata
	if isAdmin(c) {
		cats, metadata, err = models.Models.Catagory.GetAll(*input)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	} else {
		userID := getIDFromToken(c)

		cats, metadata, err = models.Models.Catagory.GetAllActive(userID, *input)
		if err != nil {
			c.Logger().Error(err)
			return err
		}
	}
	return c.JSON(http.StatusOK, echo.Map{"categories": cats, "metadata": metadata})
}

func (s *Server) getProfilePicture(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	message := doesUserExist(c, localizer)
	if message != nil {
		return c.JSON(http.StatusBadRequest, message)
	}

	userName := c.Param("name")

	picturePath, err := models.Models.User.GetProfilePicture(userName)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorUserNotExists",
		})

		return c.JSON(http.StatusBadRequest, echo.Map{"error": message})
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/%s", picturePath))
}

func (s *Server) deleteCategory(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	name := c.Param("name")

	err := models.Models.Catagory.DeleteByName(name)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrCategoryNotExists",
				Other: "That category dose not exist",
			},
		})
		return c.JSON(http.StatusBadRequest, echo.Map{"error": message})
	}

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CategoryDeleteSuccess",
			Other: "Category removed successfully",
		},
	})
	return c.JSON(http.StatusOK, echo.Map{"message": message})
}

func (s *Server) postCategory(c echo.Context) error {
	lang := c.Request().Header.Get("Accept-Language")
	localizer := i18n.NewLocalizer(&translation.Bundle, lang)

	cat := &models.Catagory{}

	err := c.Bind(cat)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrorGenericBadRequest",
				Other: "Your request doe not match the specified format, please fix and try again",
			},
		})
		return c.JSON(http.StatusBadRequest, message)
	}

	if msgs, err := Validator.Validate(cat, lang); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"errors": msgs})
	}

	err = models.Models.Catagory.Insert(cat)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			MessageID: "ErrorGenericInternal",
		})
		return c.JSON(http.StatusBadRequest, echo.Map{"error": message})
	}

	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		DefaultMessage: &i18n.Message{
			ID:    "CategoryCreatedSuccess",
			Other: "Category created successfully",
		},
	})
	return c.JSON(http.StatusCreated, echo.Map{"message": message})
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

func doesUserExist(c echo.Context, localizer *i18n.Localizer) echo.Map {
	userName := c.Param("name")
	_, err := models.Models.User.GetUserByName(userName)
	if err != nil {
		message := localizer.MustLocalize(&i18n.LocalizeConfig{
			DefaultMessage: &i18n.Message{
				ID:    "ErrorUserNotExists",
				Other: "No user with that name has been found",
			},
		})
		return echo.Map{"error": message}
	}
	return nil
}

func getIDFromToken(c echo.Context) int {
	token := c.Get("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	userIDF := claims["id"].(float64)
	userID := int(userIDF)

	return userID
}

func getIDFromParam(c echo.Context) (int, error) {
	idString := c.Param("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		return 0, err
	}

	return id, nil
}
