//go:build unix

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
)

// acquireLock attempts to acquire a single instance lock using Unix flock
func acquireLock() (*os.File, error) {
	// Get the lock file path in temp directory
	lockPath := filepath.Join(os.TempDir(), "browski.lock")

	// Try to open the lock file
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open lock file: %w", err)
	}

	// Try to acquire an exclusive lock
	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		lockFile.Close()
		// Another instance is running
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
		syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
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
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// On Unix, FindProcess always succeeds, so we need to send signal 0 to check
	err = process.Signal(syscall.Signal(0))
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
