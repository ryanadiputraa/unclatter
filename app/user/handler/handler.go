package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/ryanadiputraa/unclatter/app/user"
	"github.com/ryanadiputraa/unclatter/app/validation"
)

type handler struct {
	userService user.UserService
}

func NewUserHandler(r *echo.Group, userService user.UserService) {
	h := &handler{
		userService: userService,
	}

	r.GET("", h.getUserInfo())
}

func (h *handler) getUserInfo() echo.HandlerFunc {
	return func(c echo.Context) error {
		// TODO: get user id from access token
		user, err := h.userService.GetUserInfo(c.Request().Context(), "")
		if err != nil {
			if vErr, ok := err.(*validation.Error); ok {
				return c.JSON(http.StatusBadRequest, map[string]any{
					"message": vErr.Message,
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]any{
				"message": "internal server error",
			})
		}

		return c.JSON(http.StatusOK, map[string]any{
			"data": user,
		})
	}
}
