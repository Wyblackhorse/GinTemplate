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
	"github.com/wangyi/GinTemplate/common"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	"github.com/wangyi/GinTemplate/logger"
	"github.com/wangyi/GinTemplate/router"
	"github.com/wangyi/GinTemplate/setting"
	"github.com/wangyi/GinTemplate/tools"
	"github.com/zh-five/xdaemon"
	"go.uber.org/zap"
	"io/ioutil"
	"log"
	"os"
)

var (
	port     string
	daemon   bool
	rootPath string
)
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "启动MgHash服务",
	Example: "go-fly server",
	Run:     run,
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&rootPath, "rootPath", "r", "", "程序根目录")
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", "7777", "监听端口号")
	serverCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "是否为守护进程模式")
}

func run(cmd *cobra.Command, args []string) {
	//初始化目录

	initDir()
	//初始化守护进程
	initDaemon()
	if noExist, _ := tools.IsFileNotExist(common.LogDirPath); noExist {
		if err := os.MkdirAll(common.LogDirPath, 0777); err != nil {
			log.Println(err.Error())
		}
	}
	isMainUploadExist, _ := tools.IsFileExist(common.UploadDirPath)
	if !isMainUploadExist {
		os.Mkdir(common.UploadDirPath, os.ModePerm)
	}

	//加载配置
	if err := setting.Init(); err != nil {
		fmt.Println("配置文件初始化事变", err)
		return
	}
	//初始化日志
	if err := logger.Init(); err != nil {
		fmt.Println("日志初始化失败", err)
		return
	}
	defer zap.L().Sync() //缓存日志追加到日志文件中

	//链接数据库
	if err := mysql.Init(); err != nil {
		fmt.Println("mysql 链接失败,", err)
		return
	}
	defer mysql.Close()
	//redis 初始化
	//4.初始化redis连接
	if err := redis.Init(); err != nil {
		fmt.Println("redis文件初始化失败：", err)
		return
	}
	defer redis.Close()

	router.Setup()
}

//初始化目录
func initDir() {

	if rootPath == "" {
		rootPath = tools.GetRootPath()
	}
	log.Println("程序运行路径:" + rootPath)
	common.RootPath = rootPath
	common.LogDirPath = rootPath + "/logs/"
	common.ConfigDirPath = rootPath + "/config/"
	common.StaticDirPath = rootPath + "/static/"
	common.UploadDirPath = rootPath + "/static/upload/"

}

//初始化守护进程
func initDaemon() {
	if daemon == true {
		d := xdaemon.NewDaemon(common.LogDirPath + "MgHash.log")
		d.MaxError = 10
		d.Run()
	}
	//记录pid
	ioutil.WriteFile(common.RootPath+"/MgHash.sock", []byte(fmt.Sprintf("%d,%d", os.Getppid(), os.Getpid())), 0666)
}
