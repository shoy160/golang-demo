package main

import (
	"flag"

	"github.com/spf13/pflag"
	"shay.cn/m/demo/cmd"
)

func main() {
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Set("logtostderr", "true")
	cmd.Execute()
}
