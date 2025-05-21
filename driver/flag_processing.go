package driver

import (
	"fmt"
	"strings"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var legacyDefaultImages = [...]string{
	defaultImage,
	"ubuntu-18.04",
	"ubuntu-16.04",
	"debian-9",
}

func isDefaultImageName(imageName string) bool {
	for _, defaultImage := range legacyDefaultImages {
		if imageName == defaultImage {
			return true
		}
	}
	return false
}

func (d *Driver) setImageArch(arch string) error {
	switch arch {
	case "":
		d.ImageArch = emptyImageArchitecture
	case string(hcloud.ArchitectureARM):
		d.ImageArch = hcloud.ArchitectureARM
	case string(hcloud.ArchitectureX86):
		d.ImageArch = hcloud.ArchitectureX86
	default:
		return fmt.Errorf("unknown architecture %v", arch)
	}
	return nil
}

func (d *Driver) verifyImageFlags() error {
	if d.ImageID != 0 && d.Image != "" && !isDefaultImageName(d.Image) /* support legacy behaviour */ {
		return d.flagFailure("--%v and --%v are mutually exclusive", flagImage, flagImageID)
	} else if d.ImageID != 0 && d.ImageArch != "" {
		return d.flagFailure("--%v and --%v are mutually exclusive", flagImageArch, flagImageID)
	} else if d.ImageID == 0 && d.Image == "" {
		d.Image = defaultImage
	}
	return nil
}

func (d *Driver) verifyNetworkFlags() error {
	if !d.UsePrivateNetwork {
		return d.flagFailure("--%v must be used if public networking is disabled (hint: implicitly set by --%v)",
			flagUsePrivateNetwork, flagDisablePublic)
	}

	return nil
}

func (d *Driver) setLabelsFromFlags(opts drivers.DriverOptions) error {
	d.ServerLabels = make(map[string]string)
	for _, label := range opts.StringSlice(flagServerLabel) {
		split := strings.SplitN(label, "=", 2)
		if len(split) != 2 {
			return d.flagFailure("server label %v is not in key=value format", label)
		}
		d.ServerLabels[split[0]] = split[1]
	}
	d.keyLabels = make(map[string]string)
	for _, label := range opts.StringSlice(flagKeyLabel) {
		split := strings.SplitN(label, "=", 2)
		if len(split) != 2 {
			return fmt.Errorf("key label %v is not in key=value format", label)
		}
		d.keyLabels[split[0]] = split[1]
	}
	return nil
}
