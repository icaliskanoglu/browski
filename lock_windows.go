//go:build windows

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// acquireLock attempts to acquire a single instance lock using Windows exclusive file access
func acquireLock() (*os.File, error) {
	// Get the lock file path in temp directory
	lockPath := filepath.Join(os.TempDir(), "browski.lock")

	// Try to open the lock file with exclusive access
	// On Windows, we use FILE_SHARE_NONE to prevent other processes from opening the file
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR|os.O_EXCL, 0666)
	if err != nil {
		// If file already exists and is locked, try opening without O_EXCL to detect if locked
		lockFile, err = os.OpenFile(lockPath, os.O_RDWR, 0666)
		if err != nil {
			// Another instance is running
			return nil, fmt.Errorf("another instance of Browski is already running")
		}
		lockFile.Close()
		return nil, fmt.Errorf("another instance of Browski is already running")
	}

	// Write our PID to the lock file
	lockFile.Truncate(0)
	lockFile.Seek(0, 0)
	fmt.Fprintf(lockFile, "%d", os.Getpid())
	lockFile.Sync()

	return lockFile, nil
}

// releaseLock releases the single instance lock
func releaseLock(lockFile *os.File) {
	if lockFile != nil {
		lockFile.Close()
		os.Remove(filepath.Join(os.TempDir(), "browski.lock"))
	}
}

// checkExistingInstance checks if another instance is running and tries to communicate with it
func checkExistingInstance(urlArg string) bool {
	lockPath := filepath.Join(os.TempDir(), "browski.lock")

	// Try to read the PID from the lock file
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return false
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false
	}

	// Check if the process is actually running
	// On Windows, FindProcess returns an error if the process doesn't exist
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// Try to send signal 0 to check if process is alive (Windows will return error if process is dead)
	err = process.Signal(os.Signal(nil))
	if err != nil {
		return false
	}

	// Another instance is running
	log.Printf("Another instance of Browski (PID %d) is already running", pid)

	// Try to send the URL to the existing instance
	if urlArg != "" {
		if sendURLToExistingInstance(urlArg) {
			return true
		}
		log.Printf("Failed to communicate with existing instance")
	}

	return true
}
