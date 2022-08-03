package lnx

type Search struct {
	Query   Query  `json:"query"`
	Offset  int    `json:"offset"`
	OrderBy string `json:"order_by"`
	Sort    string `json:"sort"`
}
