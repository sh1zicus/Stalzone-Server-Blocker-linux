package tui

import "github.com/charmbracelet/lipgloss"

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("247")).
			Align(lipgloss.Center)

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("239"))

	searchStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226")).
			Bold(true)

	// Курсор
	cursorMarker = lipgloss.NewStyle().
			Foreground(lipgloss.Color("39")).
			Bold(true)

	// Подсветка строки под курсором
	highlightLine = lipgloss.NewStyle().
			Background(lipgloss.Color("240"))

	// Пулы
	poolIcon = lipgloss.NewStyle().
			Foreground(lipgloss.Color("247"))

	poolName = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("255"))

	poolNameActive = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("46"))

	poolCount = lipgloss.NewStyle().
			Foreground(lipgloss.Color("247"))

	poolCountActive = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	// Квадраты статуса сервера
	statusOn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	statusOff = lipgloss.NewStyle().
			Foreground(lipgloss.Color("248"))

	// Имя туннеля
	tunnelName = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255"))

	tunnelAddr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("247"))

	// Пинг
	pingGreen = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46"))

	pingYellow = lipgloss.NewStyle().
			Foreground(lipgloss.Color("226"))

	pingOrange = lipgloss.NewStyle().
			Foreground(lipgloss.Color("208"))

	pingRed = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196"))

	pingGray = lipgloss.NewStyle().
			Foreground(lipgloss.Color("244"))

	// Статус внизу
	statusOK = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	statusErr = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	// Хелп-бар
	helpKey = lipgloss.NewStyle().
			Foreground(lipgloss.Color("255")).
			Bold(true)

	helpDesc = lipgloss.NewStyle().
			Foreground(lipgloss.Color("247"))

	// Индикатор процесса
	processOn = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	processOff = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
)
