package token

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type Token struct {
	client HTTPClient
	token  string
	valid  bool
}

type tokenResponse struct {
	Valid bool `json:"valid"`
}

func NewToken(c HTTPClient) *Token {
	return &Token{c, "", false}
}

func (t *Token) Validate(token string) (bool, error) {
	// client := resty.New()

	// client.SetHeaders(map[string]string{
	// 	"Content-Type": "application/json",
	// 	"User-Agent":   "Azion Orchestrator",
	// })

	req, err := http.NewRequest("GET", "api.azion.net", nil)
	if err != nil {
		return false, err
	}
	q := req.URL.Query()
	q.Add("token", token)
	req.URL.RawQuery = q.Encode()
	req.Header.Add("Accept", "application/json")

	resp, err := t.client.Do(req)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	res := &tokenResponse{}
	json.Unmarshal(body, res)

	if !res.Valid {
		return false, nil
	}

	t.token = token
	t.valid = true

	return true, nil
	// resp, err := client.R().
	// 	SetQueryString("token="+token).
	// 	SetHeader("Accept", "application/json").
	// 	Get("http://localhost/apitoken")

	// if err != nil {
	// 	return false
	// }
	// fmt.Println("  Error      :", err)
	// fmt.Println("  Status Code:", resp.StatusCode())
	// fmt.Println("  Body       :\n", resp)
	// fmt.Println()

	// if resp.StatusCode() == 200 {
	// 	fmt.Println("Token valid")
	// 	return true
	// }

	// fmt.Println("Token invalid")
	// return false
}

func (t *Token) Save() {
	fbyte := []byte(t.token + "\n")
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	dirname = dirname + "/.azion/"
	err = os.MkdirAll(dirname, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	dirname = dirname + "credentials"
	err = os.WriteFile(dirname, fbyte, 0600)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Token saved in " + dirname)
}

func ReadDisk() string {
	dirname, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return ""
	}

	dirname = dirname + "/.azion/credentials"
	file, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return ""
	}

	return scanner.Text()
}
