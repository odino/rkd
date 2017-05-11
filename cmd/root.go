package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "rkd",
	Short: "rkd is docker-compose for rkt containers",
	Long:  `Development environments powered by rkt containers with ease.`,
}
