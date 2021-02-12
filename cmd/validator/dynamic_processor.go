package main

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"reflect"
	"sort"
)

// DynamicProcessor implements process address by dynamically normalizing the JSON and ensuring the results are equal.
type DynamicProcessor struct {
}

func mustEncode(data interface{}) string {
	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		errorLog.Fatalf("failed to encode data: %v", err)
	}
	return string(result)
}

// ProcessAddress is the entrypoint for the DynamicProcessor.
func (gp DynamicProcessor) ProcessAddress(addr string, config Params) (Result, error) {
	indexerURL := fmt.Sprintf("%s:/v2/accounts/%s", config.indexerURL, addr)
	indexerAcct, err := getResponse(indexerURL, config.indexerToken)
	if err != nil {
		return Result{}, errors.Wrap(err, fmt.Sprintf("unable to lookup indexer acct %s", addr))
	}
	algodURL := fmt.Sprintf("%s:/v2/accounts/%s", config.algodURL, addr)
	algodAcct, err := getResponse(algodURL, config.algodToken)
	if err != nil {
		return Result{}, errors.Wrap(err, fmt.Sprintf("unable to lookup algod acct %s", addr))
	}
	indexerNorm, err := normalize(indexerAcct)
	if err != nil {
		errorLog.Printf("failed to normalize indexer: %v\n", err)
	}
	algodNorm, err := normalize(algodAcct)
	if err != nil {
		errorLog.Printf("failed to normalize algod: %v\n", err)
	}

	if !reflect.DeepEqual(indexerNorm, algodNorm) {
		return Result{
			Equal:   false,
			Retries: 0,
			Details: &ErrorDetails{
				address: addr,
				algod:   fmt.Sprintf("RawJson\n%s\nNormalizedJson\n%s\n", mustEncode(algodAcct), mustEncode(algodNorm)),
				indexer:   fmt.Sprintf("RawJson\n%s\nNormalizedJson\n%s\n", mustEncode(indexerAcct), mustEncode(indexerNorm)),
				diff:    nil,
			},
		}, nil
	}
	return Result{Equal:   true}, nil
}


// getResponse queries the url for a response and marshals it into a generic type
func getResponse(url, token string) (map[string]interface{}, error) {
	var target map[string]interface{}
	auth := fmt.Sprintf("Bearer %s", token)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", auth)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			errorLog.Fatalf("failed to close body: %v", err)
		}
	}()

	err = json.NewDecoder(resp.Body).Decode(&target)
	if err != nil {
		return nil, err
	}
	return target, nil
}


// isDeleted is a simple helper to check for boolean values
func isDeleted(val interface{}) (bool, error) {
	switch val.(type) {
	case bool:
		return val.(bool), nil
	default:
		return false, fmt.Errorf("unable to parse value as boolean (%v)", val)
	}
}

// isEmpty checks if the value should be culled (zero, nil, "", empty object, empty array).
func isEmpty(val interface{}) (bool, error) {
	if val == nil {
		return true, nil
	}

	switch t := val.(type) {
	case int:
		i := val.(int)
		if i == 0 {
			return true, nil
		}
	case float64:
		f := val.(float64)
		if f == 0 {
			return true, nil
		}
	case string:
		s := val.(string)
		if s == "" {
			return true, nil
		}
	case bool:
		// don't cull booleans
		return false, nil
	case map[string]interface{}:
		m := val.(map[string]interface{})
		if len(m) == 0 {
			return true, nil
		}
	case []interface{}:
		arr := val.([]interface{})
		if len(arr) == 0 {
			return true, nil
		}
	default:
		var r = reflect.TypeOf(t)
		return false, fmt.Errorf("unknown leaf type:%v", r)
	}

	return false, nil
}

func sortHelper(arr []interface{}, key string) {
	sort.Slice(arr, func(i, j int) bool {
		// This will panic if things are not what they need to be
		v1 := arr[i].(map[string]interface{})[key].(float64)
		v2 := arr[j].(map[string]interface{})[key].(float64)
		return v1 < v2
	})
}

// normalize is the top level parsing function once the top level things are done it recurses.
func normalize(data map[string]interface{}) (interface{}, error) {
	var result map[string]interface{}

	// If there is a top level account field (indexer), use it.
	if val, ok := data["account"]; ok {
		result = val.(map[string]interface{})
	} else {
		result = data
	}

	// A zero amount indicates algod is reporting an unused account or indexer a deleted account.
	val, ok := result["amount"]
	if !ok || val == 0 {
		return nil, nil
	}

	return normalizeRecurse("", result)
}

func normalizeRecurse(key string, data interface{}) (interface{}, error) {
	result := data

	// If the data is a complex type, recursively normalize result
	switch data.(type) {
	case map[string]interface{}:
		// Handle objects.
		object := data.(map[string]interface{})

		if deletedVal, ok := object["deleted"]; ok {
			deleted, err := isDeleted(deletedVal)
			if err != nil {
				return nil, err
			}

			// Trim this branch if it is deleted
			if deleted {
				return nil, nil
			}
		}

		// Remove stuff that isn't in the union of indexer and algod responses

		// only in indexer
		delete(object, "deleted")

		// at-round fields are only in indexer
		delete(object, "created-at-round")
		delete(object, "deleted-at-round")
		delete(object, "destroyed-at-round")
		delete(object, "optin-at-round")
		delete(object, "opted-in-at-round")
		delete(object, "opted-out-at-round")
		delete(object, "closeout-at-round")
		delete(object, "closed-out-at-round")
		delete(object, "close-out-at-round")

		// indexer does not attach creator to asset holdings
		delete(object, "creator")

		// indexer adds a special sig-type field to the top level (this could be moved to the non-recursive call)
		delete(object, "sig-type")

		// indexer and algod are usually off by 1 round, so don't bother checking this.
		delete(object, "round")

		// algod seems to look this up on demand, indexers value is the base at the last time the account changed
		delete(object, "reward-base")

		// indexer doesn't have this field yet.
		delete(object, "apps-total-schema")

		for key, val := range object {
			normalized, err := normalizeRecurse(key, val)
			if err != nil {
				return nil, err
			}
			if normalized != nil {
				// Update with the normalized value
				object[key] = normalized
			} else {
				// Omit empty nested values
				delete(object, key)
			}
		}

		result = object
	case []interface{}:
		// Normalize each element of array
		resultArray := make([]interface{}, 0)
		for _, arrVal := range data.([]interface{}) {
			normalized, err := normalizeRecurse("", arrVal)
			if err != nil {
				return nil, err
			}
			// Only add normalized value if it wasn't culled
			if normalized != nil {
				resultArray = append(resultArray, normalized)
			}
		}

		// sort arrays
		switch key {
		case "created-assets":
			sortHelper(resultArray, "index")
		case "assets":
			sortHelper(resultArray, "asset-id")
		case "created-apps":
			sortHelper(resultArray, "id")
		case "apps-local-state":
			sortHelper(resultArray, "id")
		}

		result = resultArray
	}

	// Omit empty things from normalized result.
	empty, err := isEmpty(result)
	if empty {
		return nil, err
	}
	return result, err
}
