package main

import (
	"github.com/alexflint/go-arg"
)

type Args struct {
	Push *PushCmd `arg:"subcommand:push"`
	Pull *PullCmd `arg:"subcommand:pull"`
	Trim *TrimCmd `arg:"subcommand:trim"`
}

type PushCmd struct {
	Repos         []string `arg:"positional"`
	LocalStorage  string   `arg:"-l,--local" help:"Path to local patches storage"`
	RemoteStorage string   `arg:"-r,--remote" help:"Path to root remote (smb) patches storage"`
}

type PullCmd struct {
	Repos         []string `arg:"positional"`
	LocalStorage  string   `arg:"-l,--local" help:"Path to local patches storage"`
	RemoteStorage string   `arg:"-r,--remote" help:"Path to root remote (smb) patches storage"`
	Strip         int      `arg:"-s,--strip" help:"number of leading path components to strip from paths parsed from the patch file"`
}

type TrimCmd struct {
	DaysDiff      int    `arg:"positional" help:"Days for delete old patches"`
	LocalStorage  string `arg:"-l,--local" help:"Path to local patches storage"`
	RemoteStorage string `arg:"-r,--remote" help:"Path to root remote (smb) patches storage"`
}

var args Args

// ln -s /var/run/user/1000/gvfs/smb-share\:server\=storage\,share\=public/ ~/etc/sync/patches
func main() {
	arg.MustParse(&args)

	if args.Push != nil {
		push(args.Push)
	} else if args.Pull != nil {
		pull(args.Pull)
	} else if args.Trim != nil {
		trim(args.Trim)
	}
}
