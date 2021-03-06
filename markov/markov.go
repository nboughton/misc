/*
Package markov is a minor adaptation of https://golang.org/doc/codewalk/markov.go
with two additional functions for loading data from a file or url.
*/
package markov

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	supportedFile = regexp.MustCompile(`.*\.txt$`)
)

// Prefix is a Markov chain prefix of one or more words.
type Prefix []string

// String returns the Prefix as a string (for use as a map key).
func (p Prefix) String() string {
	return strings.Join(p, " ")
}

// Shift removes the first word from the Prefix and appends the given word.
func (p Prefix) Shift(word string) {
	copy(p, p[1:])
	p[len(p)-1] = word
}

// Chain contains a map ("chain") of prefixes to a list of suffixes.
// A prefix is a string of prefixLen words joined with spaces.
// A suffix is a single word. A prefix can have multiple suffixes.
type Chain struct {
	chain     map[string][]string
	prefixLen int
}

// NewChain returns a new Chain with prefixes of prefixLen words.
func NewChain(prefixLen int) *Chain {
	return &Chain{make(map[string][]string), prefixLen}
}

// Build reads text from the provided Reader and
// parses it into prefixes and suffixes that are stored in Chain.
func (c *Chain) Build(r io.Reader) {
	br := bufio.NewReader(r)
	p := make(Prefix, c.prefixLen)
	for {
		var s string
		if _, err := fmt.Fscan(br, &s); err != nil {
			break
		}
		key := p.String()
		c.chain[key] = append(c.chain[key], s)
		p.Shift(s)
	}
}

// Generate returns a string of at most n words generated from Chain.
func (c *Chain) Generate(n int) string {
	r, p := rand.New(rand.NewSource(time.Now().UnixNano())), make(Prefix, c.prefixLen)
	var words []string
	for i := 0; i < n; i++ {
		choices := c.chain[p.String()]
		if len(choices) == 0 {
			break
		}
		next := choices[r.Intn(len(choices))]
		words = append(words, next)
		p.Shift(next)
	}
	return strings.Join(words, " ")
}

// FromURL retrieves url (if it is supported) and returns a markov chain
// based on its contents
func FromURL(url string, textLength int) (string, error) {
	if !supportedFile.MatchString(url) {
		return "", fmt.Errorf("Unsupported filetype: %v", url)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Could not retrieve file: %v", err)
	}
	defer resp.Body.Close()

	c := NewChain(2)   // Initialize a new Chain.
	c.Build(resp.Body) // Build chains from standard input.

	text := c.Generate(textLength)
	return TrimToSentence(text), nil
}

// FromFile generates a markov chain from the contents of a text file
func FromFile(file string, textLength int) (string, error) {
	if !supportedFile.MatchString(file) {
		return "", fmt.Errorf("Unsupported filetype: %v", file)
	}

	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	c := NewChain(2)
	c.Build(f)

	text := c.Generate(textLength)
	return TrimToSentence(text), nil
}

// TrimToSentence returns the text trimmed to the last full stop
func TrimToSentence(text string) string {
	for i := len(text) - 1; i > 0; i-- {
		if string(text[i]) == "." {
			return string(text[:i+1])
		}
	}

	return ""
}
