package main

import (
	"context"
	"fmt"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/remotecommand"
	"strings"
	"time"
	"k8s.io/klog/v2"
)

func main() {

	klog.InitFlags(nil)

	kubeConfig := kubeConfig()

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if nil != err {
		panic(err.Error())
	}

	namespace, _, err := kubeConfig.Namespace()
	if err != nil {
		panic(err.Error())
	}

	pods, err := client.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%d pods in namespace %s\n", len(pods.Items), namespace)

	for _, pod := range pods.Items {
		fmt.Println(pod.Name)
	}

	fmt.Println("starting KRM Function in kubernetes")

	pod, err := client.CoreV1().Pods(namespace).Create(context.TODO(), &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "ciis",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:      "ciis",
					Image:     "minikube:5000/echo-func:3",
					Stdin:     true,
					StdinOnce: true,
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}, metav1.CreateOptions{})

	defer func(){
		fmt.Println("cleaning up KRM function")
		err = client.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
		if err != nil {
			panic(err)
		}
	}()

	if err != nil {
		panic(err.Error())
	}

	// TODO: implement proper wait
	time.Sleep(1 * time.Second) // as this is all within minikube 1 sec should be enough

	req := client.CoreV1().RESTClient().Post().Namespace(pod.Namespace).Resource("pods").Name(pod.Name).
		SubResource("attach").VersionedParams(&v1.PodAttachOptions{
		Stdin: true,
		Stdout: true,
		Stderr: true,
	}, scheme.ParameterCodec)

	fmt.Println("URL", req.URL())

	attach, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		panic(err)
	}

	stdin := strings.NewReader("hello\n")
	stdout := &strings.Builder{}
	stderr := &strings.Builder{}

	fmt.Println("attaching...")

	err = attach.Stream(remotecommand.StreamOptions{
		Stdin: stdin,
		Stdout: stdout,
		Stderr: stderr,
	})

	if err != nil {
		fmt.Println("err?", err)
	}

	fmt.Printf("KRM function result:\n---%s---\n", stdout.String())
	fmt.Println("KRM function errors:", stderr.String())

}

func kubeConfig() clientcmd.ClientConfig {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}

	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)

}
