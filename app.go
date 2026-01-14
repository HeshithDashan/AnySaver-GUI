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
	"github.com/wailsapp/wails/v2/pkg/runtime"
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

type progressWriter struct {
	total      int64
	downloaded int64
	ctx        context.Context
}

func (pw *progressWriter) Write(p []byte) (int, error) {
	n := len(p)
	pw.downloaded += int64(n)
	if pw.total > 0 {
		percentage := float64(pw.downloaded) / float64(pw.total) * 100
		runtime.EventsEmit(pw.ctx, "download-progress", int(percentage))
	}
	return n, nil
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
	
	fmt.Println("[DEBUG] Step 1: Fetching video info from YouTube...")
	video, err := client.GetVideo(url)
	if err != nil {
		fmt.Println("[ERROR] Step 1 Failed - YouTube Info Error:", err)
		return "වැරදියි: වීඩියෝ විස්තර ලබාගත නොහැක."
	}

	fmt.Printf("[DEBUG] Step 2: Found Video - Title: %s\n", video.Title)

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		fmt.Println("[ERROR] No suitable formats found.")
		return "වැරදියි: ගැලපෙන වීඩියෝ Format එකක් හමු නොවීය."
	}

	fmt.Println("[DEBUG] Step 3: Getting video stream...")
	stream, size, err := client.GetStream(video, &formats[0])
	if err != nil {
		fmt.Println("[ERROR] Step 3 Failed - Stream Error:", err)
		return "වැරදියි: Stream එක ආරම්භ කළ නොහැක."
	}
	defer stream.Close()

	fmt.Printf("[DEBUG] Step 4: Stream ready. File size: %d bytes\n", size)

	homeDir, _ := os.UserHomeDir()
	saveName := "AnySaver_" + time.Now().Format("150405") + ".mp4"
	fileName := filepath.Join(homeDir, "Desktop", saveName)
	
	fmt.Println("[DEBUG] Step 5: Creating file on Desktop...")
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Println("[ERROR] Step 5 Failed - File Creation Error:", err)
		return "වැරදියි: File එක සෑදිය නොහැක."
	}
	defer file.Close()

	pw := &progressWriter{
		total: size,
		ctx:   a.ctx,
	}

	fmt.Println("[DEBUG] Step 6: Starting data copy (TeeReader)...")
	_, err = io.Copy(file, io.TeeReader(stream, pw))
	if err != nil {
		fmt.Println("[ERROR] Step 6 Failed - Copy Error:", err)
		return "වැරදියි: Download වීම අතරමග නැවතුණි."
	}

	fmt.Println("[DEBUG] SUCCESS: Download completed!")
	return "සාර්ථකයි! වීඩියෝව Desktop එකේ " + saveName + " නමින් සේව් වුණා."
}

func (a *App) downloadDirectFile(url string) string {
	fmt.Println("[DEBUG] Starting direct file download...")
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

	pw := &progressWriter{
		total: resp.ContentLength,
		ctx:   a.ctx,
	}

	_, err = io.Copy(file, io.TeeReader(resp.Body, pw))
	if err != nil {
		return "Error: " + err.Error()
	}

	return "සාර්ථකයි! ෆයිල් එක Desktop එකට සේව් වුණා."
}