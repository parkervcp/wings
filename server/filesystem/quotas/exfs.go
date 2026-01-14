package quotas

import (
	"os"

	"github.com/g0rbe/go-chattr"
)

// EnableEXFSQuota enables quotas on a specified folder
func EnableEXFSQuota(serverPath string) (err error) {
	serverdir, err := os.OpenFile(serverPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}

	err = chattr.SetAttr(serverdir, chattr.FS_PROJINHERIT_FL)
	return
}

// DisableEXFSQuota disables quotas on a specified folder
func DisableEXFSQuota(serverPath string) (err error) {

	return
}

// SetEXFSQuota sets the quota in bytes for the specified server uuid
func SetEXFSQuota(serverUUID string, byteLimit int64) (err error) {

	return
}

// GetEXFSQuota gets the specified quotas and usage of a specified server uuid
func GetEXFSQuota(serverUUID string) (byteLimit, bytesUsed int64, err error) {

	return
}

// NewEXFSQuota
func NewEXFSQuota() (err error) {
	return
}
