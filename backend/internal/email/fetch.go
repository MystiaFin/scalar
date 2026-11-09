package email

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
)

type RawEmail struct {
	MessageID string
	Hash      string 
	Subject   string
	From      string
	Date      time.Time
	Body      string  
}

func hashMessageID(messageID string) string {
	h := sha256.Sum256([]byte(messageID))
	return hex.EncodeToString(h[:])
}

var htmlTagRe = regexp.MustCompile(`<[^>]+>`)

func stripHTML(raw string) string {
	text := htmlTagRe.ReplaceAllString(raw, " ")
	text = strings.ReplaceAll(text, "&nbsp;", " ")
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&#39;", "'")
	text = strings.ReplaceAll(text, "&quot;", `"`)
	return strings.Join(strings.Fields(text), " ")
}

func fetchCleanText(msg *imap.Message, section *imap.BodySectionName) (string, error) {
	if msg == nil {
		return "", fmt.Errorf("server didn't return message")
	}

	r := msg.GetBody(section)
	if r == nil {
		return "", fmt.Errorf("server didn't return message body")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		return "", err
	}

	var htmlFallback string

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Println("MIME parse error:", err)
			break
		}

		if h, ok := p.Header.(*mail.InlineHeader); ok {
			contentType, _, _ := h.ContentType()
			switch strings.ToLower(contentType) {
			case "text/plain":
				b, _ := io.ReadAll(p.Body)
				return string(b), nil
			case "text/html":
				b, _ := io.ReadAll(p.Body)
				htmlFallback = string(b)
			}
		}
	}

	if htmlFallback != "" {
		return stripHTML(htmlFallback), nil
	}

	return "", fmt.Errorf("no text body found")
}

func FetchReceipts(days int, alreadyProcessed func(hash string) bool) ([]RawEmail, error) {
	emailAddress := os.Getenv("SCALAR_EMAIL")
	appPassword := os.Getenv("SCALAR_APP_PASSWORD")
	if emailAddress == "" || appPassword == "" {
		return nil, fmt.Errorf("SCALAR_EMAIL atau SCALAR_APP_PASSWORD belum di-set")
	}

	log.Println("connecting to imap.gmail.com:993...")
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		return nil, fmt.Errorf("imap dial: %w", err)
	}

	defer c.Logout()

	if err := c.Login(emailAddress, appPassword); err != nil {
		return nil, fmt.Errorf("imap login: %w", err)
	}

	mbox, err := c.Select("INBOX", true) // read-only
	if err != nil {
		return nil, fmt.Errorf("imap select: %w", err)
	}

	if mbox.Messages == 0 {
		return nil, nil
	}

	criteria := imap.NewSearchCriteria()
	criteria.Since = time.Now().AddDate(0, 0, -days)

	ids, err := c.Search(criteria)
	if err != nil {
		return nil, fmt.Errorf("imap search: %w", err)
	}

	if len(ids) == 0 {
		log.Printf("tidak ada email dalam %d hari terakhir", days)
		return nil, nil
	}

	log.Printf("found %d messages", len(ids))

	seqset := new(imap.SeqSet)
	seqset.AddNum(ids...)

	envCh := make(chan *imap.Message, len(ids))
	envDone := make(chan error, 1)
	go func() {
		envDone <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, envCh)
	}()

	var msgs []*imap.Message
	for m := range envCh {
		if m != nil && m.Envelope != nil {
			msgs = append(msgs, m)
		}
	}
	if err := <-envDone; err != nil {
		return nil, fmt.Errorf("imap envelope fetch: %w", err)
	}

	sort.Slice(msgs, func(i, j int) bool {
		return msgs[i].Envelope.Date.After(msgs[j].Envelope.Date)
	})

	var results []RawEmail
	for _, msg := range msgs {
		subject := msg.Envelope.Subject
		from := ""
		if len(msg.Envelope.From) > 0 {
			from = msg.Envelope.From[0].Address()
		}

		ok, err := IsExpenseReceipt(subject, from)
		if err != nil {
			return nil, err
		}

		if !ok {
			continue
		}

		hash := hashMessageID(msg.Envelope.MessageId)
		if alreadyProcessed != nil && alreadyProcessed(hash) {
			log.Printf("sudah diproses, skip: %s", subject)
			continue
		}

		bodySet := new(imap.SeqSet)
		bodySet.AddNum(msg.SeqNum)
		section := &imap.BodySectionName{Peek: true} 
		bodyCh := make(chan *imap.Message, 1)
		bodyDone := make(chan error, 1)
		go func() {
			bodyDone <- c.Fetch(bodySet, []imap.FetchItem{section.FetchItem()}, bodyCh)
		}()

		var bodyMsg *imap.Message
		for m := range bodyCh {
			bodyMsg = m
		}
		
		if err := <-bodyDone; err != nil {
			log.Printf("body fetch error (%s): %v", subject, err)
			continue
		}

		body, err := fetchCleanText(bodyMsg, section)
		if err != nil {
			log.Printf("gagal baca body (%s): %v", subject, err)
			continue
		}

		results = append(results, RawEmail{
			MessageID: msg.Envelope.MessageId,
			Hash:      hash,
			Subject:   subject,
			From:      from,
			Date:      msg.Envelope.Date,
			Body:      body,
		})
	}

	log.Printf("%d email lolos filter", len(results))
	return results, nil
}
