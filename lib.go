package lib

import (
	"crypto/md5"
	"encoding/hex"
	"log"
	"math/rand"
	"os"
	"os/exec"
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
func FileExits(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		return err
	}
	return nil
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

/* os end */

//random number for time.Duration
func RandInt64(min, max int64) int {
	rand.Seed(time.Now().UnixNano())
	return int(min + rand.Int63n(max-min+1))
}

//create md5 string
func Strtomd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

//password hash function
func Pwdhash(str string) string {
	return Strtomd5(str)
}
