package loaders

type ServerJar struct {
	Loader  string // "vanilla" | "paper" | "fabric"
	Version string // Minecraft version, ex: "1.21.1"
	Label   string // friendly description to show on screen (ex: "Paper 1.21.1 (build 132)")
	URL     string
}

const (
	Vanilla = "vanilla"
	Paper   = "paper"
	Fabric  = "fabric"
)

func Releases(loader string, limit int) ([]string, error) {
	switch loader {
	case Vanilla:
		return VanillaReleases(limit)
	case Paper:
		return PaperReleases(limit)
	case Fabric:
		return FabricReleases(limit)
	default:
		return nil, unknownLoaderErr(loader)
	}
}

func Resolve(loader, version string) (*ServerJar, error) {
	switch loader {
	case Vanilla:
		return VanillaServerJar(version)
	case Paper:
		return PaperServerJar(version)
	case Fabric:
		return FabricServerJar(version)
	default:
		return nil, unknownLoaderErr(loader)
	}
}

func unknownLoaderErr(loader string) error {
	return &UnknownLoaderError{Loader: loader}
}

type UnknownLoaderError struct {
	Loader string
}

func (e *UnknownLoaderError) Error() string {
	return "unknown loader: " + e.Loader
}
