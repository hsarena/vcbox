package vmware

import (
	"fmt"
	"net/url"
)

func ParseCredentials(vcc *VCClient) (*url.URL, error) {

	// Check that an address was actually entered
	if vcc.Address == "" {
		return nil, fmt.Errorf("no VMware vcCredenter URL/Address has been submitted")
	}

	// Check that the URL can be parsed
	u, err := url.Parse(vcc.Address)
	if err != nil {
		return nil, fmt.Errorf("url can't be parsed, ensure it is https://username:password/<address>/sdk")
	}

	// Check if a username was entered
	if vcc.Username == "" {
		// if no username does one exist as part of the url
		if u.User.Username() == "" {
			return nil, fmt.Errorf("no VMware vcCredenter Username has been submitted")
		}
	} else {
		// A username was submitted update the url
		u.User = url.User(vcc.Username)
	}

	if vcc.Password == "" {
		_, set := u.User.Password()
		if !set {
			return nil, fmt.Errorf("no VMware vcCredenter Password has been submitted")
		}
	} else {
		u.User = url.UserPassword(u.User.Username(), vcc.Password)
	}
	return u, nil
}
