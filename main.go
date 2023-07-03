package main

import (
	"context"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
	"k8s.io/klog/v2"
)

func main() {

	config, err := clientcmd.BuildConfigFromFlags("", "/Users/anhongkui/.kube/config")
	PanicIf(err)

	client, err := clientset.NewForConfig(config)
	PanicIf(err)

	leader, err := leaderelection.NewLeaderElector(leaderelection.LeaderElectionConfig{
		ReleaseOnCancel: true,
		LeaseDuration:   60 * time.Second, // 租约到期时间
		RenewDeadline:   15 * time.Second, // 租约续期间隔？
		RetryPeriod:     5 * time.Second,

		Lock: &resourcelock.LeaseLock{
			LeaseMeta:  metav1.ObjectMeta{Namespace: "default", Name: "my-lock-test"},
			Client:     client.CoordinationV1(),
			LockConfig: resourcelock.ResourceLockConfig{Identity: "my-random-id2"}, // 本实例的唯一标识，每个实例的标识不一样
		},

		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(ctx context.Context) { klog.Info("I'm leader") },
			OnStoppedLeading: func() { klog.Info("leader lost") },
			OnNewLeader:      func(identity string) { klog.Infof("new leader elected: %s", identity) }, // 能够得知当前Leader的唯一标识
		},
	})

	PanicIf(err)
	leader.Run(context.Background()) // Never Return
}

func PanicIf(e any) {
	if e != nil {
		fmt.Println(e)
		panic(e)
	}
}
