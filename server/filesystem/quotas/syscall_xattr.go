package quotas

import (
	"os"
	"syscall"
	"unsafe"
)

// Pulled definitions from /usr/include/linux/fs.h

/*
 * Flags for the fsx_xflags field
 */
const (
	FS_XFLAG_REALTIME     = 0x00000001 /* data in realtime volume */
	FS_XFLAG_PREALLOC     = 0x00000002 /* preallocated file extents */
	FS_XFLAG_IMMUTABLE    = 0x00000008 /* file cannot be modified */
	FS_XFLAG_APPEND       = 0x00000010 /* all writes append */
	FS_XFLAG_SYNC         = 0x00000020 /* all writes synchronous */
	FS_XFLAG_NOATIME      = 0x00000040 /* do not update access time */
	FS_XFLAG_NODUMP       = 0x00000080 /* do not include in backups */
	FS_XFLAG_RTINHERIT    = 0x00000100 /* create with rt bit set */
	FS_XFLAG_PROJINHERIT  = 0x00000200 /* create with parents projid */
	FS_XFLAG_NOSYMLINKS   = 0x00000400 /* disallow symlink creation */
	FS_XFLAG_EXTSIZE      = 0x00000800 /* extent size allocator hint */
	FS_XFLAG_EXTSZINHERIT = 0x00001000 /* inherit inode extent size */
	FS_XFLAG_NODEFRAG     = 0x00002000 /* do not defragment */
	FS_XFLAG_FILESTREAM   = 0x00004000 /* use filestream allocator */
	FS_XFLAG_DAX          = 0x00008000 /* use DAX for IO */
	FS_XFLAG_COWEXTSIZE   = 0x00010000 /* CoW extent size allocator hint */
	FS_XFLAG_HASATTR      = 0x80000000 /* no DIFLAG for this   */
)

/*
#define FS_IOC_GETFLAGS                 _IOR('f', 1, long)
#define FS_IOC_SETFLAGS                 _IOW('f', 2, long)
#define FS_IOC_FSGETXATTR               _IOR('X', 31, struct fsxattr)
#define FS_IOC_FSSETXATTR               _IOW('X', 32, struct fsxattr)
*/

const (
	FS_IOC_FSGETXATTR uintptr = 0x801c581f // https://docs.rs/linux-raw-sys/latest/linux_raw_sys/ioctl/constant.FS_IOC_FSGETXATTR.html
	FS_IOC_FSSETXATTR uintptr = 0x401c5820 // https://docs.rs/linux-raw-sys/latest/linux_raw_sys/ioctl/constant.FS_IOC_FSSETXATTR.html
)

// fsXAttr is the struct defining the structure
// for FS_IOC_FSGETXATTR and FS_IOC_FSSETXATTR
type fsXAttr struct {
	XFlags    uint32
	ExtSize   uint32
	NextENTs  uint32
	ProjectID uint32
	FSXPad    byte
}

// xAttrCtl sets the
func xAttrCtl(f *os.File, request uintptr, xattr *fsXAttr) (err error) {
	attreq := uintptr(unsafe.Pointer(xattr))

	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, f.Fd(), request, attreq)

	if errno != 0 {
		return os.NewSyscallError("ioctl", errno)
	}

	return
}

// getXAttr gets the extended attributes of a file
func getXAttr(f *os.File) (attr fsXAttr, err error) {
	if err = xAttrCtl(f, FS_IOC_FSGETXATTR, &attr); err != nil {
		return
	}

	return
}

// setXAttr sets xattr values for the
func setXAttr(serverDir *os.File, fsXAttr fsXAttr) (err error) {
	xAttr, err := getXAttr(serverDir)
	if err != nil {
		return err
	}

	// bitwise add for uint32 X Attributes
	xAttr.XFlags |= fsXAttr.XFlags

	err = xAttrCtl(serverDir, FS_IOC_FSSETXATTR, &fsXAttr)

	return
}
