package seed

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jphilipstevens/web-service-gin/app/db"
	"github.com/jphilipstevens/web-service-gin/example/features/albums"
	"github.com/jphilipstevens/web-service-gin/testUtils"
)

var data = []albums.Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	{ID: "4", Title: "Kind of Blue", Artist: "Miles Davis", Price: 56.99},
	{ID: "5", Title: "Everlong", Artist: "Blink-182", Price: 19.99},
	{ID: "6", Title: "The Wall", Artist: "Blink-182", Price: 19.99},
	{ID: "7", Title: "Going to California", Artist: "Blink-182", Price: 19.99},
	{ID: "8", Title: "The One And Only", Artist: "Blink-182", Price: 19.99},
	{ID: "9", Title: "A Love Supreme", Artist: "John Coltrane", Price: 49.99},
	{ID: "10", Title: "Bitches Brew", Artist: "Miles Davis", Price: 39.99},
	{ID: "11", Title: "Take Five", Artist: "Dave Brubeck", Price: 24.99},
	{ID: "12", Title: "Giant Steps", Artist: "John Coltrane", Price: 29.99},
	{ID: "13", Title: "Ella and Louis", Artist: "Ella Fitzgerald", Price: 34.99},
	{ID: "14", Title: "What's Going On", Artist: "Marvin Gaye", Price: 28.99},
	{ID: "15", Title: "All the Things You Are", Artist: "Ella Fitzgerald", Price: 32.99},
	{ID: "16", Title: "In a Silent Way", Artist: "Miles Davis", Price: 36.99},
	{ID: "17", Title: "Untitled", Artist: "Blink-182", Price: 21.99},
	{ID: "18", Title: "Mingus Ah Um", Artist: "Charles Mingus", Price: 27.99},
	{ID: "19", Title: "Sketches of Spain", Artist: "Miles Davis", Price: 42.99},
	{ID: "20", Title: "Dookie", Artist: "Green Day", Price: 18.99},
}

func createAlbumsTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS albums (
			id     SERIAL PRIMARY KEY,
			title  TEXT NOT NULL,
			artist TEXT NOT NULL,
			price  NUMERIC(10, 2)
		)
	`)
	return err
}

func SeedAlbums(dbConn db.Database) error {
	ctx := testUtils.CreateTestContext()
	// Create the albums table if it doesn't exist
	if err := createAlbumsTable(ctx, dbConn.GetClient()); err != nil {
		return fmt.Errorf("fatal error cannot create Album Table: %w", err)
	}

	contractsRepository := albums.NewAlbumRepository(dbConn)

	_, err := dbConn.GetClient().ExecContext(ctx, "TRUNCATE TABLE contracts")
	if err != nil {
		return fmt.Errorf("failed to truncate table: %w", err)
	}

	err = contractsRepository.InsertBatch(ctx, data)
	if err != nil {
		return fmt.Errorf("failed to insert album: %w", err)
	}
	return nil
}
