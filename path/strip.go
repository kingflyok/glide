package path

import (
	"os"
	"path/filepath"

	"github.com/Masterminds/glide/godep/strip"
	"github.com/Masterminds/glide/msg"
)

// StripVcs removes VCS metadata (.git, .hg, .bzr, .svn) from the vendor/
// directory.
func StripVcs() error {
	searchPath := VendorDir + string(os.PathSeparator)
	if _, err := os.Stat(searchPath); err != nil {
		if os.IsNotExist(err) {
			msg.Debug("Vendor directory does not exist.")
		}

		return err
	}
	return filepath.Walk(searchPath, stripHandler)
}

func stripHandler(path string, info os.FileInfo, err error) error {

	name := info.Name()
	if name == ".git" || name == ".bzr" || name == ".svn" || name == ".hg" {
		if _, err := os.Stat(path); err == nil {
			if info.IsDir() {
				msg.Info("Removing: %s", path)
				return os.RemoveAll(path)
			}

			msg.Debug("%s is not a directory. Skipping removal", path)
			return nil
		}
	}
	return nil
}

// StripVendor removes nested vendor and Godeps/_workspace/ directories.
func StripVendor() error {
	searchPath := VendorDir + string(os.PathSeparator)
	if _, err := os.Stat(searchPath); err != nil {
		if os.IsNotExist(err) {
			msg.Debug("Vendor directory does not exist.")
		}

		return err
	}

	err := filepath.Walk(searchPath, func(path string, info os.FileInfo, err error) error {
		// Skip the base vendor directory
		if path == searchPath || path == VendorDir {
			return nil
		}

		name := info.Name()
		if name == "vendor" {
			if _, err := os.Stat(path); err == nil {
				if info.IsDir() {
					msg.Info("Removing: %s", path)
					return os.RemoveAll(path)
				}

				msg.Debug("%s is not a directory. Skipping removal", path)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return strip.GodepWorkspace(searchPath)
}
