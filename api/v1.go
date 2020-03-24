package api

import (
	"fmt"
	"net/http"
	"strings"

	"bitbucket.org/mine/miniurl/misc"
	"bitbucket.org/mine/miniurl/schema"
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	log "github.com/sirupsen/logrus"
)

type WriteResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type ReadResponse struct {
	Count   int         `json:"count"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

func (s *Service) getMiniFromURL(c *gin.Context) {

	var (
		message, success, code, chanRE = "", false, http.StatusInternalServerError, make(chan ResultError)
	)

	if url := c.DefaultPostForm("url", ""); url != "" {
		s.Counter.Count += 1
		hash := misc.Encode(s.Counter.Count)

		go s.persistTheMiniURLMapping(url, hash, s.Counter.Count, chanRE)
		resChan := <-chanRE
		defer close(chanRE)

		//TODO: Reconsider this later
		//go s.PersistTheCount(&s.Counter)

		if resChan.Res {
			code = http.StatusOK
			success = true
			message = fmt.Sprintf(misc.MiniURLDomain, hash)
		} else {
			message = resChan.Error.Error()
		}

		s.responseWriter(c, WriteResponse{
			Message: message,
			Success: success,
		}, code)

		return
	}
}


func (s *Service) getURLFromMini(c *gin.Context) {

	var (
		message, success, code = "", false, http.StatusInternalServerError
	)

	if url := c.DefaultQuery("url", ""); url != "" {
		id, err := misc.Decode(strings.Split(url,"/")[3])

		if err != nil {
			message = err.Error()
			log.WithField("Error", err).Error("Failed to decode the url")
		}

		cacheKey := fmt.Sprintf("mini-hash-%d", id)
		if url, sErr := s.rc.Get(cacheKey); sErr == nil && len(url) > 0 {
			var mapping schema.Mappings
	
			if rcErr := jsoniter.ConfigCompatibleWithStandardLibrary.UnmarshalFromString(url, &mapping); rcErr != nil {
				log.WithError(rcErr).WithField("key", cacheKey).Warn("Failed to decode cache value url mappings")
			}
		}

		mapping, err := s.getURLFromID(id); if err != nil {
			log.WithError(err).Warn("Failed to get the mapping")
		} else {
			message = mapping.URL
		}

		s.responseWriter(c, WriteResponse{
			Message: message,
			Success: success,
		}, code)

		return
	}
}
