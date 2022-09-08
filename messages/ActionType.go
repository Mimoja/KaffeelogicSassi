package messages

type ActionType uint8

const (
	SASSI_ACTION_CODE_SAVE ActionType = iota + 1
	SASSI_ACTION_CODE_RESTART
	SASSI_ACTION_CODE_FORMAT_FILESYS

	SASSI_ACTION_CODE_LOAD_PROFILE   = 100
	SASSI_ACTION_CODE_SET_LOG_NUMBER = 101
)

const (
	SASSI_RESTART_TYPE_POWERUP = iota + 1
	SASSI_RESTART_TYPE_RESET
	SASSI_RESTART_TYPE_INSTALL
	SASSI_RESTART_TYPE_RESCUE
	SASSI_RESTART_TYPE_BOOTSEL
)
