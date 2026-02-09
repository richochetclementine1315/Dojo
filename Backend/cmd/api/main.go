package main

import (
	"dojo/internal/config"
	"dojo/internal/handler"
	"dojo/internal/middleware"
	"dojo/internal/models"
	"dojo/internal/repository"
	"dojo/internal/routes"
	"dojo/internal/service"
	"dojo/internal/utils"
	"dojo/internal/websocket"
	"dojo/pkg/database"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load Configs
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	// Initialize Validator
	utils.InitValidator()

	// DB ka connection and other initializations go here
	isDebug := cfg.App.Env == "development" // Enable GORM debug mode in development
	db, err := database.Connect(cfg.GetDSN(), isDebug)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	database.DB = db
	log.Println("Database connected successfully")

	// AutoMigrate models
	if err := database.AutoMigrate(
		&models.User{},
		&models.AuthAccount{},
		&models.RefreshToken{},
		&models.UserProfile{},
		&models.UserPlatformStat{},
		&models.Friend{},
		&models.FriendRequest{},
		&models.BlockedUser{},
		&models.Notification{},
		&models.Problem{},
		&models.ProblemSheet{},
		&models.SheetProblem{},
		&models.UserNote{},
		&models.UserProblemProgress{},
		&models.Contest{},
		&models.ContestReminder{},
		&models.Room{},
		&models.RoomParticipant{},
		&models.CodeSession{},
		&models.WhiteboardSession{},
		&models.WhiteboardStroke{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	authRepo := repository.NewAuthRepository(db)
	problemRepo := repository.NewProblemRepository(db)
	contestRepo := repository.NewContestRepository(db)
	sheetRepo := repository.NewSheetRepository(db)
	socialRepo := repository.NewSocialRepository(db)
	roomRepo := repository.NewRoomRepository(db)

	// initialize Services
	authService := service.NewAuthService(userRepo, authRepo, cfg)
	userService := service.NewUserService(userRepo)
	problemService := service.NewProblemService(problemRepo)
	contestService := service.NewContestService(contestRepo)
	sheetService := service.NewSheetService(sheetRepo, problemRepo)
	socialService := service.NewSocialService(socialRepo, userRepo)
	roomService := service.NewRoomService(roomRepo, userRepo)

	// Initialize contest sync service
	contestSyncService := service.NewContestSyncService(contestService)
	// Sync contests every 6 hours
	go contestSyncService.Start(6 * time.Hour)
	log.Println("Contest sync service started (syncing every 6 hours)")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, cfg)
	userHandler := handler.NewUserHandler(userService)
	problemHandler := handler.NewProblemHandler(problemService)
	contestHandler := handler.NewContestHandler(contestService)
	sheetHandler := handler.NewSheetHandler(sheetService)
	socialHandler := handler.NewSocialHandler(socialService)
	roomHandler := handler.NewRoomHandler(roomService)

	// Initialize WebSocket Hub
	wsHub := websocket.NewHub()
	// Starting hub in background
	go wsHub.Run()
	// Initialize WebSocket handler
	roomWSHandler := websocket.NewRoomHandler(wsHub)

	// Create Fiber App
	app := fiber.New(fiber.Config{
		AppName:      cfg.App.Name,
		ErrorHandler: middleware.ErrorHandler,
	})
	// TODO: Set up middlewares

	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CORSMiddleware(cfg))
	app.Use(middleware.RateLimitMiddleware(cfg))

	// setup routes
	handlers := &routes.Handlers{
		Auth:    authHandler,
		User:    userHandler,
		Problem: problemHandler,
		Contest: contestHandler,
		Sheet:   sheetHandler,
		Social:  socialHandler,
		Room:    roomHandler,
		RoomWS:  roomWSHandler,
	}
	routes.SetupRoutes(app, handlers, cfg)

	// Start server
	log.Printf("Starting server on port %s", cfg.Server.Port)
	err = app.Listen("0.0.0.0:" + cfg.Server.Port)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
