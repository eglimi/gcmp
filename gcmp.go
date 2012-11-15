/*
gcmp is a tool to compare files in two directories and copy new files into a
third directory.
*/
package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

type FileInfoPath struct {
	fi   os.FileInfo
	path string
}

var verbose bool
var icase bool

func main() {

	// Command line parsing
	var orig = flag.String("orig", "", "The directory with the original file")
	var new = flag.String("new", "", "The directory with the new files")
	var store = flag.String("out", "./file_diff", "The directory to copy the files to. Defaults to './file_diff'")
	flag.BoolVar(&verbose, "verbose", false, "Print more information")
	flag.BoolVar(&icase, "icase", false, "Ignore case of filenames")

	flag.Parse()

	if len(*orig) == 0 {
		flag.PrintDefaults()
		log.Fatal("orig directory not declared")
	}

	if len(*new) == 0 {
		flag.PrintDefaults()
		log.Fatal("new directory not declared")
	}

	// Check for destination directory
	_, err := os.Stat(*store)
	if err != nil {
		if os.IsNotExist(err) {
			err := os.Mkdir(*store, 0755)
			if err != nil {
				log.Fatalf("Could not create output directory %v : %v", *store, err)
			}
		} else {
			log.Fatalf("Error when checking for destination folder. %v", err)
		}
	}

	log.Printf("Starting to compare dir %s with %s and writing files to %s", *orig, *new, *store)

	// Visit orig directory
	origMap := make(map[string]FileInfoPath, 100)
	origPath, err := filepath.Abs(*orig)
	if err != nil {
		log.Fatal(err)
	}
	origChan := make(chan int)
	go func() {
		visit(origPath, origMap)
		origChan <- 1
	}()

	// Visit new directory
	newMap := make(map[string]FileInfoPath, 100)
	newPath, err := filepath.Abs(*new)
	if err != nil {
		log.Fatal(err)
	}
	newChan := make(chan int)
	go func() {
		visit(newPath, newMap)
		newChan <- 1
	}()

	// Wait for visits to complete
	<-origChan
	<-newChan

	// Check new files and copy them
	for k, v := range newMap {
		if _, ok := origMap[k]; ok == false {
			if verbose {
				log.Printf("Found diff file %v", v.fi.Name())
			}
			copyFile(*store, v.path, v.fi.Name())
		}
	}
}

func visit(path string, data map[string]FileInfoPath) {
	if verbose {
		log.Printf("Checking directory %v", path)
	}

	d, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	fr, err := d.Readdir(0)
	d.Close()
	if err != nil {
		log.Fatal(err)
	}

	for _, fi := range fr {
		if fi.IsDir() {
			visit(filepath.Join(path, fi.Name()), data)
		} else {
			// Add file entry
			var name string
			if icase {
				name = strings.ToUpper(fi.Name())
			} else {
				name = fi.Name()
			}
			data[name] = FileInfoPath{fi, path}
		}
	}
}

func copyFile(destPath string, srcPath string, srcName string) {
	// Source file
	srcFile, err := os.Open(filepath.Join(srcPath, srcName))
	if err != nil {
		log.Printf("Cannot open %v", filepath.Join(srcPath, srcName))
		return
	}
	defer srcFile.Close()

	// Destination file
	dstFile, err := os.Create(filepath.Join(destPath, srcName))
	if err != nil {
		log.Printf("Could not open %v", filepath.Join(destPath, srcName))
		return
	}
	defer dstFile.Close()

	// Copy the file
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		log.Printf("Could not copy %v", srcFile)
	}
}

