package armada

import (
	"context"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"

	"golang.org/x/exp/slog"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonversioned "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	k8scheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
)

type Clients struct {
	Kubernetes *kubernetes.Clientset
	Tekton     *tektonversioned.Clientset
	Dynamic    dynamic.Interface
	Namespace  string
	Host       string
}

func (c *Clients) Connect() error {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kclient := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	kubeConfig, err := kclient.ClientConfig()
	if err != nil {
		return err
	}

	if c.Namespace, _, err = kclient.Namespace(); err != nil {
		return err
	}

	if c.Kubernetes, err = kubernetes.NewForConfig(kubeConfig); err != nil {
		return err
	}

	if c.Tekton, err = tektonversioned.NewForConfig(kubeConfig); err != nil {
		return err
	}

	if c.Dynamic, err = dynamic.NewForConfig(kubeConfig); err != nil {
		return err
	}

	c.Host = kubeConfig.Host
	return nil
}

type TektonTypes struct {
	PipelineRuns []*tektonv1.PipelineRun
	Pipelines    []*tektonv1.Pipeline
	TaskRuns     []*tektonv1.TaskRun
	Tasks        []*tektonv1.Task
}

type KubeTypes struct {
	Secrets    []*corev1.Secret
	ConfigMaps []*corev1.ConfigMap
}

type Types struct {
	Kube   KubeTypes
	Tekton TektonTypes
}

var yamlDocSeparatorRe = regexp.MustCompile(`(?m)^---\s*$`)

func ReadTektonTypes(ctx context.Context, log *slog.Logger, yamls []string) (Types, error) {
	types := Types{}
	decoder := k8scheme.Codecs.UniversalDeserializer()

	for _, base64data := range yamls {
		// decode base64data
		bdata, err := base64.StdEncoding.DecodeString(base64data)
		if err != nil {
			return Types{}, err
		}

		for _, doc := range yamlDocSeparatorRe.Split(string(bdata), -1) {
			if strings.TrimSpace(doc) == "" {
				continue
			}

			obj, _, err := decoder.Decode([]byte(doc), nil, nil)
			if err != nil {
				log.Error(fmt.Sprintf("Skipping document not looking like a kubernetes resources: %s", err.Error()))
				continue
			}
			switch o := obj.(type) {
			case *tektonv1.PipelineRun:
				types.Tekton.PipelineRuns = append(types.Tekton.PipelineRuns, o)
			case *tektonv1.Pipeline:
				types.Tekton.Pipelines = append(types.Tekton.Pipelines, o)
			case *tektonv1.Task:
				types.Tekton.Tasks = append(types.Tekton.Tasks, o)
			case *corev1.Secret:
				types.Kube.Secrets = append(types.Kube.Secrets, o)
			case *corev1.ConfigMap:
				types.Kube.ConfigMaps = append(types.Kube.ConfigMaps, o)
			default:
				log.Info("Skipping document not looking like a tekton resource we can Resolve.")
			}
		}
	}

	return types, nil
}

//nolint:gochecknoinits
func init() {
	_ = tektonv1.AddToScheme(k8scheme.Scheme)
	_ = corev1.AddToScheme(k8scheme.Scheme)
}
