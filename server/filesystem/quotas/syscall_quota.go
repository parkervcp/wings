package quotas

const (
	Q_GETQUOTA     = 0x0080000700
	Q_SETQUOTA     = 0x0080000800
	Q_GETNEXTQUOTA = 0x00080000900
)

const (
	USRQUOTA = 0x0000000000
	GRPQUOTA = 0x0000000001
	PRJQUOTA = 0x0000000002
)

type DQBlk struct {
	dqbBHardlimit uint64
	dqbBSoftlimit uint64
	dqbCurSpace   uint64
	dqbIHardlimit uint64
	dqbISoftlimit uint64
	dqbCurInodes  uint64
	dqbBTime      uint64
	dqbITime      uint64
	dqbValid      uint32
}
