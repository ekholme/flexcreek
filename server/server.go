package server

import (
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/ekholme/flexcreek"
)

type Server struct {
	Router *http.ServeMux
	Srvr   *http.Server

	//services
	UserService         flexcreek.UserService
	WorkoutService      flexcreek.WorkoutService
	ActivityTypeService flexcreek.ActivityTypeService

	//templates
	Templates *template.Template

	//logging
	Logger *slog.Logger

	//config
	//todo
}

func NewServer(addr string, us flexcreek.UserService, ws flexcreek.WorkoutService, ats flexcreek.ActivityTypeService, templatePaths string, logFilePath string) *Server {

	tmpl, err := template.ParseGlob(templatePaths)

	if err != nil {
		log.Fatalf("couldn't parse templates: %v", err.Error())
	}

	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		log.Fatalf("couldn't open log file: %v", err.Error())
	}

	lh := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		AddSource: true,
		Level:     slog.LevelDebug,
	})

	logger := slog.New(lh)

	return &Server{
		Router: http.NewServeMux(),
		Srvr: &http.Server{
			Addr: addr,
		},
		UserService:         us,
		WorkoutService:      ws,
		ActivityTypeService: ats,
		Templates:           tmpl,
		Logger:              logger,
	}
}

func (s *Server) registerRoutes() {
	//html routes
	//todo

	//api routes
	s.Router.HandleFunc("/api", s.handleApiIndex)
}

func (s *Server) Run() {
	s.registerRoutes()

	s.Srvr.Handler = s.Router

	s.Logger.Info("Starting server", slog.String("addr", s.Srvr.Addr))

	fmt.Printf("Starting server on %s\n", s.Srvr.Addr)

	s.Srvr.ListenAndServe()
}
