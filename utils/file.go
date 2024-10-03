package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func CopyFileOrDir(src, dst string) error {
	src = ExpandPath(src)
	dst = ExpandPath(dst)

	sourceInfo, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("error accessing source: %v", err)
	}

	if sourceInfo.IsDir() {
		return CopyDir(src, dst)
	}
	return CopyFile(src, dst, sourceInfo.Mode())
}

func CopyFile(src, dst string, mode os.FileMode) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer sourceFile.Close()

	// Ensure the destination directory exists
	dstDir := filepath.Dir(dst)
	if err := os.MkdirAll(dstDir, os.ModePerm); err != nil {
		return fmt.Errorf("could not create destination directory: %v", err)
	}

	destinationFile, err := os.OpenFile(dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	return nil
}

func CopyDir(src, dst string) error {
	srcInfo, err := os.Lstat(src)
	if err != nil {
		return fmt.Errorf("error accessing source directory: %v", err)
	}

	// If the destination doesn't exist, create it with the same name as the source
	if _, err := os.Stat(dst); os.IsNotExist(err) {
		dst = filepath.Join(dst, filepath.Base(src))
	} else {
		// If the destination exists, append the source directory name
		dst = filepath.Join(dst, filepath.Base(src))
	}

	// Create the destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("could not create destination directory: %v", err)
	}

	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("could not read source directory: %v", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		entryInfo, err := os.Lstat(srcPath)
		if err != nil {
			return fmt.Errorf("error accessing entry %s: %v", srcPath, err)
		}

		switch {
		case entryInfo.Mode()&os.ModeSymlink != 0:
			// Handle symbolic link
			linkTarget, err := os.Readlink(srcPath)
			if err != nil {
				return fmt.Errorf("error reading symlink %s: %v", srcPath, err)
			}
			if err := os.Symlink(linkTarget, dstPath); err != nil {
				return fmt.Errorf("error creating symlink %s: %v", dstPath, err)
			}
		case entryInfo.IsDir():
			// Recursive call for subdirectories
			if err := CopyDir(srcPath, dstPath); err != nil {
				return err
			}
		default:
			// Copy files
			if err := CopyFile(srcPath, dstPath, entryInfo.Mode()); err != nil {
				return err
			}
		}
	}

	return nil
}

func MoveFileOrDir(src, dst string) error {
	src = ExpandPath(src)
	dst = ExpandPath(dst)

	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error accessing source: %v", err)
	}

	// Check if the destination exists
	dstInfo, err := os.Stat(dst)
	if err == nil {
		// Destination exists
		if dstInfo.IsDir() {
			// If destination is a directory, append the source filename
			dst = filepath.Join(dst, filepath.Base(src))
		} else if srcInfo.IsDir() {
			// If source is a directory but destination is a file, it's an error
			return fmt.Errorf("cannot overwrite non-directory %s with directory %s", dst, src)
		}
	} else if !os.IsNotExist(err) {
		// If there's an error other than "not exists", return it
		return fmt.Errorf("error accessing destination: %v", err)
	}

	// Attempt to rename (move) the file or directory
	err = os.Rename(src, dst)
	if err == nil {
		return nil
	}

	// If rename fails (e.g., cross-device link), fall back to copy and delete
	if srcInfo.IsDir() {
		err = CopyDir(src, dst)
	} else {
		err = CopyFile(src, dst, srcInfo.Mode())
	}

	if err != nil {
		return fmt.Errorf("error copying: %v", err)
	}

	// After successful copy, delete the source
	err = os.RemoveAll(src)
	if err != nil {
		return fmt.Errorf("error removing source after copy: %v", err)
	}

	return nil
}
