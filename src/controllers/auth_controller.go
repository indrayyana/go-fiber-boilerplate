package controllers

import (
	"app/src/config"
	"app/src/model"
	"app/src/response"
	"app/src/services"
	"app/src/validation"
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type AuthController struct {
	AuthService  services.AuthService
	UserService  services.UserService
	TokenService services.TokenService
}

func NewAuthController(
	authService services.AuthService, userService services.UserService, tokenService services.TokenService,
) *AuthController {
	return &AuthController{
		AuthService:  authService,
		UserService:  userService,
		TokenService: tokenService,
	}
}

// @Tags         Auth
// @Summary      Register as user
// @Accept       json
// @Produce      json
// @Param        request  body  validation.Register  true  "Request body"
// @Router       /auth/register [post]
// @Success      201  {object}  example.RegisterResponse
// @Failure      409  {object}  example.DuplicateEmail  "Email already taken"
func (a *AuthController) Register(c *fiber.Ctx) error {
	req := new(validation.Register)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Register(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusCreated).
		JSON(response.SuccessWithTokens[model.User]{
			Code:    fiber.StatusCreated,
			Status:  "success",
			Message: "Register successfully",
			Data:    user,
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Login
// @Accept       json
// @Produce      json
// @Param        request  body  validation.Login  true  "Request body"
// @Router       /auth/login [post]
// @Success      200  {object}  example.LoginResponse
// @Failure      401  {object}  example.FailedLogin  "Invalid email or password"
func (a *AuthController) Login(c *fiber.Ctx) error {
	req := new(validation.Login)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	user, err := a.AuthService.Login(c, req)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens[model.User]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			Data:    user,
			Tokens:  *tokens,
		})
}

// @Tags         Auth
// @Summary      Logout
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/logout [post]
// @Success      200  {object}  example.LogoutResponse
// @Failure      404  {object}  example.NotFound  "Not found"
func (a *AuthController) Logout(c *fiber.Ctx) error {
	req := new(validation.RefreshToken)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if err := a.AuthService.Logout(c, req); err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.Common{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Logout successfully",
		})
}

// @Tags         Auth
// @Summary      Refresh auth tokens
// @Accept       json
// @Produce      json
// @Param        request  body  example.RefreshToken  true  "Request body"
// @Router       /auth/refresh-tokens [post]
// @Success      200  {object}  example.RefreshTokenResponse
// @Failure      401  {object}  example.Unauthorized  "Unauthorized"
func (a *AuthController) RefreshToken(c *fiber.Ctx) error {
	req := new(validation.RefreshToken)

	if err := c.BodyParser(req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	tokens, err := a.AuthService.RefreshToken(c, req)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.RefreshToken{
			Code:   fiber.StatusOK,
			Status: "success",
			Tokens: *tokens,
		})
}

// @Tags         Auth
// @Summary      Login with google
// @Description  This route initiates the Google OAuth2 login flow. Please try this in your browser.
// @Router       /auth/google [get]
// @Success      200  {object}  example.GoogleLoginResponse
func (a *AuthController) GoogleLogin(c *fiber.Ctx) error {
	// Generate a random state
	state := uuid.New().String()

	c.Cookie(&fiber.Cookie{
		Name:   "oauth_state",
		Value:  state,
		MaxAge: 30,
	})

	url := config.AppConfig.GoogleLoginConfig.AuthCodeURL(state)

	return c.Status(fiber.StatusSeeOther).Redirect(url)
}

func (a *AuthController) GoogleCallback(c *fiber.Ctx) error {
	state := c.Query("state")
	storedState := c.Cookies("oauth_state")

	if state != storedState {
		return fiber.NewError(fiber.StatusUnauthorized, "States don't Match!")
	}

	code := c.Query("code")
	googlecon := config.GoogleConfig()

	token, err := googlecon.Exchange(context.Background(), code)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(
		c.Context(), http.MethodGet,
		"https://www.googleapis.com/oauth2/v2/userinfo?access_token="+token.AccessToken,
		nil,
	)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	userData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	googleUser := new(validation.GoogleLogin)
	if errJSON := json.Unmarshal(userData, googleUser); errJSON != nil {
		return errJSON
	}

	user, err := a.UserService.CreateGoogleUser(c, googleUser)
	if err != nil {
		return err
	}

	tokens, err := a.TokenService.GenerateAuthTokens(c, user)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).
		JSON(response.SuccessWithTokens[model.User]{
			Code:    fiber.StatusOK,
			Status:  "success",
			Message: "Login successfully",
			Data:    user,
			Tokens:  *tokens,
		})
}
