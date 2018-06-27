// +build !nofuse

package node

import (
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
	"time"

	"context"

	core "github.com/ahester57/go-ipfs/core"
	ipns "github.com/ahester57/go-ipfs/fuse/ipns"
	mount "github.com/ahester57/go-ipfs/fuse/mount"
	namesys "github.com/ahester57/go-ipfs/namesys"

	offroute "gx/ipfs/QmWLQyLU7yopJnwMvpHM5VSMG4xmbKgcq6P246mDy9xy5E/go-ipfs-routing/offline"
	ci "gx/ipfs/QmcW4FGAt24fdK1jBgWQn3yP4R9ZLyWQqjozv9QK7epRhL/go-testutil/ci"
)

func maybeSkipFuseTests(t *testing.T) {
	if ci.NoFuse() {
		t.Skip("Skipping FUSE tests")
	}
}

func mkdir(t *testing.T, path string) {
	err := os.Mkdir(path, os.ModeDir|os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
}

// Test externally unmounting, then trying to unmount in code
func TestExternalUnmount(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}

	// TODO: needed?
	maybeSkipFuseTests(t)

	node, err := core.NewNode(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	err = node.LoadPrivateKey()
	if err != nil {
		t.Fatal(err)
	}

	node.Routing = offroute.NewOfflineRouter(node.Repo.Datastore(), node.RecordValidator)
	node.Namesys = namesys.NewNameSystem(node.Routing, node.Repo.Datastore(), 0)

	err = ipns.InitializeKeyspace(node, node.PrivateKey)
	if err != nil {
		t.Fatal(err)
	}

	// get the test dir paths (/tmp/fusetestXXXX)
	dir, err := ioutil.TempDir("", "fusetest")
	if err != nil {
		t.Fatal(err)
	}

	ipfsDir := dir + "/ipfs"
	ipnsDir := dir + "/ipns"
	mkdir(t, ipfsDir)
	mkdir(t, ipnsDir)

	err = Mount(node, ipfsDir, ipnsDir)
	if err != nil {
		t.Fatal(err)
	}

	// Run shell command to externally unmount the directory
	cmd := "fusermount"
	args := []string{"-u", ipnsDir}
	if err := exec.Command(cmd, args...).Run(); err != nil {
		t.Fatal(err)
	}

	// TODO(noffle): it takes a moment for the goroutine that's running fs.Serve to be notified and do its cleanup.
	time.Sleep(time.Millisecond * 100)

	// Attempt to unmount IPNS; check that it was already unmounted.
	err = node.Mounts.Ipns.Unmount()
	if err != mount.ErrNotMounted {
		t.Fatal("Unmount should have failed")
	}

	// Attempt to unmount IPFS; it should unmount successfully.
	err = node.Mounts.Ipfs.Unmount()
	if err != nil {
		t.Fatal(err)
	}
}
