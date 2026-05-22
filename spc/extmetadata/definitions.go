package extmetadata

import (
	audiofilespb "pyrorhythm.dev/libspot/gen/spotify/extendedmetadata/audiofiles"
	metadatapb "pyrorhythm.dev/libspot/gen/spotify/metadata"
)

var (
	AudioFiles = Define[audiofilespb.AudioFilesExtensionResponse](KindAudioFiles)
	Track      = Define[metadatapb.Track](KindTrackV4)
	Album      = Define[metadatapb.Album](KindAlbumV4)
	Artist     = Define[metadatapb.Artist](KindArtistV4)
)
