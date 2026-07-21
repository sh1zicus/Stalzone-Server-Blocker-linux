package tui

import (
	"log/slog"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
	"github.com/sh1zicus/stalzone-server-blocker/internal/ping"
	"github.com/sh1zicus/stalzone-server-blocker/internal/nft"
	"github.com/sh1zicus/stalzone-server-blocker/internal/config"
)

type statusClearMsg struct{}
type tickProcessMsg struct{}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		func() tea.Msg {
			ping.Refresh(m.state.Pools)
			return nil
		},
		tickProcess(),
	)
}

func tickProcess() tea.Cmd {
	return tea.Tick(2*time.Second, func(_ time.Time) tea.Msg {
		return tickProcessMsg{}
	})
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

		case statusClearMsg:
			m.status = ""

		case tickProcessMsg:
			m.stalzoneRunning = isStalzoneRunning()
			return m, tickProcess()

		case tea.WindowSizeMsg:

			m.width = msg.Width
			m.height = msg.Height

			if m.viewport.Width == 0 {

				m.viewport = viewport.New(msg.Width, msg.Height-7)

			} else {

				m.viewport.Width = msg.Width
				m.viewport.Height = msg.Height - 7
			}

		case tea.KeyMsg:

			if m.searchMode {

				switch msg.String() {

					case "esc":

						m.searchMode = false
						m.search = ""

						m.state.Filter = ""
						m.state.Rebuild()

					case "enter":

						m.searchMode = false

					case "backspace":

						if len(m.search) > 0 {
							m.search = m.search[:len(m.search)-1]
						}

						m.state.Filter = m.search
						m.state.Rebuild()

					default:

						if len(msg.Runes) == 1 {

							m.search += string(msg.Runes)

							m.state.Filter = m.search
							m.state.Rebuild()
						}
				}

				m.viewport.SetContent(renderList(m))

				return m, nil
			}


			switch msg.String() {

				case "ctrl+c", "q":
					return m, tea.Quit

				case "up":
					m.state.MoveUp()

				case "down":
					m.state.MoveDown()

				case "left":

					row := m.state.Current()

					if row != nil && row.Type == model.RowPool {

						if m.state.Expanded[row.Pool] {

							m.state.Expanded[row.Pool] = false
							m.state.Rebuild()
						}
					}

				case "right":

					row := m.state.Current()

					if row != nil && row.Type == model.RowPool {

						if !m.state.Expanded[row.Pool] {

							m.state.Expanded[row.Pool] = true
							m.state.Rebuild()
						}
					}

				case " ":
					m.state.Toggle()
					config.UpdateSelection(m.cfg, m.state.Pools)
					_ = config.Save(m.cfg)

				case "enter":
					m.state.TogglePool()

				case "pgdown":
					m.viewport.LineDown(10)

				case "pgup":
					m.viewport.LineUp(10)

				case "/":
					m.searchMode = true

				case "r":
					ping.Refresh(m.state.Pools)

				case "a":
					config.UpdateSelection(m.cfg, m.state.Pools)
					_ = config.Save(m.cfg)

					if err := nft.Apply(m.state.Pools); err != nil {
						m.log.Error("apply nft", "err", err)
						m.status = "Ошибка: " + err.Error()
					} else {
						m.status = "Правила применены"
					}
					m.viewport.SetContent(renderList(m))
					return m, clearStatusAfter(3 * time.Second)

				case "d":
					if err := nft.Reset(); err != nil {
						m.log.Error("reset nft", "err", err)
						m.status = "Ошибка: " + err.Error()
					} else {
						m.status = "Правила сброшены"
					}
					m.viewport.SetContent(renderList(m))
					return m, clearStatusAfter(3 * time.Second)
			}
	}

	m.viewport.SetContent(renderList(m))

	cursor := m.state.Cursor

	top := m.viewport.YOffset
	bottom := top + m.viewport.Height - 1

	if cursor < top {
		m.viewport.SetYOffset(cursor)
	}

	if cursor > bottom {
		m.viewport.SetYOffset(cursor - m.viewport.Height + 1)
	}

	return m, nil
}

func Run(cfg *model.Config, log *slog.Logger, pools []model.Pool) error {

	m := Model{
		cfg: cfg,
		log: log,
		state: model.NewState(pools),
	}

	m.viewport = viewport.New(100, 30)

	m.viewport.SetContent(renderList(m))

	p := tea.NewProgram(m, tea.WithAltScreen())

	_, err := p.Run()

	return err
}

func clearStatusAfter(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(_ time.Time) tea.Msg {
		return statusClearMsg{}
	})
}

func isStalzoneRunning() bool {
	cmd := exec.Command("pgrep", "-f", "stalzone.exe")
	return cmd.Run() == nil
}
