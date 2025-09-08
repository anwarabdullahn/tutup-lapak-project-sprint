package handlers

import (
	presenter "auth-service/api/presenters"
	"auth-service/pkg/entities"
	"auth-service/pkg/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateAuth = validator.New()

// Register is handler/controller which create new user
// @Summary      Create new user
// @Description  Create new user
// @Tags         Autentifikasi
// @Accept       json
// @Produce      json
// @Param        user  body      entities.AuthRequest   true  "User "
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/register [post]
func Register(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.AuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidasi := validateAuth.Struct(req); errValidasi != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidasi.Error()))
		}

		// cek email apakah sudah terdaftar
		existingUser, _ := service.FindByEmail(req.Email)
		if existingUser != nil {
			return c.Status(http.StatusConflict).
				JSON(presenter.ErrorResponse("email already registered"))
		}

		postData := &entities.User{
			Email:    req.Email,
			Password: string(req.Password),
		}

		user, err := service.Register(postData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(presenter.UserSuccessResponse(user))
	}
}

// Login is handler/controller login
// @Summary      Login
// @Description  Login
// @Tags         Autentifikasi
// @Accept       json
// @Produce      json
// @Param        user  body      entities.AuthRequest   true  "User "
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/login [post]
func Login(service user.Service, jwtManager *user.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.AuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidasi := validateAuth.Struct(req); errValidasi != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidasi.Error()))
		}

		postData := &entities.User{
			Email:    req.Email,
			Password: req.Password,
		}

		user, err := service.Login(postData)
		if err != nil {
			return c.Status(http.StatusUnauthorized).
				JSON(presenter.ErrorResponse("invalid email or password"))
		}

		userID := user.ID.String()
		// generate token
		token, err := jwtManager.Generate(userID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("failed to generate token"))
		}

		return c.JSON(fiber.Map{
			"email": user.Email,
			"token": token,
		})
	}
}
