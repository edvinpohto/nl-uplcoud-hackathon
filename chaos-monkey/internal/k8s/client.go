package k8s

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
)

type Client struct {
	cs     *kubernetes.Clientset
	config *rest.Config
}

// NewFromBase64 creates a Kubernetes client from a base64-encoded kubeconfig.
func NewFromBase64(kubeconfigB64 string) (*Client, error) {
	decoded, err := base64.StdEncoding.DecodeString(kubeconfigB64)
	if err != nil {
		return nil, fmt.Errorf("decoding kubeconfig: %w", err)
	}

	cfg, err := clientcmd.RESTConfigFromKubeConfig(decoded)
	if err != nil {
		return nil, fmt.Errorf("building rest config: %w", err)
	}

	cs, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("creating clientset: %w", err)
	}

	return &Client{cs: cs, config: cfg}, nil
}

// ListRunningPods returns all running pods in the given namespace.
func (c *Client) ListRunningPods(ctx context.Context, namespace string) ([]corev1.Pod, error) {
	pods, err := c.cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		FieldSelector: "status.phase=Running",
	})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// DeletePod deletes a pod by name.
func (c *Client) DeletePod(ctx context.Context, namespace, name string) error {
	return c.cs.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
}

// ListNodes returns all nodes in the cluster.
func (c *Client) ListNodes(ctx context.Context) ([]corev1.Node, error) {
	nodes, err := c.cs.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return nodes.Items, nil
}

// ListDaemonSetPods returns pods from a DaemonSet matching the given label selector.
func (c *Client) ListDaemonSetPods(ctx context.Context, namespace, labelSelector string) ([]corev1.Pod, error) {
	pods, err := c.cs.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{
		LabelSelector: labelSelector,
	})
	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

// ScaleDeployment patches a Deployment's replica count.
func (c *Client) ScaleDeployment(ctx context.Context, namespace, name string, replicas int32) error {
	scale, err := c.cs.AppsV1().Deployments(namespace).GetScale(ctx, name, metav1.GetOptions{})
	if err != nil {
		return fmt.Errorf("get scale: %w", err)
	}
	scale.Spec.Replicas = replicas
	_, err = c.cs.AppsV1().Deployments(namespace).UpdateScale(ctx, name, scale, metav1.UpdateOptions{})
	return err
}

// GetDeploymentReplicas returns current desired replica count.
func (c *Client) GetDeploymentReplicas(ctx context.Context, namespace, name string) (int32, error) {
	d, err := c.cs.AppsV1().Deployments(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return 0, err
	}
	if d.Spec.Replicas == nil {
		return 1, nil
	}
	return *d.Spec.Replicas, nil
}

// ExecInPod runs a command inside a container and returns stdout+stderr.
func (c *Client) ExecInPod(ctx context.Context, namespace, podName, container string, command []string) (string, string, error) {
	req := c.cs.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(podName).
		Namespace(namespace).
		SubResource("exec").
		VersionedParams(&corev1.PodExecOptions{
			Container: container,
			Command:   command,
			Stdin:     false,
			Stdout:    true,
			Stderr:    true,
			TTY:       false,
		}, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(c.config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("create executor: %w", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.StreamWithContext(ctx, remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	})
	return stdout.String(), stderr.String(), err
}
