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
	clients   *KClients
}

func NewArmada(logger *slog.Logger, targetURL string) (*armada, error) {
	clients, err := NewKClients()
	if err != nil {
		return nil, err
	}
	return &armada{
		targetURL: targetURL,
		logger:    logger,
		clients:   clients,
	}, nil
}

// title returns a copy of the string s with all Unicode letters that begin words
// mapped to their Unicode title case.
func title(source string) string {
	return cases.Title(language.Und, cases.NoLower).String(source)
}

// apply this function is stupid, i should just use runtime.Object and kubectl
// apply code, but laziness you know
func (a armada) apply(ctx context.Context, job Job) error {
	alltypes, err := ReadTektonTypes(ctx, a.logger, job.Yamls)
	if err != nil {
		return fmt.Errorf("Error reading templates: %w", err)
	}

	for _, configmap := range alltypes.Kube.ConfigMaps {
		if _, err := a.clients.Kubernetes.CoreV1().ConfigMaps(a.clients.Namespace).Get(ctx, configmap.GetName(), metav1.GetOptions{}); err == nil {
			if err := a.clients.Kubernetes.CoreV1().ConfigMaps(a.clients.Namespace).Delete(ctx, configmap.GetName(), metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("error deleting configmap %s on %s: %w", configmap.GetName(), a.clients.Host, err)
			} else {
				a.logger.Info(fmt.Sprintf("configmap %s has been deleted on %s", configmap.GetName(), a.clients.Host))
			}
		}
		if cp, err := a.clients.Kubernetes.CoreV1().ConfigMaps(a.clients.Namespace).Create(ctx, configmap, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("error creating configmap %s on %s: %w", cp.GetName(), a.clients.Host, err)
		} else {
			a.logger.Info(fmt.Sprintf("configmap %s has been created on %s", cp.GetName(), a.clients.Host))
		}
	}

	for _, secret := range alltypes.Kube.Secrets {
		if _, err := a.clients.Kubernetes.CoreV1().Secrets(a.clients.Namespace).Get(ctx, secret.GetName(), metav1.GetOptions{}); err == nil {
			if err := a.clients.Kubernetes.CoreV1().Secrets(a.clients.Namespace).Delete(ctx, secret.GetName(), metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("error deleting secret %s on %s: %w", secret.GetName(), a.clients.Host, err)
			} else {
				a.logger.Info(fmt.Sprintf("secret %s has been deleted on %s", secret.GetName(), a.clients.Host))
			}
		}
		if cp, err := a.clients.Kubernetes.CoreV1().Secrets(a.clients.Namespace).Create(ctx, secret, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("error creating secret %s on %s: %w", cp.GetName(), a.clients.Host, err)
		} else {
			a.logger.Info(fmt.Sprintf("secret %s has been created on %s", cp.GetName(), a.clients.Host))
		}
	}

	for _, pipeline := range alltypes.Tekton.Pipelines {
		if _, err := a.clients.Tekton.TektonV1().Pipelines(a.clients.Namespace).Get(ctx, pipeline.GetName(), metav1.GetOptions{}); err == nil {
			if err := a.clients.Tekton.TektonV1().Pipelines(a.clients.Namespace).Delete(ctx, pipeline.GetName(), metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("error deleting pipeline %s on %s: %w", pipeline.GetName(), a.clients.Host, err)
			} else {
				a.logger.Info(fmt.Sprintf("pipeline %s has been deleted on %s", pipeline.GetName(), a.clients.Host))
			}
		}
		if cp, err := a.clients.Tekton.TektonV1().Pipelines(a.clients.Namespace).Create(ctx, pipeline, metav1.CreateOptions{}); err != nil {
			a.logger.Error(fmt.Sprintf("error creating pipeline%s on %s: %s", cp.GetName(), a.clients.Host, err.Error()))
		} else {
			a.logger.Info(fmt.Sprintf("pipeline %s has been created on %s", cp.GetName(), a.clients.Host))
		}
	}

	for _, pipelinerun := range alltypes.Tekton.PipelineRuns {
		if _, err := a.clients.Tekton.TektonV1().PipelineRuns(a.clients.Namespace).Get(ctx, pipelinerun.GetName(), metav1.GetOptions{}); err == nil {
			if err := a.clients.Tekton.TektonV1().PipelineRuns(a.clients.Namespace).Delete(ctx, pipelinerun.GetName(), metav1.DeleteOptions{}); err != nil {
				return fmt.Errorf("error deleting pipelinerun %s on %s: %w", pipelinerun.GetName(), a.clients.Host, err)
			} else {
				a.logger.Info(fmt.Sprintf("pipelinerun %s has been deleted on %s", pipelinerun.GetName(), a.clients.Host))
			}
		}

		if cp, err := a.clients.Tekton.TektonV1().PipelineRuns(a.clients.Namespace).Create(ctx, pipelinerun, metav1.CreateOptions{}); err != nil {
			return fmt.Errorf("error creating pipelinerun %s on %s: %v", cp.GetName(), a.clients.Host, err)
		} else {
			a.logger.Info(fmt.Sprintf("pipelinerun %s has been created on %s", cp.GetName(), a.clients.Host))
		}
	}
	return nil
}

func (a armada) clientSetup(ctx context.Context) error {
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
			s := fmt.Sprintf("%s Listening to %s on cluster %s", nowStr, a.targetURL, a.clients.Host)
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
		if err := a.apply(ctx, job); err != nil {
			a.logger.Error(err.Error())
			return
		}
	})
	return err
}
