package main

import (
	"flag"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"

	"github.com/kkdai/youtube/v2"
)

func main() {
	playlistUrl := flag.String("url", "", "Your Playlist URL")
	flag.Parse()

	if *playlistUrl == "" {
		panic("URL is required")
	}

	urlString := *playlistUrl
	u, er := url.Parse(urlString)
	if er != nil {
		fmt.Println(er.Error())
		return
	}
	playlistId := u.Query()["list"][0]

	downloadedFromPlaylist(playlistId)
}

func downloadedFromPlaylist(playListId string) {
	client := youtube.Client{}

	playlist, err := client.GetPlaylist(playListId)
	if err != nil {
		panic(err)
	}

	/* ----- Enumerating playlist videos ----- */
	header := fmt.Sprintf("Playlist %s by %s", playlist.Title, playlist.Author)
	println(header)
	println(strings.Repeat("=", len(header)) + "\n")

	for k, v := range playlist.Videos {
		fmt.Printf("(%d) %s - '%s [%v]' \n", k+1, v.Author, v.Title, v.Duration)

		/* ----- Downloading video ----- */
		entry := v
		video, err := client.VideoFromPlaylistEntry(entry)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Downloading %s by '%s'!\n", video.Title, video.Author)

		stream, _, err := client.GetStream(video, video.Formats.FindByQuality("medium"))
		// stream, _, err := client.GetStream(video, &video.Formats[0])
		if err != nil {
			panic(err)
		}

		file, err := os.Create("downloads/" + video.Title + ".mp4")
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(file, stream)

		if err != nil {
			panic(err)
		}

		println("Downloaded " + video.Title)
		file.Close()
	}

}
