package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/html/v2"
	"github.com/momokii/go-sso-web/internal/database"
	"github.com/momokii/go-sso-web/internal/handlers"
	"github.com/momokii/go-sso-web/internal/middlewares"
	"github.com/momokii/go-sso-web/internal/repository/session"
	"github.com/momokii/go-sso-web/internal/repository/user"
	"github.com/momokii/go-sso-web/pkg/worker"

	_ "github.com/joho/godotenv/autoload"
)

const (
	WORKER_SESSION_CHECKER_DURATION = 25 * time.Second
)

func main() {
	// db and session storage init
	database.InitDB()
	middlewares.InitSession()

	// repo init
	userRepo := user.NewUserRepo()
	sessionRepo := session.NewSessionRepo()

	// handler init
	authHandler := handlers.NewAuthHandler(*userRepo, *sessionRepo)
	userHandler := handlers.NewUserHandler(*userRepo)

	// worker init
	sessionChecker := worker.NewSessionChecker(*sessionRepo)

	// start worker
	sessionChecker.StartChecker(WORKER_SESSION_CHECKER_DURATION)

	// app server setup and init
	engine := html.New("./web", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return c.Status(code).Render("error", fiber.Map{
				"Code":    code,
				"Message": err.Error(),
			})
		},
	})

	api := app.Group("/api") // for api endpoint

	app.Use(cors.New())
	app.Use(logger.New())
	app.Static("/web", "./web")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("dashboard", fiber.Map{
			"Title": "Dashboard - Klan SSO",
		})
	})
	// dashboard api get data
	api.Get("/dashboard", authHandler.CheckAuthDashboard)

	// auth
	api.Get("/redirect", middlewares.IsAuth, authHandler.RedirectRequest)

	app.Get("/login", middlewares.IsNotAuth, authHandler.LoginView)
	api.Post("/login", middlewares.IsNotAuth, authHandler.Login)

	app.Get("/signup", middlewares.IsNotAuth, authHandler.SignUpView)
	api.Post("/signup", middlewares.IsNotAuth, authHandler.SignUp)

	api.Post("/logout", middlewares.IsAuth, authHandler.Logout)

	// user
	api.Patch("/users", middlewares.IsAuth, userHandler.ChangeUsername)
	api.Patch("/users/password", middlewares.IsAuth, userHandler.ChangePassword)

	devMode := os.Getenv("APP_ENV")
	if devMode != "development" && devMode != "production" {
		log.Println("APP_ENV not set")
	} else {
		log.Println("Running on: " + devMode)
		if err := app.Listen(":3001"); err != nil {
			log.Println("Error running Server: ", err)
		}
	}
}
