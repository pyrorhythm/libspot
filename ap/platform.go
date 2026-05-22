package ap

import (
	"fmt"
	"runtime"

	pb "pyrorhythm.dev/libspot/gen/spotify"
)

// SpotifyVersionCode is the client version advertised during the handshake.
const SpotifyVersionCode = 127700358

func getOS() pb.Os {
	switch runtime.GOOS {
	case "android":
		return pb.Os_OS_ANDROID
	case "darwin":
		return pb.Os_OS_OSX
	case "freebsd":
		return pb.Os_OS_FREEBSD
	case "ios":
		return pb.Os_OS_IPHONE
	case "linux":
		return pb.Os_OS_LINUX
	case "windows":
		return pb.Os_OS_WINDOWS
	default:
		return pb.Os_OS_UNKNOWN
	}
}

func getCpuFamily() pb.CpuFamily {
	switch runtime.GOARCH {
	case "386":
		return pb.CpuFamily_CPU_X86
	case "amd64":
		return pb.CpuFamily_CPU_X86_64
	case "arm", "arm64":
		return pb.CpuFamily_CPU_ARM
	case "mips", "mips64":
		return pb.CpuFamily_CPU_MIPS
	case "ppc64":
		return pb.CpuFamily_CPU_PPC_64
	default:
		return pb.CpuFamily_CPU_UNKNOWN
	}
}

func getPlatform() pb.Platform {
	switch runtime.GOOS {
	case "android":
		return pb.Platform_PLATFORM_ANDROID_ARM
	case "darwin":
		switch runtime.GOARCH {
		case "386":
			return pb.Platform_PLATFORM_OSX_X86
		case "amd64":
			return pb.Platform_PLATFORM_OSX_X86_64
		case "ppc64":
			return pb.Platform_PLATFORM_OSX_PPC
		}
	case "freebsd":
		switch runtime.GOARCH {
		case "386":
			return pb.Platform_PLATFORM_FREEBSD_X86
		case "amd64":
			return pb.Platform_PLATFORM_FREEBSD_X86_64
		}
	case "ios":
		switch runtime.GOARCH {
		case "arm":
			return pb.Platform_PLATFORM_IPHONE_ARM
		case "arm64":
			return pb.Platform_PLATFORM_IPHONE_ARM64
		}
	case "linux":
		switch runtime.GOARCH {
		case "386":
			return pb.Platform_PLATFORM_LINUX_X86
		case "amd64":
			return pb.Platform_PLATFORM_LINUX_X86_64
		case "mips", "mips64":
			return pb.Platform_PLATFORM_LINUX_MIPS
		case "arm", "arm64":
			return pb.Platform_PLATFORM_LINUX_ARM
		}
	case "windows":
		switch runtime.GOARCH {
		case "386":
			return pb.Platform_PLATFORM_WIN32_X86
		case "amd64":
			return pb.Platform_PLATFORM_WIN32_X86_64
		case "arm", "arm64":
			return pb.Platform_PLATFORM_WINDOWS_CE_ARM
		}
	}

	return pb.Platform_PLATFORM_GENERIC_PARTNER
}

func versionString() string {
	return "libspot 0.0.0"
}

func systemInfoString() string {
	return fmt.Sprintf("%s; Go %s (%s %s)", versionString(), runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

// obfuscateUsername hides all but the first and last character of a username
// for logging.
func obfuscateUsername(username string) string {
	switch len(username) {
	case 0, 1, 2:
		return username
	default:
		return fmt.Sprintf("%c%s%c", username[0], "***", username[len(username)-1])
	}
}
