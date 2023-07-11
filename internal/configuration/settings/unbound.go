package settings

import (
	"errors"
	"net/netip"

	"github.com/qdm12/dns/v2/pkg/provider"
	"github.com/qdm12/gosettings"
	"github.com/qdm12/gotree"
)

// Unbound is settings for the Unbound program.
type Unbound struct {
	Providers []string `json:"providers"`
	Caching   *bool    `json:"caching"`
	IPv6      *bool    `json:"ipv6"`
}

func (u *Unbound) setDefaults() {
	u.Providers = gosettings.DefaultSlice(u.Providers, []string{
		provider.Cloudflare().Name,
	})
	u.Caching = gosettings.DefaultPointer(u.Caching, true)
	u.IPv6 = gosettings.DefaultPointer(u.IPv6, false)
}

func (u Unbound) validate() (err error) {
	for _, s := range u.Providers {
		_, err := provider.Parse(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func (u Unbound) copy() (copied Unbound) {
	return Unbound{
		Providers: gosettings.CopySlice(u.Providers),
		Caching:   gosettings.CopyPointer(u.Caching),
		IPv6:      gosettings.CopyPointer(u.IPv6),
	}
}

func (u *Unbound) mergeWith(other Unbound) {
	u.Providers = gosettings.MergeWithSlice(u.Providers, other.Providers)
	u.Caching = gosettings.MergeWithPointer(u.Caching, other.Caching)
	u.IPv6 = gosettings.MergeWithPointer(u.IPv6, other.IPv6)
}

func (u *Unbound) overrideWith(other Unbound) {
	u.Providers = gosettings.OverrideWithSlice(u.Providers, other.Providers)
	u.Caching = gosettings.OverrideWithPointer(u.Caching, other.Caching)
	u.IPv6 = gosettings.OverrideWithPointer(u.IPv6, other.IPv6)
}

var (
	ErrConvertingNetip = errors.New("converting net.IP to netip.Addr failed")
)

func (u Unbound) GetFirstPlaintextIPv4() (ipv4 netip.Addr) {
	s := u.Providers[0]
	provider, err := provider.Parse(s)
	if err != nil {
		// Settings should be validated before calling this function,
		// so an error happening here is a programming error.
		panic(err)
	}

	return provider.DNS.IPv4[0]
}

func (u Unbound) String() string {
	return u.toLinesNode().String()
}

func (u Unbound) toLinesNode() (node *gotree.Node) {
	node = gotree.New("DNS over TLS settings:")

	authServers := node.Appendf("Authoritative servers:")
	for _, provider := range u.Providers {
		authServers.Appendf(provider)
	}

	node.Appendf("Caching: %s", gosettings.BoolToYesNo(u.Caching))
	node.Appendf("IPv6: %s", gosettings.BoolToYesNo(u.IPv6))

	return node
}
