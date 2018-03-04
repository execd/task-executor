package executor

import (
	. "github.com/onsi/ginkgo"
	"k8s.io/client-go/kubernetes/fake"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/pkg/errors"
	"github.com/wayofthepie/task-executor/mocks/manager"
)

var context = GinkgoT()

var _ = Describe("KubernetesImpl", func() {
	Describe("ExecuteTask", func() {
		var name string
		var image string
		var init string
		var initArgs []string
		var spec *task.TaskSpec
		var clientSet *fake.Clientset
		var managerMock *manager.Service
		var service *KubernetesImpl

		BeforeEach(func() {
			name = "test"
			image = "alpine"
			init = "init.sh"
			initArgs = []string{"test"}
			spec = &task.TaskSpec{Image: image, Name: name, Init: init, InitArgs: initArgs}
			clientSet = fake.NewSimpleClientset()
			managerMock = &manager.Service{}
			service = NewKubernetesClientImpl(clientSet, managerMock)
		})

		It("should manage the execution of the task", func() {
			// Arrange
			managerMock.On("ManageExecutingTask", "", mock.Anything).Return(nil, nil)

			// Act
			service.ExecuteTask(spec)

			// Assert
			managerMock.AssertNumberOfCalls(context, "ManageExecutingTask", 1)
		})

		It("should return error if management of task fails", func() {
			// Arrange
			expectedErr := errors.New("error")
			managerMock.On("ManageExecutingTask", "", mock.Anything).Return(nil, expectedErr)

			// Act
			_, err := service.ExecuteTask(spec)

			// Assert
			assert.NotNil(context, err)
			assert.Equal(context, expectedErr, err)
		})
	})
})
