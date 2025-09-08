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

// RegisterEmail handles email registration
// @Summary      Register with email
// @Description  Register new user with email
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      entities.EmailAuthRequest   true  "Email Registration Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/register/email [post]
func RegisterEmail(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.EmailAuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		// Check if email already registered
		existingUser, _ := service.FindByEmail(req.Email)
		if existingUser != nil {
			return c.Status(http.StatusConflict).
				JSON(presenter.ErrorResponse("email already registered"))
		}

		postData := &entities.User{
			Email:    req.Email,
			Password: req.Password,
		}

		user, err := service.Register(postData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(presenter.UserSuccessResponse(user))
	}
}

// RegisterPhone handles phone registration
// @Summary      Register with phone
// @Description  Register new user with phone
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      entities.PhoneAuthRequest   true  "Phone Registration Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/register/phone [post]
func RegisterPhone(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.PhoneAuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		// Check if phone already registered
		existingUser, _ := service.FindByPhone(req.Phone)
		if existingUser != nil {
			return c.Status(http.StatusConflict).
				JSON(presenter.ErrorResponse("phone already registered"))
		}

		postData := &entities.User{
			Phone:    req.Phone,
			Password: req.Password,
		}

		user, err := service.Register(postData)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(presenter.UserSuccessResponse(user))
	}
}

// LoginEmail handles email login
// @Summary      Login with email
// @Description  Login with email and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      entities.EmailAuthRequest   true  "Email Login Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/login/email [post]
func LoginEmail(service user.Service, jwtManager *user.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.EmailAuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
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
		token, err := jwtManager.Generate(userID)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse("failed to generate token"))
		}

		return c.JSON(fiber.Map{
			"email": user.Email,
			"phone": user.Phone,
			"token": token,
		})
	}
}

// LoginPhone handles phone login
// @Summary      Login with phone
// @Description  Login with phone and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        user  body      entities.PhoneAuthRequest   true  "Phone Login Data"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      401   {object}  map[string]interface{}
// @Router       /api/v1/login/phone [post]
func LoginPhone(service user.Service, jwtManager *user.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req entities.PhoneAuthRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body"))
		}

		if errValidation := validateAuth.Struct(req); errValidation != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errValidation.Error()))
		}

		postData := &entities.User{
			Phone:    req.Phone,
			Password: req.Password,
		}

		user, err := service.LoginByPhone(postData)
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

		return c.JSON(fiber.Map{
			"email": user.Email,
			"phone": user.Phone,
			"token": token,
		})
	}
}
