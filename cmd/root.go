package cmd

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/vmware/govmomi"

	i "github.com/hsarena/vcbox/internal"
)

type vcCred struct {
	address  string
	username string
	password string
	insecure bool
}

var (
	vcc vcCred
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vcbox",
	Short: "VMware vCenter Text User Interface as a Box",
	Run: func(cmd *cobra.Command, args []string) {
		u, err := parseCredentials(&vcc)
		if err != nil {
			log.Printf("%s", err.Error())
		}

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		c, err := govmomi.NewClient(ctx, u, vcc.insecure)
		if err != nil {
			log.Printf("%s", err.Error())
		}

		if err := i.NewUI(c).Run(); err != nil {
			panic(err)
		}

		discovery := i.NewDiscoveryService(c)
		dc, err := discovery.DiscoverDatacenters()
		if err != nil {
			log.Printf("%s", err.Error())
		}

		cr, err := discovery.DiscoverComputeResource(dc[0])
		if err != nil {
			log.Printf("%s", err.Error())
		}

		log.Printf("Compute Resource: %s\n", cr)
		defer c.Logout(ctx)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by mai.main(). It only needs to happen once to the rootCmd.
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
}

func parseCredentials(vcc *vcCred) (*url.URL, error) {

	// Check that an address was actually entered
	if vcc.address == "" {
		return nil, fmt.Errorf("no VMware vcCredenter URL/Address has been submitted")
	}

	// Check that the URL can be parsed
	u, err := url.Parse(vcc.address)
	if err != nil {
		return nil, fmt.Errorf("url can't be parsed, ensure it is https://username:password/<address>/sdk")
	}

	// Check if a username was entered
	if vcc.username == "" {
		// if no username does one exist as part of the url
		if u.User.Username() == "" {
			return nil, fmt.Errorf("no VMware vcCredenter Username has been submitted")
		}
	} else {
		// A username was submitted update the url
		u.User = url.User(vcc.username)
	}

	if vcc.password == "" {
		_, set := u.User.Password()
		if !set {
			return nil, fmt.Errorf("no VMware vcCredenter Password has been submitted")
		}
	} else {
		u.User = url.UserPassword(u.User.Username(), vcc.password)
	}
	return u, nil
}
