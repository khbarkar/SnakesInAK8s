package k8s

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// PodInfo holds the minimal info we need from a running pod.
type PodInfo struct {
	Name      string
	Namespace string
}

// Client wraps the Kubernetes clientset for pod operations.
type Client struct {
	clientset   *kubernetes.Clientset
	rawConfig   api.Config
	clusterName string
	namespace   string // empty string means all namespaces
}

// ResolveKubeconfig returns the kubeconfig path by checking, in order:
// 1. Explicit flag value (if non-empty)
// 2. KUBECONFIG environment variable
// 3. Default ~/.kube/config
func ResolveKubeconfig(flagValue string) string {
	if flagValue != "" {
		return flagValue
	}
	if env := os.Getenv("KUBECONFIG"); env != "" {
		return env
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".kube", "config")
}

// NewClient builds a Client from the given kubeconfig path.
// Pass an empty namespace to operate across all namespaces.
func NewClient(kubeconfigPath, namespace string) (*Client, error) {
	loadingRules := &clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath}
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

	rawConfig, err := kubeConfig.RawConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load kubeconfig: %w", err)
	}

	restConfig, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to build rest config: %w", err)
	}

	cs, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create clientset: %w", err)
	}

	clusterName := rawConfig.CurrentContext
	if ctx, ok := rawConfig.Contexts[rawConfig.CurrentContext]; ok && ctx.Cluster != "" {
		clusterName = ctx.Cluster
	}

	return &Client{
		clientset:   cs,
		rawConfig:   rawConfig,
		clusterName: clusterName,
		namespace:   namespace,
	}, nil
}

// ClusterName returns the name of the current cluster context.
func (c *Client) ClusterName() string {
	return c.clusterName
}

// Namespace returns the configured namespace filter (empty = all).
func (c *Client) Namespace() string {
	return c.namespace
}

// SetNamespace updates the namespace filter.
func (c *Client) SetNamespace(ns string) {
	c.namespace = ns
}

// ListNamespaces returns all namespace names in the cluster.
func (c *Client) ListNamespaces(ctx context.Context) ([]string, error) {
	nsList, err := c.clientset.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list namespaces: %w", err)
	}
	names := make([]string, 0, len(nsList.Items))
	for _, ns := range nsList.Items {
		names = append(names, ns.Name)
	}
	return names, nil
}

// RandomPod picks a random running pod, filtered by the client's namespace.
// If namespace is empty, picks from all namespaces.
// Only picks pods with the label app=snakefood to avoid killing real workloads.
// Pods whose names appear in exclude are skipped.
func (c *Client) RandomPod(ctx context.Context, exclude map[string]bool) (*PodInfo, error) {
	ns := c.namespace // empty string = all namespaces in the API

	pods, err := c.clientset.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{
		LabelSelector: "app=snakefood",
		FieldSelector: "status.phase=Running",
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}

	// Filter out excluded pods
	var candidates []PodInfo
	for _, p := range pods.Items {
		if !exclude[p.Name] {
			candidates = append(candidates, PodInfo{Name: p.Name, Namespace: p.Namespace})
		}
	}

	if len(candidates) == 0 {
		return nil, nil
	}

	pick := candidates[rand.Intn(len(candidates))]
	return &pick, nil
}

// KillPod force-deletes the given pod. Brutal.
func (c *Client) KillPod(ctx context.Context, name, namespace string) error {
	gracePeriod := int64(0)
	err := c.clientset.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{
		GracePeriodSeconds: &gracePeriod,
	})
	if err != nil {
		return fmt.Errorf("failed to kill pod %s/%s: %w", namespace, name, err)
	}
	return nil
}
