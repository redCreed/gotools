package main

import (
	"github.com/spf13/cobra"
)

var AdminCmd *cobra.Command

type AdminFlags struct {
	Host          string
	Username      string
	Authorization string
}

func init() {
	var flags AdminFlags
	AdminCmd = &cobra.Command{
		Use:   "admin",
		Short: "Storj 卫星admin 操作",
	}
	AdminCmd.PersistentFlags().StringVarP(&flags.Host, "host", "", "http://36.138.1.13:10005", "admin地址")
	AdminCmd.PersistentFlags().StringVarP(&flags.Username, "username", "", "", "新创建的用户名")
	AdminCmd.PersistentFlags().StringVarP(&flags.Authorization, "authorization", "", "abcd", "授权码")
	AdminCmd.AddCommand(
		&cobra.Command{
			Use:   "create-user",
			Short: "创建用户",
			//RunE:  flags.CreateUser,
			Run: func(cmd *cobra.Command, args []string) {
				cmd.AddCommand(&cobra.Command{
					Use:   "get-user",
					Short: "创建用户",
					//RunE:  flags.GetUser,
				},
				)
			},
		},

	)
}
