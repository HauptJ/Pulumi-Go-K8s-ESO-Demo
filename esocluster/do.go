package main

import (
	"github.com/pulumi/pulumi-digitalocean/sdk/v4/go/digitalocean"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func doK8sCluster(ctx *pulumi.Context) (*digitalocean.KubernetesCluster, error) {
	//var errr error = nil
	cfg := config.New(ctx, "")
	doK8sClusterName := cfg.Require("doK8sClusterName")
	doK8sNodePoolName := cfg.Require("doK8sNodePoolName")
	doK8sRegion := cfg.Require("doK8sRegion")
	doK8sVersion := cfg.Require("doK8sVersion")
	doK8sNodeSize := cfg.Require("doK8sNodeSize")
	doK8sAutoScale := cfg.RequireBool("doK8sAutoScale")
	doK8sMinNodes := cfg.RequireInt("doK8sMinNodes")
	doK8sMaxNodes := cfg.RequireInt("doK8sMaxNodes")

	k8s, err := digitalocean.NewKubernetesCluster(ctx, doK8sClusterName, &digitalocean.KubernetesClusterArgs{
		Name:    pulumi.String(doK8sClusterName),
		Region:  pulumi.String(doK8sRegion),
		Version: pulumi.String(doK8sVersion),
		NodePool: &digitalocean.KubernetesClusterNodePoolArgs{
			Name:      pulumi.String(doK8sNodePoolName),
			AutoScale: pulumi.BoolPtr(doK8sAutoScale),
			NodeCount: pulumi.Int(doK8sMinNodes),
			MinNodes:  pulumi.Int(doK8sMinNodes),
			MaxNodes:  pulumi.Int(doK8sMaxNodes),
			Size:      pulumi.String(doK8sNodeSize),
			Tags: pulumi.StringArray{
				pulumi.String(doK8sClusterName),
				pulumi.String(doK8sNodePoolName),
			},
		},
	}, pulumi.Timeouts(&pulumi.CustomTimeouts{Create: "15m"}))
	if err != nil {
		//errr = err
		return nil, err
	}

	//clusterOutput := k8s.ToKubernetesClusterOutput()
	//k8s.NodePool.Nodes().Index(pulumi.Int(0)).DropletId()
	//ctx.Export("kubeconfigs", myKube.KubeConfigs)
	//ctx.Export("doKubeConfig", clusterOutput.KubeConfigs().Index(pulumi.Int(0)).RawConfig())

	return k8s, nil
}

func doK8sClusterLb(ctx *pulumi.Context) (*digitalocean.LoadBalancer, error) {
	var errr error = nil
	cfg := config.New(ctx, "")
	//doK8sClusterName := cfg.Require("doK8sClusterName")
	doK8sNodePoolName := cfg.Require("doK8sNodePoolName")
	doK8sRegion := cfg.Require("doK8sRegion")
	doK8sLbSize := cfg.Require("doK8sLbSize")

	k8sLbName := doK8sNodePoolName + "-lb"
	k8sLb, err := digitalocean.NewLoadBalancer(ctx, k8sLbName, &digitalocean.LoadBalancerArgs{
		Name:   pulumi.String(k8sLbName),
		Region: pulumi.String(doK8sRegion),
		Size:   pulumi.String(doK8sLbSize),
		ForwardingRules: digitalocean.LoadBalancerForwardingRuleArray{
			&digitalocean.LoadBalancerForwardingRuleArgs{
				EntryPort:      pulumi.Int(80),
				EntryProtocol:  pulumi.String("http"),
				TargetPort:     pulumi.Int(80),
				TargetProtocol: pulumi.String("http"),
			},
		},
		Healthcheck: &digitalocean.LoadBalancerHealthcheckArgs{
			Port:     pulumi.Int(22),
			Protocol: pulumi.String("tcp"),
		},
		DropletTag: pulumi.String(doK8sNodePoolName),
	})
	if err != nil {
		errr = err
	}

	ctx.Export("lb_ID", k8sLb.ID())
	ctx.Export("lb_ipv4", k8sLb.Ip)
	ctx.Export("lb_ipv6", k8sLb.Ipv6)

	return k8sLb, errr

}
