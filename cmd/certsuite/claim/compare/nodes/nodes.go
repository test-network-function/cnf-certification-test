package nodes

import (
	"github.com/test-network-function/cnf-certification-test/cmd/certsuite/claim/compare/diff"
	"github.com/test-network-function/cnf-certification-test/cmd/certsuite/pkg/claim"
)

// Structure that holds a summary of nodes roles and a slice of NodeDiffReports,
// one per node found in both claim files. In case one node only exists in one
// claim file, it will be marked as "not found in claim[1|2]".
type DiffReport struct {
	Nodes    *diff.Diffs `json:"nodes"`
	CNI      *diff.Diffs `json:"CNI"`
	CSI      *diff.Diffs `json:"CSI"`
	Hardware *diff.Diffs `json:"hardware"`
}

// Stringer method to show in a table the the differences found on each node
// appearing on both claim files. If a node only appears in one claim file, it
// will be flagged as "not found in claim[1|2]".
func (d DiffReport) String() string {
	str := "CLUSTER NODES DIFFERENCES\n"
	str += "-------------------------\n\n"

	if d.Nodes != nil {
		str += d.Nodes.String() + "\n"
	}

	if d.CNI != nil {
		str += d.CNI.String() + "\n"
	}

	if d.CSI != nil {
		str += d.CSI.String() + "\n"
	}

	if d.Hardware != nil {
		str += d.Hardware.String() + "\n"
	}

	return str
}

// Generates a DiffReport from two pointers to claim.Nodes. The report consists
// of a diff.Diffs object per node's section (CNIs, CSIs & Hardware).
func GetDiffReport(claim1Nodes, claim2Nodes *claim.Nodes) *DiffReport {
	return &DiffReport{
		Nodes:    diff.Compare("Nodes", claim1Nodes.NodesSummary, claim2Nodes.NodesSummary, []string{"labels", "annotations"}),
		CNI:      diff.Compare("CNIs", claim1Nodes.CniNetworks, claim2Nodes.CniNetworks, nil),
		CSI:      diff.Compare("CSIs", claim1Nodes.CsiDriver, claim2Nodes.CsiDriver, nil),
		Hardware: diff.Compare("Hardware", claim1Nodes.NodesHwInfo, claim2Nodes.NodesHwInfo, nil),
	}
}
