package api

import (
	"sync"
	"time"

	"bitbucket.org/smartclean/routines-go/config"
	"bitbucket.org/smartclean/routines-go/db"
	"bitbucket.org/smartclean/routines-go/schema"
	gintemplate "github.com/foolin/gin-template"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// Service HTTP server info
type Service struct {
	shutdownChan chan bool

	router    *gin.Engine
	wg        sync.WaitGroup
	rc        *db.RedisConn
	channel   map[string]schema.Data
	AppName   string
	Version   string
	BuildTime string
}

// NewService Create a new service
func InitService(conf *config.Config) (*Service, error) {

	rc, err := db.NewRedis(&conf.Redis)
	if err != nil {
		log.WithError(err).Error("Failed to connect to redis")

		return nil, err
	}

	s := &Service{
		router:       gin.New(),
		shutdownChan: make(chan bool),
		rc:           rc,
		channel :    make(map[string]schema.Data),
	}

	s.router.Use(gin.Logger())

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS", "PUT", "POST"},
		AllowHeaders:     []string{"origin"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.router.HTMLRender = gintemplate.New(gintemplate.TemplateConfig{
		Root:      "templates",
		Extension: ".html",
		DisableCache: true,
	})

	s.router.GET("/", s.index)
	v1 := s.router.Group("v1")
	{
		v1.GET("/_create", s.create)
		v1.GET("/_check", s.check)
		v1.PUT("/_pause", s.pause)
		v1.PUT("/_clear", s.clear)
		v1.GET("/_render", s.render)
		v1.POST("/_snapshot", s.snapshot)
	}
	return s, nil
}

// Start the web service
func (s *Service) Start(address string) error {
	return s.router.Run(address)
}

// Close all threads and free up resources
func (s *Service) Close() {
	close(s.shutdownChan)

	s.wg.Wait()

}
