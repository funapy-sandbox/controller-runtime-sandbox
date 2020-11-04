package main

import (
	"context"

	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func run() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	scheme := runtime.NewScheme()
	if err := batchv1.AddToScheme(scheme); err != nil {
		return err
	}

	c, err := client.New(cfg, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		return err
	}

	err = c.Create(context.Background(), &batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind: "batch/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      "sleep",
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			Completions: func() *int32 {
				var n int32 = 5
				return &n
			}(),
			Parallelism: func() *int32 {
				var n int32 = 3
				return &n
			}(),
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:  "sleep",
							Image: "alpine",
							Command: []string{
								"sh", "-c",
							},
							Args: []string{
								"sleep 3",
							},
						},
					},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}
