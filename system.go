package libspot

import (
	"runtime"

	pbdata "github.com/pyrorhythm/libspot/api/spotify/clienttoken/data/v0"
)

func platformSpecificData() *pbdata.PlatformSpecificData {
	switch runtime.GOOS {
	case "android":
		return &pbdata.PlatformSpecificData{
			Data: &pbdata.PlatformSpecificData_Android{
				Android: &pbdata.NativeAndroidData{},
			},
		}
	case "darwin":
		return &pbdata.PlatformSpecificData{
			Data: &pbdata.PlatformSpecificData_Mac{
				Mac: &pbdata.NativeDesktopMacOSData{},
			},
		}
	case "ios":
		return &pbdata.PlatformSpecificData{
			Data: &pbdata.PlatformSpecificData_Ios{
				Ios: &pbdata.NativeIOSData{},
			},
		}
	case "linux", "freebsd":
		return &pbdata.PlatformSpecificData{
			Data: &pbdata.PlatformSpecificData_Linux{
				Linux: &pbdata.NativeDesktopLinuxData{},
			},
		}
	case "windows":
		return &pbdata.PlatformSpecificData{
			Data: &pbdata.PlatformSpecificData_Win{
				Win: &pbdata.NativeDesktopWindowsData{},
			},
		}
	}

	return nil
}
