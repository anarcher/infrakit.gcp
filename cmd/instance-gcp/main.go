package main

import (
	"os"

	log "github.com/Sirupsen/logrus"
	//"github.com/anarcher/infrakit.gcp/plugin/instance"
	"github.com/docker/infrakit/pkg/cli"
	//instance_plugin "github.com/docker/infrakit/pkg/rpc/instance"
	"github.com/spf13/cobra"
	//"strings"
)

func main() {

	var logLevel int
	var name string
	var namespaceTags []string

	cmd := &cobra.Command{}

	cmd.Flags().IntVar(&logLevel, "log", cli.DefaultLogLevel, "Logging level. 0 is least verbose. Max is 5")
	cmd.Flags().StringVar(&name, "name", "instance-gcp", "Plugin name to advertise for discovery")
	cmd.Flags().StringSliceVar(
		&namespaceTags,
		"namespace-tags",
		[]string{},
		"A list of key=value resource tags to namespace all resources created")

	cmd.AddCommand(cli.VersionCommand())

	err := cmd.Execute()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
}
