package mount

import (
	// "os/exec"
	// "fmt"
	"time"
	"testing"
	"github.com/elrichgro/protofuse/test"
)

func TestInvalidMount(t *testing.T) {
	err := Mount(nil, nil, "", "", "invalid_mount_point")
	if err == nil {
		t.FailNow()
	}
}

func TestUnmount(t *testing.T) {
	c := make(chan bool)
	var mountpoint string = "../test/mp"

	buf, fDesc, packageName, messageName, err := test.GenerateFull()
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		err = Mount(buf, fDesc, packageName, messageName, mountpoint)
		if err != nil {
			t.Fatal(err)
		}
		c <- true
	}()	
	time.Sleep(100*time.Millisecond)
	err = Unmount(mountpoint)
	if err != nil {
		t.Fatal(err)
	}
	unmounted := <-c
	if !unmounted {
		t.FailNow()
	}
}

func TestMountLarge(t *testing.T) {
	// c := make(chan bool)
	var mountpoint string = "../test/mp"

	buf, fDesc, packageName, messageName, err := test.GenerateLarge()
	if err != nil {
		t.Fatal(err)
	}

	// fmt.Println(len(buf))

	go func() {
		err = MountList(buf, fDesc, packageName, messageName, mountpoint)
		if err != nil {
			t.Fatal(err)
		}
		// c <- true
	}()	
	// time.Sleep(1000*time.Millisecond)
	// err = Unmount(mountpoint)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// unmounted := <-c
	// if !unmounted {
	// 	t.FailNow()
	// }
}