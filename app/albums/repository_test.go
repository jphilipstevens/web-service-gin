package albums

import (
	"context"
	"example/web-service-gin/app/db"
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

	testCases := []struct {
		name     string
		artist   string
		mockFunc func()
		expected *db.Paginated[Album]
		err      error
	}{
		{
			name:   "Get all albums",
			artist: "",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
					AddRow("1", "Album 1", "Artist 1", 9.99).
					AddRow("2", "Album 2", "Artist 2", 14.99)
				mock.ExpectQuery("SELECT \\* FROM albums").WillReturnRows(rows)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))
			},
			expected: &db.Paginated[Album]{
				Items: []Album{
					{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
					{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
				},
				Total: 2,
			},
			err: nil,
		},
		{
			name:   "Get albums by artist",
			artist: "Artist 1",
			mockFunc: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
					AddRow("1", "Album 1", "Artist 1", 9.99)
				mock.ExpectQuery("SELECT \\* FROM albums WHERE artist ILIKE \\$1").WithArgs("Artist 1").WillReturnRows(rows)
				mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums WHERE artist SIMILAR TO \\$1").WithArgs("Artist 1").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			expected: &db.Paginated[Album]{
				Items: []Album{
					{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
				},
				Total: 1,
			},
			err: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.mockFunc()
			result, err := repo.GetAlbums(context.Background(), tc.artist)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, result)
		})
	}
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
