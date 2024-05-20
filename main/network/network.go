package network

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	fl "bytebrother/main/filer"
	g "bytebrother/main/global"
	st "bytebrother/main/settings"
)

func Start() error {
	if g.ChosenIndex == -1 {
		return nil
	}

	// Find all available network interfaces
	devices, err := pcap.FindAllDevs()
	if err != nil {
		return fmt.Errorf("failed to find all devices: %w", err)
	}

	if g.ChosenIndex == 69420 {

		// Print information about each network interface
		fmt.Println("Available network interfaces:")
		for i, device := range devices {
			fmt.Printf("%d: %s (%s)\n", i, device.Name, device.Description)
		}

		// Prompt user to choose a network interface
		for {
			fmt.Print("\nEnter the index of the network you want to monitor (Check in task manager if you are unsure)\nif you don't know what unsafe network logging is, type '-1': ")
			var input string
			_, err := fmt.Scanln(&input)
			if err != nil {
				log.Fatal("Error reading input:", err)
			}

			index, err := strconv.Atoi(input)

			if err != nil || index < 1 && index != -1 || index > len(devices) && index != -1 {
				fmt.Println("Invalid input. Please enter a valid index.")
				continue
			}

			g.ChosenIndex = index

			settings := st.LoadSettings()

			settings.NetworkIndexToMonitor = g.ChosenIndex
			st.SaveSettings(settings)

			if g.ChosenIndex == -1 {
				fmt.Println("Unsafe network logging is disabled.")
				return nil
			}

			break
		}
	}

	if g.IsSetup {
		fmt.Println("Network logging is set up. You can later change the network interface to monitor in the bytebrother/settings.json file.\n if youre still unsure and want to try again, delete the entry 'network_index_to_monitor' in the settings.json file.")
		os.Exit(0)
	}

	chosenInterface := devices[g.ChosenIndex].Name

	// Open chosen network interface for capturing
	handle, err := pcap.OpenLive(chosenInterface, 65536, true, pcap.BlockForever)
	if err != nil {
		return fmt.Errorf("failed to open live network interface: %w", err)
	}
	defer handle.Close()

	// Set filter to capture TCP packets
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		return fmt.Errorf("failed to set BPF filter: %w", err)
	}

	// Define a regular expression to match URLs either http or https
	urlRegex := regexp.MustCompile(`(?i)\b((?:https?://|www\d{0,3}[.]|[a-z0-9.\-]+[.][a-z]{2,4}/)(?:[^\s()<>]+|\(([^\s()<>]+|(\([^\s()<>]+\)))*\))+(?:\(([^\s()<>]+|(\([^\s()<>]+\)))*\)|[^\s]))`)

	// Start capturing packets
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		// Check if the packet contains application layer data
		appLayer := packet.ApplicationLayer()
		if appLayer != nil {
			// Extract payload from the packet
			payload := appLayer.Payload()

			// fmt.Println("Packet found:" + string(payload))
			// Find URLs in the payload
			urls := urlRegex.FindAll(payload, -1)
			for _, url := range urls {
				// fmt.Printf("URL found: %s\n", url)
				fl.AppendToFile(fl.LogFolder+fl.NetworkLogs, "["+fl.CurrentTime(":")+"] "+string(url)+"\n")
			}
		}
	}

	return nil
}
