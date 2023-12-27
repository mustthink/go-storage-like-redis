package internal

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"

	"github.com/mustthink/go-storage-like-redis/config"
	"github.com/mustthink/go-storage-like-redis/internal/handlers"
	"github.com/mustthink/go-storage-like-redis/internal/storage"
)

type Application struct {
	config  *config.Config
	storage storage.Storage
	logger  *logrus.Logger
}

func NewApplication(configPath string) *Application {
	log := logrus.New()
	log.SetLevel(logrus.DebugLevel)
	log.Debug("logger created")

	appConfig, err := config.New(configPath)
	if err != nil {
		log.Fatalf("couldn't create config w err: %s", err.Error())
	}
	log.Debug("config created")

	appStorage := storage.New(appConfig.StorageConfig)
	log.Debug("storage created")
	return &Application{
		config:  appConfig,
		logger:  log,
		storage: appStorage,
	}
}

func (a *Application) Run() {
	r := mux.NewRouter()

	mainHandler := func(writer http.ResponseWriter, request *http.Request) {
		handlers.Handler(writer, request, a.storage)
	}
	r.HandleFunc("/", handlers.BaseAuth(mainHandler, a.config.ServerConfig.Auth))

	server := &http.Server{
		Addr:         a.config.ServerConfig.URL(),
		Handler:      r,
		ReadTimeout:  a.config.ServerConfig.ReadTimeout * time.Millisecond,
		WriteTimeout: a.config.ServerConfig.WriteTimeout * time.Millisecond,
	}
	a.logger.Debug("start listening and serve")
	a.logger.Fatal(server.ListenAndServe())
}
