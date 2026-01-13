package main

import (
	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
)

func initNginx(ctx *pulumi.Context, provider *k8s.Provider, resource pulumi.Resource, lbId pulumi.StringOutput) (*helmv3.Release, error) {
	cfg := config.New(ctx, "")
	nginxHelmVer := cfg.Require("nginxHelmVer")

	/*helmValues := pulumi.Map{
		"controller": pulumi.Map{
			"publishService": pulumi.Map{
				"enabled": pulumi.Bool(true),
			},
		},
	}*/

	helmValues := pulumi.Map{
		"image": pulumi.Map{
			"repository": pulumi.String("bitnamilegacy/nginx-ingress-controller"),
		},
		"defaultBackend": pulumi.Map{
			"image": pulumi.Map{
				"repository": pulumi.String("bitnamilegacy/nginx"),
			},
		},
		"service": pulumi.Map{
			"type": pulumi.String("LoadBalancer"),
			"annotations": pulumi.Map{
				"kubernetes.digitalocean.com/load-balancer-id":  lbId,
				"nginx.ingress.kubernetes.io/enable-access-log": pulumi.String("true"),
				"nginx.ingress.kubernetes.io/enable-error-log":  pulumi.String("true"),
			},
		},
	}

	helmInfos := HelmChartInfo{
		name:            "nginx-ingress-controller",
		version:         nginxHelmVer,
		namespace:       "nginx-ingress-controller",
		createNamespace: true,
		//url:		"https://kubernetes.github.io/ingress-nginx",
		url:      "https://charts.bitnami.com/bitnami",
		skipCrds: false,
		values:   helmValues,
	}
	nginxHelmRelease, err := initHelm(ctx, provider, helmInfos, resource)
	if err != nil {
		return nil, err
	}
	return nginxHelmRelease, nil
}
