package libspot

import (
	"runtime"

	datav0 "github.com/pyrorhythm/libspot/gen/spotify/clienttoken/data/v0"
)

func platformSpecificData() *datav0.PlatformSpecificData {
	psd := datav0.PlatformSpecificData_builder{}

	switch runtime.GOOS {
	case "android":
		psd.Android = &datav0.NativeAndroidData{}
	case "darwin":
		psd.Mac = &datav0.NativeDesktopMacOSData{}
	case "ios":
		psd.Ios = &datav0.NativeIOSData{}
	case "linux", "freebsd":
		psd.Linux = &datav0.NativeDesktopLinuxData{}
	case "windows":
		psd.Win = &datav0.NativeDesktopWindowsData{}
	}

	return psd.Build()
}

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
	case AppPlatformLinux, AppPlatformOSX, AppPlatformOSXArm, AppPlatformWin32, AppPlatformWin32X86, AppPlatformWin32Arm:
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
