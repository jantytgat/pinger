package cmd

//import "C"
import (
	"github.com/go-ping/ping"
	"github.com/jantytgat/pinger/pinger"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"sync"
	"time"
)

var rootCmd = &cobra.Command{
	Use:   "pinger",
	Short: "Pinger allows availability checks for multiple hosts at the same time",
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

func Execute() {
	clearDisplay()

	hosts := []string{
		"1.1.1.1",
		"1.1.1.2",
		"1.1.1.3",
		"8.8.8.8",
		"8.8.4.4",
		"9.9.9.9",
	}
	wg := &sync.WaitGroup{}
	ch := make(chan pinger.Host)

	wg.Add(1)
	go consolePrinter(ch, wg)

	for _, host := range hosts {
		wg.Add(1)
		go pingHost(host, ch, wg)
		time.Sleep(time.Millisecond * 100)
	}

	wg.Wait()
	close(ch)
}

func clearDisplay() {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}

	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func pingHost(address string, ch chan<- pinger.Host, wg *sync.WaitGroup) {

	host := pinger.Host{
		Address:         address,
		PacketsSent:     0,
		PacketsReceived: 0,
		Available:       false,
		PacketLoss:      0,
	}

	hostPinger, err := ping.NewPinger(address)
	if err != nil {
		host.Comment = err.Error()
		ch <- host

	} else {
		timer := time.Second + time.Millisecond*time.Duration(250*(rand.Int63n(4)))
		hostPinger.Interval = timer

		if runtime.GOOS == "windows" {
			hostPinger.SetPrivileged(true)
		}

		hostPinger.OnSend = func(pkt *ping.Packet) {
			host = updateHostData(host, hostPinger.Statistics())
			ch <- pinger.Host{
				Address:         host.Address,
				PacketsSent:     host.PacketsSent,
				PacketsReceived: host.PacketsReceived,
				Available:       host.Available,
				PacketLoss:      host.PacketLoss,
			}
		}

		err = hostPinger.Run()
		if err != nil {
			host.Comment = err.Error()
			ch <- host
		}
	}
	wg.Done()
}

func updateHostData(data pinger.Host, stats *ping.Statistics) pinger.Host {
	if stats.PacketLoss > data.PacketLoss || stats.PacketLoss == 100 {
		data.Available = false
	} else {
		data.Available = true
	}
	if !math.IsNaN(stats.PacketLoss) {
		data.PacketLoss = stats.PacketLoss
	}

	data.PacketsSent = stats.PacketsSent
	data.PacketsReceived = stats.PacketsRecv

	return data
}

func consolePrinter(ch <-chan pinger.Host, wg *sync.WaitGroup) {
	completed := false

	area, _ := pterm.DefaultArea.Start()
	var hostsData = make(map[string]pinger.Host)
	for !completed {
		select {
		case data, ok := <-ch:
			if !ok {
				completed = true
			}
			hostsData[data.Address] = data
		default:
			area.Update(updateOutput(hostsData))
			time.Sleep(time.Millisecond * 500)
		}
		err := area.Stop()
		if err != nil {
			return
		}
	}
	wg.Done()
}

func updateOutput(hosts map[string]pinger.Host) string {
	// Sort keys
	keys := make([]string, 0, len(hosts))
	for k := range hosts {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var data [][]string
	data = append(data, []string{"Address", "Status", "TX", "RX", "Lost", "Details"})

	for _, k := range keys {
		data = append(data, hosts[k].GetOutputData())
	}
	table, tableErr := pterm.DefaultTable.WithHasHeader().WithData(data).Srender()
	if tableErr != nil {
		panic(tableErr)
	}
	return table
}
