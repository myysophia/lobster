package platform

import (
	"os"
	"runtime"
)

type OS string

const (
	Windows OS = "windows"
	Darwin  OS = "darwin"
	Linux   OS = "linux"
	Unknown OS = "unknown"
)

type Info struct {
	OS         OS
	Arch       string
	HasDesktop bool
}

func Detect() Info {
	info := Info{
		OS:   normalizeOS(runtime.GOOS),
		Arch: runtime.GOARCH,
	}

	switch info.OS {
	case Darwin, Windows:
		info.HasDesktop = true
	case Linux:
		info.HasDesktop = os.Getenv("DISPLAY") != "" || os.Getenv("WAYLAND_DISPLAY") != ""
	default:
		info.HasDesktop = false
	}

	return info
}

func (i Info) String() string {
	return string(i.OS) + "/" + i.Arch
}

func normalizeOS(goos string) OS {
	switch goos {
	case "windows":
		return Windows
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	default:
		return Unknown
	}
}
