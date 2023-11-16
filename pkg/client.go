package andromeda

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	_ "embed"

	"github.com/r3labs/sse/v2"
	"golang.org/x/exp/slog"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var pmEventRe = regexp.MustCompile(`(\w+|\d+|_|-|:)`)

const (
	defaultTimeout = 5
	Version        = "0.0.1"
)

const smeeChannel = "messages"

const tsFormat = "2006-01-02T15.04.01.000"

type andromeda struct {
	targetURL string
	logger    *slog.Logger
}

// title returns a copy of the string s with all Unicode letters that begin words
// mapped to their Unicode title case.
func title(source string) string {
	return cases.Title(language.Und, cases.NoLower).String(source)
}

func (a andromeda) clientSetup() error {
	version := strings.TrimSpace(string(Version))
	s := fmt.Sprintf("Starting andromeda version: %s", version)
	a.logger.Info(s)
	client := sse.NewClient(a.targetURL, sse.ClientMaxBufferSize(1<<20))
	client.Headers["User-Agent"] = fmt.Sprintf("andromeda/%s", version)
	// this is to get nginx to work
	client.Headers["X-Accel-Buffering"] = "no"
	channel := filepath.Base(a.targetURL)
	err := client.Subscribe(channel, func(msg *sse.Event) {
		now := time.Now().UTC()
		nowStr := now.Format(tsFormat)

		if string(msg.Event) == "ready" || string(msg.Data) == "ready" {
			s := fmt.Sprintf("%s Listening to %s", nowStr, a.targetURL)
			a.logger.Info(s)
			return
		}

		if string(msg.Event) == "ping" {
			return
		}

		if string(msg.Data) != "{}" {
			fmt.Printf("msg.Data: %v\n", msg.Data)
		}
	})
	return err
}
