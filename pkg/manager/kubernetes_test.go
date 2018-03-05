package manager

import (
	. "github.com/onsi/ginkgo"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-executor/pkg/model/k8s"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"k8s.io/api/core/v1"
	"log"
	"time"
)

var context = GinkgoT()

var _ = Describe("KubernetesImpl", func() {
	Describe("ManageExecutingTask", func() {
		name := "test"
		image := "alpine"
		init := "init.sh"
		initArgs := []string{"test"}
		spec := &task.TaskSpec{Image: image, Name: name, Init: init, InitArgs: initArgs}
		timeout := time.Millisecond * 10

		var clientSet *fake.Clientset

		It("should return error if job is deleted before execution completes", func() {
			// Arrange
			container := v1.Container{
				Name:    spec.Name,
				Image:   spec.Image,
				Command: []string{spec.Init},
				Args:    spec.InitArgs,
			}

			k8sJob := k8s.Job(spec.Name, []v1.Container{container})
			clientSet = fake.NewSimpleClientset(k8sJob)
			manager := NewKubernetesImpl(clientSet)

			// Act
			successChan := make(chan bool, 2)

			go func() {
				defer GinkgoRecover()
				// Act
				_, err := manager.ManageExecutingTask(k8sJob.Name, make(chan int))

				// Assert
				assert.NotNil(context, err)
				assert.Equal(context, "job has been deleted before execution completed", err.Error())
				successChan <- true
			}()

			time.Sleep(timeout) // we need to wait so the watch job gets the event
			err := clientSet.BatchV1().Jobs(v1.NamespaceDefault).Delete(k8sJob.Name, &v12.DeleteOptions{})
			failOnError(err)

			waitForResult(time.After(timeout), successChan)
		})
	})
})

func waitForResult(timeout <-chan time.Time, successChan <-chan bool) {
WaitLoop:
	for {
		select {
		case <-timeout:
			assert.Fail(context, "Timeout waiting for channel result reached!")
		case ok := <-successChan:
			if ok {
				break WaitLoop
			}
		}
	}
}

func failOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
