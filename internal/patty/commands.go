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

	if err := writeLoaderFile(); err != nil {
		return fmt.Errorf("failed to create patty_loader.lua: %w\n  This file is needed so Lua can find packages installed by patty", err)
	}
	fmt.Println("Created patty_loader.lua")

	lock := LockFile{
		Meta:     LockMetaSection{PattyVersion: "0.1.0"},
		Packages: []LockPackageSection{},
	}
	if err := SaveLockFile(lock); err != nil {
		return fmt.Errorf("failed to create patty.lock: %w\n  The lockfile pins your dependency versions for reproducible builds", err)
	}
	fmt.Println("Created patty.lock")

	_ = ensureGitIgnoreLine(".patty/")
	fmt.Println("Updated .gitignore")

	fmt.Println("\nProject initialized! Run 'patty install <package>' to add dependencies.")
	return nil
}

// CmdInstall handles both:
//   - patty install <pkg>    → add package to manifest, then install everything
//   - patty install          → install everything from manifest
func CmdInstall(pkgArgs []string) error {
	if err := CheckLuarocksInstalled(); err != nil {
		return err
	}

	m, err := LoadManifest()
	if err != nil {
		return err
	}

	for _, arg := range pkgArgs {
		name, version := parsePackageArg(arg)
		if name == "" {
			return fmt.Errorf("invalid package name in '%s'\n  Usage: patty install <name> or patty install <name>@<version>", arg)
		}
		if version == "" {
			version = "latest"
		}
		m.Dependencies[name] = version
		fmt.Printf("Added %s@%s to patty.toml\n", name, version)
	}

	if len(pkgArgs) > 0 {
		if err := SaveManifest(m); err != nil {
			return fmt.Errorf("failed to save patty.toml: %w\n  Check that the file isn't open in another program or read-only", err)
		}
	}

	if len(m.Dependencies) == 0 {
		fmt.Println("No dependencies to install.")
		fmt.Println("  Add packages with: patty install <package>")
		fmt.Println("  Example: patty install luafilesystem")
		return nil
	}

	if err := EnsurePattyDir(); err != nil {
		return fmt.Errorf("failed to create .patty directory: %w\n  Patty stores installed packages in .patty/ — check folder permissions", err)
	}

	names := make([]string, 0, len(m.Dependencies))
	for k := range m.Dependencies {
		names = append(names, k)
	}
	sort.Strings(names)

	fmt.Printf("\nInstalling %d package(s)...\n\n", len(names))

	for _, name := range names {
		ver := m.Dependencies[name]
		spinner := NewSpinner(fmt.Sprintf("Installing %s@%s", name, ver))
		spinner.Start()

		err := InstallPackage(name, ver)
		spinner.Stop()

		if err != nil {
			fmt.Printf("✗ Failed %s@%s\n", name, ver)
			return fmt.Errorf("could not install '%s': %w\n\n  Possible causes:\n  - Package name is misspelled (check https://luarocks.org)\n  - Version '%s' does not exist for this package\n  - Network issue — check your internet connection\n  - Native module that failed to compile (see error above)", name, err, ver)
		}

		fmt.Printf("✓ Installed %s@%s\n", name, ver)
	}

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
		return fmt.Errorf("failed to write patty.lock: %w\n  Packages were installed but the lockfile couldn't be saved", err)
	}

	if err := writeLoaderFile(); err != nil {
		return fmt.Errorf("failed to write patty_loader.lua: %w\n  Packages were installed but the loader couldn't be generated", err)
	}

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
		available := make([]string, 0, len(m.Dependencies))
		for k := range m.Dependencies {
			available = append(available, k)
		}
		sort.Strings(available)

		msg := fmt.Sprintf("package '%s' is not in your dependencies", name)
		if len(available) > 0 {
			msg += fmt.Sprintf("\n  Installed packages: %s", strings.Join(available, ", "))
		} else {
			msg += "\n  You have no dependencies installed"
		}
		msg += "\n  Did you mean a different package name?"
		return fmt.Errorf("%s", msg)
	}

	delete(m.Dependencies, name)
	if err := SaveManifest(m); err != nil {
		return fmt.Errorf("failed to save patty.toml after removing '%s': %w", name, err)
	}

	fmt.Printf("✓ Removed %s from patty.toml\n", name)
	fmt.Println("  Run 'patty install' to update your .patty/ directory")
	return nil
}

func CmdUpdate() error {
	fmt.Println("Reinstalling all dependencies...")
	return CmdInstall(nil)
}

func parsePackageArg(arg string) (string, string) {
	if i := strings.Index(arg, "@"); i != -1 {
		return arg[:i], arg[i+1:]
	}
	return arg, ""
}

func ensureGitIgnoreLine(line string) error {
	const path = ".gitignore"

	b, err := os.ReadFile(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if strings.Contains(string(b), line) {
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
