package lib

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/lonng/nanoserver/pkg/errutil"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

/*os  start*/

//create a dir
func Create_dir(path string) error {
	err := os.MkdirAll(path, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

//remove multi file
func Removetempimages(filenames []string) {
	for _, name := range filenames {
		os.Remove(name)
	}
}

//check file exits
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}


// windows kill process
func WinKill(process string) {
	kill := exec.Command("taskkill", "/f", "/im", process)
	//kill.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	kill.Run()
}

//windows list process
func WinListProcess() {

}

func Utf8ToGBK(utf8str string) string {
	result, _, _ := transform.String(simplifiedchinese.GBK.NewEncoder(), utf8str)
	return result
}

func CopyFile(dst, src string) error {
	if dst == "" || src == "" {
		return errutil.ErrIllegalParameter
	}

	srcDir, _ := filepath.Split(src)
	// get properties of source dir
	srcDirInfo, err := os.Stat(srcDir)
	if err != nil {
		return err
	}

	dstDir, _ := filepath.Split(dst)
	if err != nil {
		return err
	}

	MakeDirIfNeed(dstDir, srcDirInfo.Mode())

	sf, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sf.Close()

	df, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer df.Close()

	_, err = io.Copy(df, sf)
	return err
}

func CopyDir(dst string, src string) error {
	// get properties of source dir
	srcDirInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	// create dest dir
	err = MakeDirIfNeed(dst, srcDirInfo.Mode())
	if err != nil {
		return err
	}

	srcDir, _ := os.Open(src)
	objs, err := srcDir.Readdir(-1)
	if err != nil {
		return err
	}

	const sep = string(filepath.Separator)
	for _, obj := range objs {
		srcFile := src + sep + obj.Name()
		dstFile := dst + sep + obj.Name()

		if obj.IsDir() {
			// create sub-directories - recursively
			if err = CopyDir(dstFile, srcFile); err != nil {
				return err
			}
			continue
		}

		err = CopyFile(dstFile, srcFile)
		if err != nil {
			return err
		}

	}
	return err
}

func MakeDirIfNeed(dir string, mode os.FileMode) error {
	dir = strings.TrimRight(dir, "/")

	if FileExists(dir) {
		return nil
	}

	err := os.MkdirAll(dir, mode)
	return err
}


func RunCmd(cmdName string, workingDir string, args ...string) (string, error) {
	const duration = time.Second * 7200

	cmd := exec.Command(cmdName, args...)

	if workingDir != "" {
		cmd.Dir = workingDir
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	buf := &bytes.Buffer{}

	chanErr := make(chan error)
	go func() {
		multiReader := io.MultiReader(stdout, stderr)
		in := bufio.NewScanner(multiReader)
		for in.Scan() {
			buf.Write(in.Bytes())
			buf.WriteString("\n")
		}

		if err := in.Err(); err != nil {
			chanErr <- err
			return
		}

		close(chanErr)

	}()

	// wait or timeout
	chanDone := make(chan error)

	go func() {
		chanDone <- cmd.Wait()
	}()
	select {
	case <-time.After(duration):
		cmd.Process.Kill()
		return "", fmt.Errorf("run command: %s failed with timeout", cmdName)

	case err, ok := <-chanErr:
		if ok {
			return "", err
		}

	case e := <-chanDone:
		fmt.Printf("error %+v\n", e)
	}

	return buf.String(), nil
}
/* os end */

//random number for time.Duration
func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}

//TimeRange adjust the time range.
func TimeRange(start, end int64) (int64, int64) {
	if start < 0 {
		start = 0
	}
	if end < 0 || end > time.Now().Unix() {
		end = time.Now().Unix()
	}

	if start > end {
		start, end = end, start
	}

	return start, end
}

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}



/* web */
func AccessControl(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")

		h.ServeHTTP(w, r)
	})
}

//HTTPGet http's get method
func HTTPGet(url string) (string, error) {
	var body []byte
	rspn, err := http.Get(url)

	if err != nil {
		return "", err
	}
	defer rspn.Body.Close()
	body, err = ioutil.ReadAll(rspn.Body)

	return string(body), err
}
/* web */
