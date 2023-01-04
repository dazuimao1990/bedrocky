/*
Copyright © 2022 dazuimao1990<gx_927@163.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

// hostsCmd represents the hosts command
var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hosts called")
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd)
}

// Build a new hostSpec object
func NewHostSpec(host string, port int, password string, username string) TaskSpec {
	return TaskSpec{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		Role:     []string{},
	}
}

func (h *TaskSpec) SshConfigInit() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: h.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(h.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func HostSpecCompletion(h *TaskSpec) {
	// 对于 Port Username Publickey Password
	// 优先取命令行参数
	// 其次取配置文件中为每个 host 定义的值
	// 再次则取配置文件中的全局值
	// 最后取命令行参数默认值
	// 如无，则置空
	// 垃圾函数，迟早优化

	// fmt.Println(h.Host, h.Port, h.Password, h.Publickey)
	if port != 22 {
		h.Port = port
	} else if h.Port != 0 {
		// 当命令行没有输入时，且有单独定义的值时，pass
	} else if BedrockyCfg.GlobelPort != 0 {
		h.Port = BedrockyCfg.GlobelPort
	} else {
		h.Port = 22
	}

	if username != "root" {
		h.Username = username
	} else if h.Username != "" {
		// 当命令行没有输入时，且有单独定义的值时，pass
	} else if BedrockyCfg.GlobelUsername != "" {
		h.Username = BedrockyCfg.GlobelUsername
	} else {
		h.Username = "root"
	}

	if password != "" {
		h.Password = password
	} else if h.Password != "" {
		// 当命令行没有输入时，且有单独定义的值时，pass
	} else if BedrockyCfg.GlobelPassword != "" {
		h.Password = BedrockyCfg.GlobelPassword
	} else {
		h.Password = ""
	}

	if publickey != "" {
		h.Publickey = publickey
	} else if h.Publickey != "" {
		// 当命令行没有输入时，且有单独定义的值时，pass
	} else if BedrockyCfg.GlobelPublickey != "" {
		h.Publickey = BedrockyCfg.GlobelPublickey
	} else {
		h.Publickey = "~/.ssh/id_rsa.pub"
	}
}
