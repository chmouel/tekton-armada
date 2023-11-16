package andromeda

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/r3labs/sse/v2"
	"github.com/urfave/cli/v2"
)

var (
	defaultServerPort    = 3344
	defaultServerAddress = "localhost"
)

type Job struct {
	Tags []string `json:"tags"`
	URL  string   `json:"url"`
}

func errorIt(w http.ResponseWriter, _ *http.Request, status int, err error) {
	w.WriteHeader(status)
	_, _ = w.Write([]byte(err.Error()))
}

func serve(c *cli.Context) error {
	publicURL := c.String("public-url")
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	events := sse.New()
	events.AutoReplay = false
	events.AutoStream = true
	events.OnSubscribe = (func(sid string, _ *sse.Subscriber) {
		events.Publish(sid, &sse.Event{
			Data: []byte("ready"),
		})
	})

	router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		// redirect to /new
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte("<img src=\"https://cdn.futura-sciences.com/sources/images/glossaire/and.jpg\">"))
		w.WriteHeader(http.StatusFound)
	})

	router.Post("/job", func(w http.ResponseWriter, r *http.Request) {
		// grab current time stamp before we take any further actions
		now := time.Now().UTC()
		// check if we have content-type json
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			errorIt(w, r, http.StatusBadRequest, fmt.Errorf("content-type must be application/json"))
			return
		}
		// try to json decode body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			errorIt(w, r, http.StatusInternalServerError, err)
			return
		}
		var job Job
		if err := json.Unmarshal(body, &job); err != nil {
			errorIt(w, r, http.StatusBadRequest, err)
			return
		}
		channel := strings.Join(job.Tags, "-")
		events.CreateStream(channel)
		events.Publish(channel, &sse.Event{
			Data: body,
		})
		w.WriteHeader(http.StatusAccepted)

		fmt.Fprintf(w, "{\"status\": %d, \"channel\": \"%s\", \"message\": \"ok\"}\n", http.StatusAccepted, channel)
		fmt.Fprintf(os.Stdout, "%s Published %s on channel %s\n", now.Format("2006-01-02T15.04.01.000"), middleware.GetReqID(r.Context()), channel)
	})

	portAddr := fmt.Sprintf("%s:%d", c.String("address"), c.Int("port"))
	if publicURL == "" {
		publicURL = "http://"
		publicURL = fmt.Sprintf("%s%s", publicURL, portAddr)
	}

	fmt.Fprintf(os.Stdout, "Serving for jobs on %s\n", publicURL)

	//nolint:gosec
	return http.ListenAndServe(portAddr, router)
}
