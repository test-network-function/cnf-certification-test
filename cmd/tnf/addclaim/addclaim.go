package claim

import (
	"fmt"
	"os"
	"path/filepath"

	"encoding/json"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/test-network-function/cnf-certification-test/pkg/junit"
	"github.com/test-network-function/test-network-function-claim/pkg/claim"
)

var (
	Reportdir string
	Claim     string
	Claim1    string
	Claim2    string

	addclaim = &cobra.Command{
		Use:   "claim",
		Short: "The test suite generates a \"claim\" file",
		RunE:  claimUpdate,
	}
	claimAddFile = &cobra.Command{
		Use:   "add",
		Short: "The test suite generates a \"claim\" file",
		RunE:  claimUpdate,
	}
	claimCompareFiles = &cobra.Command{
		Use:   "compare",
		Short: "Compare 2 \"claim\" file",
		RunE:  claimCompare,
	}
)

const (
	claimFilePermissions = 0o644
)

type cniplugin struct {
	name   string
	plugin interface{}
}
type cnistruct []struct {
	Name    string        "json:\"name\""
	Plugins []interface{} "json:\"plugins\""
}
type Cni struct {
	Claim struct {
		Nodes struct {
			CniPlugins map[string]cnistruct `json:"cniPlugins"`
		} `json:"nodes"`
	} `json:"claim"`
}
type Csi struct {
	Claim struct {
		Nodes struct {
			CsiDriver interface{} `json:"csiDriver"`
		} `json:"nodes"`
	} `json:"claim"`
}

type HwInfo struct {
	Claim struct {
		Nodes struct {
			NodesHwInfo map[string]interface{} `json:"nodesHwInfo"`
		} `json:"nodes"`
	} `json:"claim"`
}

type RawResult struct {
	Claim struct {
		RawResults struct {
			Cnfcertificationtest struct {
				Testsuites struct {
					Testsuite struct {
						Testcase testcase `json:"testcase"`
					} `json:"testsuite"`
				} `json:"testsuites"`
			} `json:"cnf-certification-test"`
		} `json:"rawResults"`
	} `json:"claim"`
}

func claimCompare(cmd *cobra.Command, args []string) error {
	//var cnis []cniplugin
	//var plugin interface{}
	//var name string
	var nodes, nodes2 []string
	claimFileTextPtr := &Claim1
	dat, err := os.ReadFile(*claimFileTextPtr)
	if err != nil {
		log.Fatalf("Error reading claim1 file:%v", err)
	}
	claimFileTextPtr2 := &Claim2
	dat2, err2 := os.ReadFile(*claimFileTextPtr2)
	if err != nil {
		log.Fatalf("Error reading claim2 file :%v", err2)
	}
	// cniclaimRoot.Claim.RawResults
	var cni2 Cni
	err = json.Unmarshal(dat2, &cni2)
	var cni Cni
	err = json.Unmarshal(dat, &cni)
	// csi
	var csi, csi2 Csi
	err = json.Unmarshal(dat2, &csi)
	err = json.Unmarshal(dat2, &csi2)
	// HwInfo
	var hwinfo, hwinfo2 HwInfo
	err = json.Unmarshal(dat, &hwinfo)
	err = json.Unmarshal(dat2, &hwinfo2)
	// rawResult
	var rawResult, rawResult2 RawResult
	err = json.Unmarshal(dat, &rawResult)
	err = json.Unmarshal(dat2, &rawResult2)
	for node, val := range cni.Claim.Nodes.CniPlugins {
		for node2, val2 := range cni2.Claim.Nodes.CniPlugins {
			if node == node2 {
				c, s := compare2cnis(val, val2)
				if len(s) != 0 {
					log.Info("node ", node2, " cnis found in claim1 but not present in claim2: ", s)
				}
				if len(c) != 0 {
					log.Info("node ", node2, " cnis present in both claim 1 and 2 but with different plugins: ", c)
				}
			}
		}
	}
	for key := range hwinfo.Claim.Nodes.NodesHwInfo {
		nodes = append(nodes, key)
	}
	for key := range hwinfo2.Claim.Nodes.NodesHwInfo {
		nodes2 = append(nodes2, key)
	}
	fmt.Println("nodes2 and nodes diffs", missing(nodes2, nodes))
	slist, r := compare2TestCaseResults(rawResult.Claim.RawResults.Cnfcertificationtest.Testsuites.Testsuite.Testcase,
		rawResult2.Claim.RawResults.Cnfcertificationtest.Testsuites.Testsuite.Testcase)
	log.Info("claim1 and claim2 has diff RawResults ", slist)
	log.Info("test name that claim1 has but claim 2 dont has", r)
	return nil
}

type testcase []struct {
	Name   string `json:"-name"`
	Status string `json:"-status"`
}

func compare2TestCaseResults(testcaseResult1, testcaseResult2 testcase) (testcase, []string) {
	var diffresult testcase
	var notFoundtest []string
	for _, result1 := range testcaseResult1 {
		findeName := false
		for _, result2 := range testcaseResult2 {
			if result2.Name == result1.Name {
				findeName = true
				if (result2.Status) != (result1.Status) {
					diffresult = append(diffresult, result1)
				}
				break
			}
		}
		if !findeName {
			notFoundtest = append(notFoundtest, result1.Name)
		}
	}
	return diffresult, notFoundtest
}

// empty struct (0 bytes)
type void struct{}

// missing compares two slices and returns slice of differences
func missing(a, b []string) []string {
	// create map with length of the 'a' slice
	ma := make(map[string]void, len(a))
	diffs := []string{}
	// Convert first slice to map with empty struct (0 bytes)
	for _, ka := range a {
		ma[ka] = void{}
	}
	// find missing values in a
	for _, kb := range b {
		if _, ok := ma[kb]; !ok {
			diffs = append(diffs, kb)
		}
	}
	return diffs
}

func compare2cnis(cniList1, cniList2 cnistruct) (cnistruct, []string) {
	var diffplugins cnistruct
	var notFoundNames []string
	for _, plugin1 := range cniList1 {
		findeName := false
		for _, plugin2 := range cniList2 {
			if plugin2.Name == plugin1.Name {
				findeName = true
				if plugin2.Plugins != nil {
					if len(plugin2.Plugins) != len(plugin1.Plugins) {
						diffplugins = append(diffplugins, plugin1)
					}
				}

				break
			}
		}
		if !findeName {
			notFoundNames = append(notFoundNames, plugin1.Name)
		}
	}
	return diffplugins, notFoundNames
}

//nolint:funlen
func claimUpdate(cmd *cobra.Command, args []string) error {
	claimFileTextPtr := &Claim
	reportFilesTextPtr := &Reportdir
	fileUpdated := false
	dat, err := os.ReadFile(*claimFileTextPtr)
	if err != nil {
		log.Fatalf("Error reading claim file :%v", err)
	}
	claimRoot := readClaim(&dat)
	junitMap := claimRoot.Claim.RawResults
	items, err := os.ReadDir(*reportFilesTextPtr)
	if err != nil {
		log.Fatalf("Error reading directory: %v", err)
	}
	for _, item := range items {
		fileName := item.Name()
		extension := filepath.Ext(fileName)
		reportKeyName := fileName[0 : len(fileName)-len(extension)]
		if _, ok := junitMap[reportKeyName]; ok {
			log.Printf("Skipping: %s already exists in supplied `%s` claim file", reportKeyName, *claimFileTextPtr)
		} else {
			junitMap[reportKeyName], err = junit.ExportJUnitAsMap(fmt.Sprintf("%s/%s", *reportFilesTextPtr, item.Name()))
			if err != nil {
				log.Fatalf("Error reading JUnit XML file into JSON: %v", err)
			}
			fileUpdated = true
		}
	}
	claimRoot.Claim.RawResults = junitMap
	payload, err := json.MarshalIndent(claimRoot, "", "  ")
	if err != nil {
		log.Fatalf("Failed to generate the claim: %v", err)
	}
	err = os.WriteFile(*claimFileTextPtr, payload, claimFilePermissions)
	if err != nil {
		log.Fatalf("Error writing claim data:\n%s", string(payload))
	}
	if fileUpdated {
		log.Printf("Claim file `%s` updated\n", *claimFileTextPtr)
	} else {
		log.Printf("No changes were applied to `%s`\n", *claimFileTextPtr)
	}
	return nil
}

func readClaim(contents *[]byte) *claim.Root {
	var claimRoot claim.Root
	err := json.Unmarshal(*contents, &claimRoot)
	if err != nil {
		log.Fatalf("Error reading claim constents file into type: %v", err)
	}
	return &claimRoot
}

func NewCommand() *cobra.Command {
	claimAddFile.Flags().StringVarP(
		&Reportdir, "reportdir", "r", "",
		"dir of JUnit XML reports. (Required)",
	)

	err := claimAddFile.MarkFlagRequired("reportdir")
	if err != nil {
		return nil
	}

	claimAddFile.Flags().StringVarP(
		&Claim, "claim", "c", "",
		"existing claim file. (Required)",
	)
	err = claimAddFile.MarkFlagRequired("claim")
	if err != nil {
		return nil
	}
	addclaim.AddCommand(claimAddFile)
	claimCompareFiles.Flags().StringVarP(
		&Claim1, "claim1", "1", "",
		"existing claim1 file. (Required) first file to compare",
	)
	claimCompareFiles.Flags().StringVarP(
		&Claim2, "claim2", "2", "",
		"existing claim2 file. (Required) seconed file to compare with",
	)
	err = claimAddFile.MarkFlagRequired("claim")
	if err != nil {
		return nil
	}
	addclaim.AddCommand(claimCompareFiles)
	return addclaim
}
