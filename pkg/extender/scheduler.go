package extender

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/comail/colog"
	"github.com/julienschmidt/httprouter"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	extenderapi "k8s.io/kube-scheduler/extender/v1"
)

const (
	versionPath      = "/version"
	apiPrefix        = "/scheduler"
	bindPath         = apiPrefix + "/bind"
	preemptionPath   = apiPrefix + "/preemption"
	predicatesPrefix = apiPrefix + "/predicates"
	prioritiesPrefix = apiPrefix + "/priorities"
)

var (
	version string // injected via ldflags at build time

	Graviton2Filter = Predicate{
		Name: "filter_graviton2",
		Func: func(pod v1.Pod, node v1.Node) (bool, error) {
			podArchitectureCounts, err := architectureForPod(pod)
			if err != nil {
				panic(err)
			}

			nodeArchitecture, err := getNodeArchitecture(node)
			if err != nil {
				panic(err)
			}

			// TODO There are additional architectures and OS considerations here
			log.Print("info: ", "filter_graviton2", "podArchitectureCounts = ", podArchitectureCounts, "nodeArchitecture = ", nodeArchitecture)

			if nodeArchitecture == "amd64" {
				if podArchitectureCounts.CountTotal > podArchitectureCounts.CountAmd64 {
					return false, nil
				}
			}

			if nodeArchitecture == "arm64" {
				if podArchitectureCounts.CountTotal > podArchitectureCounts.CountArm64 {
					return false, nil
				}
			}

			return true, nil
		},
	}

	Graviton2Priority = Prioritize{
		Name: "prefer_graviton2",
		Func: func(pod v1.Pod, nodes []v1.Node) (*extenderapi.HostPriorityList, error) {
			var priorityList extenderapi.HostPriorityList
			priorityList = make([]extenderapi.HostPriority, len(nodes))
			for i, node := range nodes {
				priorityList[i] = extenderapi.HostPriority{
					Host:  node.Name,
					Score: 0,
				}
			}
			return &priorityList, nil
		},
	}

	NoBind = Bind{
		Func: func(podName string, podNamespace string, podUID types.UID, node string) error {
			return fmt.Errorf("This extender doesn't support Bind.  Please make 'BindVerb' be empty in your ExtenderConfig.")
		},
	}

	EchoPreemption = Preemption{
		Func: func(_ v1.Pod, _ map[string]*extenderapi.Victims, nodeNameToMetaVictims map[string]*extenderapi.MetaVictims) map[string]*extenderapi.MetaVictims {
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
		log.Printf("warning: LOG_LEVEL=\"%s\" is empty or invalid, falling back to \"INFO\".\n", level)
		return colog.LInfo
	}
}

func Run() {
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
	AddVersion(router)

	predicates := []Predicate{Graviton2Filter}
	for _, p := range predicates {
		AddPredicate(router, p)
	}

	priorities := []Prioritize{Graviton2Priority}
	for _, p := range priorities {
		AddPrioritize(router, p)
	}

	AddBind(router, NoBind)

	log.Print("info: server starting on the port :80")
	if err := http.ListenAndServe(":80", router); err != nil {
		log.Fatal(err)
	}
}
