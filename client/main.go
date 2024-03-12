package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Hunter-club/cloudman/view"
	"github.com/Hunter-club/cloudman/xui"
	"github.com/imroc/req/v3"
	"github.com/spf13/cobra"
	"github.com/tidwall/gjson"
)

var (
	cloudmanURL string = os.Getenv("CLOUDMAN_URL")
)

func main() {
	cmd := &cobra.Command{
		Use: "cloudman",
		Run: func(cmd *cobra.Command, args []string) {
			os.Exit(-1)
		},
	}
	cmd.AddCommand(PreAllocateLineCommand())
	cmd.AddCommand(ConfigXuiCommand())
	cmd.AddCommand(GenSubCommand())

	err := cmd.Execute()
	if err != nil {
		panic(err)
	}
}

// 预分配线路
func PreAllocateLineCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "allocate_line",
		Run: func(cmd *cobra.Command, args []string) {
			client := GetClient()
			resp, err := client.NewRequest().SetBody(&view.AllocateRequest{
				Lines:   linemap,
				OrderID: orderID,
			}).Post("/api/v1/line")

			if err != nil {
				panic(err)
			}

			fmt.Println(string(resp.Bytes()))
		},
	}
	cmd.Flags().Var(&linemap, "linemap", "use it by linemap=us:4,linemap=eng:3")
	cmd.Flags().StringVar(&orderID, "order_id", "test", "use it by order_id=`test`")
	cmd.Flags().Parse(nil)
	return cmd
}

type LineMap map[string]int

var linemap LineMap = make(map[string]int)

func (l LineMap) String() string {
	var res string
	flag := true
	for k, v := range l {
		if flag {
			res += fmt.Sprintf("%s:%d", k, v)
			flag = false
		} else {
			res += fmt.Sprintf(" %s:%d", k, v)
		}
	}
	return res
}

func (l LineMap) Set(value string) error {
	res := strings.Split(value, ":")
	if len(res) != 2 {
		return errors.New("unexcepeted linemap params")
	}
	count, err := strconv.Atoi(res[1])
	if err != nil {
		return err
	}
	l[res[0]] = count
	return nil
}

func (l LineMap) Type() string {
	return ""
}

var orderID string

// 配置线路
func ConfigXuiCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "config_xui",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(orderID)

			client := GetClient()
			resp, err := client.NewRequest().SetBody(&view.XrayRequest{
				OrderID: orderID,
			}).Post("/api/v1/xray")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			inbounds := make([]xui.Inbound, 0)
			err = json.Unmarshal([]byte(gjson.GetBytes(resp.Bytes(), "data").String()), &inbounds)
			if err != nil {
				fmt.Println(string(resp.Bytes()))
				fmt.Println(err.Error())
				return
			}

			results := view.SubRequest{
				OrderID: orderID,
				Entries: make([]view.SubConfigTransferEntry, 0),
			}

			for _, inbound := range inbounds {
				subID := xui.GetInboundSubId(&inbound)
				addr := inbound.Remark[strings.Index(inbound.Remark, "-")+1:]
				fmt.Printf("addr: %s, subid: %s\n", addr, subID)
				results.Entries = append(results.Entries, view.SubConfigTransferEntry{
					Transfer: view.Transfer{
						Addr: "127.0.0.1",
						Port: 7890,
					},
					TargetHost: view.TargetHost{
						Addr:  addr,
						SubID: subID,
					},
				})
			}

			data, err := json.MarshalIndent(results, "", " ")

			if err != nil {
				panic(err)
			}

			err = os.WriteFile("subrequest.json", data, os.ModePerm)

			if err != nil {
				panic(err)
			}

		},
	}

	cmd.Flags().Var(&linemap, "linemap", "use it by linemap=us:4,linemap=eng:3")
	cmd.Flags().StringVar(&orderID, "order_id", "test", "use it by order_id=`test`")
	cmd.Flags().Parse(nil)

	return cmd
}

func GenSubCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use: "gen_sub",
		Run: func(cmd *cobra.Command, args []string) {
			data, err := os.ReadFile("subrequest.json")
			if err != nil {
				panic(err)
			}
			req := &view.SubRequest{}
			err = json.Unmarshal(data, req)
			if err != nil {
				fmt.Println(err.Error())
			}
			client := GetClient()

			resp, err := client.NewRequest().SetBody(req).Post("/api/v1/sub")
			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(string(resp.Bytes()))
		},
	}
	return cmd

}

func GetClient() *req.Client {
	return req.C().SetBaseURL(cloudmanURL).SetCommonHeader("secret", "FXf4nzFzax8A.k-a")
}
