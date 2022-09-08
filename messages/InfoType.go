package messages

type InfoType uint8

const (
	SASSI_INFO_CODE_FOLDER_UPDATED     InfoType = iota + 1 //Unknown?
	SASSI_INFO_CODE_FAULT                                  // Software Failrue
	SASSI_INFO_CODE_SYSTEM_INFO                            // Roaster System Info
	SASSI_INFO_CODE_FILESYS_INFO                           // Roaster FileSystem Info
	SASSI_INFO_CODE_TECHNICAL_INFO                         // Roaster Techincal Info
	SASSI_INFO_CODE_APPLIANCE_BUSY                         // Roaster is busy
	SASSI_INFO_CODE_APPLIANCE_NOT_BUSY                     // Roaster is not busy!?
	SASSI_INFO_CODE_FILE_UPDATED                           // Unknown?
	SASSI_INFO_CODE_OPERATIONAL_STATUS                     // Roaster System Status
)
