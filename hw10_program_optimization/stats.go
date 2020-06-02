package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type User struct {
	ID       int    `json:"-"`
	Name     string `json:"-"`
	Username string `json:"-"`
	Email    string
	Phone    string `json:"-"`
	Password string `json:"-"`
	Address  string `json:"-"`
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return getUsers(r, domain)
}

func getUsers(r io.Reader, domain string) (DomainStat, error) {

	var result = make(DomainStat, 1000)
	var user User
	reader := bufio.NewReader(r)
	for i := 0; ; i++ {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			return result, err
		}
		if err := json.Unmarshal(line, &user); err != nil {
			return result, err
		}
		if strings.HasSuffix(user.Email, "."+domain) {
			result[strings.ToLower(strings.Split(user.Email, "@")[1])]++
		}
	}
	return result, nil
}
