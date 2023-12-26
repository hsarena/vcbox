package cmd

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vmware/govmomi"

	"github.com/hsarena/vcbox/pkg/vcbox"
)

type vcCred struct {
	address  string
	username string
	password string
	insecure bool
}

var (
	dcName string
	vcc    vcCred
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vcbox",
	Short: "VMware vCenter Text User Interface as a Box",
	Run: func(cmd *cobra.Command, args []string) {
		u, err := parseCredentials(&vcc)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c, err := govmomi.NewClient(ctx, u, vcc.insecure)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}

		defer c.Logout(ctx)
		dcs, err := vcbox.DCDiscovery(c)
		if err != nil {
			fmt.Printf("%s", err.Error())
		}
		vms, err := vcbox.VMDiscovery(c, "AFR")
		vcbox.NewUi(dcs,vms)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
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

	rootCmd.Flags().StringVar(&vcc.address, "address", os.Getenv("VCURL"), "The Url/address of a VMware vCenter server")
	rootCmd.Flags().StringVar(&vcc.username, "user", os.Getenv("VCUSER"), "The Url/address of a VMware vCenter server")
	rootCmd.Flags().StringVar(&vcc.password, "password", os.Getenv("VCPASS"), "The Url/address of a VMware vCenter server")
	rootCmd.Flags().BoolVar(&vcc.insecure, "insecure", boolVar, "The Url/address of a VMware vCenter server")
	rootCmd.Flags().StringVar(&dcName, "datacenter", os.Getenv("VC_DATACENTER"), "The Datacenter Name")
}

func parseCredentials(vcc *vcCred) (*url.URL, error) {

	// Check that an address was actually entered
	if vcc.address == "" {
		return nil, fmt.Errorf("No VMware vcCredenter URL/Address has been submitted")
	}

	// Check that the URL can be parsed
	u, err := url.Parse(vcc.address)
	if err != nil {
		return nil, fmt.Errorf("URL can't be parsed, ensure it is https://username:password/<address>/sdk")
	}

	// Check if a username was entered
	if vcc.username == "" {
		// if no username does one exist as part of the url
		if u.User.Username() == "" {
			return nil, fmt.Errorf("No VMware vcCredenter Username has been submitted")
		}
	} else {
		// A username was submitted update the url
		u.User = url.User(vcc.username)
	}

	if vcc.password == "" {
		_, set := u.User.Password()
		if !set {
			return nil, fmt.Errorf("No VMware vcCredenter Password has been submitted")
		}
	} else {
		u.User = url.UserPassword(u.User.Username(), vcc.password)
	}
	return u, nil
}
