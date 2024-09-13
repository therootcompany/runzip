package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/nwaples/rardecode"
)

func main() {
	nArgs := len(os.Args)
	if nArgs < 2 || nArgs > 3 {
		fmt.Printf("USAGE\n\tunrar <archive.rar> [./dst/]\n\n")
		fmt.Printf("EXAMPLES\n\tunrar ./archive.rar\n\tunrar ./archive.rar ./unpacked/\n\n")
		os.Exit(1)
		return
	}

	// note: we only need to know that Getwd works once
	//       (otherwise something is wrong in the universe)
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "sanity fail: %v\n", err)
		os.Exit(1)
		return
	}

	useInnerName := false
	rarFile := os.Args[1]
	finalDir := "."
	if nArgs > 2 {
		finalDir = os.Args[2]
	}
	finalDir, _ = filepath.Abs(finalDir)

	// encapsulate
	if _, err := os.Stat(finalDir); nil == err {
		useInnerName = true
		name := filepath.Base(rarFile)
		ext := filepath.Ext(name)
		name = strings.TrimSuffix(name, ext)
		finalDir = filepath.Join(finalDir, name)
	}

	tmpDir, _ := filepath.Abs(finalDir + ".tmp")
	tmpRel, _ := filepath.Rel(cwd, tmpDir)

	if err := os.MkdirAll(tmpDir, os.ModePerm); err != nil {
		fmt.Fprintf(os.Stderr, "error: could not create destination: %v\n", err)
		os.Exit(1)
		return
	}

	fmt.Fprintf(os.Stderr, "extracting to temporary path '%s/'...\n", tmpRel)

	topLevelFiles, err := unrar(rarFile, tmpDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not unarchive: %v\n", err)
		os.Exit(1)
		return
	}

	finalDir, err = finalizeDestination(tmpDir, finalDir, topLevelFiles, useInnerName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: could not rename directory: %v\n", err)
		os.Exit(1)
		return
	}

	finalRel, _ := filepath.Rel(cwd, finalDir)
	fmt.Fprintf(os.Stderr, "extracted to '%s/'\n", finalRel)
}

func unrar(rarFile, absRoot string) ([]string, error) {
	var err error
	topLevelFiles := []string{}

	f, err := os.Open(rarFile)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r, err := rardecode.NewReader(f, "")
	if err != nil {
		return nil, err
	}

	for {
		header, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		// TODO: might this store symlinks that we should delay?
		// header.Attributes

		cleanPath, err := sanitizePath(absRoot, header.Name)
		if err != nil {
			return nil, err
		}
		cleanRel, _ := filepath.Rel(absRoot, cleanPath)

		parts := strings.Split(cleanRel, string(filepath.Separator))
		if len(parts) == 1 {
			topLevelFiles = append(topLevelFiles, parts[0])
		}

		if header.IsDir {
			err = os.MkdirAll(cleanPath, os.ModePerm)
			if err != nil {
				return nil, err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(cleanPath), os.ModePerm); err != nil {
			return nil, err
		}

		outFile, err := os.Create(cleanPath)
		if err != nil {
			return nil, err
		}

		_, err = io.Copy(outFile, r)
		if err != nil {
			outFile.Close()
			return nil, err
		}
		outFile.Close()
	}

	return topLevelFiles, nil
}

// sanitizePath prevents overwriting ../../../etc/passwd and the like
func sanitizePath(absRoot, filePath string) (string, error) {
	dstPath := filepath.Join(absRoot, filePath)
	absPath, _ := filepath.Abs(dstPath)

	// guard against strings.HasPrefix(filepath.Abs("/etc/yo/../yolo"), "/etc/yo")
	checkPath := absPath
	if checkPath != "/" {
		checkPath += "/"
	}
	if absRoot != "/" {
		absRoot += "/"
	}
	if !strings.HasPrefix(checkPath, absRoot) {
		return "", fmt.Errorf("potential path traversal detected: %s", absPath)
	}

	return absPath, nil
}

// finalizeDestination moves and/or unnests the extracted director to its final destination
func finalizeDestination(tmpDir, finalDir string, topLevelFiles []string, useInnerName bool) (string, error) {
	if len(topLevelFiles) != 1 {
		return finalDir, os.Rename(tmpDir, finalDir)
	}

	topLevelName := topLevelFiles[0]
	topLevelPath := filepath.Join(tmpDir, topLevelName)
	topLevelEntryInfo, err := os.Stat(topLevelPath)
	if err != nil {
		return finalDir, err
	}

	preserveFileName := false
	if !topLevelEntryInfo.IsDir() {
		preserveFileName = true
	}

	if useInnerName {
		finalDir = filepath.Dir(finalDir)
		finalDir = filepath.Join(finalDir, topLevelName)
	} else if preserveFileName {
		finalDir = filepath.Join(finalDir, topLevelName)
	}

	if err := os.Rename(topLevelPath, finalDir); err != nil {
		return finalDir, err
	}

	// if the user has begun browsing this folder, it may .DS_Store,
	// desktop.ini, thumbs.db, etc - or even files they've added.
	_ = os.Remove(tmpDir)
	return finalDir, nil
}
