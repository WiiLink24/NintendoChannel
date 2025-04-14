package main

import (
	"NintendoChannel/v3/wc24data"
	"NintendoChannel/v6/csdata"
	"NintendoChannel/v6/dllist"
	"NintendoChannel/v6/thumbnail"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	var selection int

	// Allow for arguments to be passed for autogenerate.
	if len(os.Args) > 1 {
		var err error
		selection, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatalln("Invalid selection")
		}
		selector(selection)
	}

	fmt.Println("WiiLink Nintendo Channel File Generator")
	fmt.Println()
	fmt.Println("1. DLList and game info\n2. DLList only\n3. Thumbnails\n4. CSData\n5. v3 List")
	fmt.Println()
	fmt.Printf("Choose: ")

	fmt.Scanln(&selection)
	selector(selection)
}

func selector(selection int) {
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
