package tui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/sh1zicus/stalzone-server-blocker/internal/model"
)

const (
	indCursor  = "▸ "
	indExpand  = "▼ "
	indCollapse = "▶ "
	indOn      = "■"
	indOff     = "□"
	indentPool = "  "
	indentTun  = "   "
)

func renderHeader(m Model) string {

	w := m.width
	if w == 0 {
		w = 80
	}

	var b strings.Builder

	// ---- Заголовок ----
	title := titleStyle.Render("STALZONE SERVER BLOCKER")
	b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, title))
	b.WriteString("\n")

	// ---- Статус процесса ----
	var processLine string
	if m.stalzoneRunning {
		processLine = processOn.Render("● stalzone.exe активен")
	} else {
		processLine = processOff.Render("● stalzone.exe не активен")
	}
	b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, processLine))
	b.WriteString("\n")

	// ---- Топ-3 (одна строка) ----
	type server struct {
		name string
		ping int
	}
	var best []server
	for _, p := range m.state.Pools {
		for _, t := range p.Tunnels {
			if t.PingOK {
				best = append(best, server{t.Name, t.PingMS()})
			}
		}
	}
	sort.Slice(best, func(i, j int) bool { return best[i].ping < best[j].ping })

	if len(best) > 0 {
		n := min(3, len(best))
		parts := make([]string, n)
		for i := 0; i < n; i++ {
			parts[i] = pingGreen.Render(best[i].name)
		}
		topLine := subtitleStyle.Render("Оптимальные: ") + strings.Join(parts, "  ·  ")
		b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, topLine))
		b.WriteString("\n")
	}

	// ---- Разделитель ----
	sep := dividerStyle.Render(strings.Repeat("─", min(w-4, 70)))
	b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, sep))
	b.WriteString("\n")

	// ---- Поиск ----
	if m.searchMode {
		search := searchStyle.Render("Поиск: " + m.search + "█")
		b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, search))
	} else {
		hint := subtitleStyle.Render("[ / ] поиск")
		b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, hint))
	}
	b.WriteString("\n")

	return b.String()
}

func renderList(m Model) string {

	w := m.width
	if w == 0 {
		w = 80
	}

	var b strings.Builder

	// ---- Колонки ----
	nw, pw, aw, cw := calcColWidths(m)

	// ---- Список ----
	for i, row := range m.state.Rows {

		isCursor := i == m.state.Cursor

		cur := indentPool
		if isCursor {
			cur = indCursor
		}

		// ---------- Пул ----------
		if row.Type == model.RowPool {

			p := m.state.Pools[row.Pool]

			icon := indCollapse
			if m.state.Expanded[row.Pool] {
				icon = indExpand
			}

			hasActive := p.SelectedCount() > 0
			nameStyle := poolName
			countStyle := poolCount
			if hasActive {
				nameStyle = poolNameActive
				countStyle = poolCountActive
			}

			count := countStyle.Render(padded(fmt.Sprintf("[%d/%d]", p.SelectedCount(), p.IPCount()), cw))

			if isCursor {
				bg := highlightLine
				nameStyled := nameStyle.Render(p.Name)
				namePad := strings.Repeat(" ", nw+2-lipgloss.Width(nameStyled))
				content := bg.Render(cur+icon+" ") + bg.Render(nameStyled) + bg.Render(namePad) + bg.Render(count)
				b.WriteString(padHighlightLeft(content, w, bg))
			} else {
				name := nameStyle.Render(p.Name)
				raw := cur + icon + " " + padVis(name, nw+2) + count
				b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, raw))
			}
			b.WriteString("\n")
			continue
		}

		// ---------- Туннель ----------
		t := m.state.Pools[row.Pool].Tunnels[row.Tunnel]

		var sq string
		if t.Selected {
			sq = statusOn.Render(indOn)
		} else {
			sq = statusOff.Render(indOff)
		}

		if isCursor {
			bg := highlightLine
			nameStyled := tunnelName.Render(padded(t.Name, nw))
			pingStyled := pingTextPadded(t, pw)
			addrStyled := tunnelAddr.Render(padded(t.Address, aw))
			content := bg.Render(indentTun) + bg.Render(sq) + bg.Render("  ") +
				bg.Render(nameStyled) + bg.Render("  ") +
				bg.Render(pingStyled) + bg.Render("  ") +
				bg.Render(addrStyled)
			b.WriteString(padHighlightLeft(content, w, bg))
		} else {
			name := tunnelName.Render(padded(t.Name, nw))
			ping := pingTextPadded(t, pw)
			addr := tunnelAddr.Render(padded(t.Address, aw))
			raw := indentTun + sq + "  " + name + "  " + ping + "  " + addr
			b.WriteString(lipgloss.PlaceHorizontal(w, lipgloss.Center, raw))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func calcColWidths(m Model) (nameW, pingW, addrW, countW int) {

	nameW = 16
	pingW = 7
	addrW = 18
	countW = 0

	for _, row := range m.state.Rows {
		if row.Type == model.RowPool {
			p := m.state.Pools[row.Pool]
			cw := lipgloss.Width(poolCount.Render(fmt.Sprintf("[%d/%d]", p.SelectedCount(), p.IPCount())))
			if cw > countW {
				countW = cw
			}
			continue
		}

		t := m.state.Pools[row.Pool].Tunnels[row.Tunnel]

		if len(t.Name) > nameW {
			nameW = len(t.Name)
		}
		if len(t.Address) > addrW {
			addrW = len(t.Address)
		}
		pw := lipgloss.Width(pingText(t))
		if pw > pingW {
			pingW = pw
		}
	}

	nameW += 2
	pingW += 2
	addrW += 1

	return nameW, pingW, addrW, countW
}

func (m Model) View() string {

	w := m.width
	if w == 0 {
		w = 80
	}

	// ---- Статус ----
	var statusLine string
	if m.status != "" {
		if strings.HasPrefix(m.status, "Ошибка") {
			statusLine = statusErr.Render(m.status)
		} else {
			statusLine = statusOK.Render(m.status)
		}
	}

	// ---- Хелп-бар ----
	bindings := []struct{ key, desc string }{
		{"↑↓", "движение"},
		{"←→", "раскрыть"},
		{"␣", "перекл"},
		{"enter", "пул"},
		{"R", "пинг"},
		{"A", "применить"},
		{"D", "выкл"},
		{"/", "поиск"},
		{"q", "выход"},
	}

	var helpParts []string
	for _, b := range bindings {
		helpParts = append(helpParts, helpKey.Render(b.key)+" "+helpDesc.Render(b.desc))
	}
	help := strings.Join(helpParts, "   ")
	helpLine := lipgloss.PlaceHorizontal(w, lipgloss.Center, help)

	bottom := helpLine
	if statusLine != "" {
		bottom = lipgloss.JoinVertical(lipgloss.Center, statusLine, helpLine)
	}

	return renderHeader(m) + m.viewport.View() + "\n\n" + bottom
}

func pingText(t model.Tunnel) string {
	if !t.PingOK {
		return pingGray.Render("таймаут")
	}
	ms := t.PingMS()
	switch {
	case ms <= 30:
		return pingGreen.Render(t.PingString())
	case ms <= 60:
		return pingYellow.Render(t.PingString())
	case ms <= 90:
		return pingOrange.Render(t.PingString())
	default:
		return pingRed.Render(t.PingString())
	}
}

func pingTextPadded(t model.Tunnel, width int) string {
	pt := pingText(t)
	vw := lipgloss.Width(pt)
	if vw >= width {
		return pt
	}
	return pt + strings.Repeat(" ", width-vw)
}

func padded(s string, width int) string {
	if len(s) >= width {
		return s
	}
	return s + strings.Repeat(" ", width-len(s))
}

func padVis(styled string, width int) string {
	vw := lipgloss.Width(styled)
	if vw >= width {
		return styled
	}
	return styled + strings.Repeat(" ", width-vw)
}

func padCenterPlain(styled string, width int) string {
	vw := lipgloss.Width(styled)
	if vw >= width {
		return styled
	}
	left := (width - vw) / 2
	right := width - vw - left
	return strings.Repeat(" ", left) + styled + strings.Repeat(" ", right)
}

func padHighlightLeft(styled string, width int, bg lipgloss.Style) string {
	vw := lipgloss.Width(styled)
	if vw >= width {
		return styled
	}
	left := (width - vw) / 2
	right := width - vw - left
	return bg.Render(strings.Repeat(" ", left)) + styled + bg.Render(strings.Repeat(" ", right))
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
