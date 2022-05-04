/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/

package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os/exec"
	"runtime"
	"strings"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止客服http服务",
	Run: func(cmd *cobra.Command, args []string) {
		pids, err := ioutil.ReadFile("MgHash.sock")
		if err != nil {
			return
		}
		pidSlice := strings.Split(string(pids), ",")
		var command *exec.Cmd
		for _, pid := range pidSlice {
			if runtime.GOOS == "windows" {
				command = exec.Command("taskkill.exe", "/f", "/pid", pid)
			} else {
				fmt.Printf("关闭pid %s", pid)
				command = exec.Command("kill", pid)
			}
			command.Start()
		}
	},
}
