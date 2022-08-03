package lnx

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Service struct {
	host   string
	client http.Client
}

func NewService(host string, port int) Service {
	return Service{
		host: fmt.Sprintf("%s:%d", host, port),
	}
}

func (s *Service) Search(board string, ctx string, offset int) (CondensedSearchResult, error) {
	search := Search{
		Query: Query{
			Normal: NormalQuery{Ctx: ctx},
		},
		Offset:  offset,
		OrderBy: "post_number",
		Sort:    "desc",
	}

	pipeReader, pipeWriter := io.Pipe()

	go func() {
		encoder := json.NewEncoder(pipeWriter)
		err := encoder.Encode(&search)
		pipeWriter.CloseWithError(err)
	}()

	resp, err := s.client.Post(fmt.Sprintf("%s/indexes/post_%s/search", s.host, board), "application/json", pipeReader)

	if err != nil {
		return CondensedSearchResult{}, err
	}

	if resp.StatusCode != 200 {
		responseBody, _ := io.ReadAll(resp.Body)
		log.Printf("Received error from Lnx: Status - %s, Body - %s", resp.Status, responseBody)
		resp.Body.Close()

		return CondensedSearchResult{}, err
	}

	var result SearchResult
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	resp.Body.Close()

	if err != nil {
		return CondensedSearchResult{}, err
	}

	hits := make([]int64, 0, len(result.Data.Hits))
	for _, hit := range result.Data.Hits {
		hits = append(hits, hit.Doc.PostNumber)
	}

	condensedSearchResult := CondensedSearchResult{
		Count: result.Data.Count,
		Hits:  hits,
	}

	return condensedSearchResult, nil
}
