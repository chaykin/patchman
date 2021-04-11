package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func trim(trimCmd *TrimCmd) {
	mountRemoteStorageSafe(trimCmd.RemoteStorage)

	thresholdDate := formatTime(time.Now().Add(-1 * time.Hour * 24 * time.Duration(trimCmd.DaysDiff)))

	trimPatchNames := make([]string, 0)
	err := filepath.WalkDir(trimCmd.LocalStorage, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if strings.HasSuffix(path, info.Name()) {
				return nil
			}
			return filepath.SkipDir
		} else {
			if len(info.Name()) < datePatternLength {
				// skip other files
				return nil
			}
			date := info.Name()[:datePatternLength]
			if date < thresholdDate {
				trimPatchNames = append(trimPatchNames, path)
			}
		}
		return nil
	})
	check(err)

	for _, patchName := range trimPatchNames {
		fmt.Printf("Trim patch %s...\n", patchName)
		check(os.Remove(patchName))
	}
	println("Done")
}
