package api

import (
	"bitbucket.org/mine/miniurl/schema"
	log "github.com/sirupsen/logrus"
)

func (s *Service) persistTheMiniURLMapping(url, hash string, count int, re chan ResultError) {
	select {
	case <-re:
		return
	default:
		var err error
		if err = s.pg.Insert(&schema.Mappings{
			URL:   url,
			Hash:  hash,
			Count: count,
		}); err != nil {
			log.WithError(err).Error("Failed to insert the data!")
			re <- ResultError{
				Res:   false,
				Error: err,
			}
		} else {
			re <- ResultError{
				Res:   true,
				Error: nil,
			}
		}
		return
	}
}

func (s *Service) PersistTheCount(counter *schema.Counter) {
	if err := s.pg.Insert(&schema.Counter{
		MachineID: counter.MachineID,
		Count:     counter.Count,
	}); err != nil {
		log.WithError(err).Error("Failed to insert the counter!!")
	}
}
