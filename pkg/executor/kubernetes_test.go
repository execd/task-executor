package executor

import (
	. "github.com/onsi/ginkgo"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/wayofthepie/task-executor/mocks/manager"
	"github.com/wayofthepie/task-executor/pkg/model/task"
	"k8s.io/client-go/kubernetes/fake"
)

var context = GinkgoT()

var _ = Describe("KubernetesImpl", func() {
	Describe("Execute", func() {
		var name string
		var image string
		var init string
		var initArgs []string
		var spec *task.Spec
		var clientSet *fake.Clientset
		var managerMock *manager.Service
		var service *KubernetesImpl
		var executable Executable

		BeforeEach(func() {
			name = "test"
			image = "alpine"
			init = "init.sh"
			initArgs = []string{"test"}
			spec = &task.Spec{Image: image, Name: name, Init: init, InitArgs: initArgs}
			executable = NewExecutableImpl(spec)
			clientSet = fake.NewSimpleClientset()
			managerMock = &manager.Service{}
			service = NewKubernetesClientImpl(clientSet, managerMock)
		})

		It("should manage the execution of the task", func() {
			// Arrange
			managerMock.On("ManageExecutingTask", mock.AnythingOfType("Info"), mock.Anything).Return(nil, nil)

			// Act
			service.ExecuteTask(executable)

			// Assert
			managerMock.AssertNumberOfCalls(context, "ManageExecutingTask", 1)
		})

		It("should return error if management of task fails", func() {
			// Arrange
			expectedErr := errors.New("error")
			managerMock.On("ManageExecutingTask", mock.AnythingOfType("Info"), mock.Anything).Return(nil, expectedErr)

			// Act
			_, err := service.ExecuteTask(executable)

			// Assert
			assert.NotNil(context, err)
			assert.Equal(context, expectedErr, err)
		})
	})
})
