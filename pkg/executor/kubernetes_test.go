package executor

import (
	. "github.com/onsi/ginkgo"
	"k8s.io/client-go/kubernetes/fake"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"github.com/stretchr/testify/assert"
	"k8s.io/api/batch/v1"
	"fmt"
)

var context = GinkgoT()

var _ = Describe("KubernetesImpl", func() {
	Describe("ExecuteTask", func() {
		It("should call kubernetes with correct Job details", func() {
			// Arrange
			name := "test"
			image := "alpine"
			init := "init.sh"
			initArgs := []string{"test"}
			spec := &task.TaskSpec{Image: image, Name: name, Init: init, InitArgs: initArgs}
			clientSet := fake.NewSimpleClientset()
			service := NewKubernetesClientImpl(clientSet)

			// Act
			info, err := service.ExecuteTask(spec)

			// Assert

			assert.Nil(context, err)
			assert.IsType(context, &v1.Job{}, info.Metadata)
			// Verify a "-" is added to GenerateName, the server will add a uid post fix to this
			// and set the field Name. I can't figure out how to test the Name  field directly here.
			assert.Equal(context, fmt.Sprintf("%s-", name), info.Metadata.(*v1.Job).GenerateName)
			assert.Equal(context, image, info.Metadata.(*v1.Job).Spec.Template.Spec.Containers[0].Image)
			assert.Equal(context, init, info.Metadata.(*v1.Job).Spec.Template.Spec.Containers[0].Command[0])
		})
	})
})
