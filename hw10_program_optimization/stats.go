package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (users, error) {
	result := make(users, 0, 1000)
	decoder := json.NewDecoder(r)

	for decoder.More() {
		var u User
		if err := decoder.Decode(&u); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return result, err
		}

		result = append(result, u)
	}

	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat, 100)
	suffix := "." + domain

	for i := 0; i < len(u); i++ {
		if strings.HasSuffix(u[i].Email, suffix) {
			at := strings.LastIndex(u[i].Email, "@")
			if at != -1 {
				d := strings.ToLower(u[i].Email[at+1:])
				result[d]++
			}
		}
	}
	return result, nil
}
