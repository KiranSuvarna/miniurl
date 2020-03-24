package api

import (
	"fmt"

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

func (s *Service) persistTheCount(counter *schema.Counter) {
	if err := s.pg.Insert(&schema.Counter{
		MachineID: counter.MachineID,
		Count:     counter.Count,
	}); err != nil {
		log.WithError(err).Error("Failed to insert the counter!!")
	}
}

func (s *Service) getURLFromID(id int) (*schema.Mappings, error) {
	var err error
	var mapping schema.Mappings
	if err = s.pg.Query(map[string]interface{}{"count": id}, fmt.Sprintf("%s", "url"), 0, 1, &mapping); err == nil {
		return &mapping, nil
	}
	log.WithError(err).Error("Failed to get the mapping")
	return nil, err
}
