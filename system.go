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
