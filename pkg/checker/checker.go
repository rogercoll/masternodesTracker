package checker


import (
	"regexp"
	"net/http"
	"io/ioutil"
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