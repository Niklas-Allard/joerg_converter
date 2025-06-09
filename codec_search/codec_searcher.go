package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type CodecDB map[string][]string

func main() {
	var rootDir string
	fmt.Println("Bitte gib den Pfad zum Verzeichnis ein:")
	fmt.Scanln(&rootDir)

	codecDB := make(CodecDB)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			codec, err := probeVideoCodec(path)
			if err == nil && codec != "" {
				fmt.Printf("Datei: %s | Codec: %s\n", path, codec)
				if !contains(codecDB[codec], path) {
					codecDB[codec] = append(codecDB[codec], path)
				}
			}
		}
		return nil
	})
	if err != nil {
		fmt.Printf("Fehler beim Durchlaufen des Verzeichnisses: %v\n", err)
	}

	jsonFile, err := os.Create("codecs.json")
	if err != nil {
		fmt.Printf("Fehler beim Erstellen der JSON-Datei: %v\n", err)
		return
	}
	defer jsonFile.Close()

	encoder := json.NewEncoder(jsonFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(codecDB); err != nil {
		fmt.Printf("Fehler beim Schreiben der JSON-Datei: %v\n", err)
	}

	fmt.Println("Fertig! Die Codecs wurden in codecs.json gespeichert.")
}

func probeVideoCodec(filePath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_streams", "-select_streams", "v", "-print_format", "json", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}
	var info struct {
		Streams []struct {
			CodecName string `json:"codec_name"`
		} `json:"streams"`
	}
	if err := json.Unmarshal(out.Bytes(), &info); err != nil {
		return "", err
	}
	if len(info.Streams) > 0 {
		return info.Streams[0].CodecName, nil
	}
	return "", nil
}

func contains(list []string, item string) bool {
	for _, v := range list {
		if v == item {
			return true
		}
	}
	return false
}
