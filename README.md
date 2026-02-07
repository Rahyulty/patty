# Patty

> ⚠️ **Work in Progress** — Patty is in early development (v0.1.0 MVP). Features may change, break, or be incomplete. Contributions and feedback are welcome!

A modern dependency manager for Lua. Think npm, but for Lua projects.

Patty wraps [LuaRocks](https://luarocks.org) to give you a simple, familiar workflow for managing Lua packages with a local project tree — no global installs, no permission headaches.

## Features

- `patty install <package>` — adds and installs in one step (like `npm install`)
- Local `.patty/` directory for project-scoped dependencies
- `patty.toml` manifest for declaring dependencies
- `patty.lock` lockfile for reproducible builds
- Auto-generated `patty_loader.lua` so Lua can find your packages
- Automatic `.gitignore` management
- Windows support with automatic Visual Studio Build Tools detection

## Prerequisites

- **Lua** (5.1+)
- **LuaRocks** — Patty uses LuaRocks under the hood to fetch and build packages
- **Windows only:** Visual Studio Build Tools with C++ workload (needed for native C modules)

### Installing prerequisites

**Windows:**
```
choco install lua luarocks
```

**macOS:**
```
brew install lua luarocks
```

**Linux (Debian/Ubuntu):**
```
sudo apt install lua5.4 luarocks
```

## Getting Started

### 1. Initialize a project

```
patty init
```

This creates:
- `patty.toml` — your project manifest
- `patty_loader.lua` — auto-generated loader for your packages
- `patty.lock` — lockfile (empty until you install packages)
- `.gitignore` — adds `.patty/` so installed packages aren't committed

### 2. Install a package

```
patty install luafilesystem
```

This adds `luafilesystem` to `patty.toml` and installs it into `.patty/`.

You can also install multiple packages at once:

```
patty install lua-json luasocket
```

Or pin a specific version:

```
patty install lua-json@1.0
```

### 3. Use your packages

Add this line to the top of your Lua entrypoint:

```lua
require("patty_loader")
```

Then use your packages normally:

```lua
require("patty_loader")
local lfs = require("lfs")

print(lfs.currentdir())
```

### 4. Reinstall from manifest

If you clone a project or want to reinstall everything from `patty.toml`:

```
patty install
```

## Commands

| Command | Description |
|---|---|
| `patty init` | Initialize a new patty project |
| `patty install [packages...]` | Install packages (adds to patty.toml and installs) |
| `patty remove <package>` | Remove a dependency from the project |
| `patty update` | Reinstall all dependencies |
| `patty help` | Show help |

### Aliases

| Alias | Command |
|---|---|
| `patty i` | `patty install` |
| `patty rm` | `patty remove` |

## Project Structure

After running `patty init` and installing packages, your project will look like:

```
my-project/
├── patty.toml          # your dependencies
├── patty.lock          # pinned versions
├── patty_loader.lua    # auto-generated, tells Lua where to find packages
├── .patty/             # installed packages (gitignored)
├── .gitignore
└── main.lua            # your code
```

## patty.toml

The manifest file where your project info and dependencies live:

```toml
[package]
name = "my-project"
version = "0.1.0"
lua = ">=5.1"

[dependencies]
luafilesystem = "latest"
lua-json = "1.0"

[dev_dependencies]
```

## Building from Source

Patty is written in Go. To build:

```
git clone https://github.com/Rahyulty/patty.git
cd patty
go build ./cmd/patty
```

## Known Limitations

- **MVP stage** — core functionality works but edge cases may not be handled
- Lockfile only tracks direct dependencies (no dependency tree resolution yet)
- `patty remove` updates `patty.toml` but doesn't delete files from `.patty/` (run `patty install` after to rebuild)
- Version resolution is basic — `latest` fetches whatever LuaRocks considers latest
- Windows requires Visual Studio Build Tools for packages with native C code

## License

MIT
