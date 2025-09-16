package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"profile-service/api/presenter"
	"profile-service/pkg/entities"
	"profile-service/pkg/user"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validateUser = validator.New()

// GetMe is handler/controller which lists current user
// @Summary      Get current user
// @Description Get user profile from JWT token
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Success      200   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/user [get]
func GetMe(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDStr := c.Get("X-User-ID")

		data, err := service.GetByID(userIDStr)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(err.Error())
		}

		return c.JSON(presenter.ProfileSuccessResponse(data))

	}
}

// UpdateProfile is handler/controller which updates data of current user
// @Summary      Update a profile
// @Description  Update an existing profile (partial updates allowed)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      entities.UpdateUserRequest   true  "User update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/user [put]
func UpdateProfile(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.UpdateUserRequest
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid user id"))
		}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body: " + err.Error()))
		}

		user, err := service.UpdateProfile(userIDStr, requestBody)
		if err != nil {
			// Mapping error → HTTP response
			if errors.Is(err, entities.ErrInvalidFileID) {
				return c.Status(http.StatusBadRequest).
					JSON(presenter.ErrorResponse(err.Error()))
			}
			if errors.Is(err, entities.ErrFileNotFound) {
				return c.Status(http.StatusNotFound).
					JSON(presenter.ErrorResponse(err.Error()))
			}
			if errors.Is(err, entities.ErrInvalidUserID) {
				return c.Status(http.StatusBadRequest).
					JSON(presenter.ErrorResponse(err.Error()))
			}

			// default: internal error
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(presenter.ProfileSuccessResponse(user))
	}
}

// UpdateEmail is handler/controller which updates data of current user
// @Summary      Update a email
// @Description  Update an existing email (partial updates allowed)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      entities.EmailRequest   true  "User update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/user/link/email [post]
func UpdateEmail(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.EmailRequest

		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid user id"))
		}
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body: " + err.Error()))
		}

		if errVal := validateUser.Struct(requestBody); errVal != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errVal.Error()))
		}

		user, err := service.UpdateEmail(userIDStr, requestBody.Email)
		fmt.Println(user)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(presenter.ProfileSuccessResponse(user))

	}
}

// UpdatePhone is handler/controller which updates data of current user
// @Summary      Update a phone
// @Description  Update an existing phoe (partial updates allowed)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      entities.PhoneRequest   true  "User update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/user/link/phone [post]
func UpdatePhone(service user.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.PhoneRequest

		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid user id"))
		}
		// 2. Parse request body
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body: " + err.Error()))
		}

		// 3. Validasi request
		if errVal := validateUser.Struct(requestBody); errVal != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errVal.Error()))
		}

		//cek panjang phone
		if len(requestBody.Phone) < 8 || len(requestBody.Phone) > 16 {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid phone"))
		}

		user, err := service.UpdatePhone(userIDStr, requestBody.Phone)
		if err != nil {
			return c.Status(http.StatusInternalServerError).
				JSON(presenter.ErrorResponse(err.Error()))
		}

		return c.JSON(presenter.ProfileSuccessResponse(user))

	}
}
