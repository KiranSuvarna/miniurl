package api

import (
	"fmt"
	"net/http"

	"bitbucket.org/mine/miniurl/misc"
	"github.com/gin-gonic/gin"
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

func (s *Service) getMini(c *gin.Context) {

	var (
		message, success, code, chanRE = "", false, http.StatusInternalServerError, make(chan ResultError)
	)

	if url := c.DefaultPostForm("url", ""); url != "" {
		//TODO: Get count from cache
		s.Counter.Count += 1

		hash := misc.GenerateBase62Hash(s.Counter.Count)

		go s.persistTheMiniURLMapping(url, hash, s.Counter.Count, chanRE)
		resChan := <-chanRE
		defer close(chanRE)

		go s.PersistTheCount(&s.Counter)

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
