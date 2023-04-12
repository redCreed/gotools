package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {
	//cmd := cobra.Command{
	//	Use:   "test",
	//	Short: "short test",
	//}
	//var com Com
	//cmd.PersistentFlags().StringVarP(&com.Test, "h", "", "test", "test mode type")
	//
	//cmd.Execute()
	//fmt.Println(com.Test)

	cobra.EnableCommandSorting = false

	rootCmd := &cobra.Command{
		Use:   "tools",
		Short: "Storj 工具包",
	}

	rootCmd.AddCommand(AdminCmd)
	rootCmd.SilenceUsage = true
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}

	flag := AdminCmd.Flags().Lookup("host")
	if flag == nil {
		fmt.Println("err nil")
	}

	fmt.Println("ttt:", flag.Value.String(), flag.DefValue, flag.Usage, flag.Name)
}
