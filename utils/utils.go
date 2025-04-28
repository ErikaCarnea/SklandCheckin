package utils

import (
	"os"
	"path/filepath"
	"time"
)

func getLastRunFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return "skland_lastRun"
	}
	exeDir := filepath.Dir(exePath)
	return filepath.Join(exeDir, ".skland_lastRun")
}

func HasRunToday() bool {
	filename := getLastRunFilePath()
	info, err := os.Stat(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		return false
	}
	modTime := info.ModTime().Local()
	now := time.Now().Local()
	return modTime.Year() == now.Year() &&
		modTime.Month() == now.Month() &&
		modTime.Day() == now.Day()
}

func MarkRun() error {
	filename := getLastRunFilePath()
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}
