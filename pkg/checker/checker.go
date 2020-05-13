package checker


import (
	"io"
	"log"
	"fmt"
	"time"
	"regexp"
	"net/http"
	"encoding/hex"
	"io/ioutil"
	"crypto/sha256"
	"github.com/rogercoll/masternodesTracker/pkg/db"

)

const (
	InfoColor    = "\033[1;34m%s\033[0m"
	NoticeColor  = "\033[1;36m%s\033[0m"
	WarningColor = "\033[1;33m%s\033[0m"
	ErrorColor   = "\033[1;31m%s\033[0m"
	DebugColor   = "\033[0;36m%s\033[0m"
)

const (
    // See http://golang.org/pkg/time/#Parse
    timeFormat = "2006-01-02 15:04 MST"
)

func Check(url, regex string) (bool, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return false, err
	}
	resp, err := http.Get(url)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	match := re.FindAllString(string(body),-1)
	if match != nil {
		return true, nil
	}
	return false, nil
}

//We return the hash and the balance
func getActual(url, regex string) (string, string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "","", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "","", err
	}
	h := sha256.New()
	io.WriteString(h, string(body))
	s := hex.EncodeToString(h.Sum(nil))

	re, err := regexp.Compile(regex)
	if err != nil {
		return "","", err
	}

	//returns the match in the first position, in the second the captured substring
	match := re.FindStringSubmatch(string(body))
	if len(match) > 0 {
		return match[1], s, nil
	}

	return "", s, nil
}

func Check2() {
	c, err := db.NewMongoClient()
	if err != nil { log.Fatal(err)}
	//mcoins, err := db.GetCoinInfo(c, "eska")
	mcoins, err := db.GetCoinInfo(c, "iqcash")
	if err != nil { log.Fatal(err)}
	for _, master := range *mcoins {
		fmt.Printf("\033[1;33mCoin: %s     PublicKey: %s\033[0m\n", master.Coin, master.PublicKey)
		diff := time.Now().Sub(time.Unix(master.LastCheck,0))
		fmt.Printf("Time since last check: %2.f days %.f hours %.f minutes\n", diff.Hours()/24, diff.Hours(), diff.Minutes())
		balance,actualHash, err := getActual(master.ApiEndpoint, master.RegexBalance)
		if err != nil { log.Fatal(err)}
		fmt.Printf("Actual balance: %s\n", balance)
		if actualHash !=  master.LastHash {
			fmt.Printf("Masternode with public key %v has new transactions\n", master.PublicKey)
			master.LastHash = actualHash
			master.LastCheck = int64(time.Now().Unix())
			err = db.UpdateCoinInfo(c, master.PublicKey, &master)
			if err != nil { log.Fatal(err)}
			fmt.Printf(NoticeColor + "\n", "MongoDB document updated")
		} else {
			fmt.Printf(ErrorColor + "\n", "No changes since last check")
		}
		fmt.Println()
	}
}