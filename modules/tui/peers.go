package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) peersView() string {
	var b strings.Builder

	header := lipgloss.JoinHorizontal(
		lipgloss.Top,
		tabStyle.Width(17).Render("Status"),
		tabStyle.Width(17).Render("Peer"),
		tabStyle.Width(17).Render("IP"),
		tabStyle.Width(22).Render("RX"),
		tabStyle.Width(22).Render("TX"),
	)
	b.WriteString(sectionStyle.Render("Peers"))
	b.WriteString("\n")
	b.WriteString(fmt.Sprintf("  %s\n", header))
	b.WriteString(strings.Repeat("─", m.width-10))
	b.WriteString("\n")

	if len(m.peers) == 0 {
		b.WriteString("  No peers connected.\n")
	} else {
		for _, p := range m.peers {
			icon := "○"
			if p.Connected {
				icon = connectedStyle.Render("●")
			} else {
				icon = disconnectedStyle.Render("○")
			}

			row := lipgloss.JoinHorizontal(
				lipgloss.Top,
				tabStyle.Width(17).Render(icon),
				peerNameStyle.Width(17).Render(p.Name),
				peerIPStyle.Width(17).Render(p.IP),
				peerTrafficStyle.Width(22).Render(p.Rx),
				peerTrafficStyle.Width(22).Render(p.Tx),
			)
			fmt.Fprintf(&b, "  %s\n", row)
		}
	}

	return boxStyle.Width(m.width - 8).Render(b.String())
}
