package albums

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAlbumModel(t *testing.T) {
	expectedJSON := `{"id":"1","title":"Blue Train","artist":"John Coltrane","price":56.99}`
	album := Album{
		ID:     "1",
		Title:  "Blue Train",
		Artist: "John Coltrane",
		Price:  56.99,
	}

	t.Run("Successful Album Model Marshalling for JSON", func(t *testing.T) {

		jsonData, err := json.Marshal(album)
		assert.NoError(t, err)

		assert.JSONEq(t, expectedJSON, string(jsonData))

	})

	t.Run("Successful Album Model Unmarshalling for JSON", func(t *testing.T) {

		marshalledData := []byte(expectedJSON)
		var unmarshaledAlbum Album
		err := json.Unmarshal(marshalledData, &unmarshaledAlbum)
		assert.NoError(t, err)
		assert.Equal(t, album, unmarshaledAlbum)
	})
}
