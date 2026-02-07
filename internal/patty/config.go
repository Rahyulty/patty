package patty

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const ManifestPath = "patty.toml"

type Manifest struct {
	Package         PackageSection    `toml:"package"`
	Dependencies    map[string]string `toml:"dependencies"`
	DevDependencies map[string]string `toml:"dev_dependencies"`
}

type PackageSection struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Lua     string `toml:"lua"`
}

func LoadManifest() (Manifest, error) {
	var m Manifest

	if _, err := os.Stat(ManifestPath); os.IsNotExist(err) {
		return Manifest{}, fmt.Errorf("patty.toml not found in the current directory\n  Run 'patty init' to create a new project first")
	}

	if _, err := toml.DecodeFile(ManifestPath, &m); err != nil {
		return Manifest{}, fmt.Errorf("failed to read patty.toml: %w\n  Check that patty.toml has valid TOML syntax", err)
	}

	if m.Dependencies == nil {
		m.Dependencies = map[string]string{}
	}

	if m.DevDependencies == nil {
		m.DevDependencies = map[string]string{}
	}

	return m, nil
}

func SaveManifest(m Manifest) error {
	file, err := os.Create(ManifestPath)
	if err != nil {
		return fmt.Errorf("cannot write to patty.toml: %w\n  Check that the file isn't read-only or locked by another program", err)
	}
	defer file.Close()

	enc := toml.NewEncoder(file)
	return enc.Encode(m)
}

func EnsureManifestExists() error {
	if _, err := os.Stat(ManifestPath); err == nil {
		return fmt.Errorf("a patty project already exists in this directory (patty.toml found)\n  Delete patty.toml if you want to start fresh")
	}

	m := Manifest{
		Package: PackageSection{
			Name:    "my-project",
			Version: "0.1.0",
			Lua:     ">=5.1",
		},
		Dependencies:    map[string]string{},
		DevDependencies: map[string]string{},
	}

	return SaveManifest(m)
}
