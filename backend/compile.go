package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

const blurThumbnail = "convert files/thumbnail.jpg -blur 0x8 -fill black -colorize 50% files/thumbnail.jpg"
const addTextThumbnail = "convert files/thumbnail.jpg -fill black -font Norwester -pointsize 90 -gravity center -annotate 0 \"%s\n\n\nJuke Daily Highlights\" -fill white -annotate +2+2 \"%s\n\n\nJuke Daily Highlights\" files/thumbnail.jpg"
const firstFrameFromVid = "ffmpeg -ss 00:00:00 -i files/clip_files/%s.mp4 -vframes 1 -q:v 1 files/thumbnail.jpg"
const createTsFile = "ffmpeg -i files/clip_files/%s.mp4 -c copy -bsf:v h264_mp4toannexb -f mpegts files/clip_files/%s.ts"
const removeTSFiles = "rm files/clip_files/*.ts"

// InitCompileClips deletes all things needed by CompileClips
func InitCompileClips() {
	exec.Command("bash", "-c", removeTSFiles).Run()

	// If files/juke_highlights_vid.mp4 exists, then delete it
	if _, err := os.Stat("files/juke_highlights_vid.mp4"); !os.IsNotExist(err) {
		os.Remove("files/juke_highlights_vid.mp4")
	}

	// If thumbnail.jpg file exists, then delete it
	if _, err := os.Stat("files/thumbnail.jpg"); !os.IsNotExist(err) {
		os.Remove("files/thumbnail.jpg")
	}

	// If clips shuffled exists, then delete it
	if _, err := os.Stat(ClipsShuffledFile); !os.IsNotExist(err) {
		os.Remove(ClipsShuffledFile)
	}
}

// CreateThumbnail creates the video thumbnail
func CreateThumbnail(orderedClips *ClipJSON, customVidName string, clipIndex int) error {
	clipTitle := ""

	if customVidName != "" {
		clipTitle = customVidName
	} else {
		clipTitle = orderedClips.Clips[clipIndex].Title
	}
	clipID := orderedClips.Clips[clipIndex].TrackingID

	// Creates the thumbnail
	err := exec.Command("bash", "-c", fmt.Sprintf(firstFrameFromVid, clipID)).Run()

	if err != nil {
		fmt.Println(err)
		return errors.New("Error with executing creating thumbnail command")
	}

	//Blur thumbnail
	err = exec.Command("bash", "-c", blurThumbnail).Run()

	if err != nil {
		return errors.New("Error with blurring thumbnail command for thumbnail")
	}

	//Use the ImageMagick tool to add text
	err = exec.Command("bash", "-c", fmt.Sprintf(addTextThumbnail, strings.Title(clipTitle), strings.Title(clipTitle))).Run()

	if err != nil {
		return errors.New("Error with executing add text command for thumbnail")
	}

	return nil
}

// ConcatClips concatenates the clips along with the outro
func ConcatClips(orderedClips *ClipJSON, gameType string) error {
	for _, currentClip := range orderedClips.Clips {
		// Creates a ts file from the original mp4
		err := exec.Command("bash", "-c", fmt.Sprintf(createTsFile, currentClip.TrackingID, currentClip.TrackingID)).Run()

		if err != nil {
			return fmt.Errorf(fmt.Sprintf("Unable to convert %s.mp4 to ts", currentClip.TrackingID))
		}
	}

	err := exec.Command("bash", "-c", fmt.Sprintf(concatVidCommand, printFiles(orderedClips), gameType)).Run()

	if err != nil {
		return errors.New("Unable to concatenate ts files")

	}

	return nil
}

func printFiles(newClips *ClipJSON) string {
	output := ""

	for _, val := range newClips.Clips {
		output += fmt.Sprintf("files/clip_files/%s.ts|", val.TrackingID)
	}

	return output
}

// ProcessClips will put the clip at index firstClipIndex first and then shuffle the rest of the clips
func ProcessClips(orderedClips *ClipJSON, firstClipIndex int) (*ClipJSON, error) {
	clipToPutFirst := orderedClips.Clips[firstClipIndex]

	// Remove the element at index firstClipIndex
	orderedClips.Clips = append(orderedClips.Clips[:firstClipIndex], orderedClips.Clips[firstClipIndex+1:]...)

	//Next shuffle the array
	for i := range orderedClips.Clips {
		j := rand.Intn(i + 1)
		orderedClips.Clips[i], orderedClips.Clips[j] = orderedClips.Clips[j], orderedClips.Clips[i]
	}
	orderedClips.Clips = append([]*Clip{clipToPutFirst}, orderedClips.Clips...)

	file, err := os.OpenFile(ClipsShuffledFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)

	if err != nil {
		return nil, errors.New("Failed to open clips_shuffled.json")
	}

	jsonParser := json.NewEncoder(file)

	err = jsonParser.Encode(*orderedClips)

	if err != nil {
		return nil, errors.New("Error with sending to json")
	}

	return orderedClips, nil
}
