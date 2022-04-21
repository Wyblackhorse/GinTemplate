/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package cmd

import (
	"github.com/spf13/cobra"
	"log"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "MgHash",
	Short: "MgHash",
	Long:  `MgHash 启动程序`,
}

func init() {
	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(stopCmd)
}
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Println("执行命令参数错误:", err)
		os.Exit(1)
	}
}