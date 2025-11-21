package cmd

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{ //定义根命令（应用），默认情况下根命令不执行任何操作，仅作为其他子命令的父级，因此不需要属性Run
	Use:   "pScan",
	Short: "Fast TCP port scanner",
	Long: `pScan - short for Port Scanner - executes TCP port scan on a list of hosts.

pScan allows you to add, list, and delete hosts from the list.

pScan executes a port scan on specified TCP ports. You can customize the target ports using a command line flag.`,
	Version: "v0.1.0", //定义应用程序版本，可使用命令行标志-v、--version查看
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() { //在main()函数之前运行，通常用来添加命令行标志
	//设置打印版本信息时的模板
	versionTemplate := `{{printf "%s: %s - version %s\n" .Name .Short .Version}}` //Go模板字符串，用反引号`包裹，支持多行，且模板语法{{ }}会被cobra解析，然后使用printf格式化输出
	rootCmd.SetVersionTemplate(versionTemplate)                                   //将自定义的版本模板绑定到根命令
}
