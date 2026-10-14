# go-snapctl
[![Go Reference](https://pkg.go.dev/badge/github.com/canonical/go-snapctl.svg)](https://pkg.go.dev/github.com/canonical/go-snapctl)

Go wrapper library for the [snapctl](https://snapcraft.io/docs/using-snapctl) tool.

Wrappers for the following subcommands are partially implemented:

- [ ] `fde-setup-request`: Obtain full disk encryption setup request
- [ ] `fde-setup-result`: Set result for full disk encryption
- [x] `get`: The get command prints configuration and interface connection settings.                
- [x] `is-connected`: Return success if the given plug or slot is connected, and failure otherwise
- [ ] `reboot`: Control the reboot behavior of the system          
- [x] `restart`: Restart services    
- [x] `services`: Query the status of services      
- [x] `set`: Changes configuration options
- [ ] `set-health`: Report the health status of a snap
- [x] `start`: Start services 
- [x] `stop`: Stop services
- [ ] `system-mode`: Get the current system mode and associated details
- [x] `unset`: Remove configuration options
- [x] `install`: Install a snap component
- [x] `remove`: Remove a snap component

The commands and descriptions are from `snapctl --help`.

# Usage

```go
package main

import (
	"fmt"

	"github.com/canonical/go-snapctl"
)

func main() {
	// unset
	err := snapctl.Unset("http").Run()
	if err != nil {
		panic(err)
	}

	// set values
	err = snapctl.Set("http.bind-address", "0.0.0.0").Run()
	if err != nil {
		panic(err)
	}
	err = snapctl.Set("http.bind-port", "8080").Run()
	if err != nil {
		panic(err)
	}

	// set values with a JSON object
	err = snapctl.Set("http.tls",
		`{
			"enabled":"true",
			"cert":"./cert.pem",
			"privkey":"./key.pem"
		}`).Document().Run()
	if err != nil {
		panic(err)
	}

	// get one value
	value, err := snapctl.Get("http.bind-port").Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	// Outputs:
	// 8080

	// get values as JSON object
	value, err = snapctl.Get("http").Document().Run()
	if err != nil {
		panic(err)
	}
	fmt.Println(value)
	// Outputs:
	// {
	//   "bind-address": "0.0.0.0",
	//   "bind-port": 8080,
	//   "tls": {
	//     "cert": "./cert.pem",
	// 	   "enabled": "true",
	// 	   "privkey": "./key.pem"
	//   }
	// }

	// start and enable a service
	err := snapctl.Start("snap-name.service-name").Enable().Run()
	if err != nil {
		panic(err)
	}

	// start all services
	err := snapctl.Start("snap.name").Run()
	if err != nil {
		panic(err)
	}
	
	// install components
	err := snapctl.InstallComponents("component-name", "another-component").Run()
	if err != nil {
		panic(err)
	}
}
```

### Testing
The tests need to run in a snap environment:

Build and install the tester snap:
```bash
make build # or clean-build
make install 
```

Run all tests:
```bash
make test
```

The above runs the all the tests with sudo, required to test privileged operations such as `snapctl set`.

To manually run tests, use: 
```bash
sudo go-snapctl-tester.test
```

This app copies the project files to a writable data location inside the snap.
This is to allow running tests (which required file locking) that are in the user's home via the root user.

The `test` app accepts arguments that are supported by both `go test` and `go vet`, appended to the end of the commands.

For example, to run tests in one package, i.e. `log`:
```bash
sudo go-snapctl-tester.test ./log
```
