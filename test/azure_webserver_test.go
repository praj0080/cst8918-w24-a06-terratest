package test

import (
	"testing"
	

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/stretchr/testify/assert"
)

func TestAzureLinuxVMCreation(t *testing.T) {
	subscriptionId := "e41f28e7-a9fc-46c3-98cc-3a82627e79c8"
	labelPrefix := "praj0080A05VM"

	terraformOptions := &terraform.Options{
		TerraformDir: "../", // Path to your Terraform code
		Vars: map[string]interface{}{
			"labelPrefix": labelPrefix,
		},
		EnvVars: map[string]string{
			"ARM_SUBSCRIPTION_ID": subscriptionId,
		},
	}

	// Wait before destroying resources to avoid Azure dependency issues
	// defer func() {
		// time.Sleep(20 * time.Second)
		// terraform.Destroy(t, terraformOptions)
	// }() 

	terraform.InitAndApply(t, terraformOptions)

	// Outputs to validate
	vmName := terraform.Output(t, terraformOptions, "vm_name")
	nicName := terraform.Output(t, terraformOptions, "nic_name") // used for extra validation
	resourceGroup := terraform.Output(t, terraformOptions, "resource_group_name")

	// ✅ Check VM name is not empty
	assert.NotEmpty(t, vmName)

	// ✅ Confirm NIC exists and is attached to the VM
	nicList, err := azure.GetVirtualMachineNicsE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.NotEmpty(t, nicList, "NIC list should not be empty")
	assert.Contains(t, nicList, nicName, "NIC from output should be attached to the VM")

	t.Logf("NIC(s) attached to VM: %v", nicList)

	// ✅ Confirm the VM is running the correct Ubuntu version
	vmImage, err := azure.GetVirtualMachineImageE(vmName, resourceGroup, subscriptionId)
	assert.NoError(t, err)
	assert.Contains(t, vmImage.SKU, "22_04", "Expected Ubuntu version not found") // or "20_04-lts"
}
