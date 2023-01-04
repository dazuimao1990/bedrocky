/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "bedrocky",
	Short:   "A brief description of your application",
	Long:    `Bedrocky is a tool designed to ease the deployment of a range of infrastructure software`,
	Version: "0.1",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

var (
	host        string
	port        int
	username    string
	password    string
	publickey   string
	BedrockyCfg *Config
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		log.Fatalln(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&host, "host", "H", "localhost", "Specifies the host to be operated on over ssh")
	rootCmd.PersistentFlags().IntVarP(&port, "port", "p", 22, "Specifies the host port to operate over ssh")
	rootCmd.PersistentFlags().StringVarP(&username, "username", "u", "root", "Specifies the username of the host to be operated on over ssh")
	rootCmd.PersistentFlags().StringVarP(&password, "password", "P", "", "Specifies the password of the host to be operated on over ssh")
	rootCmd.PersistentFlags().StringVarP(&publickey, "key", "i", "~/.ssh/id_rsa.pub", "Specifies the publickey filepath of the host to be operated on over ssh")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("bedrocky")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// 从配置文件中解析出 hosts 列表
	err := viper.Unmarshal(&BedrockyCfg)
	if err != nil {
		log.Fatalln(err)
	}

}
