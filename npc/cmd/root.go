package cmd

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"natproxy/npc/imp"
	"os"
)

var (
	server  string
	vkey    string
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "npc",
	Short: "NPC",
	Long:  `NPC (NAT Proxy Client，NAT 代理穿透工具). `,
	Run: func(cmd *cobra.Command, args []string) {
		if len(server) == 0 {
			cmd.Help()
			return
		}

		imp.Show(server, vkey)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&server, "server", "s", "", "公网服务器地址(ip:port)")
	rootCmd.Flags().StringVarP(&vkey, "key", "k", "", "认证秘钥字符串")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		viper.AddConfigPath(home)
		viper.SetConfigName(".npc")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
