package patty

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func CheckLuarocksInstalled() error {
	_, err := exec.LookPath("luarocks")
	if err != nil {
		return errors.New("luarocks was not found on your system\n\n" +
			"  Patty uses luarocks to fetch Lua packages. To install it:\n" +
			"    Windows:  choco install luarocks\n" +
			"    macOS:    brew install luarocks\n" +
			"    Linux:    sudo apt install luarocks (or your package manager)\n\n" +
			"  After installing, restart your terminal and try again")
	}
	return nil
}

func InstallPackage(pkg string, version string) error {
	var args []string
	if version == "" || version == "latest" {
		args = []string{"install", pkg, "--tree=.patty", "--local"}
	} else {
		args = []string{"install", pkg, version, "--tree=.patty", "--local"}
	}

	// on Windows, native C modules need cl.exe from Visual Studio Build Tools
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("cl"); err != nil {
			vcvars := findVcvarsall()
			if vcvars == "" {
				return errors.New("this package contains native C code and needs a C compiler to build\n\n" +
					"  Visual Studio Build Tools with C++ workload is required.\n" +
					"  Install it by running (in an admin terminal):\n\n" +
					"    winget install Microsoft.VisualStudio.2022.BuildTools --override \"--quiet --add Microsoft.VisualStudio.Workload.VCTools --includeRecommended\"\n\n" +
					"  After installing, restart your terminal and try again")
			}
			// run vcvarsall.bat to set up compiler environment, then luarocks
			// write a temp .bat file to avoid Go escaping quotes in cmd /C args
			luarocksCmd := fmt.Sprintf("luarocks %s", joinArgs(args))
			arch := detectLuaArch()
			// suppress compiler output (cl.exe, link.exe messages) but keep errors
			script := fmt.Sprintf("@call \"%s\" %s >nul 2>&1\n@if errorlevel 1 exit /b 1\n@%s 2>&1\n", vcvars, arch, luarocksCmd)

			tmp, err := os.CreateTemp("", "patty-*.bat")
			if err != nil {
				return fmt.Errorf("failed to create temp script: %w", err)
			}
			defer os.Remove(tmp.Name())

			if _, err := tmp.WriteString(script); err != nil {
				tmp.Close()
				return fmt.Errorf("failed to write temp script: %w", err)
			}
			tmp.Close()

			cmd := exec.Command("cmd", "/C", tmp.Name())
			output, err := cmd.CombinedOutput()
			if err != nil {
				// filter output to show only errors, not verbose compiler messages
				filtered := filterCompilerOutput(string(output))
				if filtered != "" {
					fmt.Fprint(os.Stderr, filtered)
				}
				return err
			}
			return nil
		}
	}

	// cl.exe available (or not Windows), run luarocks directly
	cmd := exec.Command("luarocks", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		filtered := filterCompilerOutput(string(output))
		if filtered != "" {
			fmt.Fprint(os.Stderr, filtered)
		}
		return err
	}
	return nil
}

func findVcvarsall() string {
	roots := []string{
		`C:\Program Files (x86)\Microsoft Visual Studio`,
		`C:\Program Files\Microsoft Visual Studio`,
	}
	years := []string{"2022", "2019", "2017"}
	editions := []string{"BuildTools", "Community", "Professional", "Enterprise"}

	for _, root := range roots {
		for _, year := range years {
			for _, edition := range editions {
				p := filepath.Join(root, year, edition, "VC", "Auxiliary", "Build", "vcvarsall.bat")
				if _, err := os.Stat(p); err == nil {
					return p
				}
			}
		}
	}
	return ""
}

func joinArgs(args []string) string {
	result := ""
	for i, arg := range args {
		if i > 0 {
			result += " "
		}
		if containsSpace(arg) {
			result += fmt.Sprintf(`"%s"`, arg)
		} else {
			result += arg
		}
	}
	return result
}

func containsSpace(s string) bool {
	for _, r := range s {
		if r == ' ' {
			return true
		}
	}
	return false
}

// detectLuaArch figures out if the luarocks Lua is 32-bit or 64-bit.
// checks the luarocks executable path and Lua install path for clues.
func detectLuaArch() string {
	// find where luarocks lives
	luarocksPath, _ := exec.LookPath("luarocks")

	// find where lua lives
	luaPath, _ := exec.LookPath("lua")

	// check both paths for 32-bit indicators
	for _, p := range []string{luarocksPath, luaPath} {
		if p == "" {
			continue
		}
		lower := toLower(p)
		if containsStr(lower, "win32") ||
			containsStr(lower, "x86") ||
			containsStr(lower, "program files (x86)") ||
			containsStr(lower, "32bit") {
			return "x86"
		}
	}

	return "x64"
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		b[i] = c
	}
	return string(b)
}

func containsStr(s, sub string) bool {
	if len(sub) > len(s) {
		return false
	}
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

// filterCompilerOutput removes verbose compiler messages but keeps errors and important info
func filterCompilerOutput(output string) string {
	var filtered []string
	scanner := bufio.NewScanner(strings.NewReader(output))

	for scanner.Scan() {
		line := scanner.Text()
		lower := strings.ToLower(line)

		// skip verbose compiler output
		if strings.Contains(lower, "microsoft (r) incremental linker") ||
			strings.Contains(lower, "copyright (c) microsoft corporation") ||
			strings.Contains(lower, "creating library") ||
			strings.Contains(lower, "cl /nologo") ||
			strings.Contains(lower, "link -dll") ||
			strings.Contains(lower, "no existing manifest") ||
			strings.Contains(lower, "visual studio") ||
			strings.Contains(lower, "vcvarsall.bat") ||
			strings.Contains(lower, "environment initialized") ||
			strings.HasPrefix(strings.TrimSpace(line), "**") ||
			strings.TrimSpace(line) == "" {
			continue
		}

		// keep errors and important messages
		if strings.Contains(lower, "error") ||
			strings.Contains(lower, "warning") ||
			strings.Contains(lower, "installing") ||
			strings.Contains(lower, "is now installed") ||
			strings.Contains(lower, "failed") {
			filtered = append(filtered, line)
		}
	}

	return strings.Join(filtered, "\n")
}
