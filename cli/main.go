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
	fmt.Println("1. DLList and game info\n2. DLList only\n3. Thumbnails\n4. CSData")
	fmt.Println()
	fmt.Printf("Choose: ")

	var selection int
	fmt.Scanln(&selection)

	switch selection {
	case 1:
		dllist.MakeDownloadList(true)
		break
	case 2:
		dllist.MakeDownloadList(false)
		break
	case 3:
		thumbnail.WriteThumbnail()
		break
	case 4:
		csdata.CreateCSData()
		break
	default:
		fmt.Println("\nInvalid Selection")
		main()
	}
}
