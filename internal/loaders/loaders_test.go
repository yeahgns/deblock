package loaders

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestVanillaParsing(t *testing.T) {
	manifestJSON := `{
		"latest": {"release": "1.21.1", "snapshot": "24w40a"},
		"versions": [
			{"id": "1.21.1", "type": "release", "url": "%s/v/1.21.1.json"},
			{"id": "24w40a", "type": "snapshot", "url": "%s/v/24w40a.json"},
			{"id": "1.21", "type": "release", "url": "%s/v/1.21.json"}
		]
	}`
	versionJSON := `{"downloads": {"server": {"url": "https://piston-data.mojang.com/v1/objects/abc/server.jar"}}}`

	var srv *httptest.Server
	mux := http.NewServeMux()
	mux.HandleFunc("/manifest.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strings.ReplaceAll(manifestJSON, "%s", srv.URL)))
	})
	mux.HandleFunc("/v/1.21.1.json", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(versionJSON))
	})
	srv = httptest.NewServer(mux)
	defer srv.Close()

	original := vanillaManifestURL
	vanillaManifestURL = srv.URL + "/manifest.json"
	defer func() { vanillaManifestURL = original }()

	releases, err := VanillaReleases(5)
	if err != nil {
		t.Fatalf("VanillaReleases failed: %v", err)
	}
	if len(releases) != 2 || releases[0] != "1.21.1" || releases[1] != "1.21" {
		t.Fatalf("unexpected releases: %v", releases)
	}

	jar, err := VanillaServerJar("latest")
	if err != nil {
		t.Fatalf("VanillaServerJar failed: %v", err)
	}
	if jar.Version != "1.21.1" {
		t.Fatalf("incorrectly resolved version: %s", jar.Version)
	}
	if jar.URL != "https://piston-data.mojang.com/v1/objects/abc/server.jar" {
		t.Fatalf("incorrect download URL: %s", jar.URL)
	}
}

func TestPaperParsing(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"versions": ["1.20.6", "1.21", "1.21.1"]}`))
	})
	mux.HandleFunc("/versions/1.21.1", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"builds": [130, 131, 132]}`))
	})
	mux.HandleFunc("/versions/1.21.1/builds/132", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"build": 132, "downloads": {"application": {"name": "paper-1.21.1-132.jar"}}}`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	original := paperAPIBase
	paperAPIBase = srv.URL
	defer func() { paperAPIBase = original }()

	releases, err := PaperReleases(2)
	if err != nil {
		t.Fatalf("PaperReleases failed: %v", err)
	}
	if len(releases) != 2 || releases[0] != "1.21.1" || releases[1] != "1.21" {
		t.Fatalf("unexpected releases: %v", releases)
	}

	jar, err := PaperServerJar("1.21.1")
	if err != nil {
		t.Fatalf("PaperServerJar failed: %v", err)
	}
	want := srv.URL + "/versions/1.21.1/builds/132/downloads/paper-1.21.1-132.jar"
	if jar.URL != want {
		t.Fatalf("incorrect download URL:\n got: %s\nwant: %s", jar.URL, want)
	}
}

func TestFabricParsing(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/game", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"version": "1.21.1", "stable": true}, {"version": "24w40a", "stable": false}, {"version": "1.21", "stable": true}]`))
	})
	mux.HandleFunc("/loader", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"version": "0.16.9", "stable": true}]`))
	})
	mux.HandleFunc("/installer", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`[{"version": "1.0.1", "stable": true}]`))
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	original := fabricMetaBase
	fabricMetaBase = srv.URL
	defer func() { fabricMetaBase = original }()

	releases, err := FabricReleases(5)
	if err != nil {
		t.Fatalf("FabricReleases falhou: %v", err)
	}
	if len(releases) != 2 || releases[0] != "1.21.1" || releases[1] != "1.21" {
		t.Fatalf("unexpected releases (snapshot should have been filtered out): %v", releases)
	}

	jar, err := FabricServerJar("1.21.1")
	if err != nil {
		t.Fatalf("FabricServerJar falhou: %v", err)
	}
	want := srv.URL + "/loader/1.21.1/0.16.9/1.0.1/server/jar"
	if jar.URL != want {
		t.Fatalf("incorrect download URL:\n got: %s\nwant: %s", jar.URL, want)
	}
}
