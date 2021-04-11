package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"time"
)

const datePatternLength = 10

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func mountRemoteStorageSafe(remoteStorage string) {
	println("Mounting remote storage...")
	check(mountRemoteStorage(remoteStorage))
	println("Mounted")
}

func checkMountRemoteStorage(remoteStorage string) (bool, error) {
	out, err := exec.Command("gio", "mount", "-l", remoteStorage).CombinedOutput()
	if err != nil {
		return false, err
	}

	r := regexp.MustCompile(`Mount\(\d+\):\s+.+?\s+->\s+` + regexp.QuoteMeta(remoteStorage))
	return r.Match(out), nil
}

func mountRemoteStorage(remoteStorage string) error {
	mounted, err := checkMountRemoteStorage(remoteStorage)
	if err != nil {
		return err
	}

	if !mounted {
		return exec.Command("gio", "mount", "-a", remoteStorage).Run()
	}
	return nil
}

func getBranchFromInfoCommand(info []byte) (string, string, error) {
	r := regexp.MustCompile(`URL:\s+.+/(.+)/(.+)\n`)
	findRes := r.FindStringSubmatch(string(info))
	if len(findRes) <= 2 {
		return "", "", fmt.Errorf("URL not specified in svn info")
	}

	return findRes[1], findRes[2], nil
}

func formatTime(time time.Time) string {
	return time.Format("2006-01-02")
}
