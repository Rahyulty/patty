package patty

import (
	"fmt"
	"os"
	"sort"
	"strings"
)

func CmdInit() error {
	if err := EnsureManifestExists(); err != nil {
		return err
	}
	fmt.Println("Created patty.toml")
	return nil
}

func CmdAdd(arg string) error {
	m, err := LoadManifest()
	if err != nil {
		return err
	}
	name, version := parsePackageArg(arg)
	if version == "" {
		version = "latest"
	}
	m.Dependencies[name] = version
	if err := SaveManifest(m); err != nil {
		return err
	}
	fmt.Printf("Added %s@%s to patty.toml\n", name, version)
	return nil
}

func CmdInstall() error {
	if err := CheckLuarocksInstalled(); err != nil {
		return err
	}

	m, err := LoadManifest()
	if err != nil {
		return err
	}
	if len(m.Dependencies) == 0 {
		fmt.Println("No dependencies found in patty.toml")
		return nil
	}

	if err := EnsurePattyDir(); err != nil {
		return err
	}

	// Deterministic order for nicer output + lockfile stability.
	names := make([]string, 0, len(m.Dependencies))
	for k := range m.Dependencies {
		names = append(names, k)
	}
	sort.Strings(names)

	for _, name := range names {
		ver := m.Dependencies[name]
		fmt.Printf("Installing %s %s...\n", name, ver)
		if err := InstallPackage(name, ver); err != nil {
			return err
		}
	}

	// Minimal lockfile (v0.1): just the direct deps you asked for.
	lock := LockFile{
		Meta: LockMetaSection{PattyVersion: "0.1.0"},
	}
	for _, name := range names {
		lock.Packages = append(lock.Packages, LockPackageSection{
			Name:    name,
			Version: m.Dependencies[name],
			Source:  "luarocks",
		})
	}
	if err := SaveLockFile(lock); err != nil {
		return err
	}

	if err := writeLoaderFile(); err != nil {
		return err
	}

	// Suggest gitignore for .patty/
	_ = ensureGitIgnoreLine(".patty/")
	PrintPostInstallHint()
	fmt.Println("Wrote patty.lock and patty_loader.lua")
	return nil
}

func CmdRemove(name string) error {
	m, err := LoadManifest()
	if err != nil {
		return err
	}
	if _, ok := m.Dependencies[name]; !ok {
		return fmt.Errorf("dependency not found: %s", name)
	}
	delete(m.Dependencies, name)
	if err := SaveManifest(m); err != nil {
		return err
	}
	fmt.Printf("Removed %s from patty.toml\n", name)
	return nil
}

func CmdUpdate() error {
	return CmdInstall()
}

func parsePackageArg(arg string) (string, string) {
	name := arg
	version := ""
	if strings.Contains(arg, "@") {
		parts := strings.SplitN(arg, "@", 2)
		name = parts[0]
		version = parts[1]
	}
	return name, version
}

func ensureGitIgnoreLine(line string) error {
	const path = ".gitignore"
	b, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if string(b) != "" && containsLine(string(b), line) {
		return nil
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, line)
	return err
}

func containsLine(content, line string) bool {
	// Simple contains check good enough for MVP.
	return len(content) > 0 && (content == line+"\n" || (len(content) > 0 && (stringContains(content, "\n"+line+"\n") || stringContains(content, "\n"+line+"\r\n") || stringContains(content, "\n"+line))))
}

func stringContains(s, sub string) bool {
	// Avoid importing strings twice in this file.
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
