package pinger

import (
	"fmt"
	"github.com/pterm/pterm"
)

type HostOutput interface {
	getInitOutputStatus() string
	getRunOutputStatus() string
	getOutputStatus() string
	GetOutputData() []string
}

func (h *Host) getInitOutputStatus() string {
	return pterm.NewStyle(pterm.BgLightBlue, pterm.Bold).Sprintf(" INIT ")
}

func (h *Host) getRunOutputStatus() string {
	if h.Available {
		return pterm.NewStyle(pterm.BgLightGreen, pterm.Bold).Sprintf(" LIVE ")
	} else {
		return pterm.NewStyle(pterm.BgLightRed, pterm.Bold).Sprintf(" DOWN ")
	}
}

func (h Host) getOutputStatus() string {
	if h.PacketsSent == 0 {
		return h.getInitOutputStatus()
	} else {
		return h.getRunOutputStatus()
	}
}

func (h Host) GetOutputData() []string {
	return []string{
		h.Address,
		h.getOutputStatus(),
		fmt.Sprintf("%d", h.PacketsSent),
		fmt.Sprintf("%d", h.PacketsReceived),
		fmt.Sprintf("%3.0f%%", h.PacketLoss),
		h.Comment,
	}
}
