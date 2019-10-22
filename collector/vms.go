package collector

//import "encoding/json"
import (
	"github.com/JacobLube/nutanix-exporter/nutanix"

	"github.com/prometheus/client_golang/prometheus"
	//	"github.com/prometheus/log"
)

type VMStat struct {
	HelpText string
	Labels   []string
}

var (
	VMNameSpace string   = "nutanix_vm"
	VMLabels    []string = []string{"name", "uuid"}
)

var VMMetadata map[string]string = map[string]string{}

var VMEntity map[string]string = map[string]string{
	"num_cores_per_vcpu": "...",
	"num_vcpus":          "...",
	"memory_mb":          "...",
	"power_state":        "...",
}

type VMExporter struct {
	Metadata map[string]*prometheus.GaugeVec
	Entity   map[string]*prometheus.GaugeVec
}

func (e *VMExporter) Describe(ch chan<- *prometheus.Desc) {
	e.Metadata = make(map[string]*prometheus.GaugeVec)
	for metricName, helpMsg := range VMMetadata {
		name := normalizeFQN(metricName)
		e.Metadata[metricName] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: VMNameSpace, Name: name, Help: helpMsg}, VMLabels)
		e.Metadata[metricName].Describe(ch)
	}

	e.Entity = make(map[string]*prometheus.GaugeVec)
	for metricName, helpMsg := range VMEntity {
		name := normalizeFQN(metricName)
		e.Entity[metricName] = prometheus.NewGaugeVec(prometheus.GaugeOpts{Namespace: VMNameSpace, Name: name, Help: helpMsg}, VMLabels)
		e.Entity[metricName].Describe(ch)
	}
}

func (e *VMExporter) Collect(ch chan<- prometheus.Metric) {
	vms := nutanixApi.GetVms()
	for vm := range vms.Entities {
		g := e["num_cores_per_vcpu"].WithLabelValues(vm.Name, vm.UUID)
		g.Set(vm.NumCoresPerVCpu)
		g.Collect(ch)
		// for i, k := range e.Entity {
		// 	v, _ := strconv.ParseFloat(s.Entity[i], 64)
		// 	g := k.WithLabelValues(s.Name)
		// 	g.Set(v)
		// 	g.Collect(ch)
		// }
	}
}

func NewVMExporter(api *nutanix.Nutanix) *VMExporter {
	nutanixApi = api
	return &VMExporter{}
}
