package loaders

import "fmt"

var vanillaManifestURL = "https://launchermeta.mojang.com/mc/game/version_manifest_v2.json"

type vanillaManifest struct {
	Latest struct {
		Release  string `json:"release"`
		Snapshot string `json:"snapshot"`
	} `json:"latest"`
	Versions []struct {
		ID   string `json:"id"`
		Type string `json:"type"`
		URL  string `json:"url"`
	} `json:"versions"`
}

type vanillaVersionMeta struct {
	Downloads struct {
		Server struct {
			URL string `json:"url"`
		} `json:"server"`
	} `json:"downloads"`
}

func VanillaReleases(limit int) ([]string, error) {
	var manifest vanillaManifest
	if err := getJSON(vanillaManifestURL, &manifest); err != nil {
		return nil, fmt.Errorf("Failed to fetch the list of Minecraft versions from Mojang: %w", err)
	}

	releases := make([]string, 0, limit)
	for _, v := range manifest.Versions {
		if v.Type != "release" {
			continue
		}
		releases = append(releases, v.ID)
		if len(releases) >= limit {
			break
		}
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("Mojang hasn't released any version")
	}
	return releases, nil
}

func VanillaServerJar(version string) (*ServerJar, error) {
	var manifest vanillaManifest
	if err := getJSON(vanillaManifestURL, &manifest); err != nil {
		return nil, fmt.Errorf("Unable to retrieve the Minecraft version manifest: %w", err)
	}

	resolved := version
	if version == "latest" {
		resolved = manifest.Latest.Release
	}

	var versionURL string
	for _, v := range manifest.Versions {
		if v.ID == resolved {
			versionURL = v.URL
			break
		}
	}
	if versionURL == "" {
		return nil, fmt.Errorf("Version %q not found in the official Minecraft manifest", resolved)
	}

	var meta vanillaVersionMeta
	if err := getJSON(versionURL, &meta); err != nil {
		return nil, fmt.Errorf("Could not retrieve metadata for version %s: %w", resolved, err)
	}
	if meta.Downloads.Server.URL == "" {
		return nil, fmt.Errorf("Version %s does not have a server.jar available (it might be too old)", resolved)
	}

	return &ServerJar{
		Loader:  Vanilla,
		Version: resolved,
		Label:   fmt.Sprintf("Vanilla %s", resolved),
		URL:     meta.Downloads.Server.URL,
	}, nil
}
