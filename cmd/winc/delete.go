package main

import (
	"code.cloudfoundry.org/winc/container"
	"code.cloudfoundry.org/winc/container/hcsprocess"
	"code.cloudfoundry.org/winc/container/mount"
	"code.cloudfoundry.org/winc/hcs"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var deleteCommand = cli.Command{
	Name:  "delete",
	Usage: "delete a container and the resources it holds",
	ArgsUsage: `<container-id>

Where "<container-id>" is the name for the instance of the container.

EXAMPLE:
For example, if the container id is "windows01" and winc list currently shows the
status of "windows01" as "stopped" the following will delete resources held for
"windows01" removing "windows01" from the winc list of containers:

       # winc delete windows01`,
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name:  "force, f",
			Usage: "Do not return an error if <container-id> does not exist",
		},
	},

	Action: func(context *cli.Context) error {
		if err := checkArgs(context, 1, exactArgs); err != nil {
			return err
		}

		rootDir := context.GlobalString("root")
		containerId := context.Args().First()
		force := context.Bool("force")

		logger := logrus.WithFields(logrus.Fields{
			"containerId": containerId,
		})
		logger.Debug("deleting container")

		client := hcs.Client{}
		cm := container.NewManager(logger, &client, &mount.Mounter{}, &hcsprocess.Process{}, containerId, rootDir)
		return cm.Delete(force)
	},
}
