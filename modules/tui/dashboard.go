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
	b.WriteString(sectionStyle.Render("Peers Overview"))
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

	return boxStyle.Width(m.width - 8).Render(b.String())
}
