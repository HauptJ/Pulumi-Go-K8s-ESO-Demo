package main

import (
	"fmt"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		ipv4, ipv6, err := doVM(ctx)
		if err != nil {
			fmt.Printf("Error in doVaultVM: %s\n", err)
		}

		err = cfDns(ctx, ipv4, ipv6)
		if err != nil {
			fmt.Printf("Error in cfDns: %s\n", err)
		}

		err = provVaultVM(ctx, ipv6)
		if err != nil {
			fmt.Printf("Error in provVaultVM: %s\n", err)
		}

		return nil
	})
}
