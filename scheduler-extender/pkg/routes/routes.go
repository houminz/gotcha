package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/SimpCosm/godemo/scheduler-extender/pkg/scheduler"

	extender "k8s.io/kube-scheduler/extender/v1"
)

const (
	versionPath      = "/version"
	apiPrefix        = "/scheduler"
	bindPath         = apiPrefix + "/bind"
	preemptionPath   = apiPrefix + "/preemption"
	predicatesPrefix = apiPrefix + "/predicates"
	prioritiesPrefix = apiPrefix + "/priorities"
)

var version string // injected via ldflags at config time

func checkBody(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
}

func PredicateRoute(predicate scheduler.Predicate) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: ", predicate.Name, " ExtenderArgs = ", buf.String())

		var extenderArgs extender.ExtenderArgs
		var extenderFilterResult *extender.ExtenderFilterResult

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			extenderFilterResult = &extender.ExtenderFilterResult{
				Nodes:       nil,
				FailedNodes: nil,
				Error:       err.Error(),
			}
		} else {
			extenderFilterResult = predicate.Handler(extenderArgs)
		}

		if resultBody, err := json.Marshal(extenderFilterResult); err != nil {
			panic(err)
		} else {
			log.Print("info: ", predicate.Name, " extenderFilterResult = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

func PrioritizeRoute(prioritize scheduler.Prioritize) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: ", prioritize.Name, " ExtenderArgs = ", buf.String())

		var extenderArgs extender.ExtenderArgs
		var hostPriorityList *extender.HostPriorityList

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			panic(err)
		}

		if list, err := prioritize.Handler(extenderArgs); err != nil {
			panic(err)
		} else {
			hostPriorityList = list
		}

		if resultBody, err := json.Marshal(hostPriorityList); err != nil {
			panic(err)
		} else {
			log.Print("info: ", prioritize.Name, " hostPriorityList = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

func BindRoute(bind scheduler.Bind) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: extenderBindingArgs = ", buf.String())

		var extenderBindingArgs extender.ExtenderBindingArgs
		var extenderBindingResult *extender.ExtenderBindingResult

		if err := json.NewDecoder(body).Decode(&extenderBindingArgs); err != nil {
			extenderBindingResult = &extender.ExtenderBindingResult{
				Error: err.Error(),
			}
		} else {
			extenderBindingResult = bind.Handler(extenderBindingArgs)
		}

		if resultBody, err := json.Marshal(extenderBindingResult); err != nil {
			panic(err)
		} else {
			log.Print("info: extenderBindingResult = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

func PreemptionRoute(preemption scheduler.Preemption) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: extenderPreemptionArgs = ", buf.String())

		var extenderPreemptionArgs extender.ExtenderPreemptionArgs
		var extenderPreemptionResult *extender.ExtenderPreemptionResult

		if err := json.NewDecoder(body).Decode(&extenderPreemptionArgs); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
		} else {
			extenderPreemptionResult = preemption.Handler(extenderPreemptionArgs)
		}

		if resultBody, err := json.Marshal(extenderPreemptionResult); err != nil {
			panic(err)
		} else {
			log.Print("info: extenderPreemptionResult = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

func VersionRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, fmt.Sprint(version))
}

func AddVersion(router *httprouter.Router) {
	router.GET(versionPath, DebugLogging(VersionRoute, versionPath))
}

func DebugLogging(h httprouter.Handle, path string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		log.Print("debug: ", path, " request body = ", r.Body)
		h(w, r, p)
		log.Print("debug: ", path, " response=", w)
	}
}

func AddPredicate(router *httprouter.Router, predicate scheduler.Predicate) {
	path := predicatesPrefix + "/" + predicate.Name
	router.POST(path, DebugLogging(PredicateRoute(predicate), path))
}

func AddPrioritize(router *httprouter.Router, prioritize scheduler.Prioritize) {
	path := prioritiesPrefix + "/" + prioritize.Name
	router.POST(path, DebugLogging(PrioritizeRoute(prioritize), path))
}

func AddBind(router *httprouter.Router, bind scheduler.Bind) {
	if handle, _, _ := router.Lookup("POST", bindPath); handle != nil {
		log.Print("warning: AddBind was called more then once!")
	} else {
		router.POST(bindPath, DebugLogging(BindRoute(bind), bindPath))
	}
}

func AddPreemption(router *httprouter.Router, preemption scheduler.Preemption) {
	if handle, _, _ := router.Lookup("POST", preemptionPath); handle != nil {
		log.Print("warning: AddPreemption was called more then once!")
	} else {
		router.POST(preemptionPath, DebugLogging(PreemptionRoute(preemption), preemptionPath))
	}
}