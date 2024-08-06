package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/test-network-function/cnf-certification-test/pkg/configuration"
	"gopkg.in/yaml.v3"
)

type configOption struct {
	Option string
	Help   string
}

func NewCommand() *cobra.Command {
	return generateConfigCmd
}

var (
	// generateConfiCmd is a helper tool to generate a CNF config YAML file.
	generateConfigCmd = &cobra.Command{
		Use:   "config",
		Short: "Generates a CNF config YAML file with user input.",
		Run: func(cmd *cobra.Command, args []string) {
			generateConfig()
		},
	}
)

var tnfConfig = configuration.TestConfiguration{}

var templates = &promptui.SelectTemplates{
	Label:    "{{ . }}",
	Active:   "\U00002B9E {{ .Option | cyan }}",
	Inactive: "  {{ .Option | cyan }}",
	Details: `
--------- {{ .Option | faint }}  ----------
{{ .Help }}`,
}

func generateConfig() {
	mainMenu := []configOption{
		{Option: create, Help: createConfigHelp},
		{Option: show, Help: showConfigHelp},
		{Option: save, Help: saveConfigHelp},
		{Option: quit, Help: exitHelp},
	}

	var exit bool
	for !exit {
		mainPrompt := promptui.Select{
			Label:        "",
			Items:        mainMenu,
			Templates:    templates,
			Size:         4,
			HideSelected: true,
		}

		opt, _, err := mainPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		switch mainMenu[opt].Option {
		case create:
			createConfiguration()
		case show:
			showConfiguration(&tnfConfig)
		case save:
			saveConfiguration(&tnfConfig)
		case quit:
			exit = true
		}
	}
}

func createConfiguration() {
	createMenu := []configOption{
		{Option: cnfResources, Help: cnfResourcesHelp},
		{Option: exceptions, Help: exceptionsdHelp},
		// {Option: collector, Help: collectordHelp},
		{Option: settings, Help: settingsHelp},
		{Option: previousMenu, Help: backHelp},
	}

	createPrompt := promptui.Select{
		Label:        "",
		Items:        createMenu,
		Templates:    templates,
		Size:         5,
		HideSelected: true,
	}

	var exit bool
	for !exit {
		i, _, err := createPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		switch createMenu[i].Option {
		case cnfResources:
			createCnfResourcesConfiguration()
		case exceptions:
			createExceptionsConfiguration()
		// case collector:
		// 	createCollectorConfiguration()
		case settings:
			createSettingsConfiguration()
		case previousMenu:
			exit = true
		}
	}
}

func showConfiguration(config *configuration.TestConfiguration) {
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("could not marshal the YAML file, err: %v", err)
		return
	}
	fmt.Println("================= CNF CONFIGURATION =================")
	fmt.Println(string(configYaml))
	fmt.Println("=====================================================")
}

func saveConfiguration(config *configuration.TestConfiguration) {
	configYaml, err := yaml.Marshal(config)
	if err != nil {
		log.Printf("could not marshal the YAML file, err: %v", err)
		return
	}

	saveConfigPrompt := promptui.Prompt{
		Label:   "CNF config file",
		Default: defaultConfigFileName,
	}

	configFileName, err := saveConfigPrompt.Run()
	if err != nil {
		log.Printf("could not read config file name, err: %v\n", err)
		return
	}

	err = os.WriteFile(configFileName, configYaml, defaultConfigFilePermissions)
	if err != nil {
		log.Printf("could not write file, err: %v", err)
		return
	}

	fmt.Println(color.GreenString("Configuration saved"))
}

func createCnfResourcesConfiguration() {
	cnfResourcesOptions := []configOption{
		{Option: namespaces, Help: namespacesHelp},
		{Option: pods, Help: podLabelsHelp},
		{Option: operators, Help: operatorLabelsHelp},
		{Option: crdFilters, Help: crdFiltersHelp},
		{Option: managedDeployments, Help: managedDeploymentsHelp},
		{Option: managedStatefulSets, Help: managedStatefulSetsHelp},
		{Option: previousMenu, Help: backHelp},
	}
	cnfResourcesSearcher := func(input string, index int) bool {
		basicOption := cnfResourcesOptions[index]
		name := strings.ReplaceAll(strings.ToLower(basicOption.Option), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

		return strings.Contains(name, input)
	}
	cnfResourcesPrompt := promptui.Select{
		Label:        "",
		Items:        cnfResourcesOptions,
		Templates:    templates,
		Size:         7,
		Searcher:     cnfResourcesSearcher,
		HideSelected: true,
	}
	var exit bool
	for !exit {
		i, _, err := cnfResourcesPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		switch cnfResourcesOptions[i].Option {
		case namespaces:
			loadNamespaces(getAnswer(namespacePrompt, namespaceSyntax, namespaceExample))
		case pods:
			loadPodLabels(getAnswer(podsPrompt, podsSyntax, podsExample))
		case operators:
			loadOperatorLabels(getAnswer(operatorsPrompt, operatorsSyntax, operatorsExample))
		case crdFilters:
			loadCRDfilters(getAnswer(crdFiltersPrompt, crdFiltersSyntax, crdFiltersExample))
		case managedDeployments:
			loadManagedDeployments(getAnswer(managedDeploymentsPrompt, managedDeploymentsSyntax, managedDeploymentsExample))
		case managedStatefulSets:
			loadManagedStatefulSets(getAnswer(managedStatefulSetsPrompt, managedStatefulSetsSyntax, managedStatefulSetsExample))
		case previousMenu:
			exit = true
		}
	}
}

func createExceptionsConfiguration() {
	exceptionsOptions := []configOption{
		{Option: kernelTaints, Help: kernelTaintsHelp},
		{Option: helmCharts, Help: helmChartsHelp},
		{Option: protocolNames, Help: protocolNamesHelp},
		{Option: services, Help: servicesHelp},
		{Option: nonScalableDeployments, Help: nonScalableDeploymentsHelp},
		{Option: nonScalableStatefulSets, Help: nonScalableStatefulSetsHelp},
		{Option: previousMenu, Help: backHelp},
	}
	exceptionsSearcher := func(input string, index int) bool {
		exceptionOption := exceptionsOptions[index]
		name := strings.ReplaceAll(strings.ToLower(exceptionOption.Option), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

		return strings.Contains(name, input)
	}
	exceptionsPrompt := promptui.Select{
		Label:        "",
		Items:        exceptionsOptions,
		Templates:    templates,
		Size:         7,
		Searcher:     exceptionsSearcher,
		HideSelected: true,
	}
	var exit bool
	for !exit {
		i, _, err := exceptionsPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		switch exceptionsOptions[i].Option {
		case kernelTaints:
			loadAcceptedKernelTaints(getAnswer(kernelTaintsPrompt, kernelTaintsSyntax, kernelTaintsExample))
		case helmCharts:
			loadHelmCharts(getAnswer(helmChartsPrompt, helmChartsSyntax, helmChartsExample))
		case protocolNames:
			loadProtocolNames(getAnswer(protocolNamesPrompt, protocolNamesSyntax, protocolNamesExample))
		case services:
			loadServices(getAnswer(servicesPrompt, servicesSyntax, servicesExample))
		case nonScalableDeployments:
			loadNonScalableDeployments(getAnswer(nonScalableDeploymentsPrompt, nonScalableDeploymentsSyntax, nonScalableDeploymentsExample))
		case nonScalableStatefulSets:
			loadNonScalableStatefulSets(getAnswer(nonScalableStatefulSetsPrompt, nonScalableStatefulSetsSyxtax, nonScalableStatefulSetsExample))
		case previousMenu:
			exit = true
		}
	}
}

//nolint:unused
func createCollectorConfiguration() {
	collectorOptions := []configOption{
		{Option: appEndPoint, Help: ""},
		{Option: executedBy, Help: ""},
		{Option: partnerName, Help: ""},
		{Option: appPassword, Help: ""},
		{Option: previousMenu, Help: backHelp},
	}
	collectorSearcher := func(input string, index int) bool {
		collectorOption := collectorOptions[index]
		name := strings.ReplaceAll(strings.ToLower(collectorOption.Option), " ", "")
		input = strings.ReplaceAll(strings.ToLower(input), " ", "")

		return strings.Contains(name, input)
	}
	collectorPrompt := promptui.Select{
		Label:        "",
		Items:        collectorOptions,
		Templates:    templates,
		Size:         5,
		Searcher:     collectorSearcher,
		HideSelected: true,
	}
	var exit bool
	for !exit {
		i, _, err := collectorPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		switch collectorOptions[i].Option {
		case appEndPoint:
			// TODO: to be implemented
		case executedBy:
			// TODO: to be implemented
		case partnerName:
			// TODO: to be implemented
		case appPassword:
			// TODO: to be implemented
		case previousMenu:
			exit = true
		}
	}
}

func createSettingsConfiguration() {
	settingsOptions := []configOption{
		{Option: debugDaemonSet, Help: debugDaemonSetHelp},
		{Option: previousMenu, Help: backHelp},
	}
	settingsPrompt := promptui.Select{
		Label:        "",
		Items:        settingsOptions,
		Templates:    templates,
		Size:         2,
		HideSelected: true,
	}
	var exit bool
	for !exit {
		i, _, err := settingsPrompt.Run()
		if err != nil {
			log.Printf("Prompt failed %v\n", err)
			return
		}
		switch settingsOptions[i].Option {
		case debugDaemonSet:
			loadDebugDaemonSetNamespace(getAnswer(debugDaemonSetPrompt, debugDaemonSetSyntax, debugDaemonSetExample))
		case previousMenu:
			exit = true
		}
	}
}

func getAnswer(prompt, syntax, example string) []string {
	fullPrompt := color.HiCyanString("%s\n", prompt) +
		color.CyanString("Syntax: ") + color.WhiteString("%s\n", syntax) +
		color.CyanString("Example: ") + color.WhiteString("%s\n", example) + color.HiCyanString("> ")
	fmt.Print(fullPrompt)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	err := scanner.Err()
	if err != nil {
		log.Printf("could not read user input, err: %v", err)
		return nil
	}

	// Split CSV string by ',' and remove any whitespace
	fields := strings.Split(scanner.Text(), ",")
	for i, field := range fields {
		fields[i] = strings.TrimSpace(field)
	}

	return fields
}

func loadNamespaces(namespaces []string) {
	tnfConfig.TargetNameSpaces = nil
	for _, namespace := range namespaces {
		tnfNamespace := configuration.Namespace{Name: namespace}
		tnfConfig.TargetNameSpaces = append(tnfConfig.TargetNameSpaces, tnfNamespace)
	}
}

func loadPodLabels(podLabels []string) {
	tnfConfig.PodsUnderTestLabels = nil
	tnfConfig.PodsUnderTestLabels = podLabels
}

func loadOperatorLabels(operatorLabels []string) {
	tnfConfig.OperatorsUnderTestLabels = nil
	tnfConfig.OperatorsUnderTestLabels = operatorLabels
}

func loadCRDfilters(crdFilters []string) {
	tnfConfig.CrdFilters = nil
	for _, crdFilterStr := range crdFilters {
		crdFilter := strings.Split(crdFilterStr, "/")
		crdFilterName := crdFilter[0]
		crdFilterScalable, err := strconv.ParseBool(crdFilter[1])
		if err != nil {
			log.Printf("could not parse CRD filter, err: %v", err)
			return
		}
		tnfCrdFilter := configuration.CrdFilter{NameSuffix: crdFilterName, Scalable: crdFilterScalable}
		tnfConfig.CrdFilters = append(tnfConfig.CrdFilters, tnfCrdFilter)
	}
}

func loadManagedDeployments(deployments []string) {
	tnfConfig.ManagedDeployments = nil
	for _, deployment := range deployments {
		tnfManagedDeployment := configuration.ManagedDeploymentsStatefulsets{Name: deployment}
		tnfConfig.ManagedDeployments = append(tnfConfig.ManagedDeployments, tnfManagedDeployment)
	}
}

func loadManagedStatefulSets(statefulSets []string) {
	tnfConfig.ManagedStatefulsets = nil
	for _, statefulSet := range statefulSets {
		tnfManagedStatefulSet := configuration.ManagedDeploymentsStatefulsets{Name: statefulSet}
		tnfConfig.ManagedStatefulsets = append(tnfConfig.ManagedStatefulsets, tnfManagedStatefulSet)
	}
}

func loadAcceptedKernelTaints(taints []string) {
	tnfConfig.AcceptedKernelTaints = nil
	for _, taint := range taints {
		tnfKernelTaint := configuration.AcceptedKernelTaintsInfo{Module: taint}
		tnfConfig.AcceptedKernelTaints = append(tnfConfig.AcceptedKernelTaints, tnfKernelTaint)
	}
}

func loadHelmCharts(helmCharts []string) {
	tnfConfig.SkipHelmChartList = nil
	for _, chart := range helmCharts {
		tnfHelmChart := configuration.SkipHelmChartList{Name: chart}
		tnfConfig.SkipHelmChartList = append(tnfConfig.SkipHelmChartList, tnfHelmChart)
	}
}

func loadProtocolNames(protocolNames []string) {
	tnfConfig.ValidProtocolNames = nil
	tnfConfig.ValidProtocolNames = protocolNames
}

func loadServices(services []string) {
	tnfConfig.ServicesIgnoreList = nil
	tnfConfig.ServicesIgnoreList = services
}

func loadNonScalableDeployments(nonScalableDeployments []string) {
	tnfConfig.SkipScalingTestDeployments = nil
	for _, nonScalableDeploymentStr := range nonScalableDeployments {
		nonScalableDeployment := strings.Split(nonScalableDeploymentStr, "/")
		const nonScalableDeploymentsFields = 2
		if len(nonScalableDeployment) != nonScalableDeploymentsFields {
			log.Println("could not parse Non-scalable Deployment")
			return
		}
		nonScalableDeploymentName := nonScalableDeployment[0]
		nonScalableDeploymentNamespace := nonScalableDeployment[1]
		tnfNonScalableDeployment := configuration.SkipScalingTestDeploymentsInfo{Name: nonScalableDeploymentName,
			Namespace: nonScalableDeploymentNamespace}
		tnfConfig.SkipScalingTestDeployments = append(tnfConfig.SkipScalingTestDeployments, tnfNonScalableDeployment)
	}
}

func loadNonScalableStatefulSets(nonScalableStatefulSets []string) {
	tnfConfig.SkipScalingTestStatefulSets = nil
	for _, nonScalableStatefulSetStr := range nonScalableStatefulSets {
		nonScalableStatefulSet := strings.Split(nonScalableStatefulSetStr, "/")
		const nonScalableStatefulSetFields = 2
		if len(nonScalableStatefulSet) != nonScalableStatefulSetFields {
			log.Println("could not parse Non-scalable StatefulSet")
			return
		}
		nonScalableStatefulSetName := nonScalableStatefulSet[0]
		nonScalableStatefulSetNamespace := nonScalableStatefulSet[1]
		tnfNonScalableStatefulSet := configuration.SkipScalingTestStatefulSetsInfo{Name: nonScalableStatefulSetName,
			Namespace: nonScalableStatefulSetNamespace}
		tnfConfig.SkipScalingTestStatefulSets = append(tnfConfig.SkipScalingTestStatefulSets, tnfNonScalableStatefulSet)
	}
}

func loadDebugDaemonSetNamespace(namespace []string) {
	tnfConfig.DebugDaemonSetNamespace = namespace[0]
}
