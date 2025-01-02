package helpers

import (
	"fmt"
	"math"
)

func BitsToGiB(bits int64) float64 {
	return float64(bits) / (8 * math.Pow(2, 30))
}

func AnyUnitToBits(value float64, unit string) int64 {
	if unit == "GiB" {
		return int64(value * 1024 * 1024 * 1024 * 8)
	} else if unit == "TiB" {
		return int64(value * 1024 * 1024 * 1024 * 1024 * 8)
	} else {
		return 0
	}
}

// takes a whole database query result and converts the relevant items
func ProcessDataTransferUnits(results interface{}) []map[string]interface{} {
	var processMap func(map[string]any) map[string]interface{}
	processMap = func(m map[string]any) map[string]interface{} {
		updated := make(map[string]any)
		for k, v := range m {
			switch v := v.(type) {
			case map[string]interface{}:
				updated[k] = processMap(v)
			case []map[string]interface{}:
				updated[k] = ProcessDataTransferUnits(v)
			default:
				if k == "uploaded" || k == "downloaded" {
					updated[k] = BitsToGiB(v.(int64))
				} else {
					updated[k] = v
				}
			}
		}
		return updated
	}

	if dataSlice, ok := results.([]map[string]any); ok {
		fmt.Println("ok")
		newData := make([]map[string]any, len(dataSlice))
		for i, item := range dataSlice {
			newData[i] = processMap(item)
		}
		return newData
	}

	// Handle other types or return an error if unexpected
	return nil
}