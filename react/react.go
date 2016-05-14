// Package react returns a random reaction parsed from a json file, see example json file for structure
package react

import (
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/nboughton/config/parser"
	"github.com/nboughton/utils"
	"golang.org/x/exp/inotify"
)

// Item is the match text and possible responses for a reacion
type Item struct {
	Text string
	Resp []string
}

// React contains reactions from a reactions.json file
type React struct {
	Items []Item
}

var (
	// Reactions is the struct that the reactions.json file is parsed into
	Reactions  React
	reactRegex *regexp.Regexp
	reactFile  *string
)

func init() {
	reactFile = flag.String("react", "reactions.json", "Path to reactions file")
	flag.Parse()

	readFile()
	genRegex()

	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Panicln(err)
	}

	err = watcher.Watch(*reactFile)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				// if modified reload
				if ev.Mask == inotify.IN_CLOSE_WRITE {
					readFile()
					genRegex()
				}

			case err := <-watcher.Error:
				log.Println("error:", err)

			}
		}
	}()
}

// Match tests to see if a string matches a listed string
func Match(s string) bool {
	return reactRegex.MatchString(s)
}

// Respond returns a random response from the first reaction to match the generated regex
func Respond(s string) (string, error) {
	str := reactRegex.FindAllString(s, 1)[0]
	for _, r := range Reactions.Items {
		if r.Text == str {
			return utils.RandS(r.Resp), nil
		}
	}
	return "", fmt.Errorf("No response found")
}

func readFile() {
	err := parser.NewParser(*reactFile).Scan(&Reactions)
	if err != nil {
		log.Panicln(err)
	}
}

func genRegex() {
	s := "("
	for idx, i := range Reactions.Items {
		s += regexp.QuoteMeta(i.Text)
		if idx < len(Reactions.Items)-1 {
			s += "|"
		}
	}
	s += ")"

	reactRegex = regexp.MustCompile(s)
}
