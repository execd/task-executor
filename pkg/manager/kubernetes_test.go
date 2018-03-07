package manager

import (
	. "github.com/onsi/ginkgo"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"

	"fmt"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/wayofthepie/task-executor/pkg/model/k8s"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	v13 "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
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
		spec := &task.Spec{Image: image, Name: name, Init: init, InitArgs: initArgs}
		timeout := time.Millisecond * 100

		var k8sJob *v13.Job
		var manager *KubernetesImpl
		var clientSet kubernetes.Interface
		var taskInfo task.Info
		BeforeEach(func() {
			container := v1.Container{
				Name:    spec.Name,
				Image:   spec.Image,
				Command: []string{spec.Init},
				Args:    spec.InitArgs,
			}

			k8sJob = k8s.Job(spec.Name, []v1.Container{container})
			clientSet = fake.NewSimpleClientset(k8sJob)
			manager = NewKubernetesImpl(clientSet)
			uuid := uuid.Must(uuid.NewV4())
			taskInfo = task.Info{
				ID:       &uuid,
				Metadata: k8sJob,
			}
		})

		It("should return error if job is deleted before execution completes", func() {
			// Arrange
			successChan := make(chan bool, 2)

			go func() {
				defer GinkgoRecover()
				// Act
				_, err := manager.ManageExecutingTask(taskInfo, make(chan int))

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

		It("should return task information if job succeeds", func() {
			// Arrange
			successChan := make(chan bool, 2)

			go func() {
				defer GinkgoRecover()
				// Act
				info, err := manager.ManageExecutingTask(taskInfo, make(chan int))

				// Assert
				retrievedJob := info.Metadata.(*v13.Job)
				assert.Nil(context, err)
				assert.Equal(context, int32(1), retrievedJob.Status.Succeeded)
				assert.Equal(context, name, retrievedJob.GenerateName)
				successChan <- true
			}()

			time.Sleep(timeout) // we need to wait so the watch job gets the event
			jobUpdate := *k8sJob
			jobUpdate.Status = v13.JobStatus{
				Succeeded: 1,
			}
			fmt.Println(jobUpdate.Name)
			_, err := clientSet.BatchV1().Jobs(v1.NamespaceDefault).UpdateStatus(&jobUpdate)
			failOnError(err)

			waitForResult(time.After(time.Second*1), successChan)
		})

		It("should return info if job fails", func() {
			// Arrange
			successChan := make(chan bool, 2)

			go func() {
				defer GinkgoRecover()
				// Act
				info, err := manager.ManageExecutingTask(taskInfo, make(chan int))

				// Assert
				retrievedJob := info.Metadata.(*v13.Job)
				assert.Nil(context, err)
				assert.Equal(context, int32(1), retrievedJob.Status.Failed)
				successChan <- true
			}()

			time.Sleep(timeout) // we need to wait so the watch job gets the event
			jobUpdate := *k8sJob
			jobUpdate.Status = v13.JobStatus{
				Failed: 1,
			}
			fmt.Println(jobUpdate.Name)
			_, err := clientSet.BatchV1().Jobs(v1.NamespaceDefault).UpdateStatus(&jobUpdate)
			failOnError(err)

			waitForResult(time.After(time.Second*1), successChan)
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
