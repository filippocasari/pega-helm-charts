package backingservices

import (
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	"github.com/stretchr/testify/require"
)

func TestCustomerAssetVolumeClaimName(t *testing.T) {
	t.Run("Only customerAssetVolumeClaimName is set", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation.enabled": "true",
				"constellation.customerAssetVolumeClaimName": "customer-claim",
			},
				[]string{"charts/constellation/templates/clln-deployment.yaml"}),
		)

		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation",
			Kind: "Deployment",
		}, &deployment)

		volumes := deployment.Spec.Template.Spec.Volumes
		found := false
		for _, vol := range volumes {
			if vol.Name == "constellation-appstatic-assets" && vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == "customer-claim" {
				found = true
				break
			}
		}
		require.True(t, found, "Expected volume with claimName 'customer-claim' not found")
	})

	t.Run("Only custom.volumes is set", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation.enabled": "true",
				"constellation.custom.volumes.name": "custom-volume",
			},
				[]string{"charts/constellation/templates/clln-deployment.yaml"}),
		)

		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation",
			Kind: "Deployment",
		}, &deployment)

		volumes := deployment.Spec.Template.Spec.Volumes
		found := false
		for _, vol := range volumes {
			if vol.Name == "custom-volume" && vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == "custom-volume" {
				found = true
				break
			}
		}
		require.True(t, found, "Expected volume with claimName 'custom-volume' not found")
	})

	t.Run("Both customerAssetVolumeClaimName and custom.volumes are set", func(t *testing.T) {
		helmChartParser := NewHelmConfigParser(
			NewHelmTestFromTemplate(t, helmChartRelativePath, map[string]string{
				"constellation.enabled": "true",
				"constellation.customerAssetVolumeClaimName": "customer-claim",
				"constellation.custom.volumes.name": "custom-volume",
			},
				[]string{"charts/constellation/templates/clln-deployment.yaml"}),
		)

		var deployment appsv1.Deployment
		helmChartParser.getResourceYAML(SearchResourceOption{
			Name: "constellation",
			Kind: "Deployment",
		}, &deployment)

		volumes := deployment.Spec.Template.Spec.Volumes
		foundCustomer := false
		foundCustom := false
		for _, vol := range volumes {
			if vol.Name == "constellation-appstatic-assets" && vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == "customer-claim" {
				foundCustomer = true
			}
			if vol.Name == "custom-volume" && vol.PersistentVolumeClaim != nil && vol.PersistentVolumeClaim.ClaimName == "custom-volume" {
				foundCustom = true
			}
		}
		require.True(t, foundCustomer, "Expected volume with claimName 'customer-claim' not found")
		require.True(t, foundCustom, "Expected volume with claimName 'custom-volume' not found")
	})
}
