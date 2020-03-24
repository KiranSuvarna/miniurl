package api

import (
	"strconv"
	"strings"
	"sync"
	"time"

	"bitbucket.org/mine/miniurl/config"
	"bitbucket.org/mine/miniurl/db"
	"bitbucket.org/mine/miniurl/schema"
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
	pg        *db.Postgres
	Counter   schema.Counter
	AppName   string
	Version   string
	BuildTime string
}

// NewService Create a new service
func InitService(conf *config.Config) (*Service, error) {

	pg, err := db.NewPostgres(&conf.Postgres)

	if err != nil {
		log.WithError(err).Error("Failed to connect to DB")

		return nil, err
	}

	rc, err := db.NewRedis(&conf.Redis)
	if err != nil {
		log.WithError(err).Error("Failed to connect to redis")

		return nil, err
	}

	mID := &conf.Counter.MachineID
	r := &conf.Counter.Range
	c, _ := strconv.Atoi(strings.Split(*r, "-")[0])

	s := &Service{
		router:       gin.New(),
		shutdownChan: make(chan bool),
		rc:           rc,
		pg:           pg,
		Counter: schema.Counter{
			MachineID: *mID,
			Count:     c,
		},
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
		Root:         "templates",
		Extension:    ".html",
		DisableCache: true,
	})

	s.router.GET("/", s.index)
	v1 := s.router.Group("v1")
	{
		v1.POST("/mini", s.getMiniURL)
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
