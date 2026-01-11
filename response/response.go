package response

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"
)

// JSON sends a JSON response
func JSON(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// If data is already JSON bytes, write directly
	if jsonBytes, ok := data.([]byte); ok {
		if _, err := w.Write(jsonBytes); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	// Otherwise, encode the data as JSON
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// Error sends an error response
func Error(w http.ResponseWriter, message string, statusCode int) {
	errorResponse := map[string]any{
		"error":  message,
		"status": statusCode,
	}

	JSON(w, errorResponse, statusCode)
}

// Redirect sends a redirect response
func Redirect(w http.ResponseWriter, url string) {
	w.Header().Set("Location", url)
	w.WriteHeader(http.StatusFound)
}

// AppendCT appends __ct to the url
func AppendCT(url string, oDic map[string]string) string {
	if strings.Contains(url, "?") {
		return url + "&__ct=" + oDic["__ct"]
	}
	return url + "?__ct=" + oDic["__ct"]
}

// RedirectWithPost sends an HTML form that auto-submits with POST method
func RedirectWithPost(w http.ResponseWriter, url string, formData map[string]string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Build form fields
	var fields string
	for name, value := range formData {
		fields += fmt.Sprintf(`<input type="hidden" name="%s" value="%s">`,
			html.EscapeString(name),
			html.EscapeString(value))
	}

	// HTML template with auto-submit form
	htmlContent := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <title>Redirecting...</title>
</head>
<body>
    <form id="redirectForm" method="POST" action="%s">
        %s
    </form>
    <script>
        document.getElementById('redirectForm').submit();
    </script>
    <noscript>
        <p>JavaScript is disabled. Please click the button below to continue.</p>
        <button type="submit" form="redirectForm">Continue</button>
    </noscript>
</body>
</html>`, html.EscapeString(url), fields)

	w.Write([]byte(htmlContent))
}
