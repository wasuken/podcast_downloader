package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Item represents a single item in the RSS feed
type Item struct {
	Title     string `xml:"title"`
	Enclosure struct {
		URL string `xml:"url,attr"`
	} `xml:"enclosure"`
}

// RSS represents the root element of the RSS feed
type RSS struct {
	Channel struct {
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

func main() {
	// XMLファイルのURL
	xmlURL := os.Args[1]

	// XMLファイルをダウンロード
	resp, err := http.Get(xmlURL)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	// XMLファイルの内容を解析
	var rss RSS
	decoder := xml.NewDecoder(resp.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		panic(err)
	}

	outputDir := "./output"
	er := downloadMp3FromRSSItemList(rss.Channel.Items, outputDir)
	if er != nil {
		panic(er)
	}
}

func downloadMp3FromRSSItemList(items []Item, outputDir string) error {
	numOfListElm := len(items)
	for i, item := range items {
		outputFilepath := fmt.Sprintf("%s/%d_%s.mp3", outputDir, i, item.Title)
		fmt.Println(fmt.Sprintf("[%d/%d] downloading %s...", i+1, numOfListElm, item.Enclosure.URL))
		outputFile, err := os.Create(outputFilepath)
		if err != nil {
			return err
		}
		defer outputFile.Close()
		if err := downloadFile(item.Enclosure.URL, outputFile); err != nil {
			fmt.Println("Error downloading file:", err)
		} else {
			fmt.Println(fmt.Sprintf("=> downloaded write: %s", outputFilepath))
		}
	}
	return nil
}

func downloadFile(url string, outputFile *os.File) error {
	// ファイルをダウンロード
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// ダウンロードした内容を出力ファイルに書き込む
	_, err = io.Copy(outputFile, resp.Body)
	return err
}
