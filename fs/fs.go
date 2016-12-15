package fs

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var (
	//ErrSrcNotExist src file doesnt exist
	ErrSrcNotExist = errors.New("Source file does not exist")

	//ErrSrcNotRegularFile src file is not a regular file
	ErrSrcNotRegularFile = errors.New("Source file is not a regular file")

	//ErrDstNotRegularFile dst file is not a regular file
	ErrDstNotRegularFile = errors.New("Destination file is not a regular file")
)

//Fs is a static class that provides Filesystem type functions
type Fs struct{}

//NewFs generates a Fs object
func NewFs() *Fs {
	myFs := &Fs{}
	return myFs
}

//DoesFileExist just like it sounds
func (fs *Fs) DoesFileExist(fullpath string) bool {
	log.Debugln("DoesFileExist ENTER")
	log.Debugln("fullpath:", fullpath)

	if _, err := os.Stat(fullpath); !os.IsNotExist(err) {
		log.Debugln(fullpath, "does exists")
		log.Debugln("DoesFileExist LEAVE")
		return true
	}

	log.Debugln(fullpath, "does not exists")
	log.Debugln("DoesFileExist LEAVE")
	return false
}

//GetFullExePath returns the fullpath of the executable including the executable
//name itself
func (fs *Fs) GetFullExePath() (string, error) {
	path, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorln("Readlink failed:", err)
		return "", nil
	}
	log.Debugln("EXE path:", path)
	return path, nil
}

//GetPathFromFullFilename returns the parent folder name
func (fs *Fs) GetPathFromFullFilename(path string) string {
	log.Debugln("GetPathFromFullFilename ENTER")
	log.Debugln("path:", path)
	last := strings.LastIndex(path, string(filepath.Separator))
	if last == -1 {
		log.Debugln("No slash. Return Path:", path)
		log.Debugln("GetPathFromFullFilename LEAVE")
		return path
	}
	tmp := path[0:last]
	log.Debugln("Final Path:", tmp)
	log.Debugln("GetPathFromFullFilename LEAVE")
	return tmp
}

//GetFullPath returns the fullpath of the executable without the executable name
func (fs *Fs) GetFullPath() (string, error) {
	path, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorln("Readlink failed:", err)
		return "", nil
	}
	log.Debugln("EXE path:", path)

	tmp := fs.GetPathFromFullFilename(path)
	return tmp, nil
}

//GetFullPathWithExec returns the fullpath of the executable
func (fs *Fs) GetFullPathWithExec() (string, error) {
	path, err := os.Readlink("/proc/self/exe")
	if err != nil {
		log.Errorln("Readlink failed:", err)
		return "", nil
	}
	log.Debugln("EXE path:", path)

	return path, nil
}

//GetFilenameFromURIOrFullPath retrieves the filename from an URI
func (fs *Fs) GetFilenameFromURIOrFullPath(path string) string {
	log.Debugln("GetFilenameFromURI ENTER")
	log.Debugln("path:", path)

	last := strings.LastIndex(path, string(filepath.Separator))
	if last == -1 {
		log.Debugln("No slash. Return Path:", path)
		log.Debugln("GetFilenameFromURI LEAVE")
		return path
	}
	pathTmp := path[last+1:]
	log.Debugln("Return Path:", pathTmp)
	log.Debugln("GetFilenameFromURI LEAVE")

	return pathTmp
}

//AppendSlash appends a slash to a path if one is needed
func (fs *Fs) AppendSlash(path string) string {
	log.Debugln("AppendSlash ENTER")
	log.Debugln("path:", path)
	if path[len(path)-1] != filepath.Separator {
		path += string(filepath.Separator)
	}
	log.Debugln("Return Path:", path)
	log.Debugln("GetFilenameFromURI LEAVE")
	return path
}

//CopyFile copies the contents of the src file to the dst file
func (fs *Fs) CopyFile(src string, dst string) error {
	return fs.CopyFileEx(src, dst, 0666)
}

//CopyFileEx copies the contents of the src file to the dst file
func (fs *Fs) CopyFileEx(src string, dst string, mode os.FileMode) error {
	log.Debugln("CopyFile ENTER")
	log.Debugln("SRC:", src)
	log.Debugln("DST:", dst)

	sfi, err := os.Stat(src)
	if err != nil {
		log.Debugln("Src Stat Failed:", err)
		log.Debugln("CopyFile LEAVE")
		return ErrSrcNotExist
	}
	if !sfi.Mode().IsRegular() {
		//cannot copy non-regular files (e.g., directories, symlinks, devices, etc.)
		log.Debugln("Src file is not regular")
		log.Debugln("CopyFile LEAVE")
		return ErrSrcNotRegularFile
	}
	dfi, err := os.Stat(dst)
	if err == nil {
		if !(dfi.Mode().IsRegular()) {
			log.Debugln("Dst file is not regular")
			log.Debugln("CopyFile LEAVE")
			return ErrDstNotRegularFile
		}
		if os.SameFile(sfi, dfi) {
			log.Debugln("Src and Dst files are the same")
			log.Debugln("CopyFile LEAVE")
			return nil
		}
	}

	//Copy the file
	in, err := os.OpenFile(src, os.O_RDONLY, 0666)
	if err != nil {
		log.Debugln("Failed to open SRC file:", err)
		log.Debugln("CopyFile LEAVE")
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_RDWR, mode)
	if err != nil {
		log.Debugln("Failed to open DST file:", err)
		log.Debugln("CopyFile LEAVE")
		return err
	}
	defer out.Close()
	if _, err = io.Copy(out, in); err != nil {
		log.Debugln("Failed to copy file:", err)
		log.Debugln("CopyFile LEAVE")
		return err
	}

	err = out.Sync()
	if err != nil {
		log.Debugln("Failed to flush file:", err)
		log.Debugln("CopyFile LEAVE")
		return err
	}

	log.Debugln("CopyFile succeeded")
	log.Debugln("CopyFile LEAVE")
	return nil
}
