package main

import (
	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func initEso(ctx *pulumi.Context, provider *k8s.Provider, resource pulumi.Resource) (*helmv3.Release, error) {

	helmInfos := HelmChartInfo{
		name:            "external-secrets",
		version:         "1.1.1",
		namespace:       "external-secrets",
		createNamespace: true,
		url:             "https://charts.external-secrets.io",
		skipCrds:        false,
	}
	esoHelmRelease, err := initHelm(ctx, provider, helmInfos, resource)
	if err != nil {
		return nil, err
	}
	return esoHelmRelease, nil
}
