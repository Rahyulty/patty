package patty

import (
	"errors"
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

	// decode the manifest file into the manifest structure
	if _, err := toml.DecodeFile(ManifestPath, &m); err != nil {
		return Manifest{}, err
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
		return err
	}

	// close the file
	defer file.Close()

	// create a new encoder for the file
	enc := toml.NewEncoder(file)
	return enc.Encode(m)
}

func EnsureManifestExists() error {
	if _, err := os.Stat(ManifestPath); err == nil {
		return errors.New("manifest file already exists")
	}

	m := Manifest {
		Package: PackageSection{
			Name: "my-project",
			Version: "0.1.0",
			Lua: ">=5.1",
		},
		Dependencies: map[string]string{},
		DevDependencies: map[string]string{},
	}
	
	return SaveManifest(m)
}
