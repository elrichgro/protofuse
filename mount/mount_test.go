package mount

import (
	"testing"
	"github.com/elrichgro/protofuse/test"
)

func TestMount(t *testing.T) {
	err := Mount(nil, nil, "", "", "invalid_mount_point")
	if err == nil {
		t.FailNow()
	}

	buf, fDesc, packageName, messageName, err := test.GenerateFull()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err = Mount(buf, fDesc, packageName, messageName, "../test/mp")
		if err != nil {
			t.Fatal(err)
		}
	}()	
	err = Unmount("../test/mp")
	if err != nil {
		t.Fatal(err)
	}
}