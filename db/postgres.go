package db

import (
	"fmt"
	"strconv"

	"bitbucket.org/mine/miniurl/config"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/lib/pq"

	log "github.com/sirupsen/logrus"
)

type Postgres struct {
	db *gorm.DB
}

func NewPostgres(conf *config.PostgresConfig) (*Postgres, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s dbname=%s sslmode=%s",
		conf.Host, strconv.Itoa(conf.Port), conf.Database, conf.SSLMode)
	if db, err := gorm.Open("postgres", psqlInfo); err == nil {
		db.SingularTable(true)
		// TODO : Change these values depending on the results of load-testing
		//db.DB().SetMaxOpenConns(100)
		//db.DB().SetMaxIdleConns(25)
		//db.DB().SetConnMaxLifetime(1 * time.Hour)
		return &Postgres{
			db: db,
		}, nil
	} else {
		return nil, err
	}
}

func (pg *Postgres) Query(query map[string]interface{}, fields interface{}, pn, pp int, result interface{}) error {
	if err := pg.db.Select(fields).Limit(pp).Offset((pn - 1) * pp).Where(query).Find(result).Error; err != nil {
		log.WithError(err).WithField("query", query).Error("Failed to query the data!")
		return err
	}
	return nil
}

func (pg *Postgres) QueryRaw(query string, result interface{}) error {
	if err := pg.db.Raw(query).Scan(result).Error; err != nil {
		log.WithError(err).WithField("query", query).Error("Failed to query the data!")
		return err
	}
	return nil
}

func (pg *Postgres) Insert(entity interface{}) error {
	if err := pg.db.Create(entity).Scan(entity).Error; err != nil {
		log.WithError(err).Error("Failed to insert the row!")
		return err
	}
	return nil
}

func (pg *Postgres) QueryTransaction(query map[string]interface{}, fields interface{}, pn, pp int, result interface{}, count *int) error {
	tx := pg.db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Select(fields).Where(query).Find(result).Count(count).Error; err != nil {
		log.WithError(err).WithField("query", query).Error("Failed to query the count!")
		tx.Rollback()
		return err
	}

	if err := tx.Select(fields).Limit(pp).Offset((pn - 1) * pp).Where(query).Find(result).Error; err != nil {
		log.WithError(err).WithField("query", query).Error("Failed to query the data!")
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
