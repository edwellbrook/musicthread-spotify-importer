package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	spotifyPlaylistId := flag.String("playlist", "", "spotify playlist id")
	spotifyClientId := flag.String("spotify_client", "", "spotify api client id")
	spotifyClientSecret := flag.String("spotify_secret", "", "spotify api client secret")
	threadKey := flag.String("thread", "", "musicthread thread key")
	token := flag.String("token", "", "musicthread api token")

	flag.Parse()

	if *spotifyPlaylistId == "" {
		fmt.Println("missing '-playlist' argument")
		os.Exit(1)
	}

	if *spotifyClientId == "" {
		fmt.Println("missing '-spotify_client' argument")
		os.Exit(1)
	}

	if *spotifyClientSecret == "" {
		fmt.Println("missing '-spotify_secret' argument")
		os.Exit(1)
	}

	if *threadKey == "" {
		fmt.Println("missing '-thread' argument")
		os.Exit(1)
	}

	if *token == "" {
		fmt.Println("missing '-token' argument")
		os.Exit(1)
	}

	// spotifyPlaylistId, err := parseSpotifyURL("https://open.spotify.com/playlist/5det7OeX5MNsE51LdZcR7U?si=uyZZnnl4SSeemNXwM-TlFg")
	// if err != nil {
	// 	log.Println(err)
	// 	return
	// }

	spotifyAuthToken, err := fetchSpotifyAuthToken(*spotifyClientId, *spotifyClientSecret)
	if err != nil {
		fmt.Println("error fetching spotify auth token: " + err.Error())
		os.Exit(1)
	}

	playlistItemURIs, err := fetchSpotifyPlaylistItems(*spotifyPlaylistId, spotifyAuthToken)
	if err != nil {
		fmt.Println("error fetching playlist items: " + err.Error())
		os.Exit(1)
	}

	for _, uri := range playlistItemURIs {
		body := "{\"url\": \"" + uri + "\", \"thread\":\"" + *threadKey + "\"}"

		req, _ := http.NewRequest("POST", "https://musicthread.app/api/v0/add-link", bytes.NewBuffer([]byte(body)))
		req.Header.Add("Authorization", "Bearer "+*token)
		req.Header.Add("Content-Type", "application/json; charset=utf-8")

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			fmt.Println("error adding link to musicthread: " + err.Error())
			continue
		}

		respBody, _ := ioutil.ReadAll(res.Body)
		fmt.Println(strings.TrimSpace(string(respBody)))
	}
}
