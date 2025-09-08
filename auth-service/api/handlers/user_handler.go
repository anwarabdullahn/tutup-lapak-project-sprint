package handlers

import (
	"auth-service/api/presenter"
	"auth-service/config"
	"auth-service/pkg/dtos"
	"auth-service/pkg/entities"
	"auth-service/pkg/user"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateAuth = validator.New()

// RegisterEmail handles email registration
// @Summary      Register with email
// @Description  Register new user with email
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      dtos.EmailRegisterRequest   true  "Email Registration Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/register/email [post]
func RegisterEmail(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dtos.EmailRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		// Check if email already registered
		existingUser, _ := service.FindByEmail(c.Context(), req.Email)
		if existingUser != nil {
			return c.Status(http.StatusConflict).
				JSON(presenter.ErrorResponse("email already registered"))
		}

		createReq := entities.CreateUserRequest{
			Email:    req.Email,
			Password: req.Password,
		}
		user, err := service.Register(c.Context(), createReq)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": dtos.UserResponse{
				ID:    user.ID.String(),
				Email: user.Email,
			},
		})
	}
}

// RegisterPhone handles phone registration
// @Summary      Register with phone
// @Description  Register new user with phone
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      dtos.PhoneRegisterRequest   true  "Phone Registration Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/register/phone [post]
func RegisterPhone(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dtos.PhoneRegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		// Check if phone already registered
		existingUser, _ := service.FindByPhone(c.Context(), req.Phone)
		if existingUser != nil {
			return c.Status(http.StatusConflict).
				JSON(presenter.ErrorResponse("phone already registered"))
		}

		createReq := entities.CreateUserRequest{
			Phone:    req.Phone,
			Password: req.Password,
		}
		user, err := service.Register(c.Context(), createReq)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(fiber.Map{
			"success": true,
			"data": dtos.UserResponse{
				ID:    user.ID.String(),
				Phone: user.Phone,
			},
		})
	}
}

// LoginEmail handles email login
// @Summary      Login with email
// @Description  Login with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      dtos.EmailLoginRequest   true  "Email Login Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/login/email [post]
func LoginEmail(service user.Service, jwtManager *config.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dtos.EmailLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		loginReq := entities.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		}
		user, err := service.Login(c.Context(), loginReq)
		if err != nil {
			return c.Status(http.StatusUnauthorized).
				JSON(presenter.ErrorResponse("invalid email or password"))
		}

		userID := user.ID.String()
		token, err := jwtManager.Generate(userID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("failed to generate token"))
		}

		return c.JSON(dtos.LoginResponse{
			User: dtos.UserResponse{
				ID:    user.ID.String(),
				Email: user.Email,
			},
			Token: token,
		})
	}
}

// LoginPhone handles phone login
// @Summary      Login with phone
// @Description  Login with phone and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      dtos.PhoneLoginRequest   true  "Phone Login Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/login/phone [post]
func LoginPhone(service user.Service, jwtManager *config.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dtos.PhoneLoginRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		loginReq := entities.LoginRequest{
			Phone:    req.Phone,
			Password: req.Password,
		}
		user, err := service.Login(c.Context(), loginReq)
		if err != nil {
			return c.Status(http.StatusUnauthorized).
				JSON(presenter.ErrorResponse("invalid phone or password"))
		}

		userID := user.ID.String()
		token, err := jwtManager.Generate(userID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("failed to generate token"))
		}

		return c.JSON(dtos.LoginResponse{
			User: dtos.UserResponse{
				ID:    user.ID.String(),
				Phone: user.Phone,
			},
			Token: token,
		})
	}
}
