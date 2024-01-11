package ssh

import (
	"fmt"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

func SSHConnect(user string, ip string, port int, ignoreHostKey bool) {

	// Connect to the SSH agent
	sshAgent, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
	if err != nil {
		fmt.Printf("Failed to connect to SSH agent: %v\n", err)
		return
	}
	defer sshAgent.Close()

	// Create an SSH client configuration
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			// Use the signer from the private key
			ssh.PublicKeysCallback(agent.NewClient(sshAgent).Signers),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Use only in a safe environment
	}

	// config := &ssh.ClientConfig{
	// 	User: "root",
	// 	Auth: []ssh.AuthMethod{
	// 		ssh.PublicKeys(signer),
	// 	},
	// 	HostKeyCallback: hostKeyCallback,
	// }
	addr := fmt.Sprintf("%s:%v", ip, port)
	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Fatal("Failed to dial: ", err)
	}

	session, err := client.NewSession()
	if err != nil {
		log.Fatal("Failed to create session: ", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if err := session.RequestPty("linux", 80, 40, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	//set input and output
	session.Stdout = os.Stdout
	session.Stdin = os.Stdin
	session.Stderr = os.Stderr

	if err := session.Shell(); err != nil {
		log.Fatal("failed to start shell: ", err)
	}

	err = session.Wait()
	if err != nil {
		log.Fatal("Failed to run: " + err.Error())
	}
}
