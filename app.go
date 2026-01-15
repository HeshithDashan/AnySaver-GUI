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

func (a *App) SelectFolder() string {
	folder, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "වීඩියෝව සේව් කරන්න ඕන තැන තෝරන්න",
	})
	if err != nil {
		fmt.Println("[ERROR] Folder Selection Error:", err)
		return ""
	}
	return folder
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

func (a *App) DownloadVideo(url string, savePath string) string {
	fmt.Printf("\n[DEBUG] Starting download for: %s\n", url)
	
	if savePath == "" {
		homeDir, _ := os.UserHomeDir()
		savePath = filepath.Join(homeDir, "Desktop")
	}

	if strings.Contains(url, "youtu") {
		return a.downloadYouTube(url, savePath)
	}
	return a.downloadDirectFile(url, savePath)
}

func (a *App) downloadYouTube(url string, savePath string) string {
	client := youtube.Client{}
	
	fmt.Println("[DEBUG] Step 1: Fetching video info from YouTube...")
	video, err := client.GetVideo(url)
	if err != nil {
		return "වැරදියි: වීඩියෝ විස්තර ලබාගත නොහැක."
	}

	fmt.Printf("[DEBUG] Step 2: Found Video - Title: %s\n", video.Title)

	formats := video.Formats.WithAudioChannels()
	if len(formats) == 0 {
		return "වැරදියි: ගැලපෙන වීඩියෝ Format එකක් හමු නොවීය."
	}

	stream, size, err := client.GetStream(video, &formats[0])
	if err != nil {
		return "වැරදියි: Stream එක ආරම්භ කළ නොහැක."
	}
	defer stream.Close()

	saveName := "AnySaver_" + time.Now().Format("150405") + ".mp4"
	fileName := filepath.Join(savePath, saveName)
	
	fmt.Printf("[DEBUG] Step 5: Creating file at %s\n", fileName)
	file, err := os.Create(fileName)
	if err != nil {
		return "වැරදියි: File එක සෑදිය නොහැක."
	}
	defer file.Close()

	pw := &progressWriter{
		total: size,
		ctx:   a.ctx,
	}

	_, err = io.Copy(file, io.TeeReader(stream, pw))
	if err != nil {
		return "වැරදියි: Download වීම අතරමග නැවතුණි."
	}

	fmt.Println("[DEBUG] SUCCESS: Download completed!")
	return "සාර්ථකයි! වීඩියෝව " + saveName + " නමින් සේව් වුණා."
}

func (a *App) downloadDirectFile(url string, savePath string) string {
	resp, err := http.Get(url)
	if err != nil {
		return "Error: " + err.Error()
	}
	defer resp.Body.Close()

	fileName := filepath.Join(savePath, "AnySaver_File.mp4")
	
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

	return "සාර්ථකයි! ෆයිල් එක සේව් වුණා."
}