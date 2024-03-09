package view

type SubRequest struct {
	OrderID string                   `json:"order_id"`
	Model   string                   `json:"model"`
	Entries []SubConfigTransferEntry `json:"entries"`
}

type SubConfigTransferEntry struct {
	Transfer   Transfer   `json:"transfer"`
	TargetHost TargetHost `json:"target_host"`
}

type TargetHost struct {
	Addr  string `json:"addr"`
	SubID string `json:"sub_id"`
}

type Transfer struct {
	Addr string `json:"addr"`
	Port int    `json:"port"`
}
