package k8s

import (
	"k8s.io/api/batch/v1"
	v13 "k8s.io/api/core/v1"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Construct a kubernetes V1 Batch Job.
func Job(jobNamePrefix string, containers []v13.Container) *v1.Job {
	return &v1.Job{
		ObjectMeta: v12.ObjectMeta{
			GenerateName: jobNamePrefix,
			Namespace:    v12.NamespaceDefault,
		},
		Spec: v1.JobSpec{
			BackoffLimit: new(int32),
			Template: v13.PodTemplateSpec{
				Spec: v13.PodSpec{
					RestartPolicy: "Never",
					Containers:    containers,
				},
			},
		},
	}
}
