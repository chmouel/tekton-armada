package armada

import (
	"context"
	"encoding/json"
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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var pmEventRe = regexp.MustCompile(`(\w+|\d+|_|-|:)`)

const (
	defaultTimeout = 5
	Version        = "0.0.1"
)

const smeeChannel = "messages"

const tsFormat = "2006-01-02T15.04.01.000"

type armada struct {
	targetURL string
	logger    *slog.Logger
}

// title returns a copy of the string s with all Unicode letters that begin words
// mapped to their Unicode title case.
func title(source string) string {
	return cases.Title(language.Und, cases.NoLower).String(source)
}

func (a armada) clientSetup() error {
	kubeClients := &Clients{}
	if err := kubeClients.Connect(); err != nil {
		return err
	}
	version := strings.TrimSpace(string(Version))
	s := fmt.Sprintf("Starting armada version: %s", version)
	a.logger.Info(s)
	client := sse.NewClient(a.targetURL, sse.ClientMaxBufferSize(1<<20))
	client.Headers["User-Agent"] = fmt.Sprintf("armada/%s", version)
	// this is to get nginx to work
	client.Headers["X-Accel-Buffering"] = "no"
	channel := filepath.Base(a.targetURL)
	err := client.Subscribe("/tag/"+channel, func(msg *sse.Event) {
		now := time.Now().UTC()
		nowStr := now.Format(tsFormat)

		if string(msg.Event) == "ready" || string(msg.Data) == "ready" {
			s := fmt.Sprintf("%s Listening to %s on cluster %s", nowStr, a.targetURL, kubeClients.Host)
			a.logger.Info(s)
			return
		}

		if string(msg.Event) == "ping" {
			return
		}

		if string(msg.Data) == "{}" {
			return
		}

		var job Job
		if err := json.Unmarshal(msg.Data, &job); err != nil {
			a.logger.Error("Error decoding json: %v", err)
			return
		}
		if len(job.Yamls) == 0 {
			a.logger.Info("no data to apply", slog.Any("payload", job))
			return
		}
		ctx := context.Background()
		alltypes, err := ReadTektonTypes(ctx, a.logger, job.Yamls)
		if err != nil {
			a.logger.Error("Error reading templates: %v", err)
			return
		}

		for _, pipeline := range alltypes.Tekton.Pipelines {
			// check if exist or delete
			if _, err := kubeClients.Tekton.TektonV1().Pipelines(kubeClients.Namespace).Get(ctx, pipeline.GetName(), metav1.GetOptions{}); err == nil {
				if err := kubeClients.Tekton.TektonV1().Pipelines(kubeClients.Namespace).Delete(ctx, pipeline.GetName(), metav1.DeleteOptions{}); err != nil {
					a.logger.Error(fmt.Sprintf("error deleting pipeline%s on %s: %s", pipeline.GetName(), kubeClients.Host, err.Error()))
				} else {
					a.logger.Info(fmt.Sprintf("pipeline %s has been deleted on %s", pipeline.GetName(), kubeClients.Host))
				}
			}
			if cp, err := kubeClients.Tekton.TektonV1().Pipelines(kubeClients.Namespace).Create(ctx, pipeline, metav1.CreateOptions{}); err != nil {
				a.logger.Error(fmt.Sprintf("error creating pipeline%s on %s: %s", cp.GetName(), kubeClients.Host, err.Error()))
			} else {
				a.logger.Info(fmt.Sprintf("pipeline %s has been created on %s", cp.GetName(), kubeClients.Host))
			}
		}

		for _, pipelinerun := range alltypes.Tekton.PipelineRuns {
			if _, err := kubeClients.Tekton.TektonV1().PipelineRuns(kubeClients.Namespace).Get(ctx, pipelinerun.GetName(), metav1.GetOptions{}); err == nil {
				if err := kubeClients.Tekton.TektonV1().PipelineRuns(kubeClients.Namespace).Delete(ctx, pipelinerun.GetName(), metav1.DeleteOptions{}); err != nil {
					a.logger.Error(fmt.Sprintf("error deleting pipelinerun %s on %s: %s", pipelinerun.GetName(), kubeClients.Host, err.Error()))
				} else {
					a.logger.Info(fmt.Sprintf("pipelinerun %s has been deleted on %s", pipelinerun.GetName(), kubeClients.Host))
				}
			}

			if cp, err := kubeClients.Tekton.TektonV1().PipelineRuns(kubeClients.Namespace).Create(ctx, pipelinerun, metav1.CreateOptions{}); err != nil {
				a.logger.Error(fmt.Sprintf("error creating pipelinerun %s on %s: %s", cp.GetName(), kubeClients.Host, err.Error()))
			} else {
				a.logger.Info(fmt.Sprintf("pipelinerun %s has been created on %s", cp.GetName(), kubeClients.Host))
			}
		}
		// dest, err := job.DownloadExtractURL()
		// if err != nil {
		// 	a.logger.Error("Error downloading %s: %v", job.URL, err)
		// 	return
		// }
		// fmt.Printf("dest: %v\n", dest)
	})
	return err
}
