package email

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
)

type Keywords struct {
	Expense []string `json:"expense"`
	Junk    []string `json:"junk"`
}

var (
	keywords     Keywords
	keywordsErr  error
	keywordsOnce sync.Once
)

func getKeywords() (Keywords, error) {
	keywordsOnce.Do(func() {
		path := os.Getenv("KEYWORDS_PATH")
		if path == "" {
			path = "keywords.json"
		}
		
		b, err := os.ReadFile(path)
		if err != nil {
			keywordsErr = fmt.Errorf("unable to read %s: %w", path, err)
			return
		}

		if err := json.Unmarshal(b, &keywords); err != nil {
			keywordsErr = fmt.Errorf("unable to parse %s: %w", path, err)
			return
		}

		for i, k := range keywords.Expense {
			keywords.Expense[i] = strings.ToLower(k)
		}

		for i, k := range keywords.Junk {
			keywords.Junk[i] = strings.ToLower(k)
		}
	})
	return keywords, keywordsErr
}

func IsExpenseReceipt(subject, from string) (bool, error) {
	kw, err := getKeywords()
	if err != nil {
		return false, err
	}

	subject = strings.ToLower(subject)
	from = strings.ToLower(from)

	for _, junk := range kw.Junk {
		if strings.Contains(subject, junk) || strings.Contains(from, junk) {
			return false, nil
		}
	}

	for _, keyword := range kw.Expense {
		if strings.Contains(subject, keyword) {
			return true, nil
		}
	}

	return false, nil
}
