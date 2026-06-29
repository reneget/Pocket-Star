package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) dashboardView() string {
	var b strings.Builder

	b.WriteString(sectionStyle.Render("Connection"))
	b.WriteString("\n")

	var connStatus string
	if m.status.Connected {
		connStatus = connectedStyle.Render("● Connected")
	} else {
		connStatus = disconnectedStyle.Render("○ Disconnected")
	}
	fmt.Fprintf(&b, "  %s%s\n", labelStyle.Render("Status:"), connStatus)
	fmt.Fprintf(&b, "  %s%s\n", labelStyle.Render("Hub:"), valueStyle.Render(m.status.HubAddr))
	fmt.Fprintf(&b, "  %s%s\n", labelStyle.Render("My IP:"), valueStyle.Render(m.status.MyIP))
	fmt.Fprintf(&b, "  %s%s\n", labelStyle.Render("Uptime:"), valueStyle.Render(m.status.Uptime))

	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Actions"))
	b.WriteString("\n\n")

	btns := lipgloss.JoinHorizontal(
		lipgloss.Top,
		buttonStyle.Render(" Connect "),
		lipgloss.NewStyle().Width(2).Render(""),
		buttonDisabledStyle.Render(" Disconnect "),
	)
	b.WriteString(fmt.Sprintf("  %s\n", btns))

	b.WriteString("\n")
	b.WriteString(m.healthView())
	b.WriteString("\n")
	b.WriteString(m.containersView())
	b.WriteString("\n")
	b.WriteString(m.peersOverview())

	return boxStyle.Width(m.width - 8).Render(b.String())
}

func (m model) healthView() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("System"))
	b.WriteString("\n\n")

	barWidth := 20

	cpuBar := progressBar(m.healthData.CPU, barWidth)
	ramBar := progressBar(m.healthData.RAM, barWidth)
	diskBar := progressBar(m.healthData.Disk, barWidth)

	fmt.Fprintf(&b, "  %s %s %5.1f%%\n", labelStyle.Render("CPU:"), cpuBar, m.healthData.CPU)
	fmt.Fprintf(&b, "  %s %s %5.1f%%  %s\n", labelStyle.Render("RAM:"), ramBar, m.healthData.RAM, m.healthData.RAMUsed)
	fmt.Fprintf(&b, "  %s %s %5.1f%%  %s\n", labelStyle.Render("Disk:"), diskBar, m.healthData.Disk, m.healthData.DiskUsed)
	if m.healthData.PubIP != "" {
		fmt.Fprintf(&b, "  %s%s\n", labelStyle.Render("Pub IP:"), valueStyle.Render(m.healthData.PubIP))
	}

	return b.String()
}

func (m model) containersView() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("Containers"))
	b.WriteString("\n\n")

	if len(m.containers) == 0 {
		b.WriteString("  No containers found.\n")
		return b.String()
	}

	for _, ct := range m.containers {
		icon := "○"
		if ct.Running {
			icon = "●"
		}
		statusColor := disconnectedStyle
		if ct.Running {
			statusColor = connectedStyle
		}

		cpuBar := minibar(ct.CPUPercent, 10)
		memBar := minibar(ct.MemPercent, 10)

		fmt.Fprintf(&b, "  %s %s %s  cpu:%s ram:%s %s\n",
			statusColor.Render(icon),
			peerNameStyle.Width(18).Render(ct.Name),
			valueStyle.Render(ct.Status),
			cpuBar,
			memBar,
			statusColor.Render(ct.Status[:min(7, len(ct.Status))]),
		)
	}

	return b.String()
}

func (m model) peersOverview() string {
	var b strings.Builder
	b.WriteString(sectionStyle.Render("Peers"))
	b.WriteString("\n\n")

	if len(m.peers) == 0 {
		b.WriteString("  No peers connected.\n")
	} else {
		for _, p := range m.peers {
			icon := "○"
			if p.Connected {
				icon = "●"
			}
			fmt.Fprintf(&b, "  %s %s  %s  rx: %s  tx: %s\n",
				icon, peerNameStyle.Render(p.Name),
				peerIPStyle.Render(p.IP),
				peerTrafficStyle.Render(p.Rx),
				peerTrafficStyle.Render(p.Tx))
		}
	}

	return b.String()
}

func progressBar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int(pct * float64(width) / 100)
	empty := width - filled

	fillChar := "█"
	emptyChar := "░"

	fill := strings.Repeat(fillChar, filled)
	emp := strings.Repeat(emptyChar, empty)

	var color lipgloss.Color
	switch {
	case pct > 80:
		color = lipgloss.Color("#e74c3c")
	case pct > 50:
		color = lipgloss.Color("#f1c40f")
	default:
		color = lipgloss.Color("#2ecc71")
	}

	return lipgloss.NewStyle().Foreground(color).Render(fill + emp)
}

func minibar(pct float64, width int) string {
	if pct < 0 {
		pct = 0
	}
	if pct > 100 {
		pct = 100
	}
	filled := int(pct * float64(width) / 100)
	if filled > width {
		filled = width
	}

	fill := strings.Repeat("█", filled)
	emp := strings.Repeat("░", width-filled)

	var color lipgloss.Color
	switch {
	case pct > 80:
		color = lipgloss.Color("#e74c3c")
	case pct > 50:
		color = lipgloss.Color("#f1c40f")
	default:
		color = lipgloss.Color("#2ecc71")
	}

	return lipgloss.NewStyle().Foreground(color).Render(fill + emp)
}
