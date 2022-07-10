package swr

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type ArchiveService struct {
	T *time.Ticker
}

var _archiver *ArchiveService

func Archiver() *ArchiveService {
	if _archiver == nil {
		_archiver = new(ArchiveService)
		_archiver.T = time.NewTicker(time.Duration(1) * time.Hour)
	}
	return _archiver
}
func BackupStart() {
	ar := Archiver()
	for {
		t := <-ar.T.C
		DoBackup(t)
		DoBackupCleanup(t)
	}
}

func DoBackup(t time.Time) {
	log.Printf("***** BACKUP STARTED *****\r\n")
	// Lock the database to prevent disk activity while we tar a backup archive
	db := DB()
	db.Lock()
	defer db.Unlock()

	stdout, err := exec.Command("tar", "-cJf", fmt.Sprintf("./backup/archive-%s.tar.xz", t.Format(time.RFC3339Nano)), "data/*").Output()

	ErrorCheck(err)

	fmt.Printf("%s\n", string(stdout))
	log.Printf("***** BACKUP COMPLETE *****\r\n")
}

func DoBackupCleanup(t time.Time) {
	archives, err := ioutil.ReadDir("./backup")
	ErrorCheck(err)
	for _, archive := range archives {
		p := archive.Name()
		p = strings.ReplaceAll(p, ".tar.xz", "")
		p = strings.ReplaceAll(p, "archive-", "")
		ar_time, err := time.Parse(time.RFC3339Nano, p)
		ErrorCheck(err)
		cut_time := t.Add(-72 * time.Hour) // 3 days worth of backups are stored.
		if ar_time.Before(cut_time) {
			err := os.Remove(fmt.Sprintf("./backup/%s", archive.Name()))
			ErrorCheck(err)
		}
	}
}
