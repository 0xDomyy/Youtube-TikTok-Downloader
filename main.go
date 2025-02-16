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
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	window := myApp.NewWindow("Video Downloader")

	window.SetIcon(theme.FyneLogo())

	title := canvas.NewText("Video Downloader v1.0", color.White)
	title.Alignment = fyne.TextAlignCenter
	title.TextStyle = fyne.TextStyle{Bold: true}

	text := canvas.NewText("Developed by @domyy.krnl", color.White)
	text.Alignment = fyne.TextAlignCenter
	text.TextStyle = fyne.TextStyle{Italic: true}

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
			browserLabel.Hide()
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
		text,
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

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("Error creating stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Error starting yt-dlp: %v", err)
	}

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		progressRegex := regexp.MustCompile(`(\d+\.\d+)%`)

		for scanner.Scan() {
			line := scanner.Text()
			match := progressRegex.FindStringSubmatch(line)
			if len(match) > 1 {
				var percent float64
				fmt.Sscanf(match[1], "%f", &percent)

				progressBar.SetValue(percent / 100)
			}
		}
	}()

	// Aspetta che il comando finisca
	err = cmd.Wait()
	progressBar.SetValue(1.0) // Completa la progress bar al 100%

	// Gestisci eventuali errori
	var stderrBuffer bytes.Buffer
	scannerErr := bufio.NewScanner(stderrPipe)
	for scannerErr.Scan() {
		stderrBuffer.WriteString(scannerErr.Text() + "\n")
	}

	if err != nil {
		return fmt.Errorf("yt-dlp error: %v\n%s", err, stderrBuffer.String())
	}

	return nil
}
