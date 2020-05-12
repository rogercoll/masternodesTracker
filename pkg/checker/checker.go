package checker


import (
	"io"
	"log"
	"time"
	"regexp"
	"net/http"
	"encoding/hex"
	"io/ioutil"
	"crypto/sha256"
	"github.com/rogercoll/masternodesTracker/pkg/db"

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
		balance,actualHash, err := getActual(master.ApiEndpoint, master.RegexBalance)
		if err != nil { log.Fatal(err)}
		if actualHash !=  master.LastHash {
			log.Printf("Masternode with public key %v has new transactions", master.PublicKey)
			log.Printf("Masternode balance: %s", balance)
			newInfo := db.Masternode{
				Coin: master.Coin,
				PublicKey: master.PublicKey,
				ApiEndpoint: master.ApiEndpoint,
				RegexBalance: master.RegexBalance,
				LastHash: actualHash,
				LastCheck: uint64(time.Now().Unix()),
			}
			err = db.UpdateCoinInfo(c, master.PublicKey, &newInfo)
			if err != nil { log.Fatal(err)}
		}
	}
}