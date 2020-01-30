package main

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var l = logrus.New()

var RootCmd = &cobra.Command{
	Use:   "service",
	Short: "boilerplate service",
}

func init() {
	RootCmd.AddCommand(Cmd)
}

func main() {
	if err := RootCmd.Execute(); err != nil {
		l.WithError(err).Fatal("something goes wrong")
	}
}
