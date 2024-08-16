// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package entityattributes

const (
	AttributeEntityServiceName           = "aws.entity.service.name"
	AttributeEntityDeploymentEnvironment = "aws.entity.deployment.environment"
	AttributeEntityCluster               = "aws.entity.k8s.cluster.name"
	AttributeEntityNamespace             = "aws.entity.k8s.namespace.name"
	AttributeEntityWorkload              = "aws.entity.k8s.workload.name"
	AttributeEntityNode                  = "aws.entity.k8s.node.name"
)

// Container Insights attributes used for scraping EKS related information
const (
	NodeName  = "NodeName"
	Namespace = "Namespace"
	// PodName in Container Insights is the workload(Deployment, Daemonset, etc) name
	PodName = "PodName"
)
