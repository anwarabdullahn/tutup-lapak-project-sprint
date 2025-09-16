package handlers

import (
	"net/http"
	"profile-service/api/presenter"
	"profile-service/pkg/auth"
	"profile-service/pkg/entities"

	"github.com/gofiber/fiber/v2"
)

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
func Register(service auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.AuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
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

		return c.JSON(presenter.SuccessRegisterResponse(user))
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
func Login(service auth.Service, jwtManager *auth.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.AuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
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

		// generate token
		token, err := jwtManager.Generate(user.ID.String())
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
