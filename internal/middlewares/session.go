package middlewares

import (
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"github.com/momokii/go-sso-web/internal/database"
	"github.com/momokii/go-sso-web/pkg/models"
	sessionRepo "github.com/momokii/go-sso-web/pkg/repository/session"
	"github.com/momokii/go-sso-web/pkg/repository/user"
)

var (
	Store *session.Store
)

func InitSession() {
	Store = session.New(session.Config{
		Expiration:     7 * time.Hour,
		CookieSecure:   true,
		CookieHTTPOnly: true,
	})

	log.Println("Session store initialized")
}

func CreateSession(c *fiber.Ctx, key string, value interface{}) error {
	sess, err := Store.Get(c)
	if err != nil {
		return err
	}
	defer sess.Save()

	sess.Set(key, value)

	return nil
}

func DeleteSession(c *fiber.Ctx) error {
	sess, err := Store.Get(c)
	if err != nil {
		return err
	}
	defer sess.Save()

	sess.Destroy()

	return nil
}

func CheckSession(c *fiber.Ctx, key string) (interface{}, error) {
	sess, err := Store.Get(c)
	if err != nil {
		return nil, err
	}

	return sess.Get(key), nil
}

func IsNotAuth(c *fiber.Ctx) error {
	userid, err := CheckSession(c, "id")
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	session_id, err := CheckSession(c, "session_id")
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	if userid != nil && session_id != nil {
		return c.Redirect("/")
	}

	return c.Next()
}

func IsAuth(c *fiber.Ctx) error {
	// new checking, check the is_logged_in status
	is_logged_in, err := CheckSession(c, "is_logged_in")
	if err != nil || is_logged_in == nil {

		DeleteSession(c)
		return c.Redirect("/login")
	}

	base_url := c.Path()
	if base_url == "/multifa" && is_logged_in == true {
		// if is_logged_in is true, redirect to dashboard
		return c.Redirect("/")
	}

	userid, err := CheckSession(c, "id")
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	session_id, err := CheckSession(c, "session_id")
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	// if session data not found, redirect to login
	if userid == nil || session_id == nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	tx, err := database.DB.Begin()
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}
	defer func() {
		database.CommitOrRollback(tx, c, err)
	}()

	// check if session is valid
	userRepo := user.NewUserRepo()
	session_repo := sessionRepo.NewSessionRepo()

	// first check if session is valid on database
	sessData, err := session_repo.FindSession(tx, session_id.(string), userid.(int))
	// if session not found or error happen, redirect to login and delete the session local data
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	// if session is deleted/ not found
	if sessData.Id == 0 && sessData.UserId == 0 && sessData.SessionId == "" {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	userData, err := userRepo.FindByID(tx, userid.(int))
	if err != nil {
		DeleteSession(c)
		return c.Redirect("/login")
	}

	userSession := models.UserSession{
		Id:               userData.Id,
		Username:         userData.Username,
		CreditToken:      userData.CreditToken,
		LastFirstLLMUsed: userData.LastFirstLLMUsed,
		PhoneNumber:      userData.PhoneNumber,
		MultiFAEnabled:   userData.MultiFAEnabled,
	}

	// store information for next data
	c.Locals("user", userSession)

	return c.Next()
}
