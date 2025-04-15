package gametdb

import (
	"NintendoChannel/common"
	"archive/zip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
)

type GameTDB struct {
	XMLName   xml.Name  `xml:"datafile"`
	Companies Companies `xml:"companies"`
	Games     []Game    `xml:"game"`
}

type Companies struct {
	Companies []Company `xml:"company"`
}

type Company struct {
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Game struct {
	XMLName     xml.Name    `xml:"game"`
	ID          string      `xml:"id"`
	Type        string      `xml:"type"`
	Region      string      `xml:"region"`
	Locale      []GameMeta  `xml:"locale"`
	ReleaseDate Date        `xml:"date"`
	Rating      Rating      `xml:"rating"`
	Publisher   string      `xml:"publisher"`
	Controllers Controllers `xml:"input"`
	Features    Features    `xml:"wi-fi"`
	Languages   string      `xml:"languages"`
	Genre       string      `xml:"genre"`
}

type GameMeta struct {
	Language string `xml:"lang,attr"`
	Title    string `xml:"title"`
	Synopsis string `xml:"synopsis"`
}

type Date struct {
	Year  string `xml:"year,attr"`
	Month string `xml:"month,attr"`
	Day   string `xml:"day,attr"`
}

type Rating struct {
	Type  string `xml:"type,attr"`
	Value string `xml:"value,attr"`
}

type Controllers struct {
	Players    uint8 `xml:"players,attr"`
	Controller []struct {
		Type string `xml:"type,attr"`
	} `xml:"control"`
}

type Features struct {
	OnlinePlayers uint8    `xml:"players,attr"`
	Feature       []string `xml:"feature"`
}

var (
	WiiTDB     *GameTDB
	DSTDB      *GameTDB
	ThreeDSTDB *GameTDB

	tdbNames = []string{"wiitdb", "dstdb", "3dstdb"}
)

func PrepareGameTDB(config common.Config) {
	fmt.Println("Downloading GameTDB XML's...")
	client := &http.Client{}

	for i, name := range tdbNames {
		zipOut := fmt.Sprintf("%s/tdb.zip", config.AssetsPath)

		fmt.Println("Downloading " + name + "...")
		req, err := http.NewRequest("GET", fmt.Sprintf("https://www.gametdb.com/%s.zip", name), nil)
		common.CheckError(err)

		req.Header.Set("User-Agent", "WiiLink Nintendo Channel File Generator/0.1")

		response, err := client.Do(req)
		common.CheckError(err)

		defer response.Body.Close()
		contents, err := io.ReadAll(response.Body)
		common.CheckError(err)

		err = os.WriteFile(zipOut, contents, 0666)
		common.CheckError(err)

		// We need to unzip before we proceed to unmarshalling
		r, err := zip.OpenReader(zipOut)
		common.CheckError(err)

		fp, err := r.Open(fmt.Sprintf("%s.xml", name))
		common.CheckError(err)

		contents, err = io.ReadAll(fp)
		common.CheckError(err)

		var gameTDB GameTDB
		err = xml.Unmarshal(contents, &gameTDB)
		common.CheckError(err)

		switch i {
		case 0:
			WiiTDB = &gameTDB
		case 1:
			DSTDB = &gameTDB
		case 2:
			ThreeDSTDB = &gameTDB
		}

		err = os.Remove(zipOut)
		common.CheckError(err)

		fmt.Println("Finished current XML.")
	}

	fmt.Println("Finished all XMLs.")
}
