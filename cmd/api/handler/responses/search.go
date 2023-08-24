package responses

type Searchable interface {
	Block | Address | Namespace | Tx
}

type SearchResponse[T Searchable] struct {
	Result T `json:"result"`
}

func NewSearchResponse[T Searchable](val T) SearchResponse[T] {
	return SearchResponse[T]{
		Result: val,
	}
}
