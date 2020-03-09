package api

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"bitbucket.org/smartclean/routines-go/schema"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

func (s *Service) create(c *gin.Context) {
	start, _ := strconv.Atoi(c.DefaultQuery("start", "0"))
	step, _ := strconv.Atoi(c.DefaultQuery("step", "0"))

	status := fmt.Sprintf("%s", "active")
	active := true

	data := schema.Data{
		Start:     start,
		StepTime:  step,
		ID:        uuid.New().String(),
		Channel:   make(chan string),
		Active:    &active,
		CreatedAt: time.Now(),
		Count:     new(int),
		Status:    &status,
	}

	go s.worker(data)

	// TODO: Consider stopping the routine after the task
	go s.persistRoutinesdata()
	
	// Store it in the map for quick access - O(1)
	s.channel[data.ID] = data

	s.responseWriter(c, schema.Response{
		Data:    data.ID,
		Success: true,
	}, http.StatusOK)

	return
}

func (s *Service) check(c *gin.Context) {
	var res interface{}
	if id := c.DefaultQuery("id", ""); id != "" {
		if v, ok := s.channel[id]; ok {
			res = schema.Response{
				Data: schema.CheckData{
					ID:        v.ID,
					StepTime:  v.StepTime,
					Count:     *v.Count,
					Status:    *v.Status,
					CreatedAt: v.CreatedAt, 
				},
				Success: true,
			}
		} else {
			log.Warning("Can't able to find the data from given id : ", id)
			res = schema.Response{
				Data:    fmt.Sprintf("No data for given id : %s", id),
				Success: false,
			}
		}
	} else {
		data := make([]schema.CheckData, 0)
		if len(s.channel) > 0 {
			for _, v := range s.channel {
				data = append(data, schema.CheckData{
					ID:        v.ID,
					StepTime:  v.StepTime,
					Count:     *v.Count,
					Status:    *v.Status,
					CreatedAt: v.CreatedAt,
				})
			}
		}
		res = schema.Response{
			Data:    data,
			Success: true,
		}
	}
	s.responseWriter(c, res, http.StatusOK)
	return
}

func (s *Service) pause(c *gin.Context) {
	var res interface{}
	if id := c.DefaultQuery("id", ""); id != "" {
		if v, ok := s.channel[id]; ok {
			status := fmt.Sprintf("%s", "paused")
			*v.Active = false
			*v.Status = status
			v.ModifiedAt = time.Now()

			go s.persistRoutinesdata()

			res = schema.Response{
				Data: schema.PauseData{
					ID:        v.ID,
					PauseTime: v.ModifiedAt,
				},
				Success: true,
			}
		} else {
			log.Warning("Can't able to find the data from given id : ", id)
			res = schema.Response{
				Data:    fmt.Sprintf("No data for given id : %s", id),
				Success: false,
			}
		}
	} else {
		log.Warning("No id : ", id)
		res = schema.Response{
			Data:    fmt.Sprintf("No input"),
			Success: false,
		}
	}
	s.responseWriter(c, res, http.StatusOK)
}

func (s *Service) clear(c *gin.Context) {
	var res interface{}
	if id := c.DefaultQuery("id", ""); id != "" {
		if v, ok := s.channel[id]; ok {
			status := fmt.Sprintf("%s", "stopped")
			*v.Status = status
			v.Channel <- v.ID
			close(v.Channel)

			go s.persistRoutinesdata()

			res = schema.Response{
				Data:    fmt.Sprintf("The go-routine with id %s is stopped", v.ID),
				Success: true,
			}
		} else {
			log.Warning("Can't able to find the data from given id : ", id)
			res = schema.Response{
				Data:    fmt.Sprintf("No data for given id : %s", id),
				Success: false,
			}
		}
	} else {
		log.Warning("No id : ", id)
		res = schema.Response {
			Data:    fmt.Sprintf("No input"),
			Success: false,
		}
	}
	s.responseWriter(c, res, http.StatusOK)
}

func (s *Service) snapshot(c *gin.Context) {
	msg := ""
	if ok := s.persistRoutinesdata(); ok {
		msg = fmt.Sprint("Fisnished the snapshot")
	} else {
		msg = fmt.Sprint("Something went wrong")		
	}
	s.responseWriter(c,schema.Response{
		Data: msg,
		Success: true,
	},http.StatusOK)
	return
}

func (s *Service) render(c *gin.Context) {
	var render schema.HTMLRender
	if len(s.channel) > 0 {
		for _, v := range s.channel {
			render.Render = append(render.Render, schema.CheckData{
				ID:        v.ID,
				StepTime:  v.StepTime,
				Count:     *v.Count,
				Status:    *v.Status,
				CreatedAt: v.CreatedAt,
			})
		}
	}
	c.HTML(http.StatusOK, "index", render)
	return
}