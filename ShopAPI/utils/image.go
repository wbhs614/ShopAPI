package utils

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Sizer interface {
	Size() int64
}

const (
	LOCAL_FILE_DIR    = "static/uploads/file"
	MIN_FILE_SIZE     = 1       // bytes
	MAX_FILE_SIZE     = 5000000 // bytes
	IMAGE_TYPES       = "(jpeg|gif|p?jpeg|(x-)?png)"
	ACCEPT_FILE_TYPES = IMAGE_TYPES
	EXPIRATION_TIME   = 300 // seconds
	THUMBNAIL_PARAM   = "=s80"
)

var (
	imageTypes      = regexp.MustCompile(IMAGE_TYPES)
	acceptFileTypes = regexp.MustCompile(ACCEPT_FILE_TYPES)
)

type FileInfo struct {
	Url          string `json:"url,omitempty"`
	ThumbnailUrl string `json:"thumbnailUrl,omitempty"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	Size         int64  `json:"size"`
	Error        string `json:"error,omitempty"`
	DeleteUrl    string `json:"deleteUrl,omitempty"`
	DeleteType   string `json:"deleteType,omitempty"`
}

func ValidateType(fi *FileInfo) bool {
	if acceptFileTypes.MatchString(fi.Type) {
		return true
	}
	fi.Error = "Filetype not allowed"
	return false
}

func ValidateSize(fi *FileInfo) bool {
	if fi.Size < MIN_FILE_SIZE {
		fi.Error = "File is to small"
	} else if fi.Size > MAX_FILE_SIZE {
		fi.Error = "File"
	} else {
		return true
	}
	return false
}

func Escape(s string) string {
	return strings.Replace(url.QueryEscape(s), "+", "%20", -1)
}

func GetFormValue(p *multipart.Part) string {
	var b bytes.Buffer
	io.CopyN(&b, p, int64(1<<20))
	return b.String()
}

func Substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func GetParentDirectory(dirctory string) string {
	return Substr(dirctory, 0, strings.LastIndex(dirctory, "/"))
}

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}
