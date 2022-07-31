#!/usr/bin/env bash

# Requires
# - kubeconfig file mounted to /usr/tnf/kubeconfig/config
#   If more than one kubeconfig needs to be used, bind
#   additional volumes for each kubeconfig, e.g.
#     - /usr/tnf/kubeconfig/config
#     - /usr/tnf/kubeconfig/config.2
#     - /usr/tnf/kubeconfig/config.3
# - TNF config files mounted into /usr/tnf/config
# - A directory to output claim into mounted at /usr/tnf/claim
# - A $KUBECONFIG environment variable passed to the TNF container
#   containing all paths to kubeconfigs located in the container, e.g.
#   KUBECONFIG=/usr/tnf/kubeconfig/config:/usr/tnf/kubeconfig/config.2

export REQUIRED_NUM_OF_ARGS=5
export REQUIRED_VARS=('LOCAL_KUBECONFIG' 'LOCAL_TNF_CONFIG' 'OUTPUT_LOC')
export REQUIRED_VARS_ERROR_MESSAGES=(
	'KUBECONFIG is invalid or not given. Use the -k option to provide path to one or more kubeconfig files.'
	'TNFCONFIG is required. Use the -t option to specify the directory containing the TNF configuration files.'
	'OUTPUT_LOC is required. Use the -o option to specify the output location for the test results.'
)

export TNF_IMAGE_NAME=test-network-function
export TNF_IMAGE_TAG=latest
export TNF_OFFICIAL_ORG=quay.io/testnetworkfunction/
export TNF_OFFICIAL_IMAGE="${TNF_OFFICIAL_ORG}${TNF_IMAGE_NAME}:${TNF_IMAGE_TAG}"
export TNF_CMD="./run-cnf-suites.sh"
export OUTPUT_ARG="-o"
export FOCUS_ARG="-f"
export SKIP_ARG="-s"
export CONTAINER_NETWORK_MODE='host'

usage() {
	read -d '' usage_prompt <<- EOF
	Usage: $0 -t TNFCONFIG -o OUTPUT_LOC [-i IMAGE] [-k KUBECONFIG] [-n NETWORK_MODE] [-d DNS_RESOLVER_ADDRESS] -f SUITE [... SUITE]

	Configure and run the containerised TNF test offering.

	Options (required)
	  -t: set the directory containing TNF config files set up for the test.
	  -o: set the output location for the test results.

	Options (optional)
	  -i: set the TNF container image. Supports local images, as well as images from external registries.
	  -k: set path to one or more local kubeconfigs, separated by a colon.
	      The -k option takes precedence, overwriting the results of local kubeconfig autodiscovery.
	      See the 'Kubeconfig lookup order' section below for more details.
	  -n: set the network mode of the container.
	  -d: set the DNS resolver address for the test containers started by docker, may be required with 
	      certain docker version if the kubeconfig contains host names
    -f: Set the suites that should be tested, multiple suites can be supplied
    -s: Set the test cases that should be skipped

	Kubeconfig lookup order
	  1. If -k is specified, use the paths provided with the -k option.
	  2. If -k is not specified, use paths defined in \$KUBECONFIG on the underlying host.
	  3. If no paths are defined, use the default kubeconfig file located in '\$HOME/.kube/config'
	     (currently: $HOME/.kube/config).

	Examples
	  $0 -t ~/tnf/config -o ~/tnf/output -f networking access-control -s access-control-host-resource-PRIVILEGED_POD

	  Because -k is omitted, $(basename $0) will first try to autodiscover local kubeconfig files.
	  If it succeeds, the networking and access-control tests will be run using the autodiscovered configuration.
	  The test results will be saved to the '~/tnf/output' directory on the host.

	  $0 -k ~/.kube/ABC:~/.kube/DEF -t ~/tnf/config -o ~/tnf/output -f access-control networking

	  The command will bind two kubeconfig files (~/.kube/ABC and ~/.kube/DEF) to the TNF container,
	  run the access-control and networking tests, and save the test results into the '~/tnf/output' directory
	  on the host.

	  $0 -i custom-tnf-image:v1.2-dev -t ~/tnf/config -o ~/tnf/output -f access-control networking

	  The command will run the access-control and networking tests as implemented in the custom-tnf-image:v1.2-dev
	  local image set by the -i parameter. The test results will be saved to the '~/tnf/output' directory.

	Test suites
	  Allowed tests are listed in the README.
	  Note: Tests must be specified after all other arguments!
	EOF

	echo -e "$usage_prompt"
}

usage_error() {
	usage
	exit 1
}

check_required_vars() {
	local var_missing=false

	for index in "${!REQUIRED_VARS[@]}"; do
		var=${REQUIRED_VARS[$index]}
		if [[ -z ${!var} ]]; then
			error_message=${REQUIRED_VARS_ERROR_MESSAGES[$index]}
			echo "$0: error: $error_message" 1>&2
			var_missing=true
		fi
	done

	if $var_missing; then
		echo ""
		usage_error
	fi
}

check_cli_required_num_of_args() {
	if (($# < $REQUIRED_NUM_OF_ARGS)); then
		usage_error
	fi;
}

perform_kubeconfig_autodiscovery() {
	if [[ -n "$KUBECONFIG" ]]; then
		export LOCAL_KUBECONFIG=$KUBECONFIG
		kubeconfig_autodiscovery_source='$KUBECONFIG'
	elif [[ -f "$HOME/.kube/config" ]]; then
		export LOCAL_KUBECONFIG=$HOME/.kube/config
		kubeconfig_autodiscovery_source="\$HOME/.kube/config ($HOME/.kube/config)"
	fi
}

display_kubeconfig_autodiscovery_summary() {
	if [[ -n "$kubeconfig_autodiscovery_source" ]]; then
		echo "Kubeconfig Autodiscovery: configuration loaded from $kubeconfig_autodiscovery_source"
	fi
}

if [ ! -z "${REQUIRED_NUM_OF_ARGS}" ]; then
	check_cli_required_num_of_args $@
fi

perform_kubeconfig_autodiscovery

# Parge args beginning with -
while [[ $1 == -* ]]; do
    case "$1" in
      -h|--help|-\?) usage; exit 0;;
      -k) if (($# > 1)); then
            export LOCAL_KUBECONFIG=$2
            unset kubeconfig_autodiscovery_source
            shift
          else
            echo "-k requires an argument" 1>&2
            exit 1
          fi
          echo "-k $LOCAL_KUBECONFIG"
          ;;
      -t) if (($# > 1)); then
            export LOCAL_TNF_CONFIG=$2; shift
          else
            echo "-t requires an argument" 1>&2
            exit 1
          fi       
          echo  "-t $LOCAL_TNF_CONFIG"
          ;;
      -o) if (($# > 1)); then
            export OUTPUT_LOC=$2; shift
          else
            echo "-o requires an argument" 1>&2
            exit 1
          fi
          echo "-o $OUTPUT_LOC"
          ;;
      -i) if (($# > 1)); then
            export TNF_IMAGE=$2; shift
          else
            echo "-i requires an argument" 1>&2
            exit 1
          fi
          echo "-i $TNF_IMAGE"
          ;;
      -n) if (($# > 1)); then
            export CONTAINER_NETWORK_MODE=$2; shift
          else
            echo "-n requires an argument" 1>&2
            exit 1
          fi
          echo "-n $CONTAINER_NETWORK_MODE"
          ;;
      -d) if (($# > 1)); then
            export DNS_ARG=$2; shift 2
          else
            echo "-d requires an argument" 1>&2
            exit 1
          fi
          echo "-d $DNS_ARGS"
          ;;
      -s)
        ONCE=true
        while (( "$#" >= 2 )) && ! [[ $2 = --* ]] && ! [[ $2 = -* ]] ; do
          if [ $ONCE = true ]; then
            TNF_SKIP_SUITES="$2"
            ONCE=false
          fi
          TNF_SKIP_SUITES="$TNF_SKIP_SUITES $2"
          shift
        done
        export TNF_SKIP_SUITES
        echo "-s $TNF_SKIP_SUITES"
        ;;
      -f)
        ONCE=true
        while (( "$#" >= 2 )) && ! [[ $2 = --* ]]  && ! [[ $2 = -* ]] ; do
          if [ $ONCE = true ]; then
            TNF_FOCUS_SUITES="$2"
            ONCE=false
          fi
          TNF_FOCUS_SUITES="$TNF_FOCUS_SUITES $2"
          shift
        done
        export TNF_FOCUS_SUITES
        echo "-f $TNF_FOCUS_SUITES"
        ;;
      --) shift; break;;
      -*) echo "invalid option: $1" 1>&2; usage_error;;
    esac
  shift
done

display_kubeconfig_autodiscovery_summary
check_required_vars

cd script

./run-cfd-container.sh

./run-container.sh  "$@"
