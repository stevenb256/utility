package utl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	l "github.com/stevenb256/log"
)

// reusable client
var httpClient = &http.Client{}

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

// URL - makes a url with arguments and escaping
func URL(scheme, host, path string, args ...string) string {
	u := &url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path,
	}
	q := u.Query()
	for i := 0; i < len(args); i += 2 {
		q.Add(args[i], args[i+1])
	}
	u.RawQuery = q.Encode()
	return u.String()
}

// Get - simple http get
func Get(cookie *http.Cookie, uri string) error {

	// Now that you have a form, you can submit it to your handler.
	request, err := http.NewRequest("GET", uri, nil)
	if l.Check(err) {
		return err
	}

	// if we have a cookie
	if nil != cookie {
		request.AddCookie(cookie)
	}

	// do the request
	r, err := httpClient.Do(request)
	if l.Check(err) {
		return err
	}
	defer r.Body.Close()

	// if not status ok
	if http.StatusOK != r.StatusCode {
		return fmt.Errorf("failed(status:%d) - %s", r.StatusCode, uri)
	}

	// done
	return nil
}

// PostForm a file to a http endpoint
func PostForm(cookie *http.Cookie, uri string, formValues ...string) error {

	// set form values
	form := url.Values{}
	for i := 0; i < len(formValues); i += 2 {
		form.Add(formValues[i], formValues[i+1])
	}

	// Now that you have a form, you can submit it to your handler.
	request, err := http.NewRequest("POST", uri, strings.NewReader(form.Encode()))
	if l.Check(err) {
		return err
	}

	// if we have a cookie
	if nil != cookie {
		request.AddCookie(cookie)
	}

	// set header
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.PostForm = form

	// do the request
	r, err := httpClient.Do(request)
	if l.Check(err) {
		return err
	}
	defer r.Body.Close()

	// if not status ok
	if http.StatusOK != r.StatusCode {
		return fmt.Errorf("failed(status:%d) - %s", r.StatusCode, uri)
	}

	// done
	return nil
}

// PostFiles a file to a http endpoint
func PostFiles(cookie *http.Cookie, uri string, files ...string) error {

	// make buffer and writer
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// loop through files
	for i := 0; i < len(files); i += 2 {

		// set fileName
		fileName := files[i]
		path := files[i+1]

		// read the file
		fileContents, err := ioutil.ReadFile(path)
		if l.Check(err) {
			return err
		}

		// add part from file
		part, err := writer.CreateFormFile(fileName, filepath.Base(path))
		if l.Check(err) {
			return err
		}

		// write the file
		_, err = part.Write(fileContents)
		if l.Check(err) {
			return err
		}
	}

	// close the writer
	err := writer.Close()
	if l.Check(err) {
		return err
	}

	// Now that you have a form, you can submit it to your handler.
	request, err := http.NewRequest("POST", uri, body)
	if l.Check(err) {
		return err
	}

	// if we have a cookie
	if nil != cookie {
		request.AddCookie(cookie)
	}

	// Don't forget to set the content type, this will contain the boundary.
	request.Header.Set("Content-Type", writer.FormDataContentType())

	// do the request
	r, err := httpClient.Do(request)
	if l.Check(err) {
		return err
	}
	defer r.Body.Close()

	// if not status ok
	if http.StatusOK != r.StatusCode {
		return fmt.Errorf("failed(status:%d) - %s", r.StatusCode, uri)
	}

	// done
	return nil
}
