package pkgmgr

import "runtime"

type OSName string

const (
	Linux   OSName = "linux"
	Mac     OSName = "mac"
	Windows OSName = "windows"
	Unknown OSName = "unknown"
)

func DetectOS() OSName {
	switch runtime.GOOS {
	case "linux":
		return Linux
	case "darwin":
		return Mac
	case "windows":
		return Windows
	default:
		return Unknown
	}
}
