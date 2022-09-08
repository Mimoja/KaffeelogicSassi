package messages

type OutcomeCode uint8

const (
	SASSI_OUTCOME_CODE_SUCCESS         OutcomeCode = 0
	SASSI_OUTCOME_CODE_NO_ROOM                     = 1
	SASSI_OUTCOME_CODE_NOT_VALID                   = 2
	SASSI_OUTCOME_CODE_NOT_FOUND                   = 3
	SASSI_OUTCOME_CODE_SYS_ERR                     = 4
	SASSI_OUTCOME_CODE_TOO_BIG                     = 5
	SASSI_OUTCOME_CODE_DIR_NOT_EMPTY               = 6
	SASSI_OUTCOME_CODE_LFS_ERR                     = 10
	SASSI_OUTCOME_CODE_UNSPECIFIED_ERR             = 20
	SASSI_OUTCOME_CODE_CANNOT_SAVE                 = 50
	SASSI_OUTCOME_CODE_CANNOT_OPEN                 = 51
	SASSI_OUTCOME_CODE_CANNOT_RENAME               = 52
	SASSI_OUTCOME_CODE_CANNOT_DELETE               = 53
	SASSI_OUTCOME_CODE_CANNOT_MKDIR                = 54
	SASSI_OUTCOME_CODE_NAME_TOO_LONG               = 100
	SASSI_OUTCOME_CODE_DATA_SEQ_ERR                = 101
	SASSI_OUTCOME_CODE_TIMEOUT                     = 102
	SASSI_OUTCOME_CODE_BUSY                        = 103
	SASSI_OUTCOME_CODE_LAST_PACKET                 = 128
)
