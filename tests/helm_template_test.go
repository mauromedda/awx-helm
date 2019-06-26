// +build kubeall helm

// **NOTE**: we have build tags to differentiate kubernetes tests from non-kubernetes tests, and further differentiate helm
// tests. This is done because minikube is heavy and can interfere with docker related tests in terratest. Similarly, helm
// can overload the minikube system and thus interfere with the other kubernetes tests. Specifically, many of the tests
// start to fail with `connection refused` errors from `minikube`. To avoid overloading the system, we run the kubernetes
// tests and helm tests separately from the others. This may not be necessary if you have a sufficiently powerful machine.
// We recommend at least 4 cores and 16GB of RAM if you want to run all the tests together.

package tests

import (
	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
	v1beta1 "k8s.io/api/extensions/v1beta1"
	v1 "k8s.io/api/core/v1"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/helm"
)

// This file uses terratest to test helm chart template logic by rendering the templates
// using `helm template`, and then reading in the rendered templates.
// There are two tests:
// - TestHelmBasicTemplateRenderedDeployment: It Read and rendered the Deployment object and check the
//   computed values.

// Tests if the rendered Deployment template object of a Helm Chart given various inputs.
func TestHelmBasicTemplateRenderedDeployment(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../awx-helm")
	require.NoError(t, err)

	// Since we aren't deploying any resources, there is no need to setup kubectl authentication, helm home, or
	// namespaces

	// Setup the args. For this test, we will set the following input values:
	// - containerImageRepo=ansible/awx_task
	// - containerImageTag=5.0.0
	options := &helm.Options{
		SetValues: map[string]string{
			"awx_task.image.repository": "ansible/awx_task",
			"awx_task.image.tag":        "300.0.0",
		},
	}

	// Run RenderTemplate to render the template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	// Additionally, although we know there is only one yaml file in the template, we deliberately path a templateFiles
	// arg to demonstrate how to select individual templates to render.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/deployment.yaml"})

	// Now we use kubernetes/client-go library to render the template output into the Deployment struct. This will
	// ensure the Deployment resource is rendered correctly.
	var deployment appsv1.Deployment
	helm.UnmarshalK8SYaml(t, output, &deployment)

	// Finally, we verify the deployment pod template spec is set to the expected container image value
	expectedContainerImage := "ansible/awx_task:300.0.0"
	deploymentContainers := deployment.Spec.Template.Spec.Containers
	require.Equal(t, len(deploymentContainers), 4)
	require.Equal(t, deploymentContainers[1].Image, expectedContainerImage)

}

// Tests if the rendered ConfigMap template object of an Helm Chart.
func TestHelmBasicTemplateRenderedConfigMap(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../awx-helm")
	require.NoError(t, err)

	// Since we aren't deploying any resources, there is no need to setup kubectl authentication, helm home, or
	// namespaces

	// Setup the args.
	options := &helm.Options{
		SetValues: map[string]string{},
	}

	// Run RenderTemplate to render the template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	// Additionally, although we know there is only one yaml file in the template, we deliberately path a templateFiles
	// arg to demonstrate how to select individual templates to render.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/configmap.yaml"})

	// Finally, we verify the configmap template to see if the number of numbers and data are correct
	var configMapList v1.ConfigMapList
	helm.UnmarshalK8SYaml(t, output, &configMapList)
	require.Equal(t, len(configMapList.Items), 2)

	testcases := []struct {
		testname     string
		expectedData string
		itemNumber int
	}{
		{
			testname:     "Check awx_settings",
			expectedData: "awx_settings",
			itemNumber: 0,
			
		},
		{
			testname:     "Check secret_key",
			expectedData: "secret_key",
			itemNumber: 0,
		},
		{
			testname:     "Check provision_awx.sh",
			expectedData: "provision_awx.sh",
			itemNumber: 0,
		},
		{
			testname:     "Check rabbitmq.conf",
			expectedData: "rabbitmq.conf",
			itemNumber: 1,
		},
		{
			testname:     "Check rabbitmq.conf",
			expectedData: "rabbitmq.conf",
			itemNumber: 1,
		},
		{
			testname:     "Check rabbitmq_definitions.json",
			expectedData: "rabbitmq_definitions.json",
			itemNumber: 1,
		},
		{
			testname:     "Check enabled_plugins",
			expectedData: "enabled_plugins",
			itemNumber: 1,
		},
	}
	for _, tt := range testcases {
		t.Run(tt.testname, func(t *testing.T) {
			if _, ok := configMapList.Items[tt.itemNumber].Data[tt.expectedData]; !ok {
				t.Fatalf("ConfigMap %s item not present", tt.expectedData)
			}
		})
	}
}

// Tests if the rendered Ingress template object of an Helm Chart.
func TestHelmBasicTemplateRenderedIngress(t *testing.T) {
	t.Parallel()

	// Path to the helm chart we will test
	helmChartPath, err := filepath.Abs("../../awx-helm")
	require.NoError(t, err)

	// Since we aren't deploying any resources, there is no need to setup kubectl authentication, helm home, or
	// namespaces
	expectedIngressHost := "test.domain.local"
	// Setup the args.
	options := &helm.Options{
		SetValues: map[string]string{
			"ingress.enabled": "true",
			"ingress.host": expectedIngressHost,
		},
	}

	// Run RenderTemplate to render the template and capture the output. Note that we use the version without `E`, since
	// we want to assert that the template renders without any errors.
	// Additionally, although we know there is only one yaml file in the template, we deliberately path a templateFiles
	// arg to demonstrate how to select individual templates to render.
	output := helm.RenderTemplate(t, options, helmChartPath, []string{"templates/ingress.yaml"})

	// Finally, we verify the ingress template to see if the number of numbers and data are correct
	var ingress v1beta1.Ingress
	helm.UnmarshalK8SYaml(t, output, &ingress)
	require.Equal(t, ingress.Spec.Rules[0].Host, expectedIngressHost)
}