package albums

import (
	"context"
	"database/sql"
	"example/web-service-gin/app/db"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewAlbumRepository(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)
	assert.NotNil(t, repo)
}

func TestGetAlbumsRepository(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	t.Run("Get all albums", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", 9.99).
			AddRow("2", "Album 2", "Artist 2", 14.99)
		mock.ExpectQuery("SELECT \\* FROM albums").WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
		result, err := repo.GetAlbums(context.Background(), GetAlbumsParams{
			Artist: "",
			Limit:  10,
			Page:   0,
		})

		expected := &db.Paginated[Album]{
			Items: []Album{
				{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
				{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
			},
			Total: 2,
		}

		assert.NoError(t, err)
		assert.Equal(t, 2, result.Total)
		assert.Equal(t, 2, len(result.Items))
		assert.Equal(t, reflect.DeepEqual(expected.Items, result.Items), true)

	})

	t.Run("Get all albums with artist", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", 9.99)

		mock.ExpectQuery("SELECT \\* FROM albums WHERE artist ILIKE \\$1 LIMIT \\$2 OFFSET \\$3").WithArgs("Artist 1", 10, 0).WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums WHERE artist SIMILAR TO \\$1").WithArgs("Artist 1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
		result, err := repo.GetAlbums(context.Background(), GetAlbumsParams{
			Artist: "Artist 1",
			Limit:  10,
			Page:   0,
		})
		expected := &db.Paginated[Album]{
			Items: []Album{
				{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
			},
			Total: 1,
		}
		assert.NoError(t, err)
		assert.Equal(t, 1, result.Total)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, reflect.DeepEqual(expected.Items, result.Items), true)
	})

	// testCases := []struct {
	// 	name     string
	// 	params   GetAlbumsParams
	// 	mockFunc func()
	// 	expected *db.Paginated[Album]
	// 	err      error
	// }{
	// 	// {
	// 	// 	name:   "Get all albums",
	// 	// 	params: GetAlbumsParams{Artist: "", Limit: 10, Page: 0},
	// 	// 	mockFunc: func() {
	// 	// 		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
	// 	// 			AddRow("1", "Album 1", "Artist 1", 9.99).
	// 	// 			AddRow("2", "Album 2", "Artist 2", 14.99)
	// 	// 		mock.ExpectQuery("SELECT \\* FROM albums").WillReturnRows(rows)
	// 	// 		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
	// 	// 	},
	// 	// 	expected: &db.Paginated[Album]{
	// 	// 		Items: []Album{
	// 	// 			{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
	// 	// 			{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
	// 	// 		},
	// 	// 		Total: 2,
	// 	// 	},
	// 	// 	err: nil,
	// 	// },
	// 	// {
	// 	// 	name:   "Get albums by artist",
	// 	// 	params: GetAlbumsParams{Artist: "Artist 1", Limit: 10, Page: 0},
	// 	// 	mockFunc: func() {
	// 	// 		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
	// 	// 			AddRow("1", "Album 1", "Artist 1", 9.99)
	// 	// 		mock.ExpectQuery("SELECT \\* FROM albums WHERE artist ILIKE \\$1").WithArgs("Artist 1").WillReturnRows(rows)
	// 	// 		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums WHERE artist SIMILAR TO \\$1").WithArgs("Artist 1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
	// 	// 	},
	// 	// 	expected: &db.Paginated[Album]{
	// 	// 		Items: []Album{
	// 	// 			{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
	// 	// 		},
	// 	// 		Total: 1,
	// 	// 	},
	// 	// 	err: nil,
	// 	// },
	// }

	// for _, tc := range testCases {
	// 	t.Run(tc.name, func(t *testing.T) {
	// 		tc.mockFunc()
	// 		result, err := repo.GetAlbums(context.Background(), tc.params)
	// 		assert.Equal(t, tc.err, err)
	// 		assert.Equal(t, tc.expected, result)
	// 	})
	// }
}

func TestInsert(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	album := Album{ID: "1", Title: "New Album", Artist: "New Artist", Price: 19.99}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs(album.ID, album.Title, album.Artist, album.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Insert(context.Background(), album)
	assert.NoError(t, err)
}

func TestInsertBatch(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	albums := []Album{
		{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
		{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
	}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs("1", "Album 1", "Artist 1", 9.99, "2", "Album 2", "Artist 2", 14.99).
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := repo.InsertBatch(context.Background(), albums)
	assert.NoError(t, err)
}

func TestGetAlbumsNoResults(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	mock.ExpectQuery("SELECT \\* FROM albums").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	result, err := repo.GetAlbums(context.Background(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, db.NotFoundError, err)
}

func TestGetAlbumsError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	mock.ExpectQuery("SELECT \\* FROM albums").WillReturnError(sql.ErrConnDone)

	result, err := repo.GetAlbums(context.Background(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestInsertError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	album := Album{ID: "1", Title: "New Album", Artist: "New Artist", Price: 19.99}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs(album.ID, album.Title, album.Artist, album.Price).
		WillReturnError(sql.ErrConnDone)

	err := repo.Insert(context.Background(), album)
	assert.Error(t, err)
}

func TestInsertBatchEmptySlice(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	err := repo.InsertBatch(context.Background(), []Album{})
	assert.NoError(t, err)
}

func TestInsertBatchError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	albums := []Album{
		{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
		{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
	}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs("1", "Album 1", "Artist 1", 9.99, "2", "Album 2", "Artist 2", 14.99).
		WillReturnError(sql.ErrConnDone)

	err := repo.InsertBatch(context.Background(), albums)
	assert.Error(t, err)
}

func TestGetAlbumsCountError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	mock.ExpectQuery("SELECT \\* FROM albums").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", 9.99))

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").
		WillReturnError(sql.ErrConnDone)

	result, err := repo.GetAlbums(context.Background(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, db.DatabaseError, err)
}

func TestGetAlbumsScanError(t *testing.T) {
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(mockDB)

	mock.ExpectQuery("SELECT \\* FROM albums").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", "invalid_price"))

	result, err := repo.GetAlbums(context.Background(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, db.DatabaseError, err)
}
