package driver

import (
	"strconv"
	"strings"
	"testing"

	"github.com/docker/machine/commands/commandstest"
	"github.com/docker/machine/libmachine/drivers"
	"github.com/hetznercloud/hcloud-go/v2/hcloud"
)

var defaultFlags = map[string]interface{}{
	flagAPIToken: "foo",
}

func makeFlags(args map[string]interface{}) drivers.DriverOptions {
	combined := make(map[string]interface{}, len(defaultFlags)+len(args))
	for k, v := range defaultFlags {
		combined[k] = v
	}
	for k, v := range args {
		combined[k] = v
	}

	return &commandstest.FakeFlagger{Data: combined}
}

func TestDisablePublic(t *testing.T) {
	d := NewDriver("test")
	err := d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagDisablePublic: true,
	}))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if !d.UsePrivateNetwork {
		t.Error("expected private network to be enabled")
	}
}

func TestImageFlagExclusions(t *testing.T) {
	// both id and name given
	d := NewDriver("test")
	err := d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImageID: "42",
		flagImage:   "answer",
	}))
	assertMutualExclusion(t, err, flagImageID, flagImage)

	// both id and arch given
	d = NewDriver("test")
	err = d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImageID:   "42",
		flagImageArch: string(hcloud.ArchitectureX86),
	}))
	assertMutualExclusion(t, err, flagImageID, flagImageArch)
}

func TestImageArch(t *testing.T) {
	// no explicit arch
	d := NewDriver("test")
	err := d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImage: "answer",
	}))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if d.ImageArch != emptyImageArchitecture {
		t.Errorf("expected empty architecture, but got %v", d.ImageArch)
	}

	// existing architectures
	testArchFlag(t, hcloud.ArchitectureARM)
	testArchFlag(t, hcloud.ArchitectureX86)

	// invalid
	d = NewDriver("test")
	err = d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImage:     "answer",
		flagImageArch: "hal9000",
	}))
	if err == nil {
		t.Fatal("expected error, but invalid arch was accepted")
	}
}

func TestBogusId(t *testing.T) {
	d := NewDriver("test")
	err := d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImageID: "answer",
	}))
	if err == nil {
		t.Fatal("expected error, but invalid arch was accepted")
	}
}

func TestLongId(t *testing.T) {
	var testId int64 = 79871865169581

	d := NewDriver("test")
	err := d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImageID: strconv.FormatInt(testId, 10),
	}))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if d.ImageID != testId {
		t.Errorf("expected %v id, but got %v", testId, d.ImageArch)
	}
}

func testArchFlag(t *testing.T, arch hcloud.Architecture) {
	d := NewDriver("test")
	err := d.setConfigFromFlagsImpl(makeFlags(map[string]interface{}{
		flagImage:     "answer",
		flagImageArch: string(arch),
	}))
	if err != nil {
		t.Fatalf("unexpected error, %v", err)
	}

	if d.ImageArch != arch {
		t.Errorf("expected %v architecture, but got %v", arch, d.ImageArch)
	}
}

func assertMutualExclusion(t *testing.T, err error, flag1, flag2 string) {
	if err == nil {
		t.Errorf("expected mutually exclusive flags to fail, but no error was thrown: %v %v", flag1, flag2)
		return
	}

	errstr := err.Error()
	if !(strings.Contains(errstr, flag1) && strings.Contains(errstr, flag2) && strings.Contains(errstr, "mutually exclusive")) {
		t.Errorf("expected mutually exclusive flags to fail, but message differs: %v %v %v", flag1, flag2, errstr)
	}
}
