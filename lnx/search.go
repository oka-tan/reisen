package lnx

type search struct {
	Query   query  `json:"query"`
	Offset  int    `json:"offset"`
	OrderBy string `json:"order_by"`
	Sort    string `json:"sort"`
}
