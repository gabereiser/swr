/*  Star Wars Role-Playing Mud
 *  Copyright (C) 2022 @{See Authors}
 *
 *  This program is free software: you can redistribute it and/or modify
 *  it under the terms of the GNU General Public License as published by
 *  the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
 *
 *  This program is distributed in the hope that it will be useful,
 *  but WITHOUT ANY WARRANTY; without even the implied warranty of
 *  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 *  GNU General Public License for more details.
 *
 *  You should have received a copy of the GNU General Public License
 *  along with this program.  If not, see <https://www.gnu.org/licenses/>.
 *
 */
package swr

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
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
		_archiver.T = time.NewTicker(1 * time.Hour)
	}
	return _archiver
}
func StartBackup() {
	ar := Archiver()
	go func() {
		for {
			t := <-ar.T.C
			DoBackup(t)
			DoBackupCleanup(t)
		}
	}()
	log.Printf("Backup service started.\n")
	DoBackupCleanup(time.Now())
	DoBackup(time.Now())
}

func DoBackup(t time.Time) {
	log.Printf("***** BACKUP STARTED *****\r\n")

	DB().Save()
	create_backup(t)
	runtime.GC()
	log.Printf("***** BACKUP COMPLETE *****\r\n")
}

func DoBackupCleanup(t time.Time) {
	archives, err := os.ReadDir("backup")
	ErrorCheck(err)
	for _, archive := range archives {
		if strings.HasSuffix(archive.Name(), ".tar.gz") {
			p := archive.Name()
			p = strings.ReplaceAll(p, ".tar.gz", "")
			ar_time, err := time.Parse("2006_01_02_15_04_05", p)
			ErrorCheck(err)
			cut_time := t.Add(-72 * time.Hour) // 3 days worth of backups are stored.
			if ar_time.Before(cut_time) {
				if runtime.GOOS == "windows" {
					err := os.Remove(fmt.Sprintf("backup\\%s", archive.Name()))
					ErrorCheck(err)
				} else {
					err := os.Remove(fmt.Sprintf("backup/%s", archive.Name()))
					ErrorCheck(err)
				}

			}
		}

	}
}
func create_backup(t time.Time) {
	files := []string{}
	filepath.Walk("data", func(path string, info fs.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	// Create output file
	out, err := os.Create(sprintf("backup/%s.tar.gz", t.Format("2006_01_02_15_04_05")))
	if err != nil {
		log.Fatalln("Error writing archive:", err)
	}
	defer out.Close()

	// Create the archive and write the output to the "out" Writer
	err = create_archive(files, out)
	if err != nil {
		log.Fatalln("Error creating archive:", err)
	}

	fmt.Println("Archive created successfully")
}
func create_archive(files []string, buf io.Writer) error {
	// Create new Writers for gzip and tar
	// These writers are chained. Writing to the tar writer will
	// write to the gzip writer which in turn will write to
	// the "buf" writer
	gw := gzip.NewWriter(buf)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	// Iterate over files and add them to the tar archive
	for _, file := range files {
		err := archive_add(tw, file)
		if err != nil {
			return err
		}
	}

	return nil
}

func archive_add(tw *tar.Writer, filename string) error {
	// Open the file which will be written into the archive
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Get FileInfo about our file providing file size, mode, etc.
	info, err := file.Stat()
	if err != nil {
		return err
	}

	// Create a tar Header from the FileInfo data
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		return err
	}

	// Use full path as name (FileInfoHeader only takes the basename)
	// If we don't do this the directory strucuture would
	// not be preserved
	// https://golang.org/src/archive/tar/common.go?#L626
	header.Name = filename

	// Write file header to the tar archive
	err = tw.WriteHeader(header)
	if err != nil {
		return err
	}

	// Copy file content to tar archive
	_, err = io.Copy(tw, file)
	if err != nil {
		return err
	}

	return nil
}
