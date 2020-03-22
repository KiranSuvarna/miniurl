package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Name    string `json:"app"`
	Version string `json:"version"`
	Built   string `json:"built"`
	Status  bool   `json:"status"`
}

type ResultError struct {
	Res   bool
	Error error
}

func (s *Service) index(c *gin.Context) {
	c.String(http.StatusOK, "Mini URL Server")
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
