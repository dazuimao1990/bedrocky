package cmd

import (
	"encoding/base64"
	"fmt"
)

type Module interface {
	// 基于本 moudle 结构体对象特征而实现的并行任务的具体方法
	ParallelWork() error
	// 从配置文件中解析出 RoleSpec.Spec 生成本 moudle 结构体对象实例
	// ParseRoleCfg(moudle string)
}

var TaskChan chan TaskSpec

type TaskSpec struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Publickey string
	Role      []string
}

type Config struct {
	Hosts           []TaskSpec `mapstructure:"hosts"`
	GlobelPort      int        `mapstructure:"port"`
	GlobelUsername  string     `mapstructure:"username"`
	GlobelPassword  string     `mapstructure:"password"`
	GlobelPublickey string     `mapstructure:"publickey"`
}

// 解析并补全每一个 host 元素,基于 host 创建任务管道
func (c *Config) ParseHostToChan() {
	TaskChan = make(chan TaskSpec, len(c.Hosts))
	for _, host := range c.Hosts {
		HostSpecCompletion(&host)
		TaskChan <- host
	}
	close(TaskChan)
}

func (c *Config) PrintHostsInfo() {
	for k, host := range c.Hosts {
		HostSpecCompletion(&host)
		fmt.Printf("%d\t%v\t%v\t%v\t%v\t\t%v\n",
			k+1, host.Host, host.Port, host.Username, base64.StdEncoding.EncodeToString([]byte(host.Password)), host.Publickey)
	}

}
