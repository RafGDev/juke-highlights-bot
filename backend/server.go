package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

const apiKey = "hbc7t5net9h5oy4n8hj6zwphfh2kzw"

type successResponse struct {
	Response string      `json:"response"`
	Data     interface{} `json:"data"`
}

type errorResponse struct {
	Response  string `json:"response"`
	ErrorType string `json:"errorType"`
}

// deleteVideoHandler deletes a clip
func deleteClipHandler(w http.ResponseWriter, r *http.Request) {
	trackingID := r.URL.Query().Get("tracking_id")

	err := removeVideo(trackingID)

	if err != nil {
		json.NewEncoder(w).Encode(successResponse{"error", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(successResponse{"success", nil})
}

// downloadClipsHandler downloads most popular clips for that day for a specific game
func downloadClipsHandler(w http.ResponseWriter, r *http.Request) {
	numOfClips := r.URL.Query().Get("numOfClips")
	gameType := r.URL.Query().Get("gameType")

	if !Games[gameType] {
		json.NewEncoder(w).Encode(errorResponse{"error", "gameTypeError"})
		return
	}

	clipJSON, err := Download(numOfClips, gameType, apiKey)

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(successResponse{"success", *clipJSON})
}

// compileClipsHandler both concats the clips along with the outro and creates a thumbanil
// from the first clip. The user is also able to specific the first clip
func compileClipsHandler(w http.ResponseWriter, r *http.Request) {
	gameType := r.URL.Query().Get("gameType")
	customTitle := r.URL.Query().Get("customTitle")
	firstClipIndexString := r.URL.Query().Get("firstClipIndex")

	if !Games[gameType] {
		json.NewEncoder(w).Encode(errorResponse{"error", "Please choose out of the supported games please"})
		return
	}

	firstClipIndex, err := strconv.Atoi(firstClipIndexString)

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", "Please choose out of the supported games please"})
	}

	err = CompileClips(gameType, customTitle, firstClipIndex)

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(successResponse{"success", nil})
}

// uploadVidHandler is used to upload the finished video to youtube
func uploadVidHandler(w http.ResponseWriter, r *http.Request) {
	gameType := r.URL.Query().Get("gameType")
	customTitle := r.URL.Query().Get("customTitle")

	err := UploadVid(gameType, customTitle)

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(successResponse{"success", nil})
}

// getClipsHandler is used to get a list of the current clip names
func getClipsHandler(w http.ResponseWriter, r *http.Request) {
	clipJSON, err := GetClipInfo(ClipsOrderFile)

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(successResponse{"success", clipJSON})
}

// getDurationHandler will return how long the final video WILL be
func getDurationHandler(w http.ResponseWriter, r *http.Request) {
	gameType := r.URL.Query().Get("gameType")
	duration, err := GetDuration(gameType)

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", err.Error()})
	}

	json.NewEncoder(w).Encode(successResponse{"success", duration})
}

// Checks whether the final video exists, return true if it does, else otherwise
func checkVidExistsHandler(w http.ResponseWriter, r *http.Request) {
	doesExist, err := CheckVidExists()

	if err != nil {
		json.NewEncoder(w).Encode(errorResponse{"error", err.Error()})
		return
	}

	json.NewEncoder(w).Encode(successResponse{"success", doesExist})
}

// StartServer starts the server
func StartServer() {
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/api/downloadClips", downloadClipsHandler)
	r.HandleFunc("/api/concatClips", compileClipsHandler)
	r.HandleFunc("/api/getClips", getClipsHandler)
	r.HandleFunc("/api/deleteVideo", deleteClipHandler)
	r.HandleFunc("/api/uploadVid", uploadVidHandler)
	r.HandleFunc("/api/getDuration", getDurationHandler)
	r.HandleFunc("/api/checkVidExists", checkVidExistsHandler)

	r.PathPrefix("/api").Handler(http.StripPrefix("/api", http.FileServer(http.Dir("files"))))
	http.ListenAndServe(":8080", r)
}
