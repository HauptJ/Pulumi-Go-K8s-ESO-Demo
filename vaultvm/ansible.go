package main

import (
	"fmt"
	"os"

	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi-command/sdk/go/command/remote"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func provVaultVM(ctx *pulumi.Context, ipv6 pulumi.StringOutput) error {
	var errr error = nil

	cfg := config.New(ctx, "")
	doPrivateKeyPath := cfg.Require("doPrivateKeyPath")
	doVmUser := cfg.Require("doVmUser")
	vaultFqdn := cfg.Require("vaultFqdn")
	vaultTlsDisable := cfg.RequireBool("vaultTlsDisable")
	vaultTlsPath := cfg.Require("vaultTlsPath")
	vaultPort := cfg.RequireInt("vaultPort")
	vaultTokenOutFilePath := cfg.Require("vaultTokenOutFilePath")

	filePath := doPrivateKeyPath
	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		errr = err
	}
	fileContent := string(fileBytes)

	updateSysCmd, err := remote.NewCommand(ctx, "updateSysCmd", &remote.CommandArgs{
		Connection: &remote.ConnectionArgs{
			Host:       ipv6,
			Port:       pulumi.Float64(22),
			User:       pulumi.String(doVmUser),
			PrivateKey: pulumi.String(fileContent),
		},
		Create: pulumi.String("sudo lsof /var/lib/dpkg/lock-frontend" + "\n" +
			"return_cd=$?" + "\n" +
			"while [$return_cd -eq 0]" + "\n" +
			"do" + "\n" +
			"echo $return_cd" + "\n" +
			"done" + "\n" +
			"sudo apt update" + "\n" +
			"sudo apt upgrade -y"),
	})
	if err != nil {
		errr = err
	}

	ansibleReqs, err := local.NewCommand(ctx, "ansibleReqs", &local.CommandArgs{
		Create: pulumi.String("ansible-galaxy install -r requirements.yml --force"),
	}, pulumi.ReplaceWith([]pulumi.Resource{updateSysCmd}))
	if err != nil {
		errr = err
	}
	ctx.Export("ansibleReqs", ansibleReqs.Stdout)

	ansibleVault, err := local.NewCommand(ctx, "ansiblePlaybookCmd", &local.CommandArgs{
		Create: ipv6.ApplyT(func(doIpv6 string) (string, error) {
			return fmt.Sprintf("ANSIBLE_HOST_KEY_CHECKING=FALSE ansible-playbook -v "+
				"-u root "+
				"-i '%v,' "+
				"--private-key %v "+
				"--extra-vars 'vault_fqdn=%v vault_rem_tls_disable=%v vault_tls_loc_path=%v "+
				"vault_rem_token_out_file_path=%v vault_rem_port=%d' "+
				"vault_playbook.yml",
				doIpv6, doPrivateKeyPath, vaultFqdn, vaultTlsDisable, vaultTlsPath, vaultTokenOutFilePath, vaultPort), nil
		}).(pulumi.StringOutput),
	}, pulumi.DependsOn([]pulumi.Resource{
		ansibleReqs, updateSysCmd,
	}))
	if err != nil {
		errr = err
	}

	ansibleVaultOut := ansibleVault.Stdout.ApplyT(func(ansibleVaultStdOut string) string {
		return fmt.Sprintf("%s", ansibleVaultStdOut)
	}).(pulumi.StringOutput)
	fmt.Printf("Ansible output: %s\n", ansibleVaultOut)

	// Export Ansible STDOUT
	ctx.Export("AnsibleOut", ansibleVault.Stdout)

	return errr
}
