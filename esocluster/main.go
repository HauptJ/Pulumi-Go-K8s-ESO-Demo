package main

import (
	//"strconv"

	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"

	//k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	//"github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	//metav1 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/meta/v1"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func main() {

	pulumi.Run(func(ctx *pulumi.Context) error {
		cfg := config.New(ctx, "")
		doK8sClusterName := cfg.Require("doK8sClusterName")
		doK8sClusterLbEnabled := cfg.RequireBool("doK8sClusterLbEnabled")
		nginxHelmEnabled := cfg.RequireBool("nginxHelmEnabled")
		esoHelmEnabled := cfg.RequireBool("esoHelmEnabled")
		esoDnsEnabled := cfg.RequireBool("esoDnsEnabled")

		k8sCluster, err := doK8sCluster(ctx)
		if err != nil {
			return err
		}

		doKubeConfig := k8sCluster.ToKubernetesClusterOutput().KubeConfigs().Index(pulumi.Int(0)).RawConfig()
		ctx.Export("doKubeConfig", doKubeConfig)

		k8sProvider, err := k8s.NewProvider(ctx, doK8sClusterName, &k8s.ProviderArgs{
			Kubeconfig: doKubeConfig,
		}, pulumi.DependsOn([]pulumi.Resource{k8sCluster}))
		if err != nil {
			return err
		}

		if doK8sClusterLbEnabled {
			k8sLb, err := doK8sClusterLb(ctx)
			if err != nil {
				return err
			}
			ctx.Export("doK8sLbID", k8sLb.ID())
			ctx.Export("doK8sLbIpv4", k8sLb.Ip)
			ctx.Export("doK8sLbIpv6", k8sLb.Ipv6)

			if nginxHelmEnabled {
				nginx, err := initNginx(ctx, k8sProvider, k8sLb, k8sLb.ID().ToStringOutput())
				if err != nil {
					return err
				}
				ctx.Export("Nginx Helm Release ID", nginx.ID())
			}

			if esoDnsEnabled {
				err := cfDns(ctx, k8sLb.Ip, k8sLb.Ipv6)
				if err != nil {
					return err
				}
			}
		}

		if esoHelmEnabled {
			eso, err := initEso(ctx, k8sProvider, k8sCluster)
			if err != nil {
				return err
			}
			ctx.Export("Eso Helm Release ID", eso.ID())
		}

		return nil
	})
}
