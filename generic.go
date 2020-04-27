package utl

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"

	l "github.com/stevenb256/log"
)

// md5Bytes returns md5 hash of bytes
func md5Bytes(buf []byte) string {
	h := md5.Sum(buf)
	return base64.StdEncoding.EncodeToString(h[:])
}

// MinUint16 min of two uint16
func MinUint16(x, y uint16) uint16 {
	if x < y {
		return x
	}
	return y
}

// MaxUint16 Max of two uint32
func MaxUint16(x, y uint16) uint16 {
	if x > y {
		return x
	}
	return y
}

// SendError - send error over channel if not nil
func SendError(chError chan error, err error) error {
	if nil != chError {
		chError <- err
	}
	return err
}

// Execute - runs a command
func Execute(wait bool, dir, app string, args ...string) error {

	// change directory
	err := os.Chdir(dir)
	if l.Check(err) {
		return err
	}

	// start command
	command := exec.Command(app, args...)

	// give current in/out
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	// wait
	if true == wait {
		err = command.Run()
		if l.Check(err) {
			return err
		}
	} else {
		err = command.Start()
		if l.Check(err) {
			return err
		}
	}

	// done
	return nil
}

// MoveFile - copies file and then deletes the source file
func MoveFile(srcPath, dstPath string) error {
	err := CopyFile(srcPath, dstPath)
	if l.Check(err) {
		return err
	}
	return os.Remove(srcPath)
}

// CopyFileWithJoin - same as copy file but joins src/dst paths with src/dst file names
func CopyFileWithJoin(srcPath, srcFile, dstPath, dstFile string) error {
	return CopyFile(Join(srcPath, srcFile), Join(dstPath, dstFile))
}

// CopyFile copys source file path to dst file path
func CopyFile(srcPath, dstPath string) error {

	// make path
	err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if l.Check(err) {
		return err
	}

	// open source
	srcFile, err := os.Open(srcPath)
	if l.Check(err) {
		return err
	}
	defer srcFile.Close()

	// kill destination
	os.Remove(dstPath)

	// get info of source
	info, err := srcFile.Stat()
	if l.Check(err) {
		return err
	}

	// write to dest
	dstFile, err := os.OpenFile(dstPath, os.O_CREATE|os.O_RDWR, info.Mode())
	if l.Check(err) {
		return err
	}
	defer dstFile.Close()

	// make sure mode really got set
	if runtime.GOOS != "windows" {
		err = dstFile.Chmod(info.Mode())
		if l.Check(err) {
			return err
		}
	}

	// copy it
	_, err = io.Copy(dstFile, srcFile)
	if l.Check(err) {
		os.Remove(dstPath)
		return err
	}

	// done
	return nil
}

// DoesFileExist - checks to see if file exists
func DoesFileExist(filePath string) bool {
	_, err := os.Stat(filePath)
	return true != os.IsNotExist(err)
}

// IsDirectory - checks to see if path is a director or a file
func IsDirectory(path string) bool {
	stat, err := os.Stat(path)
	return nil == err && true == stat.IsDir()
}

// Join - takes a set of strings and joins them into a path
func Join(a ...string) string {
	return filepath.Join(a...)
}

// Clean - calls go filepath clean method
func Clean(path string) string {
	return filepath.Clean(path)
}

// Clone - copyies exported fields from one struct and makes new
func Clone(inter interface{}) interface{} {
	nInter := reflect.New(reflect.TypeOf(inter).Elem())
	val := reflect.ValueOf(inter).Elem()
	nVal := nInter.Elem()
	for i := 0; i < val.NumField(); i++ {
		nvField := nVal.Field(i)
		if true == nvField.CanSet() {
			nvField.Set(val.Field(i))
		}
	}
	return nInter.Interface()
}

// WriteFile writes buffer into path
func WriteFile(path string, buffer []byte) error {

	// make sure directory exists
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if l.Check(err) {
		return err
	}

	// create/lock the local file
	file, err := os.Create(path)
	if l.Check(err) {
		return err
	}
	defer file.Close()

	// write the contents
	_, err = file.Write(buffer)
	if l.Check(err) {
		return err
	}

	// truncate the file
	err = file.Truncate(int64(len(buffer)))
	if l.Check(err) {
		return err
	}

	// set the size
	return nil
}

// LaunchURL -- launches the browser to a url
func LaunchURL(home, url string) error {

	// path
	path := Join(home, "open.url")

	// write to a file that we can launch
	err := WriteFile(path, []byte(fmt.Sprintf("[InternetShortcut]\nURL=%s", url)))
	if l.Check(err) {
		return err
	}

	// setup command and launch url to do auth flow
	var verb, file string
	if runtime.GOOS == "windows" {
		verb = "cmd"
		file = "/c start " + path
	} else if runtime.GOOS == "darwin" {
		verb = "open"
		file = path
	} else if runtime.GOOS == "linux" {
		verb = "xdg-open"
		file = path
	}

	// launch command
	command := exec.Command(verb, file) // TODO: on windows use start, xdg-open on linux
	err = command.Run()
	if l.Check(err) {
		return err
	}

	// done
	return nil
}

// MakeURL - makes a url with arguments and escaping
func MakeURL(base, page string, args ...string) string {
	s := new(strings.Builder)
	fmt.Fprintf(s, "%s/%s", base, page)
	sep := "?"
	for i := 0; i < len(args); i += 2 {
		fmt.Fprintf(s, "%s%s=%s", sep, args[i], url.PathEscape(args[i+1]))
		sep = "&"
	}
	return s.String()
}

// Get - simple http get
func Get(base, page string, args ...string) error {

	// make url
	s := MakeURL(base, page, args...)

	// try to ping and get a response
	r, err := http.Get(s)
	if l.Check(err) {
		return err
	}

	// if not a 200
	if 200 != r.StatusCode {
		//b, _ := ioutil.ReadAll(r.Body)
		//l.Debug(string(b))
		err = fmt.Errorf("failed(status:%d) - %s", r.StatusCode, s)
		l.Fail(err)
		return err
	}

	// done
	return nil
}
