package view

type AllocateRequest struct {
	Lines   map[string]int `json:"line_spec"`
	OrderID string         `json:"order_id"`
}
