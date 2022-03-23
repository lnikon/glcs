package computation

import (
	"fmt"
	"strings"
)

func constructParameter(name string, prefix string) string {
	return strings.Join([]string{prefix, name}, "")
}

func (cs *ComputationService) constructLunchArguments(desc *ComputationDescription) []string {
	return []string{constructParameter(string(Replicas), "-"), fmt.Sprintf("%d", desc.Replicas), cs.pgasGraphRunner,
		constructParameter(string(VertexCount), "--"), fmt.Sprintf("%d", desc.VertexCount),
		constructParameter(string(Percentage), "--"), fmt.Sprintf("%d", desc.Density)}

	// TODO: Do we ever gonna need this?
	// constructParameter("export-path"), constructResultFileName(desc.Name)}
}
