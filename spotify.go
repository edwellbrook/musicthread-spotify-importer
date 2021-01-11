package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func parseSpotifyURL(urlStr string) (string, error) {
	playlistURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	log.Println(playlistURL.Path)

	pathParts := strings.Split(strings.Trim(playlistURL.Path, "/"), "/")

	log.Println(pathParts)

	if len(pathParts) < 2 {
		return "", errors.New("invalid_playlist_url")
	}

	if pathParts[0] != "playlist" {
		return "", errors.New("invalid_playlist_url")
	}

	playlistId := pathParts[1]

	return playlistId, nil
}

type SpotifyAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func fetchSpotifyAuthToken(clientId, clientSecret string) (string, error) {
	params := url.Values{}
	params.Set("grant_type", "client_credentials")
	body := bytes.NewBufferString(params.Encode())

	client := &http.Client{}

	req, _ := http.NewRequest("POST", "https://accounts.spotify.com/api/token", body)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=utf-8")
	req.SetBasicAuth(clientId, clientSecret)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	var tokResp SpotifyAccessTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&tokResp)
	if err != nil {
		return "", err
	}

	return tokResp.AccessToken, nil
}

type SpotifyTrack struct {
	URI string `json:"uri"`
}

type SpotifyPlaylistItem struct {
	Track SpotifyTrack `json:"track"`
}

type SpotifyPlaylistResponseContainer struct {
	Items []SpotifyPlaylistItem `json:"items"`
	Next  string                `json:"next"`
}

func fetchSpotifyPlaylistItems(playlistId string, authToken string) ([]string, error) {
	client := &http.Client{}

	// Create request
	req, _ := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlistId+"/tracks", nil)
	req.Header.Add("Authorization", "Bearer "+authToken)

	// Fetch Request
	resp, err := client.Do(req)
	if err != nil {
		return []string{}, err
	}

	var respContainer SpotifyPlaylistResponseContainer
	err = json.NewDecoder(resp.Body).Decode(&respContainer)
	if err != nil {
		return []string{}, err
	}

	var uris []string

	for _, item := range respContainer.Items {
		uris = append(uris, item.Track.URI)
	}

	return uris, nil
}
