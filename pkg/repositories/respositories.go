package repositories

import "github.com/muhammadswa/personal-library/internal/database"

type Repositories struct {
	db *database.Queries
}

func New(dbQueries *database.Queries) *Repositories {
	return &Repositories{
		db: dbQueries,
	}
}
