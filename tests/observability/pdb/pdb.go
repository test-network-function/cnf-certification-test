package pdb

import (
	"fmt"
	"math"

	policyv1 "k8s.io/api/policy/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	percentageDivisor = 100
)

// percentageToFloat converts a percentage string to a float.
func percentageToFloat(percentage string) (float64, error) {
	var percentageFloat float64
	_, err := fmt.Sscanf(percentage, "%f%%", &percentageFloat)
	if err != nil {
		return 0, err
	}
	return percentageFloat / percentageDivisor, nil
}

func CheckPDBIsValid(pdb *policyv1.PodDisruptionBudget, replicas *int32) (bool, error) {
	var replicaCount int32
	if replicas != nil {
		replicaCount = *replicas
	} else {
		replicaCount = 1 // default value
	}

	var minAvailableValue int
	var maxUnavailableValue int

	if pdb.Spec.MinAvailable != nil {
		var err error
		minAvailableValue, err = intOrStringToValue(pdb.Spec.MinAvailable, replicaCount)
		if err != nil {
			return false, err
		}

		// Tests for the minAvailable spec.
		if minAvailableValue == 0 {
			return false, fmt.Errorf("field .spec.minAvailable cannot be zero. Currently set to: %d. Replicas set to: %d", minAvailableValue, replicaCount)
		}

		if minAvailableValue > int(replicaCount) {
			return false, fmt.Errorf("minAvailable cannot be greater than replicas. Currently set to: %d. Replicas set to: %d", minAvailableValue, replicaCount)
		}
	}

	if pdb.Spec.MaxUnavailable != nil {
		var err error
		maxUnavailableValue, err = intOrStringToValue(pdb.Spec.MaxUnavailable, replicaCount)
		if err != nil {
			return false, err
		}

		// Tests for the maxUnavailable spec.
		if maxUnavailableValue >= int(replicaCount) {
			return false, fmt.Errorf("field .spec.maxUnavailable cannot be greater than or equal to the number of pods in the replica. Currently set to: %d. Replicas set to: %d", maxUnavailableValue, replicaCount)
		}
	}

	return true, nil
}

func intOrStringToValue(intOrStr *intstr.IntOrString, replicas int32) (int, error) {
	switch intOrStr.Type {
	case intstr.Int:
		return intOrStr.IntValue(), nil
	case intstr.String:
		v, err := percentageToFloat(intOrStr.StrVal)
		if err != nil {
			return 0, fmt.Errorf("invalid value %q: %v", intOrStr.StrVal, err)
		}
		return int(math.RoundToEven(v * float64(replicas))), nil
	}
	return 0, fmt.Errorf("invalid type: neither int nor percentage")
}
