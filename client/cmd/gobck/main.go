package main

import (
	"github.com/devfacet/gocmd/v3"

	"github.com/manso/gobackup/client/commands"
)

func main() {
	flags := struct {
		Help   bool `short:"h" long:"help" description:"Display usage" global:"true"`
		Backup struct {
			Path string `short:"p" long:"path" required:"true" description:"Root directory to be backed up"`
		} `command:"backup" description:"Perform backup"`
		Restore struct {
			Path string `short:"p" long:"path" required:"true" description:"Root directory where files will be restored"`
		} `command:"restore" description:"Perform restoration"`
	}{}

	gocmd.HandleFlag("Backup", func(cmd *gocmd.Cmd, args []string) error {
		commands.PerformBackup(flags.Backup.Path)
		return nil
	})

	gocmd.HandleFlag("Restore", func(cmd *gocmd.Cmd, args []string) error {
		commands.PerformRestore(flags.Restore.Path)
		return nil
	})

	gocmd.New(gocmd.Options{
		Name:        "gobck",
		Description: "Backup tool",
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
	})
}
