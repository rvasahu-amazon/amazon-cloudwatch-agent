// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package resourcestore

import (
	"reflect"
	"testing"
)

func TestInitResourceStore(t *testing.T) {
	tests := []struct {
		name string
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initResourceStore()
		})
	}
}

func TestResourceStore_EC2Info(t *testing.T) {
	tests := []struct {
		name         string
		ec2InfoInput ec2Info
		want         ec2Info
	}{
		{
			name: "happypath",
			ec2InfoInput: ec2Info{
				InstanceID:       "i-1234567890",
				AutoScalingGroup: "test-asg",
			},
			want: ec2Info{
				InstanceID:       "i-1234567890",
				AutoScalingGroup: "test-asg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceStore{
				ec2Info: tt.ec2InfoInput,
			}
			if got := r.EC2Info(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EC2Info() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceStore_EKSInfo(t *testing.T) {
	tests := []struct {
		name         string
		eksInfoInput eksInfo
		want         eksInfo
	}{
		{
			name:         "happypath",
			eksInfoInput: eksInfo{ClusterName: "test-cluster"},
			want:         eksInfo{ClusterName: "test-cluster"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceStore{
				eksInfo: tt.eksInfoInput,
			}
			if got := r.EKSInfo(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EKSInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceStore_LogFiles(t *testing.T) {
	tests := []struct {
		name         string
		logFileInput map[string]string
		want         map[string]string
	}{
		{
			name:         "happypath",
			logFileInput: map[string]string{"/opt/aws/amazon-cloudwatch-agent/logs/amazon-cloudwatch-agent.log": "cloudwatch-agent"},
			want:         map[string]string{"/opt/aws/amazon-cloudwatch-agent/logs/amazon-cloudwatch-agent.log": "cloudwatch-agent"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceStore{
				logFiles: tt.logFileInput,
			}
			if got := r.LogFiles(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LogFiles() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResourceStore_Mode(t *testing.T) {
	tests := []struct {
		name      string
		modeInput string
		want      string
	}{
		{name: "happypath", modeInput: "EC2", want: "EC2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ResourceStore{
				mode: tt.modeInput,
			}
			if got := r.Mode(); got != tt.want {
				t.Errorf("Mode() = %v, want %v", got, tt.want)
			}
		})
	}
}
