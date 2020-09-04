package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "dd",
}

// Execute add
func Execute() {
	fmt.Println("cmd.Execute")
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// rootCmd.PersistentFlags().BoolVar(&opt.Verbose, "debug", true, "Print detailed information")
}
