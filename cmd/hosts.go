package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var hostsCmd = &cobra.Command{
	Use:   "hosts",
	Short: "Manage the hosts list",
	Long: `Manages the hosts lists for pScan

Add hosts with the add command
Delete hosts with the delete command
List hosts with the list command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hosts called")
	},
}

func init() {
	rootCmd.AddCommand(hostsCmd) //将hosts子命令（hostsCmd）附加到根命令
}
