package repositories

import (
	"context"

	"github.com/muhammadswa/personal-library/internal/database"
)

type AppDB interface {
	GetBookByID(ctx context.Context, id int32) (database.Book, error)
	CreateBook(ctx context.Context, book database.CreateBookParams) (int32, error)
	GetBooks(ctx context.Context, book database.GetBooksParams) ([]database.Book, error)
	GetBooksLength(ctx context.Context) (int64, error)
	UpdateBook(ctx context.Context, book database.UpdateBookParams) error
}

//	type Repositories struct {
//		db *database.Queries
//	}
type Repositories struct {
	db AppDB
}

func New(dbQueries *database.Queries) *Repositories {
	return &Repositories{
		db: dbQueries,
	}
}
