package network

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"

	fl "bytebrother/main/filer"
	st "bytebrother/main/settings"
)

var ChosenIndex int

func Start() {
	if ChosenIndex == -1 {
		return
	}

	// Find all available network interfaces
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}

	if ChosenIndex == 69420 {

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

			ChosenIndex = index

			settings := st.LoadSettings()

			st.SaveSettings(st.Settings{
				ProcessInterval:              settings.ProcessInterval,
				NetworkIndexToMonitor:        ChosenIndex,
				NumRowsInArchive:             settings.NumRowsInArchive,
				SaveProcessInformationInFile: settings.SaveProcessInformationInFile,
			})

			if ChosenIndex == -1 {
				fmt.Println("Unsafe network logging is disabled.")
				return
			}

			break
		}
	}

	chosenInterface := devices[ChosenIndex].Name

	// Open chosen network interface for capturing
	handle, err := pcap.OpenLive(chosenInterface, 65536, true, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	// Set filter to capture TCP packets
	err = handle.SetBPFFilter("tcp")
	if err != nil {
		log.Fatal(err)
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
				fl.AppendToFile(fl.LogFolder+fl.NetworkLogs, "["+fl.CurrentTime()+"] "+string(url)+"\n")
			}
		}
	}
}
