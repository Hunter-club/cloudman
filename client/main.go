package main

import (
	"os"

	"github.com/Hunter-club/cloudman/view"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
)

var (
	cloudmanURL string = os.Getenv("CLOUDMAN_URL")
)

func main() {
	cmd := &cobra.Command{
		Use: "cloudman",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.AddCommand(PreAllocateLineCommand())
		},
	}
}

// e.POST("/api/v1/sub", Handler(handler.AllocateResource))
// e.POST("/api/v1/xray", Handler(handler.XUIConfigure))
// e.POST("/api/v1/line", Handler(handler.PreAllocateLine))

var ()

func PreAllocateLineCommand() *cobra.Command {
	return &cobra.Command{
		Use: "pre_allocate",
		Run: func(cmd *cobra.Command, args []string) {
			client := GetClient()

			client.NewRequest().SetBody(&view.AllocateRequest{
				Lines:   map[string]int{},
				OrderID: "",
			}).Post("")
		},
	}
}

func GetClient() *req.Client {
	return req.C().SetBaseURL(cloudmanURL).SetCommonHeader("secret", "FXf4nzFzax8A.k-a")
}
