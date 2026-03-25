package main

import (
	"QuickSlot/internal/handler"
	"QuickSlot/internal/middleware"
	"QuickSlot/internal/repository"
	"QuickSlot/internal/service"
	"QuickSlot/pkg/database/mysql"
	"log"
	"net/http"
)

func main() {

	cfg := loadConfig()

	dialect := mysql.NewMySQLDialect(nil, cfg)

	userRepo := repository.NewUserRepository(dialect)
	authService := service.NewAuthService(userRepo)
	authHandler := handler.NewAuthHandler(authService)

	appointmentRepo := repository.NewAppointmentRepository(dialect)
	appointmentService := service.NewAppointmentService(appointmentRepo)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)

	slotRepo := repository.NewSlotRepository(dialect)
	slotService := service.NewSlotService(slotRepo)
	slotHandler := handler.NewSlotHandler(slotService)

	mux := http.NewServeMux()

	mux.Handle("/slots/generate", middleware.AuthMiddleware(http.HandlerFunc(slotHandler.Generate)))
	mux.HandleFunc("/slots/available", slotHandler.GetAvailableByEmployee)
	mux.Handle("/appointments/book", middleware.AuthMiddleware(http.HandlerFunc(appointmentHandler.Book)))
	mux.HandleFunc("/register", authHandler.Register)
	mux.HandleFunc("/login", authHandler.Login)

	protected := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("protected route"))
	})

	mux.Handle("/me", middleware.AuthMiddleware(protected))

	log.Println("server started on :8080")
	http.ListenAndServe(":8080", mux)
}
