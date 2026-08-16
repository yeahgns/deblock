package loaders

import "fmt"

// API publica da Fabric. Documentada em https://fabricmc.net/wiki/documentation:fabric_meta
var fabricMetaBase = "https://meta.fabricmc.net/v2/versions"

type fabricVersionEntry struct {
	Version string `json:"version"`
	Stable  bool   `json:"stable"`
}

func FabricReleases(limit int) ([]string, error) {
	var games []fabricVersionEntry
	if err := getJSON(fabricMetaBase+"/game", &games); err != nil {
		return nil, fmt.Errorf("Couldn't retrieve the list of Fabric versions.: %w", err)
	}

	releases := make([]string, 0, limit)
	for _, g := range games {
		if !g.Stable {
			continue
		}
		releases = append(releases, g.Version)
		if len(releases) >= limit {
			break
		}
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("The Fabric API did not return any stable version.")
	}
	return releases, nil
}

// FabricServerJar monta a URL do jar de servidor "pronto pra rodar" do Fabric
// (a propria API ja empacota o installer junto, nao precisa de passo extra)
// para a versao do jogo pedida, usando a versao mais recente e estavel do
// loader e do installer.
func FabricServerJar(gameVersion string) (*ServerJar, error) {
	loaderVersion, err := latestStableFabricVersion(fabricMetaBase + "/loader")
	if err != nil {
		return nil, fmt.Errorf("Couldn't find the latest version of Fabric Loader: %w", err)
	}

	installerVersion, err := latestStableFabricVersion(fabricMetaBase + "/installer")
	if err != nil {
		return nil, fmt.Errorf("Couldn't find the latest version of the Fabric Installer: %w", err)
	}

	url := fmt.Sprintf("%s/loader/%s/%s/%s/server/jar", fabricMetaBase, gameVersion, loaderVersion, installerVersion)
	return &ServerJar{
		Loader:  Fabric,
		Version: gameVersion,
		Label:   fmt.Sprintf("Fabric %s (loader %s)", gameVersion, loaderVersion),
		URL:     url,
	}, nil
}

func latestStableFabricVersion(url string) (string, error) {
	var versions []fabricVersionEntry
	if err := getJSON(url, &versions); err != nil {
		return "", err
	}
	for _, v := range versions {
		if v.Stable {
			return v.Version, nil
		}
	}
	if len(versions) > 0 {
		return versions[0].Version, nil
	}
	return "", fmt.Errorf("no version found in %s", url)
}
