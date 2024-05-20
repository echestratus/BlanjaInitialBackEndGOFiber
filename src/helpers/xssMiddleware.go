package helpers

import "github.com/microcosm-cc/bluemonday"

func XSSMiddleware(param map[string]interface{}) map[string]interface{} {
	policy := bluemonday.UGCPolicy()

	for key, v := range param {
		if str, ok := v.(string); ok {
			param[key] = policy.Sanitize(str)
		}
	}
	return param
}
