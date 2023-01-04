/*
Copyright © 2022 dazuimao1990<gx_927@163.com>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var exit chan int

// playCmd represents the play command
var playCmd = &cobra.Command{
	Use:   "play -f ./bedrocky.yaml",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("play called with cfgFile", cfgFile)
		// 解析并补全每一个 host 元素,基于 host 创建任务管道
		BedrockyCfg.ParseHostToChan()
		// 1. goroutine 分发任务管道中的任务
		exit = make(chan int, len(TaskChan))
		for i := 0; i < cap(TaskChan); i++ {
			go AllotTask(TaskChan, exit)
		}
		for i := 0; i < cap(TaskChan); i++ {
			<-exit
		}

	},
}

func init() {
	rootCmd.AddCommand(playCmd)
	playCmd.PersistentFlags().StringVarP(&cfgFile, "config", "f", "./bedrocky.yaml", "Config file (default is ./bedrocky.yaml)")
}

// 从 TaskChan 管道中获取任务详情，基于 role 的切片内容，从配置文件中获取任务详情
// 并根据模块名生成对应的模块结构体实例，用多态函数按顺序执行
func AllotTask(c chan TaskSpec, exit chan int) error {
	//var roleSpec []interface{}
	task := <-c
	for _, role := range task.Role {
		//roleSpec = .([]interface{})
		for _, steps := range viper.Get(role).([]interface{}) {
			//fmt.Printf("steps v=%v t=%T\n", steps, steps)
			steps, ok := steps.(map[string]interface{})
			if !ok {
				err := fmt.Errorf("steps 类型解析失败，检查配置文件")
				return err
			}
			task.DetectModuleToRun(steps)
		}

	}
	exit <- 1
	return nil
}

// 根据输入的模块名，生成对应的模块结构体实例
func (t TaskSpec) DetectModuleToRun(steps map[string]interface{}) error {
	switch steps["module"] {
	case "shell":
		cmds := []string{}
		shellSpec, ok := steps["spec"].([]interface{})
		if !ok {
			err := fmt.Errorf("%v shell.spec 解析失败，检查配置文件", steps["name"])
			return err
		}
		for _, cmd := range shellSpec {
			cmds = append(cmds, cmd.(string))
		}
		StartTask(NewShellModule(t, cmds))
	case "copy":

	case "systemd":

	default:
		err := fmt.Errorf("%v 包含未知的模块 %v", steps["name"], steps["module"])
		return err

	}
	return nil
}

// 根据传入的模块结构体实例，来启动多态方法
func StartTask(m Module) {
	m.ParallelWork()
}
