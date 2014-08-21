// Package vagrant implements a driver based on Vagrant.
package vagrant

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	goexec "os/exec"
	"path"
	"strings"

	"github.com/mlafeldt/chef-runner/log"
	"github.com/mlafeldt/chef-runner/openssh"
	"github.com/mlafeldt/chef-runner/rsync"
)

const (
	// DefaultMachine is the name of the default Vagrant machine.
	DefaultMachine = "default"

	// ConfigPath is the path to the local directory where chef-runner
	// stores Vagrant-specific information.
	ConfigPath = ".chef-runner/vagrant"
)

// Driver is a driver based on Vagrant.
type Driver struct {
	machine     string
	sshClient   *openssh.Client
	rsyncClient *rsync.Client
}

func init() {
	os.Setenv("VAGRANT_NO_PLUGINS", "1")
}

// NewDriver creates a new Vagrant driver that communicates with the given
// Vagrant machine. Under the hood `vagrant ssh-config` is executed to get a
// working SSH configuration for the machine.
func NewDriver(machine string) (*Driver, error) {
	if machine == "" {
		machine = DefaultMachine
	}

	// TODO: reuse existing config file, but make sure it's still valid
	log.Debug("Asking Vagrant for SSH config")
	out, err := goexec.Command("vagrant", "ssh-config", machine).CombinedOutput()
	if err != nil {
		msg := fmt.Sprintf("`vagrant ssh-config` failed with output:\n\n%s", out)
		return nil, errors.New(msg)
	}

	configFile := path.Join(ConfigPath, "machines", machine, "ssh_config")
	log.Debug("Writing current SSH config to", configFile)
	if err := os.MkdirAll(path.Dir(configFile), 0755); err != nil {
		return nil, err
	}
	if err := ioutil.WriteFile(configFile, out, 0644); err != nil {
		return nil, err
	}

	sshClient := &openssh.Client{
		Host:       "default",
		ConfigFile: configFile,
	}

	rsyncClient := rsync.MirrorClient
	rsyncClient.RemoteHost = "default"
	rsyncClient.RemoteShell = strings.Join(sshClient.Shell(), " ")

	return &Driver{machine, sshClient, rsyncClient}, nil
}

// RunCommand runs the specified command on the Vagrant machine.
func (drv Driver) RunCommand(args []string) error {
	return drv.sshClient.RunCommand(args)
}

// Upload copies files to the Vagrant machine.
func (drv Driver) Upload(dst string, src ...string) error {
	return drv.rsyncClient.Copy(dst, src...)
}

// String returns the driver's name.
func (drv Driver) String() string {
	return fmt.Sprintf("Vagrant driver (machine: %s)", drv.machine)
}
