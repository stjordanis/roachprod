package local

import (
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"
	"time"

	"github.com/cockroachdb/roachprod/config"
	"github.com/cockroachdb/roachprod/install"
	"github.com/cockroachdb/roachprod/vm"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

// ProviderName is config.Local.
const ProviderName = config.Local

func init() {
	vm.Providers[ProviderName] = &Provider{}
}

// A Provider is used to create stub VM objects.
type Provider struct{}

// No-op implementation of ProviderFlags
type emptyFlags struct{}

// ConfigureCreateFlags is part of ProviderFlags.  This implementation is a no-op.
func (o *emptyFlags) ConfigureCreateFlags(flags *pflag.FlagSet) {
}

// ConfigureClusterFlags is part of ProviderFlags.  This implementation is a no-op.
func (o *emptyFlags) ConfigureClusterFlags(*pflag.FlagSet) {
}

// CleanSSH is part of the vm.Provider interface.  This implementation is a no-op.
func (p *Provider) CleanSSH() error {
	return nil
}

// ConfigSSH is part of the vm.Provider interface.  This implementation is a no-op.
func (p *Provider) ConfigSSH() error {
	return nil
}

// Create just creates fake host-info entries in the local filesystem
func (p *Provider) Create(names []string, opts vm.CreateOpts) error {
	path := filepath.Join(os.ExpandEnv(config.DefaultHostDir), config.Local)
	file, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "problem creating file %s", path)
	}
	defer file.Close()

	// Align columns left and separate with at least two spaces.
	tw := tabwriter.NewWriter(file, 0, 8, 2, ' ', 0)
	tw.Write([]byte("# user@host\tlocality\n"))
	for i := 0; i < len(names); i++ {
		tw.Write([]byte(fmt.Sprintf(
			"%s@%s\t%s\n", config.OSUser.Username, "127.0.0.1", "region=local,zone=local")))
	}
	if err := tw.Flush(); err != nil {
		return errors.Wrapf(err, "problem writing file %s", path)
	}
	return nil
}

// Delete is part of the vm.Provider interface. This implementation is a no-op.
func (p *Provider) Delete(vms vm.List) error {
	return nil
}

// Extend is part of the vm.Provider interface.  This implementation returns an error.
func (p *Provider) Extend(vms vm.List, lifetime time.Duration) error {
	return errors.New("local clusters have unlimited lifetime")
}

// FindActiveAccount is part of the vm.Provider interface. This implementation is a no-op.
func (p *Provider) FindActiveAccount() (string, error) {
	return "", nil
}

// Flags is part of the vm.Provider interface. This implementation is a no-op.
func (p *Provider) Flags() vm.ProviderFlags {
	return &emptyFlags{}
}

// List constructs N-many localhost VM instances, using SyncedCluster as a way to remember
// how many nodes we should have
func (p *Provider) List() (ret vm.List, _ error) {
	if sc, ok := install.Clusters[ProviderName]; ok {
		now := time.Now()
		for range sc.VMs {
			ret = append(ret, vm.VM{
				Name:        "localhost",
				CreatedAt:   now,
				Lifetime:    time.Hour,
				PrivateIP:   "127.0.0.1",
				Provider:    ProviderName,
				ProviderID:  ProviderName,
				PublicIP:    "127.0.0.1",
				RemoteUser:  config.OSUser.Username,
				VPC:         ProviderName,
				MachineType: ProviderName,
				Zone:        ProviderName,
			})
		}
	}
	return
}

// Name returns the name of the Provider, which will also surface in VM.Provider
func (p *Provider) Name() string {
	return ProviderName
}
