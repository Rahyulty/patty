package patty

import (
	"os"

	"github.com/BurntSushi/toml"
)

const LockPath = "patty.lock"

type LockFile struct {
	Meta LockMetaSection `toml:"meta"`
	Packages []LockPackageSection `toml:"packages"`
}

type LockMetaSection struct {
	PattyVersion string `toml:"patty_version"`
}

type LockPackageSection struct {
	Name string `toml:"name"`
	Version string `toml:"version"`
	Source string `toml:"source"`
}

func SaveLockFile(lockFile LockFile) error {
	file, err := os.Create(LockPath)

	if err != nil {
		return err
	}

	defer file.Close()
	// i dont know what .encode(1) means but its suggested by chatgpt
	return toml.NewEncoder(file).Encode(1)
}
