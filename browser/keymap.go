package browser

import "github.com/charmbracelet/bubbles/key"

type keyMap struct {
	Up        key.Binding
	Down      key.Binding
	Left      key.Binding
	Right     key.Binding
	Select    key.Binding
	Home      key.Binding
	End       key.Binding
	Sync      key.Binding
	Backspace key.Binding
	Root      key.Binding
	Help      key.Binding
	Quit      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Up, k.Down, k.Left, k.Right, k.Select},
		{k.Home, k.End, k.Backspace, k.Root},
		{k.Sync, k.Help, k.Quit},
	}
}

var keys = keyMap{
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "next"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "prev"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "next page"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "prev page"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter", " "),
		key.WithHelp("enter/space", "select"),
	),
	Home: key.NewBinding(
		key.WithKeys("g", "home"),
		key.WithHelp("g/home", "go to start"),
	),
	End: key.NewBinding(
		key.WithKeys("G", "end"),
		key.WithHelp("G/end", "go to end"),
	),
	Backspace: key.NewBinding(
		key.WithKeys("backspace"),
		key.WithHelp("backspace", "cd .."),
	),
	Root: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "cd /"),
	),
	Sync: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "synchronize"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}
