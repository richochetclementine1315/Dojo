package routes

import (
	"dojo/internal/config"
	"dojo/internal/handler"
	"dojo/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

// SetUpRoutes sets up all the routes for the application
func SetupRoutes(app *fiber.App, handlers *Handlers, cfg *config.Config) {
	// Prefix for APIs
	api := app.Group("/api")
	// Health  Check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Server is Running",
		})
	})
	// Auth Routes(PUBLIC WALEE!!)
	authRoutes := api.Group("/auth")
	{
		authRoutes.Post("/register", handlers.Auth.Register)
		authRoutes.Post("/login", handlers.Auth.Login)
		authRoutes.Get("/google", handlers.Auth.GoogleLogin)
		authRoutes.Get("/google/callback", handlers.Auth.GoogleCallback)
		authRoutes.Get("/github", handlers.Auth.GitHubLogin)
		authRoutes.Get("/github/callback", handlers.Auth.GitHubCallback)
		authRoutes.Post("/refresh", handlers.Auth.RefreshToken)
		authRoutes.Post("/logout", handlers.Auth.Logout)
	}
	// Protected routes(require authentication)
	protected := api.Group("", middleware.AuthMiddleware(cfg))
	{
		userRoutes := protected.Group("/users")
		{
			userRoutes.Get("/profile", handlers.User.GetProfile)
			userRoutes.Put("/profile", handlers.User.UpdateProfile)
			userRoutes.Patch("/account", handlers.User.UpdateUser)
			userRoutes.Post("/change-password", handlers.User.ChangePassword)
			userRoutes.Post("/sync-stats", handlers.User.SyncPlatformStats)

		}
		// Problem Routes
		problemRoutes := protected.Group("/problems")
		{
			problemRoutes.Get("", handlers.Problem.ListProblems)
			problemRoutes.Post("", handlers.Problem.CreateProblem)
			problemRoutes.Get("/:id", handlers.Problem.GetProblem)
			problemRoutes.Put("/:id", handlers.Problem.UpdateProblem)
			problemRoutes.Delete("/:id", handlers.Problem.DeleteProblem)
		}

	}
}

type Handlers struct {
	Auth    *handler.AuthHandler
	User    *handler.UserHandler
	Problem *handler.ProblemHandler
}
