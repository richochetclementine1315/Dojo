package routes

import (
	"dojo/internal/config"
	"dojo/internal/handler"
	"dojo/internal/middleware"
	"dojo/internal/websocket"

	fiberws "github.com/gofiber/contrib/websocket"
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

	// Public contest routes (no authentication required)
	contestRoutes := api.Group("/contests")
	{
		contestRoutes.Get("", handlers.Contest.ListContests)
		contestRoutes.Get("/:id", handlers.Contest.GetContest)
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
		// Protected Contest Routes (sync and reminders require auth)
		protectedContestRoutes := protected.Group("/contests")
		{
			protectedContestRoutes.Post("/sync", handlers.Contest.SyncContests)
			protectedContestRoutes.Post("/reminders", handlers.Contest.CreateReminder)
			protectedContestRoutes.Delete("/reminders/:id", handlers.Contest.DeleteReminder)
		}
		sheetRoutes := protected.Group("/sheets")
		{
			sheetRoutes.Get("/public", handlers.Sheet.GetPublicSheets)
			sheetRoutes.Get("", handlers.Sheet.GetUserSheets)
			sheetRoutes.Post("", handlers.Sheet.CreateSheet)
			sheetRoutes.Get("/:id", handlers.Sheet.GetSheet)
			sheetRoutes.Put("/:id", handlers.Sheet.UpdateSheet)
			sheetRoutes.Delete("/:id", handlers.Sheet.DeleteSheet)
			sheetRoutes.Post("/:id/problems", handlers.Sheet.AddProblemToSheet)
			sheetRoutes.Delete("/:id/problems/:problemId", handlers.Sheet.RemoveProblemFromSheet)
			sheetRoutes.Patch("/:id/problems/:problemId", handlers.Sheet.UpdateSheetProblem)
		}
		// Social Routes (protected)
		socialRoutes := protected.Group("/social")
		{
			// Friend requests
			socialRoutes.Post("/friends/requests", handlers.Social.SendFriendRequest)
			socialRoutes.Get("/friends/requests/received", handlers.Social.GetReceivedRequests)
			socialRoutes.Get("/friends/requests/sent", handlers.Social.GetSentRequests)
			socialRoutes.Patch("/friends/requests/:id", handlers.Social.RespondToFriendRequest)
			socialRoutes.Delete("/friends/requests/:id", handlers.Social.CancelFriendRequest)

			// Friends
			socialRoutes.Get("/friends", handlers.Social.GetFriends)
			socialRoutes.Delete("/friends/:id", handlers.Social.RemoveFriend)

			// Blocks
			socialRoutes.Post("/blocks", handlers.Social.BlockUser)
			socialRoutes.Get("/blocks", handlers.Social.GetBlockedUsers)
			socialRoutes.Delete("/blocks/:id", handlers.Social.UnblockUser)

			// Search
			socialRoutes.Get("/users/search", handlers.Social.SearchUsers)
		}

		// Room Routes
		roomRoutes := protected.Group("/rooms")
		{
			roomRoutes.Post("", handlers.Room.CreateRoom)
			roomRoutes.Get("", handlers.Room.GetUserRooms)
			roomRoutes.Post("/join", handlers.Room.JoinRoom)
			roomRoutes.Get("/:id", handlers.Room.GetRoom)
			roomRoutes.Post("/:id/leave", handlers.Room.LeaveRoom)
			roomRoutes.Delete("/:id", handlers.Room.DeleteRoom)
			roomRoutes.Get("/:id/code", handlers.Room.GetCodeSession)
			roomRoutes.Put("/:id/code", handlers.Room.UpdateCodeSession)

			// WebSocket Connection
			roomRoutes.Get("/:id/ws", handlers.RoomWS.UpgradeConnection, fiberws.New(handlers.RoomWS.HandleConnection))
		}
	}
}

type Handlers struct {
	Auth    *handler.AuthHandler
	User    *handler.UserHandler
	Problem *handler.ProblemHandler
	Contest *handler.ContestHandler
	Sheet   *handler.SheetHandler
	Social  *handler.SocialHandler
	Room    *handler.RoomHandler
	RoomWS  *websocket.RoomHandler
}
