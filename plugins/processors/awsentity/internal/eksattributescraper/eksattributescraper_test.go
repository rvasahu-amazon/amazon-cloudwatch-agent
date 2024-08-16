// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package eksattributescraper

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
	semconv "go.opentelemetry.io/collector/semconv/v1.22.0"

	"github.com/aws/amazon-cloudwatch-agent/plugins/processors/awsentity/internal/entityattributes"
)

func TestNewEKSAttributeScraper(t *testing.T) {
	scraper := NewEKSAttributeScraper("test")
	assert.Equal(t, "test", scraper.Cluster)
}

func Test_eksattributescraper_Scrape(t *testing.T) {

	tests := []struct {
		name        string
		clusterName string
		args        pcommon.Resource
		want        pcommon.Map
	}{
		{
			name:        "Empty",
			clusterName: "",
			args:        pcommon.NewResource(),
			want:        pcommon.NewMap(),
		},
		{
			name:        "ClusterOnly",
			clusterName: "test-cluster",
			args:        pcommon.NewResource(),
			want: getAttributeMap(map[string]any{
				entityattributes.AttributeEntityCluster: "test-cluster",
			}),
		},
		{
			name:        "AllAppSignalAttributes",
			clusterName: "test-cluster",
			args:        generateResourceMetrics(semconv.AttributeK8SNamespaceName, "test-namespace", semconv.AttributeK8SDeploymentName, "test-workload", semconv.AttributeK8SNodeName, "test-node"),
			want: getAttributeMap(map[string]any{
				semconv.AttributeK8SNamespaceName:         "test-namespace",
				semconv.AttributeK8SDeploymentName:        "test-workload",
				semconv.AttributeK8SNodeName:              "test-node",
				entityattributes.AttributeEntityCluster:   "test-cluster",
				entityattributes.AttributeEntityNamespace: "test-namespace",
				entityattributes.AttributeEntityWorkload:  "test-workload",
				entityattributes.AttributeEntityNode:      "test-node",
			}),
		},
		{
			name:        "AllContainerInsightsAttributes",
			clusterName: "test-cluster",
			args:        generateResourceMetrics(entityattributes.Namespace, "test-namespace", entityattributes.PodName, "test-workload", entityattributes.NodeName, "test-node"),
			want: getAttributeMap(map[string]any{
				entityattributes.Namespace:                "test-namespace",
				entityattributes.PodName:                  "test-workload",
				entityattributes.NodeName:                 "test-node",
				entityattributes.AttributeEntityCluster:   "test-cluster",
				entityattributes.AttributeEntityNamespace: "test-namespace",
				entityattributes.AttributeEntityWorkload:  "test-workload",
				entityattributes.AttributeEntityNode:      "test-node",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := NewEKSAttributeScraper(tt.clusterName)
			e.Scrape(tt.args)
			assert.Equal(t, tt.want.AsRaw(), tt.args.Attributes().AsRaw())
		})
	}
}

func Test_eksattributescraper_decorateEntityAttributes(t *testing.T) {
	type fields struct {
		Cluster   string
		Namespace string
		Workload  string
		Node      string
	}
	tests := []struct {
		name   string
		fields fields
		want   pcommon.Map
	}{
		{
			name:   "Empty",
			fields: fields{},
			want:   pcommon.NewMap(),
		},
		{
			name: "OneAttribute",
			fields: fields{
				Cluster: "test-cluster",
			},
			want: getAttributeMap(map[string]any{
				entityattributes.AttributeEntityCluster: "test-cluster",
			}),
		},
		{
			name: "AllAttributes",
			fields: fields{
				Cluster:   "test-cluster",
				Namespace: "test-namespace",
				Workload:  "test-workload",
				Node:      "test-node",
			},
			want: getAttributeMap(map[string]any{
				entityattributes.AttributeEntityCluster:   "test-cluster",
				entityattributes.AttributeEntityNamespace: "test-namespace",
				entityattributes.AttributeEntityWorkload:  "test-workload",
				entityattributes.AttributeEntityNode:      "test-node",
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := pcommon.NewMap()
			e := &eksattributescraper{
				Cluster:   tt.fields.Cluster,
				Namespace: tt.fields.Namespace,
				Workload:  tt.fields.Workload,
				Node:      tt.fields.Node,
			}
			e.decorateEntityAttributes(p)
			assert.Equal(t, tt.want.AsRaw(), p.AsRaw())
		})
	}
}

func Test_eksattributescraper_reset(t *testing.T) {
	type fields struct {
		Cluster   string
		Namespace string
		Workload  string
		Node      string
	}
	tests := []struct {
		name   string
		fields fields
		want   *eksattributescraper
	}{
		{
			name:   "Empty",
			fields: fields{},
			want:   &eksattributescraper{},
		},
		{
			name: "ClusterExists",
			fields: fields{
				Cluster: "test-cluster",
			},
			want: &eksattributescraper{
				Cluster: "test-cluster",
			},
		},
		{
			name: "MultipleAttributeExists",
			fields: fields{
				Cluster:   "test-cluster",
				Namespace: "test-namespace",
				Workload:  "test-workload",
				Node:      "test-node",
			},
			want: &eksattributescraper{
				Cluster: "test-cluster",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eksattributescraper{
				Cluster:   tt.fields.Cluster,
				Namespace: tt.fields.Namespace,
				Workload:  tt.fields.Workload,
				Node:      tt.fields.Node,
			}
			e.reset()
			assert.Equal(t, tt.want, e)
		})
	}
}

func Test_eksattributescraper_scrapeNamespace(t *testing.T) {
	tests := []struct {
		name string
		args pcommon.Map
		want string
	}{
		{
			name: "Empty",
			args: getAttributeMap(map[string]any{"": ""}),
			want: "",
		},
		{
			name: "AppSignalNodeExists",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SNamespaceName: "namespace-name"}),
			want: "namespace-name",
		},
		{
			name: "ContainerInsightsNodeExists",
			args: getAttributeMap(map[string]any{entityattributes.Namespace: "namespace-name"}),
			want: "namespace-name",
		},
		{
			name: "NonmatchingNamespace",
			args: getAttributeMap(map[string]any{"namespace": "namespace-name"}),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eksattributescraper{}
			e.scrapeNamespace(tt.args)
			assert.Equal(t, tt.want, e.Namespace)
		})
	}
}

func Test_eksattributescraper_scrapeNode(t *testing.T) {
	tests := []struct {
		name string
		args pcommon.Map
		want string
	}{
		{
			name: "Empty",
			args: getAttributeMap(map[string]any{"": ""}),
			want: "",
		},
		{
			name: "AppsignalNodeExists",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SNodeName: "node-name"}),
			want: "node-name",
		},
		{
			name: "ContainerInsightNodeExists",
			args: getAttributeMap(map[string]any{entityattributes.NodeName: "node-name"}),
			want: "node-name",
		},
		{
			name: "NonmatchingNode",
			args: getAttributeMap(map[string]any{"node": "node-name"}),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eksattributescraper{}
			e.scrapeNode(tt.args)
			assert.Equal(t, tt.want, e.Node)
		})
	}
}

func Test_eksattributescraper_scrapeWorkload(t *testing.T) {
	tests := []struct {
		name string
		args pcommon.Map
		want string
	}{
		{
			name: "Empty",
			args: getAttributeMap(map[string]any{"": ""}),
			want: "",
		},
		{
			name: "DeploymentWorkload",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SDeploymentName: "test-deployment"}),
			want: "test-deployment",
		},
		{
			name: "DaemonsetWorkload",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SDaemonSetName: "test-daemonset"}),
			want: "test-daemonset",
		},
		{
			name: "StatefulSetWorkload",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SStatefulSetName: "test-statefulset"}),
			want: "test-statefulset",
		},
		{
			name: "ReplicaSetWorkload",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SReplicaSetName: "test-replicaset"}),
			want: "test-replicaset",
		},
		{
			name: "ContainerWorkload",
			args: getAttributeMap(map[string]any{semconv.AttributeK8SContainerName: "test-container"}),
			want: "test-container",
		},
		{
			name: "ContainerInsightPodNameWorkload",
			args: getAttributeMap(map[string]any{entityattributes.PodName: "test-workload"}),
			want: "test-workload",
		},
		{
			name: "MultipleWorkloads",
			args: getAttributeMap(map[string]any{
				semconv.AttributeK8SDeploymentName: "test-deployment",
				semconv.AttributeK8SContainerName:  "test-container"}),
			want: "test-deployment",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &eksattributescraper{}
			e.scrapeWorkload(tt.args)
			assert.Equal(t, tt.want, e.Workload)
		})
	}
}

func getAttributeMap(attributes map[string]any) pcommon.Map {
	attrMap := pcommon.NewMap()
	attrMap.FromRaw(attributes)
	return attrMap
}

func generateResourceMetrics(resourceAttrs ...string) pcommon.Resource {
	md := pmetric.NewMetrics()
	generateResource(md, resourceAttrs...)
	return md.ResourceMetrics().At(0).Resource()
}

func generateResource(md pmetric.Metrics, resourceAttrs ...string) {
	attrs := md.ResourceMetrics().AppendEmpty().Resource().Attributes()
	for i := 0; i < len(resourceAttrs); i += 2 {
		attrs.PutStr(resourceAttrs[i], resourceAttrs[i+1])
	}
}
