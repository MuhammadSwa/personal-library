package repositories

import (
	"context"
	"fmt"

	"github.com/muhammadswa/personal-library/internal/database"
)

type BooksRepository struct {
	db *database.Queries
}

func NewBooksRepository(dbQueries *database.Queries) *BooksRepository {
	return &BooksRepository{
		db: dbQueries,
	}
}

func (br BooksRepository) GetBookByID(ctx context.Context, id int) (*database.Book, error) {
	book, err := br.db.GetBookByID(ctx, int32(id))
	fmt.Println(err)
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (br BooksRepository) CreateBook(ctx context.Context, book database.CreateBookParams) (int, error) {
	id, err := br.db.CreateBook(ctx, book)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

func (br BooksRepository) GetBooks(ctx context.Context, userId int32, query string, offset int) ([]database.Book, error) {
	if offset < 0 {
		offset = 0
	}
	var books []database.Book
	var err error
	if query == "" {
		books, err = br.db.GetBooks(ctx, database.GetBooksParams{
			Offset: int32(offset),
			UserID: userId,
		})
	} else {
		books, err = br.db.GetBooksBy(ctx, database.GetBooksByParams{
			PlaintoTsquery: query,
			Offset:         int32(offset),
			UserID:         userId,
		})
	}
	if err != nil {
		return nil, err
	}
	return books, nil
}

func (br BooksRepository) GetBooksLength(ctx context.Context) (int, error) {
	length, err := br.db.GetBooksLength(ctx)
	if err != nil {
		return 0, err
	}
	return int(length), nil
}

func (br BooksRepository) UpdateBook(ctx context.Context, book database.UpdateBookParams) error {
	err := br.db.UpdateBook(ctx, book)
	if err != nil {
		return err
	}
	return nil
}
