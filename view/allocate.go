package view

// {
//
//	"line_map": {
//		"localhost":1
//	},
//
// "order_id":"test"
// }
type AllocateRequest struct {
	Lines   map[string]int `json:"line_map"`
	OrderID string         `json:"order_id"`
}
