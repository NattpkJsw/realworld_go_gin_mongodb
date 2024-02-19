package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/NattpkJsw/real-world-api-go/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
	GetServer() *server
}

type server struct {
	app *fiber.App
	db  *sqlx.DB
	cfg config.IConfig
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	// Middlewares
	middlewares := InitMiddlewares(s)
	s.app.Use(middlewares.Logger())
	s.app.Use(middlewares.Cors())

	// Modules
	api := s.app.Group("api")
	modules := InitModule(api, s, middlewares)

	modules.MonitorModule()
	modules.UsersModule()
	modules.ProfileModule()
	modules.ArticleModule()
	modules.CommentModule()
	modules.TagModule()
	modules.UserModule()

	s.app.Use(middlewares.RouterCheck())

	// Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		//_ = <-c
		log.Println("Server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	log.Printf("server is starting on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}

func (s *server) GetServer() *server {
	return s
}
