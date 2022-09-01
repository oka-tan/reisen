package lnx

type searchResult struct {
	Data searchResultData `json:"data"`
}

type searchResultData struct {
	Hits  []hit `json:"hits"`
	Count int   `json:"count"`
}

type hit struct {
	Doc doc `json:"doc"`
}

type doc struct {
	PostNumber int64 `json:"post_number"`
}

type CondensedSearchResult struct {
	Count int
	Hits  []int64
}
