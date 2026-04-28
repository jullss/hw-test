package hw10programoptimization

import (
	"bufio"
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
	return countDomains(&u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	br := bufio.NewReaderSize(r, 64*1024)
	decoder := json.NewDecoder(br)

	for i := 0; i < len(result); i++ {
		if !decoder.More() {
			break
		}

		if err := decoder.Decode(&result[i]); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}

			return result, err
		}
	}
	return result, nil
}

func countDomains(u *users, domain string) (DomainStat, error) {
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
