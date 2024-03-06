package view

type SubRequest struct {
	OrderID string                   `json:"order_id"`
	Model   string                   `json:"model"`
	Entries []SubConfigTransferEntry `json:"entries"`
}

type SubConfigTransferEntry struct {
	IP       string   `json:"ip"`
	Transfer Transfer `json:"transfer"`
	SubID    string   `json:"sub_id"`
}

type Transfer struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}
