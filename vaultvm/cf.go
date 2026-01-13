package main

import (
	"github.com/pulumi/pulumi-cloudflare/sdk/v6/go/cloudflare"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func cfDns(ctx *pulumi.Context, ipv4 pulumi.StringOutput, ipv6 pulumi.StringOutput) error {
	var errr error = nil

	cfg := config.New(ctx, "")
	cfDnsName := cfg.Require("cfDnsName")
	cfZoneId := cfg.Require("cfZoneId")

	cfAAAARecord, err := cloudflare.NewDnsRecord(ctx, "dvvm_dns_AAAA_record", &cloudflare.DnsRecordArgs{
		ZoneId:  pulumi.String(cfZoneId),
		Name:    pulumi.String(cfDnsName),
		Ttl:     pulumi.Float64(60),
		Type:    pulumi.String("AAAA"),
		Content: pulumi.StringOutput(ipv6),
		Proxied: pulumi.Bool(false),
	})
	if err != nil {
		errr = err
	}

	cfARecord, err := cloudflare.NewDnsRecord(ctx, "dvvm_dns_A_record", &cloudflare.DnsRecordArgs{
		ZoneId:  pulumi.String(cfZoneId),
		Name:    pulumi.String(cfDnsName),
		Ttl:     pulumi.Float64(60),
		Type:    pulumi.String("A"),
		Content: pulumi.StringOutput(ipv4),
		Proxied: pulumi.Bool(false),
	})
	if err != nil {
		errr = err
	}

	ctx.Export("AAAA_dns_name", cfAAAARecord.Name)
	ctx.Export("AAAA_dns_ID", cfAAAARecord.ID())
	ctx.Export("A_dns_name", cfARecord.Name)
	ctx.Export("A_dns_ID", cfARecord.ID())

	return errr
}
