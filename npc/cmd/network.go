package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"natproxy/lib/network/lan"
	"os"
)

var (
	all bool
)

func init() {
	rootCmd.AddCommand(adaptorCmd)
	adaptorCmd.Flags().BoolVarP(&all, "all", "a", false, "获取所有启用网卡信息")
}

var adaptorCmd = &cobra.Command{
	Use:   "adaptor",
	Short: "获取本机的所有启用网卡",
	Long:  `获取本机所有的已经启用的网卡信息.`,
	Run: func(cmd *cobra.Command, args []string) {
		if all {
			adaptors, err := lan.Adaptors()
			if err != nil {
				fmt.Println(err.Error())
			}
			data, err := json.Marshal(adaptors)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Printf("%s\n", data)
		} else {
			// 获取本机所在的局域网IP
			ip, err := lan.LanIP()
			if err != nil {
				log.Println(err)
			}
			adaptor, err := lan.AdaptorByIP(ip)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			data, err := json.Marshal(adaptor)
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
			fmt.Printf("%s\n", data)
		}

	},
}
