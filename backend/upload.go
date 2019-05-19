package main

import (
	"errors"
	"fmt"
	"math"
	"os/exec"
)

var description = "We are Juke Daily Highlights. Our videos are intended for people who do not have time to watch every %s Stream. Our videos fill you in with the happenings over the past 24 hours in the Fortnite twitch community! Thanks for watching and if you enjoyed or want to be kept up to date, remember to like and subscribe to our channel to support our growth!"

const youtubeUploadCommand = "youtube-upload --title=\"%s - %s Daily Highlights\" --description=\"%s\" --category=Gaming --tags=\"%s\" --default-language=\"en\" --credentials-file=%s --client-secrets=%s --thumbnail=files/thumbnail.jpg files/juke_highlights_vid.mp4"
const vrChatTags = "VRChat, Daily, Highlights, Twitch, Juke, Virtual, Reality"
const fortniteTags = "Fortnite, Daily, Highlights, Twitch, Juke, Virtual, Reality"

// UploadConfig holds the config info needed to upload
type UploadConfig struct {
	Title             string
	GameType          string
	Description       string
	Tags              string
	CredentialsFile   string
	ClientSecretsFile string
}

// PrepareUploadData prepares the data to be uploaded to youtube
func PrepareUploadData(shuffledData *ClipJSON, gameType string, customTitle string) *UploadConfig {
	tags := ""
	credentialsFile := ""
	clientSecretsFile := ""
	vidDescription := fmt.Sprintf(description, gameType) + descriptionContents(shuffledData)

	if gameType == "Fortnite" {
		tags = fortniteTags
		credentialsFile = "fortnite_credentials.json"
		clientSecretsFile = "fortnite_secrets.json"
	} else {
		tags = vrChatTags

		credentialsFile = "vrchat_credentials.json"
		clientSecretsFile = "vrchat_secrets.json"
	}

	title := ""

	if customTitle != "" {
		title = customTitle
	} else {
		title = shuffledData.Clips[0].Title
	}

	uploadConfig := &UploadConfig{title, gameType, vidDescription, tags, credentialsFile, clientSecretsFile}

	return uploadConfig
}

// UploadToYoutube uploads the video to youtube
func UploadToYoutube(uploadConfig *UploadConfig) error {
	err := exec.Command("bash", "-c", fmt.Sprintf(youtubeUploadCommand, uploadConfig.Title, uploadConfig.GameType, uploadConfig.Description, uploadConfig.Tags, uploadConfig.CredentialsFile, uploadConfig.ClientSecretsFile)).Run()

	fmt.Printf(youtubeUploadCommand, uploadConfig.Title, uploadConfig.GameType, uploadConfig.Description, uploadConfig.Tags, uploadConfig.CredentialsFile, uploadConfig.ClientSecretsFile)

	if err != nil {
		return errors.New("Error with uploading to youtube")
	}

	return nil
}

// Contains the dynamic contents of the description which is a list of clip names, vid links, and time frame in video
func descriptionContents(shuffledData *ClipJSON) string {
	contents := `\n\n` // Start with new line to seperate message at top
	currentTime := 0

	for _, clip := range shuffledData.Clips {
		contents += fmt.Sprintf("Title: %s\nVod: %s\nTime: %s\n\n", clip.Title, clip.Vod.URL, youtubeTimify(currentTime))
		// Convert clip.Duration to milliseconds, then round in to the nearest millisecond
		currentTime += int(math.Floor(clip.Duration + 0.5))
	}

	return contents
}

// youtubeTimify converts a number of seconds to format: mm:ss
func youtubeTimify(seconds int) string {
	justSeconds := seconds % 60
	justMinutes := (seconds - justSeconds) / 60

	secondsString := ""
	minutesString := ""

	if justSeconds < 10 {
		secondsString = fmt.Sprintf("0%d", justSeconds)
	} else {
		secondsString = fmt.Sprintf("%d", justSeconds)
	}

	if justMinutes < 10 {
		minutesString = fmt.Sprintf("0%d", justMinutes)
	} else {
		minutesString = fmt.Sprintf("%d", justMinutes)
	}

	return fmt.Sprintf("%s:%s", minutesString, secondsString)
}
