package lnx

type SearchResult struct {
	Data SearchResultData `json:"data"`
}

type SearchResultData struct {
	Hits  []Hit `json:"hits"`
	Count int   `json:"count"`
}

type Hit struct {
	Doc Doc `json:"doc"`
}

type Doc struct {
	PostNumber int64 `json:"post_number"`
}

type CondensedSearchResult struct {
	Count int
	Hits  []int64
}
