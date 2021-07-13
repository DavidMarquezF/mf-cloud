package uri

const (
	FirmwareIdKey string = "firmwareid"

	API string = "/v1/fmw"

	//devices
	Device     = API + "/{" + FirmwareIdKey + "}"
	Executable = Device + "/exec"
	Detail     = Device + "/inf"
)
