package main

import (
	"github.com/mitchellh/colorstring"
	k8s "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes"
	helmv3 "github.com/pulumi/pulumi-kubernetes/sdk/v4/go/kubernetes/helm/v3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	//"github.com/pulumi/pulumi/sdk/v3/go/pulumi/config"
	//"net/url"
	"strconv"
)

type HelmChartInfo struct {
	name            string
	version         string
	namespace       string
	createNamespace bool `default:"true"`
	url             string
	skipCrds        bool `default:"false"`

	values pulumi.Map
}

func initHelm(ctx *pulumi.Context, provider *k8s.Provider, chart HelmChartInfo, resource pulumi.Resource) (*helmv3.Release, error) {

	ctx.Log.Info(colorstring.Color("\t|_ Installing with [cyan]Helm chart[default]: [dark_grey]"+chart.name+":"+chart.version), nil)

	// install with Helm
	ctx.Log.Debug("Creating new Helm Release for "+chart.name, nil)
	ctx.Log.Debug(chart.namespace+", "+strconv.FormatBool(chart.createNamespace)+", "+chart.url+", "+strconv.FormatBool(chart.skipCrds)+", "+chart.version+", "+strconv.Itoa(len(chart.values)), nil)
	helmRelease, err := helmv3.NewRelease(ctx, chart.name, &helmv3.ReleaseArgs{
		Chart:           pulumi.String(chart.name),
		Version:         pulumi.String(chart.version),
		Namespace:       pulumi.String(chart.namespace),
		CreateNamespace: pulumi.Bool(chart.createNamespace),
		RepositoryOpts: &helmv3.RepositoryOptsArgs{
			Repo: pulumi.String(chart.url),
		},
		SkipCrds: pulumi.Bool(chart.skipCrds),
		Values:   chart.values,
	}, pulumi.DependsOn([]pulumi.Resource{resource}), pulumi.Provider(provider))
	if err != nil {
		return nil, err
	}
	return helmRelease, nil
}
