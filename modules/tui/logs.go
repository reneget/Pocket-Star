package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

func (m model) logsView() string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("Logs"))
	b.WriteString("\n\n")

	lines := m.logBuf.Slice()
	if len(lines) == 0 {
		b.WriteString("  No logs yet.\n")
		return boxStyle.Width(m.width - 8).Render(b.String())
	}

	// Show last N lines fitting the viewport height
	maxLines := m.height - 16
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	now := time.Now()
	for _, line := range lines {
		ts := now.Format("15:04:05")
		fmt.Fprintf(&b, "  %s %s\n",
			logTimeStyle.Render(ts),
			logTextStyle.Render(line),
		)
	}

	return boxStyle.Width(m.width - 8).Render(b.String())
}

func (m model) helpView() string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("Help"))
	b.WriteString("\n\n")

	keys := []struct {
		key string
		desc string
	}{
		{"q / Ctrl+C", "Quit Pocket Star"},
		{"? / /h", "Toggle this help screen"},
		{"/q, /quit", "Quit (command mode)"},
		{"/c, /connect", "Connect to hub"},
		{"/d, /disconnect", "Disconnect from hub"},
		{"/r, /refresh", "Refresh connection status"},
		{"/tab 1-3", "Switch to tab N"},
		{"1, 2, 3", "Switch to tab (Dashboard, Peers, Logs)"},
		{"Tab key", "Next tab"},
		{"r", "Refresh status"},
		{"Mouse", "Click tabs and buttons"},
	}

	for _, k := range keys {
		fmt.Fprintf(&b, "  %s  %s\n",
			helpStyle.Width(18).Align(lipgloss.Left).Render(k.key),
			valueStyle.Render(k.desc),
		)
	}

	b.WriteString("\n")
	b.WriteString("Press any key to close.\n")

	return boxStyle.Width(m.width - 8).Render(b.String())
}
