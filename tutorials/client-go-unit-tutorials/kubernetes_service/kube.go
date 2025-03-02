package kubernetesservice

import (
	"flag"
	"log"
	"path/filepath"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var CLIENT_SET kubernetes.Interface
var ONCE sync.Once

// DoInit Indexer相关的初始化操作，这里确保只执行一次
func DoInit() {
	ONCE.Do(initInKubernetesEnv)
}

// GetClient 调用此方法返回clientSet对象
func GetClient() kubernetes.Interface {
	return CLIENT_SET
}

// SetClient 可以通过initInKubernetesEnv在kubernetes初始化，如果有准备好的clientSet，也可以调用SetClient直接设置，而无需初始化
func SetClient(clientSet kubernetes.Interface) {
	CLIENT_SET = clientSet
}

// initInKubernetesEnv 这里是真正的初始化逻辑
func initInKubernetesEnv() {
	log.Println("开始初始化Indexer")

	var kubeconfig *string

	// 试图取到当前账号的家目录
	if home := homedir.HomeDir(); home != "" {
		// 如果能取到，就把家目录下的.kube/config作为默认配置文件
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		// 如果取不到，就没有默认配置文件，必须通过kubeconfig参数来指定
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// 加载配置文件
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// 用clientset类来执行后续的查询操作
	CLIENT_SET, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	log.Println("kubernetes服务初始化成功")
}
