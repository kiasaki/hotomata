package hotomata

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"golang.org/x/crypto/ssh"
)

type Machine struct {
	Hostname  string
	Port      int
	SSHConfig *ssh.ClientConfig
}

func MachinesFromInventoryMachines(inventoryMachines []InventoryMachine) []*Machine {
	var rawMachines = []map[string]string{}

	// convert machines to map[string]string
	for _, inventoryMachine := range inventoryMachines {
		var inventoryMachineVars = inventoryMachine.Vars()
		var rawMachine = map[string]string{}
		var err error

		for _, param := range []string{"name", "ssh_hostname", "ssh_port", "ssh_key", "ssh_password"} {
			if v, ok := inventoryMachineVars[param]; ok {
				var value string
				var intValue int
				if err = json.Unmarshal(v, &value); err == nil {
					rawMachine[param] = value
				} else if err = json.Unmarshal(v, &intValue); err == nil {
					rawMachine[param] = strconv.Itoa(intValue)
				}
			}
		}

		rawMachines = append(rawMachines, rawMachine)
	}

	return Machines(rawMachines)
}

func Machines(hosts []map[string]string) []*Machine {
	var machines []*Machine

	for _, host := range hosts {
		hostname := host["name"]
		if h, ok := host["ssh_hostname"]; ok {
			hostname = h
		}

		port := 22
		if p, ok := host["ssh_port"]; ok {
			var err error
			port, err = strconv.Atoi(p)
			if err != nil {
				fmt.Printf("Error parsing port for host [%s]\n", hostname)
				panic(err)
			}
		}

		username := "root"
		if u, ok := host["ssh_username"]; ok {
			username = u
		}

		sshAuthMethods := []ssh.AuthMethod{}
		// Password
		if password, ok := host["ssh_password"]; ok {
			sshAuthMethods = append(sshAuthMethods, ssh.Password(password))
		}
		// Key provided
		keyLocation, ok := host["ssh_key"]
		// Try default key
		if !ok {
			defaultKeyLocation := path.Join(os.Getenv("HOME"), ".ssh", "id_rsa")
			if _, err := os.Stat(defaultKeyLocation); err != nil {
				keyLocation = defaultKeyLocation
			}
		}
		if keyLocation != "" {
			authMethod, err := clientKeyAuth(keyLocation)
			if err != nil {
				fmt.Printf("Error loading key for host [%s]\n", hostname)
				panic(err)
			}
			sshAuthMethods = append(sshAuthMethods, authMethod)
		}

		config := &ssh.ClientConfig{
			User: username,
			Auth: sshAuthMethods,
		}
		machines = append(machines, &Machine{
			Hostname:  hostname,
			Port:      port,
			SSHConfig: config,
		})
	}
	return machines
}

func clientKeyAuth(keyLocation string) (ssh.AuthMethod, error) {
	buf, err := ioutil.ReadFile(keyLocation)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}
