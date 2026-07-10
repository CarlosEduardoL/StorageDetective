// Package command implements the ":" command mode. The user types a line, the parser walks a
// commands map, and most commands mutate ctx.Config and return to the caller.
package command

import (
	"strconv"
	"strings"

	"charm.land/bubbles/v2/help"
	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/textinput"
	"github.com/SolracHQ/stex/internal/config"
	"github.com/SolracHQ/stex/internal/core"
	"github.com/SolracHQ/stex/internal/debug"

	tea "charm.land/bubbletea/v2"
)

type cmdRun func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd)

type cmdDef struct {
	args []string
	run  cmdRun
}

var commands = map[string]cmdDef{
	"quit": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			return nil, tea.Quit
		},
	},
	"save": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			if err := ctx.Config.Save(); err != nil {
				return returnTo, core.NewNotifyCmd("Save failed", err.Error(), core.NotifyError)
			}
			return returnTo, core.NewNotifyCmd("Config saved", "", core.NotifyInfo)
		},
	},
	"sort": {
		args: []string{"ascending", "descending"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "ascending":
				ctx.Config.SortOrder = config.Ascending
				return returnTo, core.NewNotifyCmd("Sort order", "Ascending", core.NotifyInfo)
			case "descending":
				ctx.Config.SortOrder = config.Descending
				return returnTo, core.NewNotifyCmd("Sort order", "Descending", core.NotifyInfo)
			default:
				return returnTo, nil
			}
		},
	},
	"sortby": {
		args: []string{"name", "size"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "name":
				ctx.Config.SortBy = config.SortByName
				return returnTo, core.NewNotifyCmd("Sorting by", "Name", core.NotifyInfo)
			case "size":
				ctx.Config.SortBy = config.SortBySize
				return returnTo, core.NewNotifyCmd("Sorting by", "Size", core.NotifyInfo)
			default:
				return returnTo, nil
			}
		},
	},
	"group": {
		args: []string{"files", "dirs", "filesonly", "dirsonly", "mixed"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "files":
				ctx.Config.Grouping = config.FilesFirst
				return returnTo, core.NewNotifyCmd("Grouping", "Files first", core.NotifyInfo)
			case "dirs":
				ctx.Config.Grouping = config.DirsFirst
				return returnTo, core.NewNotifyCmd("Grouping", "Dirs first", core.NotifyInfo)
			case "filesonly":
				ctx.Config.Grouping = config.FilesOnly
				return returnTo, core.NewNotifyCmd("Grouping", "Files only", core.NotifyInfo)
			case "dirsonly":
				ctx.Config.Grouping = config.DirsOnly
				return returnTo, core.NewNotifyCmd("Grouping", "Dirs only", core.NotifyInfo)
			case "mixed":
				ctx.Config.Grouping = config.Mixed
				return returnTo, core.NewNotifyCmd("Grouping", "Mixed", core.NotifyInfo)
			default:
				return returnTo, nil
			}
		},
	},
	"toggle": {
		args: []string{"icons", "power", "hidden", "live"},
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			switch arg {
			case "icons":
				ctx.Config.ShowIcons = !ctx.Config.ShowIcons
				return returnTo, core.NewNotifyCmd("Icons", core.BoolLabel(ctx.Config.ShowIcons), core.NotifyInfo)
			case "power":
				ctx.Config.ShowPowerGlyphs = !ctx.Config.ShowPowerGlyphs
				return returnTo, core.NewNotifyCmd("Power glyphs", core.BoolLabel(ctx.Config.ShowPowerGlyphs), core.NotifyInfo)
			case "hidden":
				ctx.Config.ShowHidden = !ctx.Config.ShowHidden
				return returnTo, core.NewNotifyCmd("Hidden files", core.BoolLabel(ctx.Config.ShowHidden), core.NotifyInfo)
			case "live":
				ctx.Config.LiveFilter = !ctx.Config.LiveFilter
				return returnTo, core.NewNotifyCmd("Live filter", core.BoolLabel(ctx.Config.LiveFilter), core.NotifyInfo)
			default:
				return returnTo, nil
			}
		},
	},
	"up": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			levels := 1
			if arg != "" {
				parsed, err := strconv.Atoi(arg)
				if err != nil || parsed < 1 {
					return returnTo, core.NewNotifyCmd("Invalid number", "\""+arg+"\" is not a positive number", core.NotifyWarn)
				}
				levels = parsed
			}
			moved := 0
			for range levels {
				if ctx.Current.ParentDir() == nil {
					break
				}
				ctx.Current = ctx.Current.ParentDir()
				moved++
			}
			core.Rebuild(ctx)
			if moved < levels {
				return returnTo, core.NewNotifyCmd("Moved up", "Only "+strconv.Itoa(moved)+" "+core.Plural("level", moved)+", at root", core.NotifyWarn)
			}
			return returnTo, core.NewNotifyCmd("Moved up", strconv.Itoa(moved)+" "+core.Plural("level", moved), core.NotifyInfo)
		},
	},
	"debug": {
		run: func(ctx *core.Context, arg string, returnTo core.Mode) (core.Mode, tea.Cmd) {
			return debug.New(returnTo), core.NewNotifyCmd("Debug mode enabled", "Press esc to return to explorer", core.NotifyWarn)
		},
	},
}

// commandVerbs is the canonical verb list shown when no argument has been typed yet.
var commandVerbs = []string{"quit", "save", "sort", "sortby", "group", "toggle", "up", "debug"}

// Command is the command line mode. It owns a textinput widget with tab completion and a return
// target that the mode transitions back to on complete or cancel.
type Command struct {
	input    textinput.Model
	returnTo core.Mode
}

// New returns a Command bound to returnTo. The caller installs it as the active mode.
func New(returnTo core.Mode) *Command {
	input := textinput.New()
	input.Prompt = ":"
	input.Placeholder = "command"
	input.ShowSuggestions = true
	input.SetSuggestions(commandVerbs)
	input.KeyMap.AcceptSuggestion = commandKeys.AcceptSuggestion
	input.KeyMap.NextSuggestion = commandKeys.NextSuggestion
	input.KeyMap.PrevSuggestion = commandKeys.PrevSuggestion

	return &Command{input: input, returnTo: returnTo}
}

// Init focuses the textinput so the user can type the command.
func (cmd *Command) Init(_ *core.Context) tea.Cmd {
	return cmd.input.Focus()
}

// Update processes the textinput, refreshes the suggestion pool, and runs or cancels the
// command.
func (cmd *Command) Update(ctx *core.Context, msg tea.Msg) (core.Mode, tea.Cmd) {
	preValue := cmd.input.Value()
	needRefresh := false

	if keyMsg, ok := msg.(tea.KeyPressMsg); ok {
		switch {
		case key.Matches(keyMsg, commandKeys.Confirm):
			if cur := cmd.input.CurrentSuggestion(); cur != "" {
				cmd.input.SetValue(cur)
			}
			return runCommand(ctx, cmd.input.Value(), cmd.returnTo)
		case key.Matches(keyMsg, commandKeys.Cancel):
			return cmd.returnTo, nil
		}

		found := strings.Contains(preValue, " ")
		afterSpace := strings.IndexByte(cmd.input.Value(), ' ')
		if (!found && afterSpace != -1) || (found && afterSpace == -1) {
			needRefresh = true
		}
	}

	var c tea.Cmd
	cmd.input, c = cmd.input.Update(msg)

	if needRefresh {
		cmd.refreshSuggestions(cmd.input.Value())
	}

	if _, ok := msg.(tea.KeyPressMsg); !ok {
		cmd.refreshSuggestions(cmd.input.Value())
	}

	return nil, c
}

// refreshSuggestions sets the suggestion pool based on whether the current value has a space
// (verb plus arguments) or not (verb only).
func (cmd *Command) refreshSuggestions(value string) {
	before, _, ok := strings.Cut(value, " ")
	if !ok {
		cmd.input.SetSuggestions(commandVerbs)
		return
	}

	verb := before
	def, exists := commands[verb]
	if !exists || len(def.args) == 0 {
		cmd.input.SetSuggestions(nil)
		return
	}

	full := make([]string, len(def.args))
	for i, arg := range def.args {
		full[i] = verb + " " + arg
	}
	cmd.input.SetSuggestions(full)
}

// Help returns the command key bindings for the help footer.
func (cmd *Command) Help() help.KeyMap { return commandKeys }
func (cmd *Command) Name() string      { return "command" }

func runCommand(ctx *core.Context, value string, returnTo core.Mode) (core.Mode, tea.Cmd) {
	parts := splitFields(value)
	if len(parts) == 0 {
		return returnTo, nil
	}

	verb, arg := parts[0], ""
	if len(parts) > 1 {
		arg = parts[1]
	}

	def, exists := commands[verb]
	if !exists {
		detail := "\"" + verb + "\" is not a valid command"
		return returnTo, core.NewNotifyCmd("Command not executed", detail, core.NotifyWarn)
	}

	if len(def.args) > 0 && arg == "" {
		detail := strings.Join(def.args, ", ")
		return returnTo, core.NewNotifyCmd("\""+verb+"\" requires an argument", detail, core.NotifyWarn)
	}

	return def.run(ctx, arg, returnTo)
}

func splitFields(value string) []string {
	var out []string
	var cur strings.Builder
	for _, r := range value {
		if r == ' ' || r == '\t' {
			if cur.Len() > 0 {
				out = append(out, cur.String())
				cur.Reset()
			}
			continue
		}
		cur.WriteRune(r)
	}
	if cur.Len() > 0 {
		out = append(out, cur.String())
	}
	return out
}
