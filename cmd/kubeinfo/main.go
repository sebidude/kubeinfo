package main

import (
	"context"
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	gitcommit    string
	appversion   string
	buildtime    string
	clientconfig *rest.Config
	clientset    *kubernetes.Clientset

	namespace string
	output    string
	name      = ""
	kind      string
)

func main() {

	app := kingpin.New(os.Args[0], "encrypt decrypt data, convert yaml maps to kubernetes secrets and edit kubernetes secrets.")
	app.Flag("namespace", "Kubernetes namespace to be used.").Short('n').Envar("KUBECRYPT_NAMESPACE").StringVar(&namespace)

	get := app.Command("image", "Get the data.")
	get.Arg("name", "Name of the object to list the image.").StringVar(&name)

	kubeconfig := os.Getenv("KUBECONFIG")
	if len(kubeconfig) < 1 {
		config, err := rest.InClusterConfig()
		if err != nil {
			panic(err.Error())
		}
		clientconfig = config
	} else {
		config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			panic(err.Error())
		}
		clientconfig = config
	}
	var err error
	clientset, err = kubernetes.NewForConfig(clientconfig)
	if err != nil {
		panic(err.Error())
	}

	operation := kingpin.MustParse(app.Parse(os.Args[1:]))

	if namespace == "" {
		namespace, _, err = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		).Namespace()
		if err != nil {
			panic(err.Error())
		}
	}

	switch operation {
	case "image":
		listObjectImage(name)

	case "version":
		fmt.Printf("kubecrypt\n version: %s\n commit: %s\n buildtime: %s\n", appversion, gitcommit, buildtime)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func listObjectImage(name string) {
	ctx := context.Background()
	defer ctx.Done()
	deployments, err := clientset.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	checkError(err)
	for _, obj := range deployments.Items {
		for _, c := range obj.Spec.Template.Spec.Containers {
			if name != "" {
				if name == obj.Name && name == c.Name {
					fmt.Println(c.Image)
				}
			} else {
				fmt.Println(c.Image)
			}
		}
	}

	daemonsets, err := clientset.AppsV1().DaemonSets(namespace).List(ctx, metav1.ListOptions{})
	checkError(err)
	for _, obj := range daemonsets.Items {
		for _, c := range obj.Spec.Template.Spec.Containers {
			if name != "" {
				if name == obj.Name && name == c.Name {
					fmt.Println(c.Image)
				}
			} else {
				fmt.Println(c.Image)
			}
		}
	}

	statefulsets, err := clientset.AppsV1().StatefulSets(namespace).List(ctx, metav1.ListOptions{})
	checkError(err)
	for _, obj := range statefulsets.Items {
		for _, c := range obj.Spec.Template.Spec.Containers {
			if name != "" {
				if name == obj.Name && name == c.Name {
					fmt.Println(c.Image)
				}
			} else {
				fmt.Println(c.Image)
			}
		}
	}

}
