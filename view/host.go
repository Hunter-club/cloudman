package view

type HostImportRequest struct {
	Hosts []HostItem `json:"hosts"`
}

type HostItem struct {
	Name      string `json:"name"`
	PrimaryIP string `json:"primary_ip"`
	Zone      string `json:"zone"`
	HostID    string `json:"host_id"`
	Domain    string `json:"domain"`
}
