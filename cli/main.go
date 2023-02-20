package main

import (
	"NintendoChannel/csdata"
	"NintendoChannel/dllist"
	"NintendoChannel/thumbnail"
	"fmt"
)

func main() {
	fmt.Println("WiiLink Nintendo Channel File Generator")
	fmt.Println()
	fmt.Println("1. DLList and game info\n2. Thumbnails\n3. CSData")
	fmt.Println()
	fmt.Printf("Choose: ")

	var selection int
	fmt.Scanln(&selection)

	switch selection {
	case 1:
		dllist.MakeDownloadList()
		break
	case 2:
		thumbnail.WriteThumbnail()
		break
	case 3:
		csdata.CreateCSData()
		break
	default:
		fmt.Println("\nInvalid Selection")
		main()
	}
}
