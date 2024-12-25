package main

import (
	"NintendoChannel/v3/wc24data"
	"NintendoChannel/v6/csdata"
	"NintendoChannel/v6/dllist"
	"NintendoChannel/v6/thumbnail"
	"fmt"
)

func main() {
	fmt.Println("WiiLink Nintendo Channel File Generator")
	fmt.Println()
	fmt.Println("1. DLList and game info\n2. DLList only\n3. Thumbnails\n4. CSData\n5. v3 List")
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
	case 5:
		wc24data.MakeList()
		break
	default:
		fmt.Println("\nInvalid Selection")
		main()
	}
}
