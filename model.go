package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var accentColors = []lipgloss.Color{
	lipgloss.Color("#e06c2a"),
	lipgloss.Color("#61afef"),
	lipgloss.Color("#98c379"),
	lipgloss.Color("#c678dd"),
}

const asciiArt = `::::::::::::::::.
::::::::::::::::.
::......::::::::
::. ...  .:::::
::...:..  .::::
:::....   .::::
::::..   .:::::
:::::. .:::::::
:::::::::::::::
::::::::::::::.`

const asciiWidth = 18

type LinkItem struct {
	Label string
	Value string
}

type Page struct {
	Title   string
	Content []string
	Links   [][2]LinkItem
}

var pages = []Page{
	{
		Title: "About",
		Content: []string{
			"Thiago Augusto is a Network & Infrastructure Analyst",
			"based in Recife, Brazil, working with routing,",
			"firewalls, VPNs and monitoring.",
			"",
			"Currently studying Systems Analysis and Development.",
			"MikroTik MTCNA and Cisco CCNA certified.",
		},
	},
	{
		Title: "Links",
		Links: [][2]LinkItem{
			{
				{"Linkedin", "/in/thtpps"},
				{"Github", "@thttp"},
			},
			{
				{"Resume", "thttp.github.io"},
				{"Website", "working on..."},
			},
		},
	},
	{
		Title: "Projects",
		Content: []string{
			"Work in progress...",
		},
	},
}

type keyMap struct {
	Left  key.Binding
	Right key.Binding
	Theme key.Binding
	Quit  key.Binding
}

var keys = keyMap{
	Left:  key.NewBinding(key.WithKeys("left", "h")),
	Right: key.NewBinding(key.WithKeys("right", "l", "tab")),
	Theme: key.NewBinding(key.WithKeys("t")),
	Quit:  key.NewBinding(key.WithKeys("q", "ctrl+c")),
}

type model struct {
	currentPage int
	accentIndex int
	width       int
	height      int
}

func initialModel() model {
	return model{0, 0, 80, 24}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, keys.Right):
			m.currentPage = (m.currentPage + 1) % len(pages)
		case key.Matches(msg, keys.Left):
			m.currentPage = (m.currentPage-1+len(pages)) % len(pages)
		case key.Matches(msg, keys.Theme):
			m.accentIndex = (m.accentIndex + 1) % len(accentColors)
		}
	}
	return m, nil
}

func col(c lipgloss.Color, s string) string {
	return lipgloss.NewStyle().Foreground(c).Render(s)
}

func sp(n int) string {
	if n <= 0 {
		return ""
	}
	return strings.Repeat(" ", n)
}

func (m model) navBar(accent, muted lipgloss.Color) string {
	activeStyle := lipgloss.NewStyle().Foreground(accent).Bold(true).PaddingRight(3)
	inactiveStyle := lipgloss.NewStyle().Foreground(muted).PaddingRight(3)

	var tabs []string
	for i, p := range pages {
		if i == m.currentPage {
			tabs = append(tabs, activeStyle.Render(p.Title))
		} else {
			tabs = append(tabs, inactiveStyle.Render(p.Title))
		}
	}
	tabsStr   := strings.Join(tabs, "")
	statusStr := col(accent, "[ construction ]")
	gap       := m.width - lipgloss.Width(tabsStr) - lipgloss.Width(statusStr) - 4
	if gap < 0 {
		gap = 0
	}
	return sp(2) + tabsStr + sp(gap) + statusStr
}

func (m model) footerBar(accent, muted lipgloss.Color) string {
	hints := fmt.Sprintf("%s navigate  %s theme  %s quit",
		col(accent, "←→"),
		col(accent, "t"),
		col(accent, "q"),
	)
	version := col(muted, "v0.0.1")
	mid     := (m.width-lipgloss.Width(hints))/2 - lipgloss.Width(version)
	if mid < 1 {
		mid = 1
	}
	return version + sp(mid) + hints
}

func (m model) viewAbout(fg lipgloss.Color) string {
	asciiLines := strings.Split(asciiArt, "\n")
	content    := pages[m.currentPage].Content

	totalLines := len(asciiLines)
	if len(content) > totalLines {
		totalLines = len(content)
	}

	gap        := 4
	blockWidth := asciiWidth + gap + 44
	leftPad    := (m.width - blockWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	bodyHeight := m.height - 3
	topPad     := (bodyHeight - totalLines) / 2
	if topPad < 0 {
		topPad = 0
	}

	rows := make([]string, 0, bodyHeight)
	for i := 0; i < topPad; i++ {
		rows = append(rows, "")
	}
	for i := 0; i < totalLines; i++ {
		var left string
		if i < len(asciiLines) {
			left = col(fg, asciiLines[i]) + sp(asciiWidth-len(asciiLines[i]))
		} else {
			left = sp(asciiWidth)
		}
		var right string
		if i < len(content) {
			right = col(fg, content[i])
		}
		rows = append(rows, sp(leftPad)+left+sp(gap)+right)
	}
	for len(rows) < bodyHeight {
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func (m model) viewLinks(muted, fg lipgloss.Color) string {
	const colWidth = 28
	const colGap   = 6
	const blockWidth = colWidth*2 + colGap

	leftPad := (m.width - blockWidth) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	page := pages[m.currentPage]
	rows := make([]string, 0, len(page.Links)*3)
	for _, pair := range page.Links {
		l, r := pair[0], pair[1]
		rows = append(rows,
			sp(leftPad)+col(fg, l.Label)+sp(colWidth-len(l.Label)+colGap)+col(fg, r.Label),
			sp(leftPad)+col(muted, l.Value)+sp(colWidth-len(l.Value)+colGap)+col(muted, r.Value),
			"",
		)
	}

	bodyHeight := m.height - 3
	topPad     := (bodyHeight - len(rows)) / 2
	if topPad < 0 {
		topPad = 0
	}

	body := make([]string, 0, bodyHeight)
	for i := 0; i < topPad; i++ {
		body = append(body, "")
	}
	body = append(body, rows...)
	for len(body) < bodyHeight {
		body = append(body, "")
	}
	return strings.Join(body, "\n")
}

func (m model) viewContent(fg lipgloss.Color) string {
	content   := pages[m.currentPage].Content
	leftPad   := (m.width - 44) / 2
	if leftPad < 0 {
		leftPad = 0
	}

	bodyHeight := m.height - 3
	topPad     := (bodyHeight - len(content)) / 2
	if topPad < 0 {
		topPad = 0
	}

	rows := make([]string, 0, bodyHeight)
	for i := 0; i < topPad; i++ {
		rows = append(rows, "")
	}
	for _, line := range content {
		rows = append(rows, sp(leftPad)+col(fg, line))
	}
	for len(rows) < bodyHeight {
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func (m model) View() string {
	accent := accentColors[m.accentIndex]
	muted  := lipgloss.Color("#555555")
	fg     := lipgloss.Color("#d4d4d4")

	var body string
	switch m.currentPage {
	case 0:
		body = m.viewAbout(fg)
	case 1:
		body = m.viewLinks(muted, fg)
	default:
		body = m.viewContent(fg)
	}

	var out strings.Builder
	out.WriteString(m.navBar(accent, muted) + "\n\n")
	out.WriteString(body + "\n")
	out.WriteString(m.footerBar(accent, muted))
	return out.String()
}