// Copyright Amazon.com, Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: MIT

package prometheus

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/confmap"

	"github.com/aws/amazon-cloudwatch-agent/internal/util/collections"
	translatorConfig "github.com/aws/amazon-cloudwatch-agent/translator/config"
	"github.com/aws/amazon-cloudwatch-agent/translator/context"
	"github.com/aws/amazon-cloudwatch-agent/translator/translate/otel/common"
)

func TestTranslator(t *testing.T) {
	type want struct {
		receivers  []string
		processors []string
		exporters  []string
		extensions []string
	}
	cit := NewTranslator()
	require.EqualValues(t, "metrics/prometheus", cit.ID().String())
	testCases := map[string]struct {
		input          map[string]interface{}
		kubernetesMode string
		want           *want
		wantErr        error
	}{
		"WithoutPrometheusKey": {
			input:   map[string]interface{}{},
			wantErr: &common.MissingKeyError{ID: cit.ID(), JsonKey: "logs::metrics_collected::prometheus"},
		},
		"WithPrometheusKeyAndOnK8s": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"prometheus": nil,
					},
				},
			},
			kubernetesMode: translatorConfig.ModeEKS,
			want: &want{
				receivers:  []string{"telegraf_prometheus"},
				processors: []string{"batch/prometheus", "awsentity/service"},
				exporters:  []string{"awsemf/prometheus"},
				extensions: []string{"agenthealth/logs"},
			},
		},
		"WithPrometheusKeyAndNotOnK8s": {
			input: map[string]interface{}{
				"logs": map[string]interface{}{
					"metrics_collected": map[string]interface{}{
						"prometheus": nil,
					},
				},
			},
			kubernetesMode: "",
			want: &want{
				receivers:  []string{"telegraf_prometheus"},
				processors: []string{"batch/prometheus", "resourcedetection"},
				exporters:  []string{"awsemf/prometheus"},
				extensions: []string{"agenthealth/logs"},
			},
		},
	}
	for name, testCase := range testCases {
		t.Run(name, func(t *testing.T) {
			conf := confmap.NewFromStringMap(testCase.input)
			context.CurrentContext().SetKubernetesMode(testCase.kubernetesMode)
			got, err := cit.Translate(conf)
			assert.Equal(t, testCase.wantErr, err)
			if testCase.want == nil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, testCase.want.receivers, collections.MapSlice(got.Receivers.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.processors, collections.MapSlice(got.Processors.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.exporters, collections.MapSlice(got.Exporters.Keys(), component.ID.String))
				assert.Equal(t, testCase.want.extensions, collections.MapSlice(got.Extensions.Keys(), component.ID.String))
			}
		})
	}
}
