package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const MaxPostFormSize = 1 << 20 // 1MB

// Enhanced helper that works with both path and query parameters
func GetParam(r *http.Request, param string) string {
	// Try path parameter first
	if pathValue := r.PathValue(param); pathValue != "" {
		return pathValue
	}
	// Fall back to query parameter
	return r.URL.Query().Get(param)
}

func GetPostParam(r *http.Request, param string) string {
	return r.PostFormValue(param)
}

func GetRequiredParam(r *http.Request, param string) (string, error) {
	value := GetParam(r, param)
	if value == "" {
		return "", fmt.Errorf("parameter '%s' is required", param)
	}
	return value, nil
}

// GetOptionalParam returns the value of the parameter from the request, first from the path, then from the query, then from the post
func GetOptionalParam(r *http.Request, param, defaultValue string) string {
	value := GetParam(r, param)
	if value == "" {
		return GetOptionalPostParam(r, param, defaultValue)
		// return defaultValue
	}
	return value
}

func GetIntParam(r *http.Request, param string, defaultValue int) int {
	value := GetParam(r, param)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

func GetRequiredIntPostParam(r *http.Request, param string) (int, error) {
	err := r.ParseForm()
	if err != nil {
		return 0, fmt.Errorf("failed to parse form: %v", err)
	}
	value := r.PostForm.Get(param)
	if value == "" {
		return 0, fmt.Errorf("POST parameter '%s' is required", param)
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue, nil
	}
	return 0, fmt.Errorf("POST parameter '%s' is required", param)
}

func GetRequiredPostParam(r *http.Request, param string) (string, error) {
	err := r.ParseForm()
	if err != nil {
		return "", fmt.Errorf("failed to parse form: %v", err)
	}
	value := r.PostForm.Get(param)
	if value == "" {
		return "", fmt.Errorf("POST parameter '%s' is required", param)
	}
	return value, nil
}

func GetOptionalPostParam(r *http.Request, param, defaultValue string) string {
	err := r.ParseForm()
	if err != nil {
		return defaultValue
	}
	value := r.PostForm.Get(param)
	if value == "" {
		return defaultValue
	}
	return value
}

func GetIntPostParam(r *http.Request, param string, defaultValue int) int {
	err := r.ParseForm()
	if err != nil {
		return defaultValue
	}
	value := r.PostForm.Get(param)
	if value == "" {
		return defaultValue
	}
	if intValue, err := strconv.Atoi(value); err == nil {
		return intValue
	}
	return defaultValue
}

// RequestParam returns the value of the parameter from the request, first from the path, then from the query, then from the post
func RequestParam(r *http.Request, param string) string {
	s := GetParam(r, param)
	if s == "" {
		s = GetPostParam(r, param)
	}
	return s
}

// DictFromRequest populates the dictionary with the parameters from the request
func DictFromRequest(r *http.Request, oDic map[string]string) {
	sPost := ""
	if r.Method == "POST" {
		for k, v := range r.PostForm {
			oDic[k] = v[0]
			sPost = sPost + k + "=" + v[0] + "&"
		}
	} else {
		for k, v := range r.URL.Query() {
			oDic[k] = v[0]
			sPost = sPost + k + "=" + v[0] + "&"
		}
	}

	// oDic["post_params"] = U.Left(sPost, len(sPost)-1)
}

// ConvertPostToQueryString converts POST form data to GET query string format
// This maintains compatibility with legacy backend code expecting GET parameters
func ConvertPostToQueryString(r *http.Request, fieldName string) string {
	if err := r.ParseForm(); err != nil {
		return ""
	}

	var params []string

	// Iterate through all form values

	for key, values := range r.Form {
		// For fields with multiple values (like checkboxes with same name)
		// add each value as a separate parameter
		for _, value := range values {
			if strings.EqualFold(key, fieldName) {
				params = append(params, fmt.Sprintf("%s=%s", key, value))
				// } else {
				// 	if value != "" {
				// 		params = append(params, fmt.Sprintf("%s=%s", key, value))
				// 	}
			}
		}
	}

	// "ID_to=1234,subject_group_id=9998,Operatorpic_logo=like,_pageTable=subject,OperatorTitle=like,ChkOption=ID,ChkOption=subject_group_id,mn_sig=e98ec16d1663fd16a2e9b3423d58b1f2,ID_from=1"

	return strings.Join(params, ",")
}

func IsParamExists(r *http.Request, key string) bool {
	// This MUST be called for POST/PUT/PATCH requests to populate r.Form
	// with body data. It also parses the URL query parameters.
	r.ParseForm()

	// r.Form is of type url.Values (map[string][]string).
	// Checking the existence of the key in this map covers both query
	// and form data.
	if _, ok := r.Form[key]; ok {
		return true
	}

	return false
}

// SetParam updates a parameter in the request
// It updates path params (if exists), query params, and POST form data
func SetParam(r *http.Request, param, value string) {
	// Update query parameter
	q := r.URL.Query()
	q.Set(param, value)
	r.URL.RawQuery = q.Encode()

	// Update POST form if it exists
	if r.PostForm != nil {
		r.PostForm.Set(param, value)
	}
	// Also update Form if it exists (combined query + post)
	if r.Form != nil {
		r.Form.Set(param, value)
	}
}

// SetParamInQuery updates only the query parameter
func SetParamInQuery(r *http.Request, param, value string) {
	q := r.URL.Query()
	q.Set(param, value)
	r.URL.RawQuery = q.Encode()
}

// SetParamInPost updates only the POST form parameter
func SetParamInPost(r *http.Request, param, value string) {
	// Parse form if not already parsed
	if r.PostForm == nil {
		r.ParseForm()
	}
	if r.PostForm != nil {
		r.PostForm.Set(param, value)
	}
	if r.Form != nil {
		r.Form.Set(param, value)
	}
}

// ParseJSON reads the body once and returns the parsed map
func ParseJSON(r *http.Request) (map[string]interface{}, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
