package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "QuickSlot/docs" // swagger documentation

	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/time/rate"

	"QuickSlot/internal/handler"
	"QuickSlot/internal/middleware"
	"QuickSlot/internal/repository"
	"QuickSlot/internal/service"
	"QuickSlot/internal/worker"
	"QuickSlot/pkg/database/mysql"
)

// @title QuickSlot API
// @version 1.0
// @description Backend for appointment booking service.
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	cfg := loadConfig()

	dialect := mysql.NewMySQLDialect(context.Background(), cfg)

	// repos
	userRepo := repository.NewUserRepository(dialect)
	appointmentRepo := repository.NewAppointmentRepository(dialect)
	slotRepo := repository.NewSlotRepository(dialect)
	orgRepo := repository.NewOrganizationRepository(dialect)
	empRepo := repository.NewEmployeeRepository(dialect)
	reviewRepo := repository.NewReviewRepository(dialect)

	// services
	authService := service.NewAuthService(userRepo)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	slotService := service.NewSlotService(slotRepo)
	orgService := service.NewOrganizationService(orgRepo)
	empService := service.NewEmployeeService(empRepo, orgRepo)
	reviewService := service.NewReviewService(reviewRepo)

	// handlers
	authHandler := handler.NewAuthHandler(authService)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	slotHandler := handler.NewSlotHandler(slotService)
	orgHandler := handler.NewOrganizationHandler(orgService)
	empHandler := handler.NewEmployeeHandler(empService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	mux := http.NewServeMux()

	// Swagger documentation route
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// auth
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	// protected test
	mux.Handle("/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("protected route"))
	})))

	// slots
	mux.Handle("/slots/generate", middleware.AuthMiddleware(
		middleware.AdminOnly(http.HandlerFunc(slotHandler.Generate))))
	mux.HandleFunc("/slots/available", slotHandler.GetAvailableByEmployee)

	// appointments
	mux.Handle("/appointments/book", middleware.AuthMiddleware(http.HandlerFunc(appointmentHandler.Book)))
	mux.Handle("/appointments/cancel", middleware.AuthMiddleware(http.HandlerFunc(appointmentHandler.Cancel)))
	mux.Handle("/appointments/history", middleware.AuthMiddleware(http.HandlerFunc(appointmentHandler.History)))

	// organizations
	mux.Handle("/organizations/create", middleware.AuthMiddleware(http.HandlerFunc(orgHandler.Create)))
	mux.HandleFunc("/organizations", orgHandler.GetAll)
	mux.HandleFunc("/organizations/get", orgHandler.GetByID)
	mux.Handle("/organizations/update", middleware.AuthMiddleware(
		middleware.AdminOnly(http.HandlerFunc(orgHandler.Update))))
	mux.Handle("/organizations/delete", middleware.AuthMiddleware(
		middleware.AdminOnly(http.HandlerFunc(orgHandler.Delete))))

	// employees
	mux.Handle("/employees/create", middleware.AuthMiddleware(
		middleware.AdminOnly(http.HandlerFunc(empHandler.Create))))
	mux.HandleFunc("/employees", empHandler.GetByOrganization)
	mux.Handle("/employees/update", middleware.AuthMiddleware(
		middleware.AdminOnly(http.HandlerFunc(empHandler.Update))))
	mux.Handle("/employees/delete", middleware.AuthMiddleware(
		middleware.AdminOnly(http.HandlerFunc(empHandler.Delete))))

	// reviews
	mux.Handle("/reviews/create", middleware.AuthMiddleware(http.HandlerFunc(reviewHandler.Create)))
	mux.HandleFunc("/reviews", reviewHandler.GetByOrganization)
	mux.Handle("/reviews/delete", middleware.AuthMiddleware(http.HandlerFunc(reviewHandler.Delete)))

	// middleware chain: RateLimit -> CORS -> Logging
	rateLimitedMux := middleware.RateLimit(rate.Limit(10), 20)(mux) // 10 req/s, 20 burst
	wrapped := middleware.LoggingMiddleware(middleware.CORSMiddleware(rateLimitedMux))

	port := getEnv("SERVER_PORT", "8080")
	server := &http.Server{
		Addr:    ":" + port,
		Handler: wrapped,
	}

	// background worker
	done := make(chan struct{})
	go worker.CleanExpiredSlots(dialect, 5*time.Minute, done)

	// graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("server started on :%s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	close(done) // stop worker

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("shutdown error: %v", err)
	}

	log.Println("server stopped")
}
