# stex

stex (Storage Explorer) is a TUI app to explore disk usage. It is heavily
inspired by ncdu and nvim (in the mode architecture). It aims to be simple to
use, visually pleasant and multiplatform.

## Install

```shell
go install github.com/SolracHQ/stex@latest
```

Or clone and build:

```shell
git clone https://github.com/SolracHQ/stex.git
cd stex
go build -o stex .
```

## Usage

```shell
stex [flags] [path]
```

Defaults to the current directory when no path is given.

### Flags

| Flag | Short | Description |
| --- | --- | --- |
| `--icons` | `-i` | start with emoji icons enabled |
| `--power-glyphs` | `-g` | start with powerline glyphs in the power bar |
| `--show-all` | `-a` | start with hidden files shown |
| `--no-live-filter` | `-L` | disable live filter (compile on enter) |
| `--help` | `-h` | show this help |

## What makes it different?

stex is built around the concept of modes. Each mode changes the app behavior,
the keybindings and the overlay, reducing the amount of keybindings per mode.

### Modes

#### Explorer Mode

Explorer mode is the main app mode. It allows you to navigate the filesystem
subtree and select items to act over. It is focused completely on movement,
allowing three kinds of movement sets, hjkl for vim users, arrows for
traditionalists and awsd for gamers. On wide terminals a right panel shows
file metadata like name, size, permissions, mod time and MIME type plus the
largest children when viewing directories.

See the full keybinding list in [docs/explorer_keymap.md](docs/explorer_keymap.md).

#### Filter Mode

Filter mode lets you search the current directory by regex. Press `/` to open a
search bar at the bottom, all keystrokes are directed to the input and normal
navigation is suspended. The list narrows in real time as you type. On slower
machines you can disable live filter with the `--no-live-filter` flag or toggle
it at runtime with `ctrl+l`, in manual mode the filter only compiles on enter.

See the full keybinding list in [docs/filter_keymap.md](docs/filter_keymap.md).

#### Settings Panel

Settings panel is the UI to customize the app behavior. You can toggle sort,
order, grouping, icons, hidden files, live filter, notify level and notify
timeout. It is the same configuration you can do via CLI flags or command
mode, just in a friendly menu.

See the full keybinding list in [docs/settings_keymap.md](docs/settings_keymap.md).

#### Command Mode

Command mode is the command centric way to configure the app, allowing
everything the settings panel does but purely through commands, like a command
palette in any code editor. The idea is for it to grow into a more powerful
tool to navigate and act on the filesystem without leaving the keyboard. I even
borrowed `:q` to quit, old habits die hard. Every command shows a notification
confirming what it did, if you are using stex in a dumb terminal you can
disable or lower the verbosity in the settings panel.

See the full command list in [docs/commands.md](docs/commands.md) and the
keybinding list in [docs/command_keymap.md](docs/command_keymap.md).

### Notifications

stex shows a small box at the top right when something happens, like a command
executed, a setting changed, an error or simply a confirmation that the config
was saved. The notification auto dismisses after a configurable timeout and you
can also dismiss it immediately with esc. The verbosity level can be set in the
settings panel or config file from all (default) to off.

## Configuration

stex reads config from `$XDG_CONFIG/stex/config.json` when it exists but CLI
flags take precedence over the file and the file takes precedence over
defaults. That way you can set your preferred defaults once in the config file
and override per context with shell aliases like `alias sxh='stex -a'` for
checking disk usage in your home directory where `.local` is usually the
culprit.

## Architecture

stex uses a mode architecture that is in essence a finite state machine of
states. Modes control three things: the overlay showing at that moment, the
event processing (the app forwards non global events to the mode), and the
mode transitions, any mode can transition to any other.

The main app manages the long running tasks like notifications, the app wise
events like global help toggle or quit, and the base view composition with the
mode overlay on top. The communication between modes and the app is through
commands and messages using the standard bubbletea mechanism. Modes handle
sync operations and mutate the base state, they delegate async operations to
the app. Each mode has its own keybindings and the app is in charge of
updating the help widget with the current mode bindings alongside the global
ones.

## License

MIT. See [LICENSE](LICENSE).
