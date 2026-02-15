#!/usr/bin/env python3
"""Fix getMeIdentity to handle JWT-as-APIToken fallback."""

with open("pkg/mythic/auth.go", "r") as f:
    content = f.read()

# Find the function boundaries
start_marker = "// getMeIdentity calls the /me REST endpoint to extract the user_id from\n// the current JWT or API token claims.\nfunc (c *Client) getMeIdentity(ctx context.Context) (int, error) {"
end_marker = "\treturn meResp.UserID, nil\n}"

start_idx = content.find(start_marker)
if start_idx == -1:
    print("ERROR: Could not find start of getMeIdentity")
    exit(1)

# Find the closing brace of the function
# We need to find the end of this specific function
search_from = start_idx + len(start_marker)
end_idx = content.find(end_marker, search_from)
if end_idx == -1:
    print("ERROR: Could not find end of getMeIdentity")
    exit(1)

end_idx += len(end_marker)

old_func = content[start_idx:end_idx]
print(f"Found function ({len(old_func)} chars)")

new_func = """// getMeIdentity calls the /me REST endpoint to extract the user_id from
// the current JWT or API token claims.
//
// The /me endpoint accepts both "apitoken: <token>" (for real Mythic API
// tokens) and "Authorization: Bearer <jwt>" headers. When the caller has
// set APIToken in the config this might actually be a JWT (e.g. passed via
// MYTHIC_API_TOKEN env var), so if the first attempt returns 401 we retry
// with Authorization: Bearer as a fallback.
func (c *Client) getMeIdentity(ctx context.Context) (int, error) {
\tscheme := "https"
\tif !c.config.SSL {
\t\tscheme = "http"
\t}
\tmeURL := fmt.Sprintf("%s://%s/me", scheme, stripScheme(c.config.ServerURL))

\t// Build the list of header sets to try.
\t// Primary: whatever getAuthHeaders() returns (apitoken or Bearer).
\t// Fallback: if APIToken is set and first attempt gets 401, try the
\t// token as a Bearer JWT in case the caller supplied a JWT instead of
\t// a real Mythic API token.
\ttype headerSet = map[string]string
\tattempts := []headerSet{c.getAuthHeaders()}
\tif c.config.APIToken != "" {
\t\tattempts = append(attempts, headerSet{
\t\t\t"Authorization": "Bearer " + c.config.APIToken,
\t\t})
\t}

\tvar lastStatus int
\tvar lastBody []byte

\tfor _, headers := range attempts {
\t\treq, err := http.NewRequestWithContext(ctx, "GET", meURL, nil)
\t\tif err != nil {
\t\t\treturn 0, WrapError("getMeIdentity", err, "failed to create /me request")
\t\t}
\t\tfor k, v := range headers {
\t\t\treq.Header.Set(k, v)
\t\t}

\t\tresp, err := c.httpClient.Do(req)
\t\tif err != nil {
\t\t\treturn 0, WrapError("getMeIdentity", err, "failed to call /me endpoint")
\t\t}

\t\tbody, readErr := io.ReadAll(resp.Body)
\t\tresp.Body.Close()
\t\tif readErr != nil {
\t\t\treturn 0, WrapError("getMeIdentity", readErr, "failed to read /me response")
\t\t}

\t\tif resp.StatusCode == http.StatusOK {
\t\t\tvar meResp struct {
\t\t\t\tUserID int `json:"user_id"`
\t\t\t}
\t\t\tif err := json.Unmarshal(body, &meResp); err != nil {
\t\t\t\treturn 0, WrapError("getMeIdentity", err, "failed to parse /me response")
\t\t\t}
\t\t\tif meResp.UserID == 0 {
\t\t\t\treturn 0, WrapError("getMeIdentity", ErrInvalidResponse, "no user_id in /me response")
\t\t\t}
\t\t\treturn meResp.UserID, nil
\t\t}

\t\tlastStatus = resp.StatusCode
\t\tlastBody = body
\t}

\treturn 0, WrapError("getMeIdentity", ErrAuthenticationFailed,
\t\tfmt.Sprintf("/me returned status %d: %s", lastStatus, string(lastBody)))
}"""

content = content[:start_idx] + new_func + content[end_idx:]

with open("pkg/mythic/auth.go", "w") as f:
    f.write(content)

print("OK: getMeIdentity updated with Bearer fallback")
