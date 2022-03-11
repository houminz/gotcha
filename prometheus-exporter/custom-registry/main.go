/*
 * Copyright 2022 The Synap Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// 创建一个自定义的注册表
	registry := prometheus.NewRegistry()
	// 可选: 添加 process 和 Go 运行时指标到我们自定义的注册表中
	registry.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	registry.MustRegister(prometheus.NewGoCollector())

	// 创建一个简单 gauge 指标
	temp := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "cosmos_home_temperature_celsius",
		Help: "The current temperature in degrees Celsius.",
	})

	// 使用我们自定义的注册表注册 gauge
	registry.MustRegister(temp)

	// 设置 gague 的值为 39
	temp.Set(39)

	// 暴露自定义指标
	http.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{Registry: registry}))
	http.ListenAndServe(":8080", nil)
}