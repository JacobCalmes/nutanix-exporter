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
	VMLabels    []string = []string{"vm_name", "uuid"}
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
	for _, vm := range vms.Entities {
		metrics := make(map[string]float64)
		metrics["num_cores_per_vcpu"] = vm.NumCoresPerVCpu
		metrics["num_vcpus"] = vm.NumVCpus
		metrics["memory_mb"] = vm.MemoryMb

		metrics["power_state"] = 0
		if vm.PowerState == "on" {
			metrics["power_state"] = 1
		}

		for k, v := range metrics {
			g := e.Entity[k].WithLabelValues(vm.Name, vm.UUID)
			g.Set(v)
			g.Collect(ch)
		}

		// g := e.Entity["num_cores_per_vcpu"].WithLabelValues(vm.Name, vm.UUID)
		// g.Set(vm.NumCoresPerVCpu)
		// g.Collect(ch)

		// g = e.Entity["num_vcpus"].WithLabelValues(vm.Name, vm.UUID)
		// g.Set(vm.NumVCpus)
		// g.Collect(ch)

		// g = e.Entity["memory_mb"].WithLabelValues(vm.Name, vm.UUID)
		// g.Set(vm.MemoryMb)
		// g.Collect(ch)

		// g = e.Entity["power_state"].WithLabelValues(vm.Name, vm.UUID)
		// if vm.PowerState == 'on'g.Set(vm.NumVCpus)
		// g.Collect(ch)
	}
}

func NewVMExporter(api *nutanix.Nutanix) *VMExporter {
	nutanixApi = api
	return &VMExporter{}
}
