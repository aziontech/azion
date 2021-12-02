package token

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
)

type Token struct {
	Client *resty.Client
}

func NewToken() Token {
	return Token{resty.New()}
}

func (t Token) Validation(token string) bool {
	client := resty.New()

	client.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		"User-Agent":   "Azion Orchestrator",
	})

	resp, err := client.R().
		SetQueryString("token="+token).
		SetHeader("Accept", "application/json").
		Get("http://localhost/apitoken")

	if err != nil {
		return false
	}
	// fmt.Println("  Error      :", err)
	// fmt.Println("  Status Code:", resp.StatusCode())
	// fmt.Println("  Body       :\n", resp)
	// fmt.Println()

	if resp.StatusCode() == 200 {
		fmt.Println("Token valid")
		return true
	}

	fmt.Println("Token invalid")
	return false
}

func Save(token string) {
	fbyte := []byte(token + "\n")
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
