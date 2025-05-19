package backingservices

import (
	"testing"

	"github.com/stretchr/testify/require"
	appsv1 "k8s.io/api/apps/v1"
)

func TestConstellationMessagingDeploymentWithAffinity(t *testing.T) {

	var affintiyBasePath = "constellation-messaging.affinity.nodeAffinity.requiredDuringSchedulingIgnoredDuringExecution.nodeSelectorTerms[0].matchExpressions[0]."

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"constellation-messaging.enabled": "true",
			"constellation-messaging.name":    "constellation-messaging",
			affintiyBasePath + "key":          "kubernetes.io/os",
			affintiyBasePath + "operator":     "In",
			affintiyBasePath + "values[0]":    "linux",
		},
			[]string{"charts/constellation-messaging/templates/messaging-deployment.yaml"}),
	)

	var cllnMessagingDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "constellation-messaging",
		Kind: "Deployment",
	}, &cllnMessagingDeploymentObj)

	deploymentAffinity := cllnMessagingDeploymentObj.Spec.Template.Spec.Affinity.NodeAffinity.RequiredDuringSchedulingIgnoredDuringExecution
	require.Equal(t, "kubernetes.io/os", deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Key)
	require.Equal(t, "In", string(deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Operator))
	require.Equal(t, "linux", deploymentAffinity.NodeSelectorTerms[0].MatchExpressions[0].Values[0])
	require.Empty(t, cllnMessagingDeploymentObj.Spec.Template.Spec.Tolerations)
}

func TestConstellationMessagingDeploymentWithTolerations(t *testing.T) {

	helmChartParser := NewHelmConfigParser(
		NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
			"constellation-messaging.enabled":                 "true",
			"constellation-messaging.name":                    "constellation-messaging",
			"constellation-messaging.tolerations[0].key":      "key1",
			"constellation-messaging.tolerations[0].value":    "value1",
			"constellation-messaging.tolerations[0].operator": "Equal",
			"constellation-messaging.tolerations[0].effect":   "NotSchedule",
		},
			[]string{"charts/constellation-messaging/templates/messaging-deployment.yaml"}),
	)

	var cllnMessagingDeploymentObj appsv1.Deployment
	helmChartParser.getResourceYAML(SearchResourceOption{
		Name: "constellation-messaging",
		Kind: "Deployment",
	}, &cllnMessagingDeploymentObj)

	deploymentTolerations := cllnMessagingDeploymentObj.Spec.Template.Spec.Tolerations
	require.Equal(t, "key1", deploymentTolerations[0].Key)
	require.Equal(t, "value1", deploymentTolerations[0].Value)
	require.Equal(t, "Equal", string(deploymentTolerations[0].Operator))
	require.Equal(t, "NotSchedule", string(deploymentTolerations[0].Effect))
	require.Empty(t, cllnMessagingDeploymentObj.Spec.Template.Spec.Affinity)
}

func TestConstellationMessagingDeploymentWithAnnotations(t *testing.T) {
	t.Run("No annotations", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation-messaging.enabled": "true",
				"constellation-messaging.name":    "constellation-messaging",
			},
				[]string{"charts/constellation-messaging/templates/messaging-deployment.yaml"}),
		)
	
		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation-messaging",
			Kind: "Deployment",
		}, &deployment)

		annotations := deployment.ObjectMeta.Annotations
		podAnnotations := deployment.Spec.Template.ObjectMeta.Annotations
		require.Empty(t, annotations, "Expected no annotations, but found some")
		require.Empty(t, podAnnotations, "Expected no pod annotations, but found some")

	})

	t.Run("With annotations", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation-messaging.enabled": "true",
				"constellation-messaging.name":    "constellation-messaging",
				"constellation-messaging.annotations.key1":   "value1",
			},
				[]string{"charts/constellation-messaging/templates/messaging-deployment.yaml"}),
		)
	
		var cllnMessagingDeploymentObj appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation-messaging",
			Kind: "Deployment",
		}, &cllnMessagingDeploymentObj)
	
		require.Equal(t, "value1", cllnMessagingDeploymentObj.ObjectMeta.Annotations["key1"])
	})
	t.Run("With multiple annotations", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation-messaging.enabled": "true",
				"constellation-messaging.name":    "constellation-messaging",
				"constellation-messaging.annotations.key0":   "value0",
				"constellation-messaging.annotations.key1":   "value1",
			},
				[]string{"charts/constellation-messaging/templates/messaging-deployment.yaml"}),
		)
	
		var cllnMessagingDeploymentObj appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation-messaging",
			Kind: "Deployment",
		}, &cllnMessagingDeploymentObj)
	
		require.Equal(t, "value0", cllnMessagingDeploymentObj.ObjectMeta.Annotations["key0"])
		require.Equal(t, "value1", cllnMessagingDeploymentObj.ObjectMeta.Annotations["key1"])
	})
	t.Run("With podAnnotations", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation-messaging.enabled": "true",
				"constellation-messaging.name":    "constellation-messaging",
				"constellation-messaging.podAnnotations.key1":   "value1",

			},
				[]string{"charts/constellation-messaging/templates/messaging-deployment.yaml"}),
		)
	
		var cllnMessagingDeploymentObj appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation-messaging",
			Kind: "Deployment",
		}, &cllnMessagingDeploymentObj)
	
		require.Equal(t, "value1", cllnMessagingDeploymentObj.Spec.Template.ObjectMeta.Annotations["key1"])
	})

	
}
