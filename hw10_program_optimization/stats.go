package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson" //nolint:all
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewReader(r)

	i := 0
	var line []byte
	var user User
	for {
		line, _, err = reader.ReadLine()
		if errors.Is(err, io.EOF) {
			err = nil
			break
		}
		if err != nil {
			return
		}

		if err = easyjson.Unmarshal(line, &user); err != nil {
			return
		}
		result[i] = user
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	var emailDomain string
	var matched bool
	dottedDomain := "." + domain
	for _, user := range u {
		matched = strings.Contains(user.Email, dottedDomain)

		if matched {
			emailDomain = strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[emailDomain]
			num++
			result[emailDomain] = num
		}
	}
	return result, nil
}
