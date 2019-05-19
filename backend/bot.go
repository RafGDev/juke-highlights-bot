package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const concatVidCommand = "ffmpeg -i \"concat:%sfiles/%s_outro.ts\" -qscale 0 -vf scale=1280:720 -bsf:a aac_adtstoasc files/juke_highlights_vid.mp4"
const removeMP4Files = "rm files/clip_files/%s.mp4"

// ClipsOrderFile is the filename of the ordered clips file
const ClipsOrderFile = "clips_order.json"

// ClipsShuffledFile is the filename of the shuffled clips file
const ClipsShuffledFile = "clips_shuffled.json"

// ClipJSON holds info of the clips
type ClipJSON struct {
	Clips []*Clip `json:"clips"`
}

// Clip holds infromation on the current clip
type Clip struct {
	TrackingID  string    `json:"tracking_id"`
	Slug        string    `json:"slug"`
	Thumbnail   Thumbnail `json:"thumbnails"`
	BroadcastID string    `json:"broadcast_id"`
	Title       string    `json:"title"`
	Vod         Vod       `json:"vod"`
	Duration    float64   `json:"duration"`
}

// Vod holds the information of the vod for the clip
type Vod struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Offset int    `json:"offset"`
}

// Thumbnail holds information of the thumbnail
type Thumbnail struct {
	Tiny   string `json:"tiny"`
	Medium string `json:"medium"`
}

// Games contains map of games supported for application
var Games = map[string]bool{
	"Fortnite": true,
	"VRChat":   true,
}

// Download downloads clips and then stores then in clips_order.json
func Download(numOfClips string, gameType string, APIKey string) (*ClipJSON, error) {
	// Mainly deletes things that may be used in this function
	InitDownloadClips()

	// Sends a request to twitch api clips for a specific game and returns response as struct
	clipsFromRequest, err := GetTwitchData(numOfClips, gameType, APIKey)

	if err != nil {
		return nil, err
	}

	// If clips link to same VOD and are 30 seconds apart, then assume duplicate clips
	clipsNoDuplicates, err := RemoveDuplicates(clipsFromRequest)

	if err != nil {
		return nil, err
	}

	fmt.Println(clipsNoDuplicates)

	// Loops through vids, downloads them, and returns map of vids successfully downloaded
	downloadedVids, err := DownloadClips(clipsNoDuplicates)

	if err != nil {
		return nil, err
	}

	// Orders the clips and verifies they exist
	orderedClips := OrderClips(clipsNoDuplicates, downloadedVids)

	// Saves ordered clips to file "clips_order.json"
	err = SaveClipsToFile(orderedClips)

	if err != nil {
		return nil, err
	}

	// If no errors, returned orderedClips
	return orderedClips, nil
}

// CompileClips creates video thumbnail and concats clips with each other along with an outro
func CompileClips(game string, customVidName string, firstClipIndex int) error {
	// Deletes things needed in function
	InitCompileClips()

	orderedClips, err := GetClipInfo(ClipsOrderFile)

	if err != nil {
		return err
	}

	err = CreateThumbnail(orderedClips, customVidName, firstClipIndex)

	if err != nil {
		return err
	}

	processedClips, err := ProcessClips(orderedClips, firstClipIndex)

	if err != nil {
		return err
	}

	err = ConcatClips(processedClips, game)

	if err != nil {
		return err
	}

	return nil
}

//GetClipInfo get an array of the clips
func GetClipInfo(filename string) (*ClipJSON, error) {
	// Will hold the order of the clips
	clipOrder := new(ClipJSON)

	file, err := os.Open(filename)

	if err != nil {
		return nil, errors.New("fileError")
	}

	err = json.NewDecoder(file).Decode(clipOrder)

	if err != nil {
		return nil, err
	}

	return clipOrder, nil
}

func removeVideo(trackingID string) error {
	err := exec.Command("bash", "-c", fmt.Sprintf(removeMP4Files, trackingID)).Run()

	if err != nil {
		return err
	}

	//Remove from clips_order.json
	clipJSON := new(ClipJSON)

	file, err := os.Open("clips_order.json")

	if err != nil {
		return errors.New("error with opening clips_order.json")
	}

	err = json.NewDecoder(file).Decode(clipJSON)

	if err != nil {
		return errors.New("Error with decoding json")
	}

	for i := 0; i < len(clipJSON.Clips); i++ {
		if clipJSON.Clips[i].TrackingID == trackingID {
			clipJSON.Clips = append(clipJSON.Clips[:i], clipJSON.Clips[i+1:]...)

			file, err := os.OpenFile("clips_order.json", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

			if err != nil {
				return errors.New("Error with opening file clips_order.json")
			}

			jsonParser := json.NewEncoder(file)

			err = jsonParser.Encode(*clipJSON)

			if err != nil {
				return errors.New("Error with sending json to clips_order.json")
			}

			return nil
		}
	}

	return errors.New("error with deleting from clips_json")
}

// UploadVid uploads the juke_highlights_vid.mp4 to youtube
func UploadVid(gameType string, customTitle string) error {
	shuffledClips, err := GetClipInfo(ClipsShuffledFile)

	if err != nil {
		return err
	}

	uploadData := PrepareUploadData(shuffledClips, gameType, customTitle)

	err = UploadToYoutube(uploadData)

	if err != nil {
		return err
	}

	return nil
}

// GetDuration gets the duration
func GetDuration(gameType string) (float64, error) {
	clipJSON := new(ClipJSON)

	file, err := os.Open("clips_order.json")

	if err != nil {
		return 0, errors.New("fileError")
	}

	err = json.NewDecoder(file).Decode(clipJSON)

	if err != nil {
		return 0, errors.New("fileError")
	}

	totalDuration := 0.0

	for _, clip := range clipJSON.Clips {
		totalDuration += clip.Duration
	}

	if gameType == "Fortnite" {
		totalDuration += 26
	} else if gameType == "VRChat" {
		totalDuration += 22
	}

	return totalDuration, nil
}

// CheckVidExists checks whether the video juke_highlights_vid.mp4 exists or not
func CheckVidExists() (bool, error) {
	if _, err := os.Stat("files/juke_highlights_vid.mp4"); err == nil {
		return true, nil
	}

	return false, nil
}
