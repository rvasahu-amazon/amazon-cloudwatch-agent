// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package k8sattributescraper

import (
	"go.opentelemetry.io/collector/pdata/pcommon"
	semconv "go.opentelemetry.io/collector/semconv/v1.22.0"

	"github.com/aws/amazon-cloudwatch-agent/plugins/inputs/prometheus"
)

var (
	namespaceAllowlist = []string{
		semconv.AttributeK8SNamespaceName,
	}

	workloadAllowlist = []string{
		semconv.AttributeK8SDeploymentName,
		semconv.AttributeK8SDaemonSetName,
		semconv.AttributeK8SStatefulSetName,
		semconv.AttributeK8SReplicaSetName,
		semconv.AttributeK8SContainerName,
	}
	nodeAllowlist = []string{
		semconv.AttributeK8SNodeName,
	}
)

type K8sAttributeScraper struct {
	Cluster   string
	Namespace string
	Workload  string
	Node      string
}

func NewK8sAttributeScraper(clusterName string) *K8sAttributeScraper {
	return &K8sAttributeScraper{
		Cluster: clusterName,
	}
}

func (e *K8sAttributeScraper) Scrape(rm pcommon.Resource) {
	resourceAttrs := rm.Attributes()
	e.relabelPrometheus(resourceAttrs)

	e.scrapeNamespace(resourceAttrs)
	e.scrapeWorkload(resourceAttrs)
	e.scrapeNode(resourceAttrs)
}

func (e *K8sAttributeScraper) relabelPrometheus(p pcommon.Map) {
	// TODO: Retrieve workload from pod label
	if podName, exists := p.Get(prometheus.EntityK8sPodLabel); exists {
		p.PutStr(semconv.AttributeK8SPodName, podName.Str())
		p.Remove(prometheus.EntityK8sPodLabel)
	}

	if namespace, exists := p.Get(prometheus.EntityK8sNamespaceLabel); exists {
		p.PutStr(semconv.AttributeK8SNamespaceName, namespace.Str())
		p.Remove(prometheus.EntityK8sNamespaceLabel)
	}

	if nodeName, exists := p.Get(prometheus.EntityK8sNodeLabel); exists {
		p.PutStr(semconv.AttributeK8SNodeName, nodeName.Str())
		p.Remove(prometheus.EntityK8sNodeLabel)
	}
}

func (e *K8sAttributeScraper) scrapeNamespace(p pcommon.Map) {
	for _, namespace := range namespaceAllowlist {
		if namespaceAttr, ok := p.Get(namespace); ok {
			e.Namespace = namespaceAttr.Str()
			return
		}
	}
}

func (e *K8sAttributeScraper) scrapeWorkload(p pcommon.Map) {
	for _, workload := range workloadAllowlist {
		if workloadAttr, ok := p.Get(workload); ok {
			e.Workload = workloadAttr.Str()
			return
		}
	}

}

func (e *K8sAttributeScraper) scrapeNode(p pcommon.Map) {
	for _, node := range nodeAllowlist {
		if nodeAttr, ok := p.Get(node); ok {
			e.Node = nodeAttr.Str()
			return
		}
	}
}

func (e *K8sAttributeScraper) Reset() {
	*e = K8sAttributeScraper{
		Cluster: e.Cluster,
	}
}
