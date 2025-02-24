package handlers

import (
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/momokii/go-sso-web/internal/database"
	"github.com/momokii/go-sso-web/internal/middlewares"
	"github.com/momokii/go-sso-web/internal/models"
	modelsPkg "github.com/momokii/go-sso-web/pkg/models"
	"github.com/momokii/go-sso-web/pkg/repository/session"
	"github.com/momokii/go-sso-web/pkg/repository/user"
	"github.com/momokii/go-sso-web/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

const (
	SESSION_DURATION_DB = 6 * time.Hour
	JWT_DURATION        = 30 * time.Second
)

type AuthHandler struct {
	userRepo    user.UserRepo
	sessionRepo session.SessionRepo
}

func NewAuthHandler(userRepo user.UserRepo, sessionRepo session.SessionRepo) *AuthHandler {
	return &AuthHandler{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
	}
}

func (h *AuthHandler) LoginView(c *fiber.Ctx) error {
	return c.Render("login", fiber.Map{
		"Title": "Login - Klan SSO",
	})
}

func (h *AuthHandler) SignUpView(c *fiber.Ctx) error {
	return c.Render("signup", fiber.Map{
		"Title": "SignUp - Klan SSO",
	})
}

func (h *AuthHandler) SignUp(c *fiber.Ctx) error {
	auth := new(models.AuthLogin)
	if err := c.BodyParser(auth); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	if err := utils.ValidateStruct(auth); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Field() {
			case "Username":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Username must be alphanumeric and between 3-25 characters")
			case "Password":
				return utils.ResponseError(c, fiber.StatusBadRequest, "Password must be alphanumeric and between 6-50 characters with minimum 1 number and 1 uppercase letter")
			}
		}
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	// check if username already exist
	user_new, err := h.userRepo.FindByUsername(tx, auth.Username)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	if user_new.Id != 0 {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Username already exist")
	}

	// hashing password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(auth.Password), 16)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// add user to database
	user_new.Password = string(hashedPass)
	user_new.Username = auth.Username
	user_new.CreditToken = user.USER_MAX_DAILY_CREDIT_TOKEN

	if err := h.userRepo.Create(tx, user_new); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Signup success")
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	auth := new(models.AuthLogin)
	if err := c.BodyParser(auth); err != nil {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	userLog, err := h.userRepo.FindByUsername(tx, auth.Username)
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// check if user exist
	if userLog.Id == 0 {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	// password checking
	if err := bcrypt.CompareHashAndPassword([]byte(userLog.Password), []byte(auth.Password)); err != nil {
		return utils.ResponseError(c, fiber.StatusUnauthorized, "Invalid username or password")
	}

	// create uuid for session
	uuid, err := utils.GenerateUUIDV4()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// create session here with id for userid and tokenid for token id
	if err := middlewares.CreateSession(c, "id", userLog.Id); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}
	if err := middlewares.CreateSession(c, "session_id", uuid); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// save session to database
	time_now := time.Now()
	session_data := &modelsPkg.SessionCreate{
		UserId:    userLog.Id,
		SessionId: uuid,
		CreatedAt: time_now.Format(time.RFC3339),
		ExpiresAt: time_now.Add(SESSION_DURATION_DB).Format(time.RFC3339), // set expires at 6 hours (1 hour less than session fiber expires setup on server for 1 hour buffer time)
	}

	if err := h.sessionRepo.Create(tx, session_data); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.ResponseMessage(c, fiber.StatusOK, "Login success")
}

func (h *AuthHandler) RedirectRequest(c *fiber.Ctx) error {
	var app_url string

	app_req := c.Query("app")
	if app_req == "" && app_req != "gochat" && app_req != "llm" {
		return utils.ResponseError(c, fiber.StatusBadRequest, "Invalid request app type")
	}

	if app_req == "gochat" {
		app_url = os.Getenv("GOCHAT_URL")
	}

	if app_req == "llm" {
		app_url = os.Getenv("LLM_URL")
	}

	// create token jwt that combine session data: id and session_id
	// get user id and session id from session
	user_id, err := middlewares.CheckSession(c, "id")
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	session_id, err := middlewares.CheckSession(c, "session_id")
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// create jwt token
	sign := jwt.New(jwt.SigningMethodHS256)
	claims := sign.Claims.(jwt.MapClaims)
	claims["user_id"] = user_id
	claims["session_id"] = session_id
	// just use minimal exp time (a minute) because this token just for redirecting
	claims["exp"] = time.Now().Add(JWT_DURATION).Unix()

	token_jwt, err := sign.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.ResponseWitData(c, fiber.StatusOK, "Redirect success", fiber.Map{
		"token":        token_jwt,
		"redirect_url": app_url,
	})
}

func (h *AuthHandler) CheckAuthDashboard(c *fiber.Ctx) error {
	var user_session_data modelsPkg.UserSession
	is_logged_in := false // default value for loggin check

	// check session on local fiber data
	user_id, err := middlewares.CheckSession(c, "id")
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	session_id, err := middlewares.CheckSession(c, "session_id")
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// if session is valid and found on fiber session data, check on database
	if user_id != nil && session_id != nil {
		tx, err := database.DB.Begin()
		if err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
		}
		defer func() {
			database.CommitOrRollback(tx, c, err)
		}()

		// check session on db
		session_check, err := h.sessionRepo.FindSession(tx, session_id.(string), user_id.(int))
		if err != nil {
			return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
		}

		// if session not found on database, delete session on local fiber session and return as not logged in
		if session_check.Id == 0 && session_check.UserId == 0 && session_check.SessionId == "" {
			if err := middlewares.DeleteSession(c); err != nil {
				return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
			}

			return utils.ResponseWitData(c, fiber.StatusOK, "success dashboard data", fiber.Map{
				"is_logged_in": is_logged_in,
				"user":         user_session_data,
			})
		}

		// if session is valid on database, set is_logged_in to true and search user data
		user_data, err := h.userRepo.FindByID(tx, user_id.(int))
		if err != nil {
			if err := middlewares.DeleteSession(c); err != nil {
				return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
			}

			return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
		}

		user_session_data.Id = user_data.Id
		user_session_data.Username = user_data.Username
		user_session_data.CreditToken = user_data.CreditToken
		user_session_data.LastFirstLLMUsed = user_data.LastFirstLLMUsed

		is_logged_in = true // set to true if session is valid
	}

	return utils.ResponseWitData(c, fiber.StatusOK, "success dashboard data", fiber.Map{
		"is_logged_in": is_logged_in,
		"user":         user_session_data,
	})
}

func (h *AuthHandler) Logout(c *fiber.Ctx) error {
	// get user id and session id from session
	user_id, err := middlewares.CheckSession(c, "id")
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	session_id, err := middlewares.CheckSession(c, "session_id")
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// delete session from data first
	tx, err := database.DB.Begin()
	if err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	if err := h.sessionRepo.Delete(tx, session_id.(string), user_id.(int)); err != nil {
		return utils.ResponseError(c, fiber.StatusInternalServerError, err.Error())
	}

	// success delete session from data, now delete session from local fiber session
	middlewares.DeleteSession(c)

	return utils.ResponseMessage(c, fiber.StatusOK, "Logout success")
}
