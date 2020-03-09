package api

import (
	"fmt"
	"net/http"
	"time"

	"bitbucket.org/smartclean/routines-go/schema"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type PingResponse struct {
	Name    string `json:"app"`
	Version string `json:"version"`
	Built   string `json:"built"`
	Status  bool   `json:"status"`
}

func (s *Service) index(c *gin.Context) {
	c.String(http.StatusOK, "Ad Request Auction System")
}

func (s *Service) ping(c *gin.Context) {
	status := true

	c.JSON(http.StatusOK, PingResponse{
		Status:  status,
		Name:    s.AppName,
		Version: s.Version,
		Built:   s.BuildTime,
	})
}

func (s *Service) responseWriter(c *gin.Context, resp interface{}, code int) {
	c.Header("Cache-Control", "no-store")
	c.Header("Connection", "close")

	origin := c.GetHeader("Origin")
	if origin != "" { // TODO: validate the origin
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,OPTIONS")
	}

	c.JSON(code, resp)
}

// This either closes the routine or increments the counter
// Ref : https://yourbasic.org/golang/stop-goroutine/
func (s *Service) worker(data schema.Data) {
	for {
		select {
		case <-data.Channel:
			return
		default:
			if *data.Active {
				*data.Count += data.Start
				time.Sleep(time.Second * time.Duration(data.StepTime))
			}
		}
	}
}

func (s *Service) persistRoutinesdata() bool{
	if len(s.channel) > 0 {
		msg := ""
		log.Info("Started taking snapshot to persist data on cache...")
		for k, v := range s.channel {
			d := schema.CheckData{
				ID:        v.ID,
				StepTime:  v.StepTime,
				Count:     *v.Count,
				Status:    *v.Status,
				CreatedAt: v.CreatedAt,
			}
			if err := s.rc.Set(fmt.Sprintf("sc-%s", k), d); err != nil {
				log.WithError(err).WithField("key", fmt.Sprintf("sc-%s", k)).Warn("Failed to set cache")
			}
		}
		msg = fmt.Sprintf("Fisnished the snapshot")
		log.Info(msg)
		return true
	}
	return false
}
