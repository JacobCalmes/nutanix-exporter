package collector

import (
	"strings"

	"github.com/JacobLube/nutanix-exporter/nutanix"
)

var nutanixApi *nutanix.Nutanix

func normalizeFQN(fqn string) string {
	var _fqn string = fqn
	_fqn = strings.Replace(_fqn, ".", "_", -1)
	_fqn = strings.Replace(_fqn, "-", "_", -1)
	_fqn = strings.ToLower(_fqn)

	return _fqn
}
