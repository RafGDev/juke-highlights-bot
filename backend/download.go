package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

const twitchClipsURL = "https://api.twitch.tv/kraken/clips/top"
const removeTsFiles = "rm files/clip_files/*.ts"

const clipDownloadURL = "https://clips-media-assets2.twitch.tv/%s.mp4"
const clipDownloadURLOffset = "https://clips-media-assets2.twitch.tv/%s-offset-%s.mp4"

// GetTwitchData sends a request to the twitch clips api to get the top
// clips for a particular game
func GetTwitchData(numOfClips string, gameType string, APIKey string) (*ClipJSON, error) {
	//Preparing for request to be sent
	req, err := http.NewRequest("GET", twitchClipsURL, nil)

	if err != nil {
		return nil, errors.New("Error with setting up request")
	}

	//Adding query headers
	req.Header.Set("Client-ID", APIKey)
	req.Header.Set("Accept", "application/vnd.twitchtv.v5+json")

	//Adding querystring values
	queryString := req.URL.Query()
	queryString.Add("game", gameType)
	queryString.Add("period", "day")
	queryString.Add("trending", "false")
	queryString.Add("limit", numOfClips)

	req.URL.RawQuery = queryString.Encode()

	client := &http.Client{}

	response, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
		return nil, errors.New("Error with sending request")
	}

	defer response.Body.Close()

	clipsFromRequest := new(ClipJSON)

	err = json.NewDecoder(response.Body).Decode(clipsFromRequest)

	fmt.Println(clipsFromRequest.Clips[0].Slug)

	if err != nil {
		return nil, errors.New("Error with decoding json")
	}

	return clipsFromRequest, nil
}

// DownloadClips loops through vidURL's and downloads the videos
func DownloadClips(clipsNoDuplicates *ClipJSON) (map[string]bool, error) {
	clipsDownloaded := make(map[string]bool)

	waitGroup := new(sync.WaitGroup)

	for _, currentClip := range clipsNoDuplicates.Clips {
		waitGroup.Add(1)
		go func(currentClip *Clip) {

			client := &http.Client{}

			response, err := client.Get(generateClipURL(currentClip.Thumbnail.Tiny, currentClip))

			if err != nil {
				panic("Err with getting video of id: " + currentClip.TrackingID)
			}

			defer response.Body.Close()

			if response.StatusCode != 200 {
				waitGroup.Done()
				return
			}

			file, err := os.Create(fmt.Sprintf("files/clip_files/%s.mp4", currentClip.TrackingID))

			if err != nil {
				waitGroup.Done()
				fmt.Println(err)
				panic("Error with creating file")
			}

			_, err = io.Copy(file, response.Body)

			if err != nil {
				waitGroup.Done()
				panic("error with copying body over")
			}

			clipsDownloaded[currentClip.TrackingID] = true
			file.Close()

			waitGroup.Done()

		}(currentClip)
	}

	waitGroup.Wait()

	return clipsDownloaded, nil
}

// OrderClips verifies and orders the clips in the order from the api request
func OrderClips(clipsFromRequest *ClipJSON, clipsDownloaded map[string]bool) *ClipJSON {
	orderedClips := new(ClipJSON)

	for _, clip := range clipsFromRequest.Clips {
		// If file does exist and was recognised as downloaded, then add to orderedClips array

		_, fileNotExists := os.Stat(fmt.Sprintf("files/clip_files/%s.mp4", clip.TrackingID))

		if fileNotExists == nil && clipsDownloaded[clip.TrackingID] {
			orderedClips.Clips = append(orderedClips.Clips, clip)
		}
	}
	return orderedClips
}

// SaveClipsToFile saves the ordered clips to clips_order.json
func SaveClipsToFile(orderedClips *ClipJSON) error {
	file, err := os.OpenFile("clips_order.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		return errors.New("Failed to open clips_order.json")
	}

	jsonParser := json.NewEncoder(file)

	err = jsonParser.Encode(*orderedClips)

	if err != nil {
		return errors.New("Error with sending to json")
	}

	return nil
}

// generateClipURL generates the clip url for the current clip
func generateClipURL(url string, currentClip *Clip) string {
	if !strings.Contains(url, "offset") {
		return fmt.Sprintf(clipDownloadURL, currentClip.TrackingID)
	}

	rafFromOffset := url[strings.Index(url, "offset"):]

	rafFromOffsetHyphen := rafFromOffset[strings.Index(rafFromOffset, "-")+1:]

	offset := rafFromOffsetHyphen[:strings.Index(rafFromOffsetHyphen, "-")]

	return fmt.Sprintf(clipDownloadURLOffset, currentClip.BroadcastID, offset)
}

// InitDownloadClips deletes files needed to download the clips
func InitDownloadClips() {
	// Deletes the files/clip_files directory
	os.RemoveAll("files/clip_files")

	// Delete clips_order.json if it doesn't exist
	if _, err := os.Stat("clips_order.json"); !os.IsNotExist(err) {
		os.Remove("clips_order.json")
	}

	// Make the clip files directory
	os.Mkdir("files/clip_files", 0700)
}

// RemoveDuplicates loops through clips in clipsFromRequest and removes duplicates
func RemoveDuplicates(clipsFromRequest *ClipJSON) (*ClipJSON, error) {
	// Map used to show which duplicates have already been added
	duplicatesMap := make(map[string]bool)
	clipsNoDuplicates := new(ClipJSON)

	for _, clip := range clipsFromRequest.Clips {
		isDuplicate := clipIsDuplicate(clipsFromRequest, clip, duplicatesMap)

		if isDuplicate == nil {
			clipsNoDuplicates.Clips = append(clipsNoDuplicates.Clips, clip)
			fmt.Println("The clip was not a duplicate")
		} else {
			fmt.Println("The clip was a duplicate")
		}

	}

	return clipsNoDuplicates, nil
}

// Returns the clip that is duplicate to currentClip
func clipIsDuplicate(clipJSON *ClipJSON, currentClip *Clip, duplicatesMap map[string]bool) *Clip {
	for _, clip := range clipJSON.Clips {
		// Check if the vod id is the same

		if clip.Vod.ID == currentClip.Vod.ID && clip.TrackingID != currentClip.TrackingID {
			//Check if the start time of a clip is 30 seconds or less apart

			//if _, found := duplicatesMap[clip.TrackingID]; found {

			if clip.Vod.Offset > currentClip.Vod.Offset {
				if clip.Vod.Offset-currentClip.Vod.Offset <= 30 {
					if _, found := duplicatesMap[clip.TrackingID]; found {
						return clip
					}

					duplicatesMap[currentClip.TrackingID] = true

					return nil

				}
			} else {
				if currentClip.Vod.Offset-clip.Vod.Offset <= 30 {

					if _, found := duplicatesMap[clip.TrackingID]; found {
						return clip
					}

					duplicatesMap[currentClip.TrackingID] = true

					return nil

				}
			}

		}
	}
	return nil
}
