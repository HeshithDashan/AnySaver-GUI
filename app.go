package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kkdai/youtube/v2"
)

type App struct {
	ctx context.Context
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) DownloadVideo(url string) string {
	fmt.Printf("\n[DEBUG] Starting download for: %s\n", url)
	
	if strings.Contains(url, "youtu") {
		return a.downloadYouTube(url)
	}
	return a.downloadDirectFile(url)
}

func (a *App) downloadYouTube(url string) string {
	client := youtube.Client{}
	
	video, err := client.GetVideo(url)
	if err != nil {
		fmt.Println("[ERROR] YouTube Info Error:", err)
		return "වැරදියි: වීඩියෝ විස්තර ලබාගත නොහැක."
	}

	fmt.Println("[DEBUG] Video Title:", video.Title)

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return "වැරදියි: ගැලපෙන වීඩියෝ Format එකක් හමු නොවීය."
	}

	stream, _, err := client.GetStream(video, &formats[0])
	if err != nil {
		fmt.Println("[ERROR] Stream Error:", err)
		return "වැරදියි: Stream එක ආරම්භ කළ නොහැක."
	}
	defer stream.Close()

	homeDir, _ := os.UserHomeDir()

	saveName := "AnySaver_" + time.Now().Format("150405") + ".mp4"
	fileName := filepath.Join(homeDir, "Desktop", saveName)
	
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[ERROR] File Creation Error:", err)
		return "වැරදියි: File එක සෑදිය නොහැක."
	}
	defer file.Close()

	fmt.Println("[DEBUG] Downloading to Desktop...")
	_, err = io.Copy(file, stream)
	if err != nil {
		fmt.Println("[ERROR] Copy Error:", err)
		return "වැරදියි: Download වීම අතරමග නැවතුණි."
	}

	fmt.Println("[DEBUG] Success!")
	return "සාර්ථකයි! වීඩියෝව Desktop එකේ " + saveName + " නමින් සේව් වුණා."
}

func (a *App) downloadDirectFile(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "Error: " + err.Error()
	}
	defer resp.Body.Close()

	homeDir, _ := os.UserHomeDir()
	fileName := filepath.Join(homeDir, "Desktop", "AnySaver_File.mp4")
	
	file, err := os.Create(fileName)
	if err != nil {
		return "Error: " + err.Error()
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "Error: " + err.Error()
	}

	return "සාර්ථකයි! ෆයිල් එක Desktop එකට සේව් වුණා."
}