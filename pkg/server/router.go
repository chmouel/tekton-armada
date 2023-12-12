package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/chmouel/armadas/pkg/app"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/r3labs/sse/v2"
	"go.uber.org/zap"
)

func (c *controller) errorIt(w http.ResponseWriter, _ *http.Request, status int, err error) {
	c.Logger.Error(err)
	w.WriteHeader(status)
	_, _ = w.Write([]byte(err.Error()))
}

func SugarLogger(logger *zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			scheme := "http"
			if r.TLS != nil {
				scheme = "https"
			}
			var msg string
			reqID := middleware.GetReqID(r.Context())
			if reqID != "" {
				msg += fmt.Sprintf("[%s] ", reqID)
			}
			msg += fmt.Sprintf("%s %s://%s%s", r.Method, scheme, r.Host, r.RequestURI)
			msg += fmt.Sprintf(" from %s", r.RemoteAddr)
			defer func() {
				ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
				msg += fmt.Sprintf(" - %d", ww.Status())
				logger.Info(msg)
			}()

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}

func (c *controller) getRouter() http.Handler {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	// router.Use(SugarLogger(c.Logger))
	router.Use(middleware.Recoverer)

	events := sse.New()
	events.AutoReplay = false
	events.AutoStream = true
	events.OnSubscribe = func(sid string, _ *sse.Subscriber) {
		events.Publish(sid, &sse.Event{
			Data: []byte("ready"),
		})
	}

	router.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		// redirect to /new
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusFound)
		_, _ = w.Write([]byte("<img src=\"https://assets.editorial.aetnd.com/uploads/2018/05/spanish-armada-gettyimages-600007851.jpg\">"))
	})

	router.Get("/tag/{tag}", func(w http.ResponseWriter, r *http.Request) {
		tag := chi.URLParam(r, "tag")
		newURL, err := r.URL.Parse(fmt.Sprintf("%s?stream=%s", r.URL.Path, "/tag/"+tag))
		if err != nil {
			c.errorIt(w, r, http.StatusInternalServerError, err)
			return
		}
		r.URL = newURL
		events.ServeHTTP(w, r)
	})
	router.Post("/job", func(w http.ResponseWriter, r *http.Request) {
		// check if we have content-type json
		if !strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			c.errorIt(w, r, http.StatusBadRequest, fmt.Errorf("content-type must be application/json"))
			return
		}
		// try to json decode body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			c.errorIt(w, r, http.StatusInternalServerError, err)
			return
		}
		var job app.Job
		if err := json.Unmarshal(body, &job); err != nil {
			c.errorIt(w, r, http.StatusBadRequest, err)
			return
		}
		for _, tag := range job.Tags {
			if tag == "" {
				c.errorIt(w, r, http.StatusBadRequest, fmt.Errorf("tag cannot be empty"))
				return
			}
			events.CreateStream(tag)
			events.Publish("/tag/"+tag, &sse.Event{
				Data: body,
			})
			c.Logger.Infof("created a job for tag %s, %s\n", tag, body)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	return router
}
