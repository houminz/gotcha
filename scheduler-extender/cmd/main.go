/*
 * Copyright 2019 THL A29 Limited, a Tencent company.
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
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/SimpCosm/godemo/scheduler-extender/pkg/routes"
	"github.com/SimpCosm/godemo/scheduler-extender/pkg/scheduler"
	"github.com/comail/colog"
	"github.com/julienschmidt/httprouter"

	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	extender "k8s.io/kube-scheduler/extender/v1"
)



var (
	TruePredicate = scheduler.Predicate{
		Name: "always_true",
		Func: func(pod v1.Pod, node v1.Node) (bool, error) {
			return true, nil
		},
	}

	ZeroPriority = scheduler.Prioritize{
		Name: "zero_score",
		Func: func(_ v1.Pod, nodes []v1.Node) (*extender.HostPriorityList, error) {
			var priorityList extender.HostPriorityList
			priorityList = make([]extender.HostPriority, len(nodes))
			for i, node := range nodes {
				priorityList[i] = extender.HostPriority{
					Host:  node.Name,
					Score: 0,
				}
			}
			return &priorityList, nil
		},
	}

	NoBind = scheduler.Bind{
		Func: func(podName string, podNamespace string, podUID types.UID, node string) error {
			return fmt.Errorf("This extender doesn't support Bind.  Please make 'BindVerb' be empty in your ExtenderConfig.")
		},
	}

	EchoPreemption = scheduler.Preemption{
		Func: func(
			_ v1.Pod,
			_ map[string]*extender.Victims,
			nodeNameToMetaVictims map[string]*extender.MetaVictims,
		) map[string]*extender.MetaVictims {
			return nodeNameToMetaVictims
		},
	}
)

func StringToLevel(levelStr string) colog.Level {
	switch level := strings.ToUpper(levelStr); level {
	case "TRACE":
		return colog.LTrace
	case "DEBUG":
		return colog.LDebug
	case "INFO":
		return colog.LInfo
	case "WARNING":
		return colog.LWarning
	case "ERROR":
		return colog.LError
	case "ALERT":
		return colog.LAlert
	default:
		log.Printf("warning: LOG_LEVEL=\"%s\" is empty or invalid, fallling back to \"INFO\".\n", level)
		return colog.LInfo
	}
}

func main() {
	colog.SetDefaultLevel(colog.LInfo)
	colog.SetMinLevel(colog.LInfo)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()
	level := StringToLevel(os.Getenv("LOG_LEVEL"))
	log.Print("Log level was set to ", strings.ToUpper(level.String()))
	colog.SetMinLevel(level)

	router := httprouter.New()
	routes.AddVersion(router)

	predicates := []scheduler.Predicate{TruePredicate}
	for _, p := range predicates {
		routes.AddPredicate(router, p)
	}

	priorities := []scheduler.Prioritize{ZeroPriority}
	for _, p := range priorities {
		routes.AddPrioritize(router, p)
	}

	routes.AddBind(router, NoBind)

	log.Print("info: server starting on the port :80")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}