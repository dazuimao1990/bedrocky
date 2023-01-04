/*
Copyright © 2022 dazuimao1990<gx_927@163.com>

*/
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
)

var (
	cmds []string
)

// runCmd represents the run command
// Calling the run command directly passes command slice based on command-line arguments
// and defines a host that can be managed by the ssh protocol.
// Finally, an ordered set of commands is completed for this single host.
// ShellModule moudle provided at the same time.
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run ssh command on the remote host",
	Long:  `Bedrocky utilizes the ssh protocol to perform an ordered set of commands synchronously on multiple remote hosts`,
	RunE: func(cmd *cobra.Command, args []string) error {
		h := NewHostSpec(host, port, password, username)
		s := NewShellModule(h, cmds)
		if err := s.ParallelWork(); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.PersistentFlags().StringSliceVarP(&cmds, "command", "c", []string{}, "You can run an ordered set of commands like -c \"cmd1,cmd2\" -c \"cmd3\"")

}

// Below for ShellModule moudle
// A moudle implement Moudle interface. Execute an ordered set of commands.
type ShellModule struct {
	HostSpec TaskSpec
	Command  []string
}

// New a ShellModule soon.For bedrocky run sub-command.
func NewShellModule(host TaskSpec, cmds []string) ShellModule {
	return ShellModule{
		HostSpec: host,
		Command:  cmds,
	}
}

// Execute an ordered set of commands for a host
// Implementation method for Moudle interface.
func (s ShellModule) ParallelWork() error {
	// Init ssh client config
	sshConfig := s.HostSpec.SshConfigInit()
	// Connect to ssh server
	conn, err := ssh.Dial("tcp", s.HostSpec.Host+":"+strconv.Itoa(s.HostSpec.Port), sshConfig)
	if err != nil {
		log.Fatalf("[ERROR] 无法面向主机 %v 建立 ssh 连接: %v\n", s.HostSpec.Host, err)
	}
	defer conn.Close()
	for _, cmd := range s.Command {
		session, err := conn.NewSession()
		if err != nil {
			log.Fatalf("[ERROR] 无法面向主机 %v 建立 session 会话: %v\n", s.HostSpec.Host, err)
		}
		res, err := session.CombinedOutput(cmd)
		if err != nil {
			log.Printf("[ERROR] 主机 %v 任务 %v 执行失败. 失败原因: %v", s.HostSpec.Host, cmd, string(res))
			continue
		}
		session.Close()
		log.Printf("[SUCCESS] 主机 %v 任务 %v 执行完成. 执行结果: %v", s.HostSpec.Host, cmd, string(res))
	}
	return nil
}
