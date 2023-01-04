/*
Copyright © 2022 dazuimao1990<gx_927@163.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("序号\t主机\t\t端口\t用户\t密码\t\t\t公钥\n")
		BedrockyCfg.PrintHostsInfo()
	},
}

func init() {
	hostsCmd.AddCommand(showCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// showCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// showCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	showCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "./bedrocky.yaml", "Config file (default is ./bedrocky.yaml)")
}
