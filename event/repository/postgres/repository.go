package postgres

import (
	"backend/models"
	//"database/sql"
	log "github.com/sirupsen/logrus"
	sql "github.com/jmoiron/sqlx"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) List() ([]*models.Event, error) {
	query := `select * from "event"`
	rows, err := s.db.Queryx(query)
	if err != nil {
		log.Error("Event : repository : postgres : List() err = ", err)
		return nil, err
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var event Event
		err := rows.StructScan(&event)
		if err != nil {
			log.Error("Event : repository : postgres : List() StructScan err = ", err)
			return nil, err
		}
		modelEvent := toModelEvent(&event)
		resultEvents = append(resultEvents, modelEvent)
	}
	log.Info("Events Postgres List = ", resultEvents)
	return resultEvents, nil
}