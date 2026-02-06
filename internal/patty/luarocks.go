package patty

import (
	"errors"
	"os/exec"
)

func CheckLuarocksInstalled() error {
	_, err := exec.LookPath("luarocks")
	if err != nil {
		return errors.New("luarocks is not found in PATH (install lua rocks first to use patty)")
	}
	return nil
}

func InstallPackage(pkg string, version string) error {
	cmd := exec.Command("luarocks", "install", pkg, version, "--tree=.patty")
	cmd.Stdout = Stdout()
	cmd.Stderr = Stderr()
	return cmd.Run()
}
