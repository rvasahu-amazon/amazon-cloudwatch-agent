// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package eksattributescraper

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/collector/semconv/v1.22.0"

	"github.com/aws/amazon-cloudwatch-agent/plugins/processors/awsentity/internal/entityattributes"
)

var (
	namespaceAllowlist = []string{
		semconv.AttributeK8SNamespaceName,
		entityattributes.Namespace,
	}

	workloadAllowlist = []string{
		semconv.AttributeK8SDeploymentName,
		semconv.AttributeK8SDaemonSetName,
		semconv.AttributeK8SStatefulSetName,
		semconv.AttributeK8SReplicaSetName,
		semconv.AttributeK8SContainerName,
		entityattributes.PodName,
	}
	nodeAllowlist = []string{
		semconv.AttributeK8SNodeName,
		entityattributes.NodeName,
	}
)

type eksattributescraper struct {
	Cluster   string
	Namespace string
	Workload  string
	Node      string
}

func NewEKSAttributeScraper(clusterName string) *eksattributescraper {
	return &eksattributescraper{
		Cluster: clusterName,
	}
}

func (e *eksattributescraper) Scrape(rm pcommon.Resource) {
	resourceAttrs := rm.Attributes()
	e.scrapeNamespace(resourceAttrs)
	e.scrapeWorkload(resourceAttrs)
	e.scrapeNode(resourceAttrs)
	e.decorateEntityAttributes(resourceAttrs)
	e.reset()
}

func (e *eksattributescraper) scrapeNamespace(p pcommon.Map) {
	for _, namespace := range namespaceAllowlist {
		if namespaceAttr, ok := p.Get(namespace); ok {
			e.Namespace = namespaceAttr.Str()
			return
		}
	}
}

func (e *eksattributescraper) scrapeWorkload(p pcommon.Map) {
	for _, workload := range workloadAllowlist {
		if workloadAttr, ok := p.Get(workload); ok {
			e.Workload = workloadAttr.Str()
			return
		}
	}

}

func (e *eksattributescraper) scrapeNode(p pcommon.Map) {
	for _, node := range nodeAllowlist {
		if nodeAttr, ok := p.Get(node); ok {
			e.Node = nodeAttr.Str()
			return
		}
	}
}

func (e *eksattributescraper) decorateEntityAttributes(p pcommon.Map) {
	addAttributeIfNonEmpty(p, entityattributes.AttributeEntityCluster, e.Cluster)
	addAttributeIfNonEmpty(p, entityattributes.AttributeEntityNamespace, e.Namespace)
	addAttributeIfNonEmpty(p, entityattributes.AttributeEntityWorkload, e.Workload)
	addAttributeIfNonEmpty(p, entityattributes.AttributeEntityNode, e.Node)
}

func (e *eksattributescraper) reset() {
	*e = eksattributescraper{
		Cluster: e.Cluster,
	}
}

func addAttributeIfNonEmpty(p pcommon.Map, key string, value string) {
	if value != "" {
		p.PutStr(key, value)
	}
}
