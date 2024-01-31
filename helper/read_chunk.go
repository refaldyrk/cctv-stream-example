package helper

import (
	"fmt"
	"github.com/canhlinh/hlsdl"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func ReadChunk(url, filename string) {
	resp, err := http.Get(url + "/playlist.m3u8")
	if err != nil {
		log.Println(err)
	}

	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}

	stringResponse := string(bytes)

	splitResponse := strings.Split(stringResponse, "\n")

	chunk := splitResponse[len(splitResponse)-2]

	download, err := hlsdl.New(url+chunk, nil, "download", "", 8, false).Download()
	if err != nil {
		return
	}

	fmt.Println("Downloaded: " + download)

	ConvertToMp4(download, "./result/"+filename)

	err = os.Remove(download)
	if err != nil {
		fmt.Println(err)
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
