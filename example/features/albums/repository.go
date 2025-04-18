package albums

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jphilipstevens/web-service-gin/app/db"
)

const (
	serviceName = "album-service-repository"
)

type AlbumRepository interface {
	GetAlbums(ctx context.Context, params GetAlbumsParams) (*db.Paginated[Album], error)
	Insert(ctx context.Context, album Album) error
	InsertBatch(ctx context.Context, album []Album) error
}

type albumRepository struct {
	dbConn db.Database
}

func NewAlbumRepository(dbConn db.Database) AlbumRepository {
	return &albumRepository{
		dbConn: dbConn,
	}
}

func (ar *albumRepository) GetAlbums(ctx context.Context, params GetAlbumsParams) (*db.Paginated[Album], error) {

	var artist = params.Artist
	var query string
	var args []interface{}
	offset := (params.Page * -1) - params.Limit
	if offset < 0 {
		offset = 0
	}

	if artist != "" {
		query = "SELECT * FROM albums WHERE artist ILIKE $1 LIMIT $2 OFFSET $3"
		args = []interface{}{artist, params.Limit, offset}
	} else {
		query = "SELECT * FROM albums LIMIT $1 OFFSET $2"
		args = []interface{}{params.Limit, offset}
	}

	rows, err := ar.dbConn.QueryContext(serviceName, ctx, query, args...)
	if err != nil {
		return nil, db.MapDBError(&err)
	}
	defer rows.Close()

	var albums []Album
	for rows.Next() {
		var album Album
		if err := rows.Scan(&album.ID, &album.Title, &album.Artist, &album.Price); err != nil {
			return nil, db.MapDBError(&err)
		}
		albums = append(albums, album)

	}

	if len(albums) == 0 {
		return nil, db.MapDBError(&sql.ErrNoRows)
	}

	return &db.Paginated[Album]{
		Items: albums,
	}, nil
}

func (ar *albumRepository) Insert(ctx context.Context, album Album) error {
	_, err := ar.dbConn.ExecContext(serviceName, ctx, "INSERT INTO albums (id, title, artist, price) VALUES ($1, $2, $3, $4)", album.ID, album.Title, album.Artist, album.Price)
	if err != nil {
		return err
	}
	return nil
}

func (ar *albumRepository) InsertBatch(ctx context.Context, albums []Album) error {
	if len(albums) == 0 {
		return nil
	}

	values := make([]string, 0, len(albums))
	args := make([]interface{}, 0, len(albums)*4)
	for i, album := range albums {
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		args = append(args, album.ID, album.Title, album.Artist, album.Price)
	}

	query := fmt.Sprintf("INSERT INTO albums (id, title, artist, price) VALUES %s ON CONFLICT (id) DO UPDATE SET title = EXCLUDED.title, artist = EXCLUDED.artist, price = EXCLUDED.price", strings.Join(values, ","))

	_, err := ar.dbConn.ExecContext(serviceName, ctx, query, args...)
	return err
}
