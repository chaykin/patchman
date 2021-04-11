package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func pull(pullCmd *PullCmd) {
	mountRemoteStorageSafe(pullCmd.RemoteStorage)

	patchDate, patchNames := getPatches(pullCmd.LocalStorage)

	thresholdDate := formatTime(time.Now().Add(-1 * time.Hour * 24 * 5))
	if patchDate < thresholdDate {
		fmt.Printf("Found patch is too old: %s\n", patchDate)
		printPatches(patchNames)
		println("Do You want to apply it [y/n]?")

		if !readUserChoice() {
			println("Skip patch applying")
			return
		}
	}

	for _, repo := range pullCmd.Repos {
		applyPatch(repo, patchDate, pullCmd)
	}
	println("Done")
}

func getPatches(localStorage string) (string, []string) {
	patchDate := "0000-00-00"
	patchNames := make([]string, 0)
	err := filepath.WalkDir(localStorage, func(path string, info fs.DirEntry, err error) error {
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
			if date > patchDate {
				patchDate = date
				patchNames = []string{path}
			} else if date == patchDate {
				patchNames = append(patchNames, path)
			}
		}
		return nil
	})
	check(err)

	return patchDate, patchNames
}

func printPatches(patchNames []string) {
	for _, patchName := range patchNames {
		_, name := filepath.Split(patchName)
		fmt.Printf(" %s\n", name)
	}
}

func readUserChoice() bool {
	reader := bufio.NewReader(os.Stdin)
	inputText, err := reader.ReadString('\n')
	check(err)

	inputText = strings.ToLower(inputText)
	return strings.HasPrefix(inputText, "y")
}

func applyPatch(repo, patchDate string, pullCmd *PullCmd) {
	info, err := exec.Command("svn", "info", repo).CombinedOutput()
	check(err)

	var name, branch string
	name, branch, err = getBranchFromInfoCommand(info)
	check(err)

	patchName := fmt.Sprintf("%s-%s-%s.patch", patchDate, name, branch)

	fmt.Printf("Applying patch %s...\n", patchName)
	patchFullPath := filepath.Join(pullCmd.LocalStorage, patchName)
	var applyResult []byte
	applyResult, err = exec.Command("svn", "patch", "--strip", strconv.Itoa(pullCmd.Strip),
		patchFullPath, repo).CombinedOutput()
	fmt.Printf("%s", applyResult)
	if err != nil {
		println("An error occurred while applying patch")
	}
}
