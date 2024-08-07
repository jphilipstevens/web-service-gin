package db

type Paginated[T any] struct {
	Items []T `json:"items"`
}
