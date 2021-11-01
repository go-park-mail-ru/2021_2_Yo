package repository

import (
	sql "github.com/jmoiron/sqlx"
	"strconv"
)

const (
	updateUserImg 		= `update "user" set imgUrl = $1 where id = $2`
)

type Repository struct {
	db *sql.DB
}

func NewRepository(database *sql.DB) *Repository {
	return &Repository{
		db: database,
	}
}

func (img *Repository) StoreImage(userId string, fileName string) error {
	filepath := "~\\go\\2021_2_Yo\\static\\images"+fileName
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		//return error2.ErrAtoi
	}
	query := updateUserImg
		_, err = img.db.Exec(query, filepath, userIdInt)
		if err != nil {
			//return error2.ErrPostgres
		}
		return nil

}