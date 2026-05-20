package libspot

import (
	"runtime"
)

type AppPlatformEnum string

const (
	AppPlatformUnknown  AppPlatformEnum = "UNKNOWN"
	AppPlatformWin32    AppPlatformEnum = "WIN32"
	AppPlatformWin32X86 AppPlatformEnum = "WIN32_X86_64"
	AppPlatformWin32Arm AppPlatformEnum = "WIN32_ARM64"
	AppPlatformOSX      AppPlatformEnum = "OSX"
	AppPlatformOSXArm   AppPlatformEnum = "OSX_ARM64"
	AppPlatformLinux    AppPlatformEnum = "LINUX"
)

func (a AppPlatformEnum) String() string {
	switch a {
	case AppPlatformLinux,
		AppPlatformOSX,
		AppPlatformOSXArm,
		AppPlatformWin32,
		AppPlatformWin32X86,
		AppPlatformWin32Arm:
		return string(a)
	default:
		return string(AppPlatformUnknown)
	}
}

func AppPlatform() AppPlatformEnum {
	switch runtime.GOOS {
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return AppPlatformOSXArm
		}
		return AppPlatformOSX
	case "linux", "freebsd":
		return AppPlatformLinux
	case "windows":
		switch runtime.GOARCH {
		case "arm64":
			return AppPlatformWin32Arm
		case "amd64":
			return AppPlatformWin32X86
		}

		return AppPlatformWin32
	}

	return AppPlatformUnknown
}
