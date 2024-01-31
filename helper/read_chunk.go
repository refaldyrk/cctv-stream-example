package helper

import (
	"encoding/json"
	"fmt"
	"github.com/canhlinh/hlsdl"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Source struct {
	Name string `json:"name"`
	Uri  string `json:"uri"`
}

func ReadChunk() {
	var wg sync.WaitGroup
	for {
		sources := ReadJSON()

		for i := 0; i < len(sources); i++ {
			wg.Add(1)
			go processSource(sources[i], &wg)
		}

		wg.Wait()

		time.Sleep(3 * time.Second)
	}
}

func processSource(source Source, wg *sync.WaitGroup) {
	defer wg.Done()

	fmt.Println("Downloading: " + source.Name)
	resp, err := http.Get(source.Uri + "/playlist.m3u8")
	if err != nil {
		log.Println("Error downloading from", source.Name, ":", err)
		return
	}
	defer resp.Body.Close()

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading response from", source.Name, ":", err)
		return
	}

	stringResponse := string(bytes)

	splitResponse := strings.Split(stringResponse, "\n")

	if len(splitResponse) >= 2 {
		chunk := splitResponse[len(splitResponse)-2]

		fmt.Println("Response: ", len(splitResponse))

		download, err := hlsdl.New(source.Uri+chunk, nil, "download", "", 8, false).Download()
		if err != nil {
			log.Println("Error downloading chunk from", source.Name, ":", err)
			return
		}

		fmt.Println("Downloaded: " + download)

		ConvertToMp4(download, "result/"+source.Name+"/"+source.Name+strconv.Itoa(int(time.Now().Unix()))+".mp4")

		time.Sleep(2 * time.Second)

		if _, err := os.Stat(download); err == nil {
			err = os.Remove(download)
			if err != nil {
				fmt.Println("Error removing downloaded chunk from", source.Name, ":", err)
			}
		} else {
			fmt.Println("File not found:", download)
		}
	} else {
		log.Println("Not enough elements in splitResponse for", source.Name)
	}
}

func ConvertToMp4(filenameIn, filenameOut string) {
	cmd := exec.Command("ffmpeg", "-i", filenameIn, "-c:v", "libx264", "-c:a", "aac", "-strict", "experimental", "-b:a", "192k", filenameOut)

	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error during conversion:", err)
		fmt.Println("Output:", string(out))
		return
	}

	fmt.Println("Conversion successful!")
}

func ReadJSON() []Source {
	open, err := os.Open("./helper/source.json")
	if err != nil {
		fmt.Println("Error opening source.json:", err)
		return nil
	}
	defer func(open *os.File) {
		err := open.Close()
		if err != nil {
			fmt.Println("Error closing source.json:", err)
		}
	}(open)

	bytes, err := io.ReadAll(open)
	if err != nil {
		fmt.Println("Error reading source.json:", err)
		return nil
	}

	var sources []Source

	err = json.Unmarshal(bytes, &sources)
	if err != nil {
		fmt.Println("Error unmarshalling source.json:", err)
		return nil
	}

	return sources
}

func InitFolder() {
	dir := ReadJSON()

	for i := 0; i < len(dir); i++ {
		err := os.Mkdir("result/"+dir[i].Name, 0755)
		if err != nil {
			fmt.Println("Error creating directory for", dir[i].Name, ":", err)
		}
	}
}
