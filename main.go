package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Benutzer gibt Quell- und Zieldateitypen an
	var sourceExt, targetExt string
	fmt.Println("Geben Sie den Quell-Dateityp (z.B. mkv) ein:")
	fmt.Scanln(&sourceExt)
	fmt.Println("Geben Sie den Ziel-Dateityp (z.B. mp4) ein:")
	fmt.Scanln(&targetExt)

	// Benutzer gibt das Verzeichnis an
	var rootDir string
	fmt.Println("Geben Sie das Verzeichnis an, das durchsucht werden soll:")
	fmt.Scanln(&rootDir)

	// Rekursive Verarbeitung starten
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Überprüfen, ob die Datei die Quell-Erweiterung hat
		if !info.IsDir() && strings.HasSuffix(info.Name(), "."+sourceExt) {
			convertFile(path, sourceExt, targetExt)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Fehler beim Durchsuchen des Verzeichnisses: %v\n", err)
	}
}

func convertFile(filePath, sourceExt, targetExt string) {
	// Zielpfad erstellen
	targetPath := strings.TrimSuffix(filePath, "."+sourceExt) + "." + targetExt

	// ffmpeg-Befehl ausführen
	cmd := exec.Command("ffmpeg", "-i", filePath, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Printf("Konvertiere %s nach %s...\n", filePath, targetPath)
	if err := cmd.Run(); err != nil {
		fmt.Printf("Fehler beim Konvertieren von %s: %v\n", filePath, err)
	} else {
		// Originaldatei löschen, wenn die Konvertierung erfolgreich war
		if err := os.Remove(filePath); err != nil {
			fmt.Printf("Fehler beim Löschen der Originaldatei %s: %v\n", filePath, err)
		}
	}
}