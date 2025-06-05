package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// User provides source and target file types
	var sourceExt, targetExt, userCodec string
	fmt.Println("Enter the source file type (e.g., mkv):")
	fmt.Scanln(&sourceExt)
	fmt.Println("Enter the target file type (e.g., mp4):")
	fmt.Scanln(&targetExt)
	fmt.Println("Enter the codec to check (e.g., avc):")
	fmt.Scanln(&userCodec)

	// User provides the directory
	var rootDir string
	fmt.Println("Enter the directory to be scanned:")
	fmt.Scanln(&rootDir)

	// Start recursive processing
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Check if the file has the source extension
		if !info.IsDir() && strings.HasSuffix(info.Name(), "."+sourceExt) {
			convertFile(path, sourceExt, targetExt, userCodec)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning the directory: %v\n", err)
	}
}

func convertFile(filePath, sourceExt, targetExt, userCodec string) {
	// Create target path
	targetPath := strings.TrimSuffix(filePath, "."+sourceExt) + "_converted." + targetExt

	// Read the codec of the video
	codec, err := getVideoCodec(filePath)
	if err != nil {
		fmt.Printf("Error reading the codec of %s: %v\n", filePath, err)
		return
	}

	// Debugging output
	fmt.Printf("Extracted codec: %s\n", codec)
	fmt.Printf("User-defined codec: %s\n", userCodec)

	// If the codec matches the user-defined codec, convert to AVC
	if codec == userCodec {
		fmt.Printf("The video %s has the codec %s, converting to AVC.\n", filePath, codec)
		cmd := exec.Command("ffmpeg", "-y", "-i", filePath, "-c:v", "libx264", "-preset", "medium", "-crf", "23", "-c:a", "aac", "-b:a", "128k", targetPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("Error converting %s to AVC: %v\n", filePath, err)
			return
		}
	} else {
		// Execute ffmpeg command
		fmt.Println("Codec not the user preferred")
	}

	// Delete the original file if the conversion was successful
	if err := os.Remove(filePath); err != nil {
		fmt.Printf("Error deleting the original file %s: %v\n", filePath, err)
	}
}

func getVideoCodec(filePath string) (string, error) {
	// Execute ffprobe command
	cmd := exec.Command("ffprobe", "-v", "error", "-show_streams", "-select_streams", "v", "-print_format", "json", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error executing ffprobe: %v\n%s", err, out.String())
	}

	// Parse JSON output
	var info struct {
		Streams []struct {
			CodecName string `json:"codec_name"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		return "", fmt.Errorf("Error parsing ffprobe output: %v", err)
	}

	// Find codec from the first video stream
	if len(info.Streams) > 0 {
		return info.Streams[0].CodecName, nil
	}

	return "", fmt.Errorf("No video stream found")
}