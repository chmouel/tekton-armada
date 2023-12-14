package armada

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

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
	Tags    []string          `json:"tags"`
	Volumes map[string]string `json:"volumes"`
	Yamls   []string          `json:"yamls"`
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
		_, _ = w.Write([]byte("<img src=\"https://cdn.futura-sciences.com/sources/images/glossaire/and.jpg\">"))
		w.WriteHeader(http.StatusFound)
	})

	router.Get("/tag/{tag}", func(w http.ResponseWriter, r *http.Request) {
		tag := chi.URLParam(r, "tag")
		newURL, err := r.URL.Parse(fmt.Sprintf("%s?stream=%s", r.URL.Path, "/tag/"+tag))
		if err != nil {
			errorIt(w, r, http.StatusInternalServerError, err)
			return
		}
		r.URL = newURL
		events.ServeHTTP(w, r)
	})
	router.Post("/job", func(w http.ResponseWriter, r *http.Request) {
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
		for _, tag := range job.Tags {
			if tag == "" {
				errorIt(w, r, http.StatusBadRequest, fmt.Errorf("tag cannot be empty"))
				return
			}
			events.CreateStream(tag)
			events.Publish("/tag/"+tag, &sse.Event{
				Data: body,
			})
			fmt.Printf("created a job for tag %s, %s\n", tag, body)
		}
		w.WriteHeader(http.StatusAccepted)
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
