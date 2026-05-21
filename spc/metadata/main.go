package metadata

type MdType string

const (
	TypeTrack  MdType = "track"
	TypeAlbum  MdType = "album"
	TypeArtist MdType = "artist"
)

type HasMetadataType interface {
	Type() MdType
}

const (
	Path = "metadata/4/{type}/{gid}"
)
