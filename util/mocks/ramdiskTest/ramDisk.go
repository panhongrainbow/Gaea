package ramdiskTest

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	_ "bazil.org/fuse/fs/fstestutil"
)

func Create(mountpoint string) {
	c, err := fuse.Mount(
		mountpoint,
		fuse.FSName("helloworld"),
		fuse.Subtype("hellofs"),
	)
	if err != nil {
		return
	}
	defer c.Close()

	err = fs.Serve(c, FS{})
	if err != nil {
		return
	}
}

func Destory(c *fuse.Conn) error {
	return c.Close()
}
