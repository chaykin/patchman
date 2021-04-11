package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func push(pushCmd *PushCmd) {
	mountRemoteStorageSafe(pushCmd.RemoteStorage)

	for _, repo := range pushCmd.Repos {
		info, err := exec.Command("svn", "info", repo).CombinedOutput()
		check(err)

		var name, branch string
		name, branch, err = getBranchFromInfoCommand(info)
		check(err)

		now := formatTime(time.Now())
		patchName := fmt.Sprintf("%s-%s-%s.patch", now, name, branch)

		fmt.Printf("Creating patch %s...\n", patchName)
		var patch []byte
		patch, err = exec.Command("svn", "diff", "--git", repo).CombinedOutput()
		check(err)

		var patchFile *os.File
		patchFile, err = os.Create(filepath.Join(pushCmd.LocalStorage, patchName))
		check(err)
		_, err = patchFile.Write(patch)
		check(err)
		check(patchFile.Close())
	}
	println("Done")
}
