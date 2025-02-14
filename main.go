package main

import (
	"bufio"
	"bytes"
	"fmt"
	"image/color"
	"os/exec"
	"regexp"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Video Downloader")

	title := canvas.NewText("Video Downloader v1.0", color.White)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	platforms := []string{"YouTube", "TikTok"}
	platformSelector := widget.NewSelect(platforms, nil)
	platformSelector.SetSelected("YouTube")

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("Enter video URL...")

	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	outputLabel := widget.NewLabel("")

	browsers := []string{"Firefox", "Chrome", "Edge", "Vivaldi", "Brave", "Opera", "Safari", "Chromium", "Whale"}
	browserSelector := widget.NewSelect(browsers, nil)
	browserSelector.SetSelected("Firefox")

	browserSelector.Hide()

	platformLabel := widget.NewLabel("Select Platform:")
	browserLabel := widget.NewLabel("Select Your Browser:")
	browserLabel.Hide()

	platformSelector.OnChanged = func(selected string) {
		if selected == "TikTok" {
			browserLabel.Show()
			browserSelector.Show()
		} else {
			browserSelector.Hide()
		}
	}

	downloadButton := widget.NewButton("Download", func() {
		url := strings.TrimSpace(urlEntry.Text)
		platform := platformSelector.Selected
		browser := browserSelector.Selected

		if url == "" {
			outputLabel.SetText("Please enter a valid URL")
			return
		}

		progressBar.SetValue(0)
		progressBar.Show()

		go func() {
			err := downloadVideo(url, platform, browser, progressBar, outputLabel)
			progressBar.Hide()

			if err != nil {
				outputLabel.SetText(fmt.Sprintf("Error: %v", err))
			} else {
				outputLabel.SetText("Download complete!")

				dialog.ShowInformation("Download Complete", "The video has been downloaded successfully!", window)
			}
		}()
	})

	content := container.NewVBox(
		title,
		platformLabel,
		platformSelector,
		urlEntry,
		browserLabel,
		browserSelector,
		downloadButton,
		progressBar,
		outputLabel,
	)

	window.SetContent(content)
	window.Resize(fyne.NewSize(400, 250))
	window.ShowAndRun()
}

func downloadVideo(url, platform, browser string, progressBar *widget.ProgressBar, outputLabel *widget.Label) error {
	var cmd *exec.Cmd

	if platform == "YouTube" {
		cmd = exec.Command("yt-dlp", "--progress", "-o", "%(title)s.%(ext)s", url)
	} else if platform == "TikTok" {
		cmd = exec.Command("yt-dlp", "--progress", "-o", "%(title)s.%(ext)s", "--cookies-from-browser", browser, url)
	} else {
		return fmt.Errorf("Invalid platform selected")
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout pipe: %v", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to get stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %v", err)
	}

	progressRegex := regexp.MustCompile(`\[\s*download\s*\]\s*([\d.]+)%`)

	scanner := bufio.NewScanner(stdoutPipe)
	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			match := progressRegex.FindStringSubmatch(line)
			if len(match) > 1 {
				progress := match[1]
				var percent float64
				fmt.Sscanf(progress, "%f", &percent)
				progressBar.SetValue(percent / 100)
			}
		}
	}()

	var stderrBuffer bytes.Buffer
	scannerErr := bufio.NewScanner(stderrPipe)
	for scannerErr.Scan() {
		stderrBuffer.WriteString(scannerErr.Text() + "\n")
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("yt-dlp error: %v\n%s", err, stderrBuffer.String())
	}

	return nil
}
