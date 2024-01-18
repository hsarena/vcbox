package cmd

import (
	"context"
	"log"
	"os"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vmware/govmomi"

	"github.com/hsarena/vcbox/pkg/ssh"
	"github.com/hsarena/vcbox/pkg/tui"
	"github.com/hsarena/vcbox/pkg/vmware"
)

var (
	vcc vmware.VCClient
	sc  ssh.SSHClient
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vcbox",
	Short: "VMware vCenter Text User Interface as a Box",
	Run: func(cmd *cobra.Command, args []string) {
		u, err := vmware.ParseCredentials(&vcc)
		if err != nil {
			log.Printf("%s", err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c, err := govmomi.NewClient(ctx, u, vcc.Insecure)
		if err != nil {
			log.Printf("%s", err.Error())
		}

		// if err := ui.NewUI(c).Run(); err != nil {
		// 	panic(err)
		// }

		if _, err := tea.NewProgram(tui.InitialModel(c), tea.WithAltScreen()).Run(); err != nil {
			log.Println("Error running program:", err)
			os.Exit(1)
		}

		defer c.Logout(ctx)
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	defaultValue := true
	boolVar, err := strconv.ParseBool(os.Getenv("VC_INSECURE"))
	if err != nil {
		// Handle the error, or use the default value in case of an error
		boolVar = defaultValue
	}

	rootCmd.Flags().StringVar(&vcc.Address, "address", os.Getenv("VCURL"), "The Url/address of a VMware vCenter server")
	rootCmd.Flags().StringVar(&vcc.Username, "user", os.Getenv("VCUSER"), "The Username VMware vCenter server")
	rootCmd.Flags().StringVar(&vcc.Password, "password", os.Getenv("VCPASS"), "The Password of VMware vCenter server")
	rootCmd.Flags().BoolVar(&vcc.Insecure, "insecure", boolVar, "Ignoring the secure connection")
	rootCmd.Flags().StringVar(&sc.Username, "remote-user", os.Getenv("VCBOX_REMOTE_USER"), "The defualt remote user")
}
