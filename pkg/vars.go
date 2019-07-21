package pkg

const (
	defaultFilePerm = 0644
)

var (
	maxNewRequest = 1000
	maxPendingRequest = 1000
	maxTries = 16
	dumpFilePath = ""
	configFile = ""
	senders map[uint64]sender
)
