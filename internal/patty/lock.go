package patty

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

const LockPath = "patty.lock"

type LockFile struct {
	Meta     LockMetaSection      `toml:"meta"`
	Packages []LockPackageSection `toml:"packages"`
}

type LockMetaSection struct {
	PattyVersion string `toml:"patty_version"`
}

type LockPackageSection struct {
	Name    string `toml:"name"`
	Version string `toml:"version"`
	Source  string `toml:"source"`
}

func SaveLockFile(lockFile LockFile) error {
	file, err := os.Create(LockPath)
	if err != nil {
		return fmt.Errorf("cannot write to patty.lock: %w\n  Check that the file isn't read-only or locked by another program", err)
	}
	defer file.Close()

	return toml.NewEncoder(file).Encode(lockFile)
}
