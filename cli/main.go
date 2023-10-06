package main

import (
	"NintendoChannel/csdata"    // Generate CSData
	"NintendoChannel/dllist"    // Generate DLList and Game Info
	"NintendoChannel/thumbnail" // Generate thumbnails
	"fmt"                       // Print text
	"os"                        // Clear screen
	"os/exec"                   // Open URL in default browser
	"runtime"                   // Get OS
	"strings"                   // Repeat string
	"time"                      // Time for current year

	"github.com/fatih/color" // Colorful text
	"golang.org/x/term"      // Terminal size
)

// Prints header for the program
func printHeader() {
	clearScreen()

	// Print header text
	currentYear := time.Now().Year()
	header := fmt.Sprintf("WiiLink Nintendo Channel File Generator - (c) %d WiiLink", currentYear)
	fmt.Println(bold(header))

	// Get console width
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}

	// Print a line of '=' across the console
	fmt.Println(bold(strings.Repeat("=", width)))
	fmt.Println()
}

// Clear the console screen
func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

// Make text bold
func bold(text string) string {
	return "\033[1m" + text + "\033[0m"
}

func main() {
	printHeader()

	// Print menu with possible options
	color.HiGreen(bold("What would you like to do?"))
	fmt.Println()

	fmt.Println("1. Generate DLList and Game Info")
	fmt.Println("2. Generate DLList Only")
	fmt.Println("3. Generate Thumbnails")
	fmt.Print("4. Generate CSData\n\n")

	// Add Github link and WiiLink website link
	fmt.Printf("5. Visit GitHub Repository  (%s)\n", color.HiGreenString("https://github.com/WiiLink24/NintendoChannel"))
	fmt.Printf("6. Visit WiiLink Website    (%s)\n\n", color.HiGreenString("https://wiilink24.com"))

	fmt.Print("7. Exit Program\n\n")

	// Get user input
	var selection int
	fmt.Print(bold("Choose: "))
	fmt.Scan(&selection)

	// Handle user input
	switch selection {
	case 1:
		printHeader()
		dllist.MakeDownloadList(true)
	case 2:
		printHeader()
		dllist.MakeDownloadList(false)
	case 3:
		printHeader()
		thumbnail.WriteThumbnail()
	case 4:
		printHeader()
		csdata.CreateCSData()
	case 5:
		// Open Github link in default browser
		err := OpenBrowser("https://github.com/WiiLink24/NintendoChannel")
		if err != nil {
			fmt.Println("Failed to open link:", err)
		}
		fmt.Println()
		fmt.Println("Opening link in default browser...")
		fmt.Println()

		time.Sleep(2 * time.Second) // Wait two seconds before looping

		main()
	case 6:
		// Open WiiLink website link in default browser
		err := OpenBrowser("https://wiilink24.com")
		if err != nil {
			fmt.Println("Failed to open link:", err)
		}
		fmt.Println()
		fmt.Println("Opening link in default browser...")
		fmt.Println()

		time.Sleep(2 * time.Second) // Wait two seconds before looping

		main()
	case 7:
		fmt.Println()
		fmt.Println("Exiting...")
		fmt.Println()

		time.Sleep(1 * time.Second) // Wait a second before exiting

		clearScreen()
	default:
		fmt.Println()
		color.HiRed(bold("Invalid selection! Please try again."))
		fmt.Println()

		time.Sleep(2 * time.Second) // Wait two seconds before looping

		main()
	}
}

// OpenBrowser opens the specified URL in the default browser of the user's system
func OpenBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("Platform is unsupported.")
	}

	return err
}
