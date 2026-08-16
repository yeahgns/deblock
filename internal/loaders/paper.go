package loaders

import "fmt"

var paperAPIBase = "https://api.papermc.io/v2/projects/paper"

type paperProject struct {
	Versions []string `json:"versions"`
}

type paperVersion struct {
	Builds []int `json:"builds"`
}

type paperBuild struct {
	Build     int `json:"build"`
	Downloads struct {
		Application struct {
			Name string `json:"name"`
		} `json:"application"`
	} `json:"downloads"`
}

func PaperReleases(limit int) ([]string, error) {
	var project paperProject
	if err := getJSON(paperAPIBase, &project); err != nil {
		return nil, fmt.Errorf("Couldn't retrieve the list of Paper versions: %w", err)
	}
	if len(project.Versions) == 0 {
		return nil, fmt.Errorf("PaperAPI did not return any version.")
	}

	versions := project.Versions
	if len(versions) > limit {
		versions = versions[len(versions)-limit:]
	}
	reversed := make([]string, len(versions))
	for i, v := range versions {
		reversed[len(versions)-1-i] = v
	}
	return reversed, nil
}

func PaperServerJar(version string) (*ServerJar, error) {
	versionURL := fmt.Sprintf("%s/versions/%s", paperAPIBase, version)
	var v paperVersion
	if err := getJSON(versionURL, &v); err != nil {
		return nil, fmt.Errorf("Couldn't fetch the Paper builds for the version %s: %w", version, err)
	}
	if len(v.Builds) == 0 {
		return nil, fmt.Errorf("There are no Paper builds available for the version %s", version)
	}

	latestBuild := v.Builds[len(v.Builds)-1]
	buildURL := fmt.Sprintf("%s/builds/%d", versionURL, latestBuild)
	var b paperBuild
	if err := getJSON(buildURL, &b); err != nil {
		return nil, fmt.Errorf("failed to retrieve details for Paper build %d: %w", latestBuild, err)
	}
	if b.Downloads.Application.Name == "" {
		return nil, fmt.Errorf("Paper build %d does not have an application JAR available", latestBuild)
	}

	downloadURL := fmt.Sprintf("%s/downloads/%s", buildURL, b.Downloads.Application.Name)
	return &ServerJar{
		Loader:  Paper,
		Version: version,
		Label:   fmt.Sprintf("Paper %s (build %d)", version, latestBuild),
		URL:     downloadURL,
	}, nil
}
