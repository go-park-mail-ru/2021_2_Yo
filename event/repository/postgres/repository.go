package postgres

import (
	log "backend/logger"
	"backend/models"
	sql "github.com/jmoiron/sqlx"
)

const logMessage = "event:repository:postgres:"

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (s *Repository) List() ([]*models.Event, error) {
	message := logMessage + "List:"
	query := `select * from "event"`
	rows, err := s.db.Queryx(query)
	if err != nil {
		log.Error(message+"err =", err)
		return nil, err
	}
	defer rows.Close()
	var resultEvents []*models.Event
	for rows.Next() {
		var event Event
		err := rows.StructScan(&event)
		if err != nil {
			log.Error(message+"err =", err)
			return nil, err
		}
		modelEvent := toModelEvent(&event)
		resultEvents = append(resultEvents, modelEvent)
	}
	log.Debug(message+"resultEvents =", resultEvents)
	return resultEvents, nil
}
