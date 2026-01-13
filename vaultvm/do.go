package main

import (
	"fmt"
	"strconv"

	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func doVM(ctx *pulumi.Context) (pulumi.StringOutput, pulumi.StringOutput, error) {
	var errr error = nil

	cfg := config.New(ctx, "")
	doPublicKeyHash := cfg.Require("doPublicKeyHash")
	doVmImage := cfg.Require("doVmImage")
	doVmRegion := cfg.Require("doVmRegion")
	doVmSize := cfg.Require("doVmSize")
	doFirewallDisable := cfg.RequireBool("doFirewallDisable")
	var doFwWebIpRangesArr []string
	cfg.RequireObject("doFwWebIpRanges", &doFwWebIpRangesArr)
	SourceAddressInPulumiStrArrFwWeb := pulumi.ToStringArray(doFwWebIpRangesArr).ToStringArrayOutput()

	// Create a new Droplet
	droplet, err := digitalocean.NewDroplet(ctx, "dvvm", &digitalocean.DropletArgs{
		Image:  pulumi.String(doVmImage),
		Name:   pulumi.String("dvvm"),
		Region: pulumi.String(doVmRegion),
		Size:   pulumi.String(doVmSize),
		Ipv6:   pulumi.Bool(true),
		SshKeys: pulumi.StringArray{
			pulumi.String(doPublicKeyHash),
		},
	})
	if err != nil {
		fmt.Printf("Error creating digitalocean droplet: %v\n", err)
		errr = err
	}

	if !doFirewallDisable {

		fw, err := digitalocean.NewFirewall(ctx, "dvvm_fw", &digitalocean.FirewallArgs{
			Name: pulumi.String("dvvm-fw"),
			DropletIds: pulumi.IntArray{
				droplet.ID().ApplyT(strconv.Atoi).(pulumi.IntOutput),
			},
			InboundRules: digitalocean.FirewallInboundRuleArray{
				&digitalocean.FirewallInboundRuleArgs{
					Protocol:        pulumi.String("tcp"),
					PortRange:       pulumi.String("0"),
					SourceAddresses: SourceAddressInPulumiStrArrFwWeb,
				},
				&digitalocean.FirewallInboundRuleArgs{
					Protocol:        pulumi.String("udp"),
					PortRange:       pulumi.String("0"),
					SourceAddresses: SourceAddressInPulumiStrArrFwWeb,
				},
				&digitalocean.FirewallInboundRuleArgs{
					Protocol:        pulumi.String("icmp"),
					SourceAddresses: SourceAddressInPulumiStrArrFwWeb,
				},
			},
			OutboundRules: digitalocean.FirewallOutboundRuleArray{
				&digitalocean.FirewallOutboundRuleArgs{
					Protocol:  pulumi.String("tcp"),
					PortRange: pulumi.String("0"),
					DestinationAddresses: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
						pulumi.String("::/0"),
					},
				},
				&digitalocean.FirewallOutboundRuleArgs{
					Protocol:  pulumi.String("udp"),
					PortRange: pulumi.String("0"),
					DestinationAddresses: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
						pulumi.String("::/0"),
					},
				},
				&digitalocean.FirewallOutboundRuleArgs{
					Protocol: pulumi.String("icmp"),
					DestinationAddresses: pulumi.StringArray{
						pulumi.String("0.0.0.0/0"),
						pulumi.String("::/0"),
					},
				},
			},
		})

		if err != nil {
			fmt.Printf("Error creating digitalocean firewall: %v\n", err)
			errr = err
		}

		ctx.Export("FW ID", fw.ID())
	}

	ipv6 := droplet.Ipv6Address.ApplyT(func(doIpv6 string) string {
		return doIpv6
	}).(pulumi.StringOutput)

	ipv4 := droplet.Ipv4Address.ApplyT(func(doIpv4 string) string {
		return doIpv4
	}).(pulumi.StringOutput)

	// Export the droplets Ipv6
	ctx.Export("ipv6", droplet.Ipv6Address)
	ctx.Export("ipv4", droplet.Ipv4Address)
	ctx.Export("Droplet ID", droplet.ID())

	return ipv4, ipv6, errr
}
