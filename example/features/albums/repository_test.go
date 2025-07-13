package albums

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/jphilipstevens/web-service-gin/app/db"
	"github.com/jphilipstevens/web-service-gin/config"
	"github.com/jphilipstevens/web-service-gin/testUtils"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewAlbumRepository(t *testing.T) {
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))
	assert.NotNil(t, repo)
}

func TestGetAlbumsRepository(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})

	t.Run("Get all albums", func(t *testing.T) {
		mockDB, mock, _ := sqlmock.New()
		defer mockDB.Close()

		repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))
		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", 9.99).
			AddRow("2", "Album 2", "Artist 2", 14.99)
		mock.ExpectQuery("SELECT \\* FROM albums").WillReturnRows(rows)
		mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		result, err := repo.GetAlbums(testUtils.CreateTestContext(), GetAlbumsParams{
			Artist: "",
			Limit:  10,
			Page:   0,
		})

		expected := &db.Paginated[Album]{
			Items: []Album{
				{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
				{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
			},
		}

		assert.NoError(t, err)
		assert.Equal(t, 2, len(result.Items))
		assert.Equal(t, reflect.DeepEqual(expected.Items, result.Items), true)

	})

	t.Run("Get all albums with artist", func(t *testing.T) {
		mockDB, mock, _ := sqlmock.New()
		defer mockDB.Close()

		repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

		rows := sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", 9.99)

		mock.ExpectQuery("SELECT \\* FROM albums WHERE artist ILIKE \\$1 LIMIT \\$2 OFFSET \\$3").WithArgs("Artist 1", 10, 0).WillReturnRows(rows)
		result, err := repo.GetAlbums(testUtils.CreateTestContext(), GetAlbumsParams{
			Artist: "Artist 1",
			Limit:  10,
			Page:   0,
		})
		expected := &db.Paginated[Album]{
			Items: []Album{
				{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
			},
		}
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, reflect.DeepEqual(expected.Items, result.Items), true)
	})

}

func TestInsert(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	album := Album{ID: "1", Title: "New Album", Artist: "New Artist", Price: 19.99}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs(album.ID, album.Title, album.Artist, album.Price).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.Insert(testUtils.CreateTestContext(), album)
	assert.NoError(t, err)
}

func TestInsertBatch(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	//	config.Init(config.ConfigOptions{Path: "example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	albums := []Album{
		{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
		{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
	}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs("1", "Album 1", "Artist 1", 9.99, "2", "Album 2", "Artist 2", 14.99).
		WillReturnResult(sqlmock.NewResult(2, 2))

	err := repo.InsertBatch(testUtils.CreateTestContext(), albums)
	assert.NoError(t, err)
}

func TestGetAlbumsNoResults(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	mock.ExpectQuery("SELECT \\* FROM albums").WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}))
	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM albums").WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	result, err := repo.GetAlbums(testUtils.CreateTestContext(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, db.NotFoundError, err)
}

func TestGetAlbumsError(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	mock.ExpectQuery("SELECT \\* FROM albums").WillReturnError(sql.ErrConnDone)

	result, err := repo.GetAlbums(testUtils.CreateTestContext(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestInsertError(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	album := Album{ID: "1", Title: "New Album", Artist: "New Artist", Price: 19.99}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs(album.ID, album.Title, album.Artist, album.Price).
		WillReturnError(sql.ErrConnDone)

	err := repo.Insert(testUtils.CreateTestContext(), album)
	assert.Error(t, err)
}

func TestInsertBatchEmptySlice(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, _, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	err := repo.InsertBatch(testUtils.CreateTestContext(), []Album{})
	assert.NoError(t, err)
}

func TestInsertBatchError(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))

	albums := []Album{
		{ID: "1", Title: "Album 1", Artist: "Artist 1", Price: 9.99},
		{ID: "2", Title: "Album 2", Artist: "Artist 2", Price: 14.99},
	}

	mock.ExpectExec("INSERT INTO albums").
		WithArgs("1", "Album 1", "Artist 1", 9.99, "2", "Album 2", "Artist 2", 14.99).
		WillReturnError(sql.ErrConnDone)

	err := repo.InsertBatch(testUtils.CreateTestContext(), albums)
	assert.Error(t, err)
}

func TestGetAlbumsScanError(t *testing.T) {
	config.Init(config.ConfigOptions{Path: "../../example/config", Name: "config", Type: "yaml"})
	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	repo := NewAlbumRepository(testUtils.NewDatabase(mockDB))
	mock.ExpectQuery("SELECT \\* FROM albums").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "artist", "price"}).
			AddRow("1", "Album 1", "Artist 1", "invalid_price"))

	result, err := repo.GetAlbums(testUtils.CreateTestContext(), GetAlbumsParams{
		Artist: "",
		Limit:  10,
		Page:   0,
	})
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, db.DatabaseError, err)
}
