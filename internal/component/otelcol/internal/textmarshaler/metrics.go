// Adapted copy from the OTLP text in the Opentelemetry collector

// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package textmarshaler

import "go.opentelemetry.io/collector/pdata/pmetric"

// MarshalMetrics pmetric.Metrics to OTLP text.
func MarshalMetrics(md pmetric.Metrics) ([]byte, error) {
	buf := dataBuffer{}
	rms := md.ResourceMetrics()
	for i := 0; i < rms.Len(); i++ {
		buf.logEntry("ResourceMetrics #%d", i)
		rm := rms.At(i)
		buf.logEntry("Resource SchemaURL: %s", rm.SchemaUrl())
		buf.logAttributes("Resource attributes", rm.Resource().Attributes())
		ilms := rm.ScopeMetrics()
		for j := 0; j < ilms.Len(); j++ {
			buf.logEntry("ScopeMetrics #%d", j)
			ilm := ilms.At(j)
			buf.logEntry("ScopeMetrics SchemaURL: %s", ilm.SchemaUrl())
			buf.logInstrumentationScope(ilm.Scope())
			metrics := ilm.Metrics()
			for k := 0; k < metrics.Len(); k++ {
				buf.logEntry("Metric #%d", k)
				metric := metrics.At(k)
				buf.logMetricDescriptor(metric)
				buf.logMetricDataPoints(metric)
			}
		}
	}

	return buf.buf.Bytes(), nil
}
