package quotas

import (
	"syscall"

	"emperror.dev/errors"
	"github.com/parkervcp/fsquota"
	"github.com/pelican-dev/wings/config"
)

const (
	FSBTRFS = 2435016766
	FSEXT4  = 61267
	FSXFS   = 1481003842
	FSZFS   = 801189825
)

var fstype string

func getFSType(mount string) (fsType uint, err error) {
	var stat syscall.Statfs_t

	if mount == "" {
		return fsType, errors.New("must specify path to check the filesystem type")
	}

	err = syscall.Statfs(mount, &stat)
	if err != nil {
		return fsType, err
	}

	switch stat.Type {
	case FSBTRFS:
		return FSBTRFS, nil
	case FSEXT4:
		return FSEXT4, nil
	case FSXFS:
		return FSXFS, nil
	case FSZFS:
		return FSZFS, nil
	default:
		return fsType, errors.New("unknown filesystem type")
	}
}

// IsSupportedFS checks if the filesystem for the data files is supported.
// currently only EXT4 and XFS are supported
func IsSupportedFS() (err error) {
	checked, err := getFSType(config.Get().System.Data)
	if err != nil {
		return err
	}

	switch checked {
	case FSEXT4 | FSXFS:
		// technically tested on EXT4 and will need to be validated for XFS
		supported, err := fsquota.ProjectQuotasSupported(config.Get().System.Data)
		if err != nil {
			return err
		}
		if !supported {
			return errors.New("project quotas not enabled")
		}

		fstype = "exfs"
		return err
	case FSBTRFS:
		fstype = "btrfs"
		return errors.New("btrfs is not supported on this filesystem")
	case FSZFS:
		fstype = "zfs"
		return errors.New("zfs is not supported on this filesystem")
	default:
		return errors.New("unknown filesystem type")
	}
}

// AddQuota adds a server to the configured quotas
func AddQuota(serverID int, serverUUID string) (err error) {
	switch fstype {
	case "exfs":
		err = exfsProject{ID: serverID, Name: serverUUID}.addProject()
	}

	return
}

// DelQuota removes a server from the configured quotas
func DelQuota(serverUUID string) (err error) {
	switch fstype {
	case "exfs":
		err = exfsProject{Name: serverUUID}.removeProject()
	}
	return
}

// SetQuota configures quotas for a specified server
func SetQuota(limit int64, serverUUID string) (err error) {
	switch fstype {
	case "exfs":
		err = exfsProject{Name: serverUUID}.setQuota(uint64(limit))
	}
	return
}

// GetQuota gets the data usage for a specified server
func GetQuota(serverUUID string) (used int64, err error) {
	switch fstype {
	case "exfs":
		used, err = exfsProject{Name: serverUUID}.getQuota()
	}
	return
}
