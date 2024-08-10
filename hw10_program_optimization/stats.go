package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/mailru/easyjson"
)

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

// Номер строки из файла -> пользователь.
type users [100_000]User

func getUsers(r io.Reader) (users, error) {
	u := users{}
	sc := bufio.NewScanner(r)
	for i := 0; sc.Scan(); i++ {
		var user User
		if err := easyjson.Unmarshal(sc.Bytes(), &user); err != nil {
			return u, err
		}
		u[i] = user
	}
	return u, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	stats := make(DomainStat, len(u))
	for _, user := range u {
		if strings.Contains(user.Email, domain) {
			subDomain := strings.ToLower(user.Email[strings.Index(user.Email, "@")+1:])
			stats[subDomain]++
		}
	}
	return stats, nil
}
