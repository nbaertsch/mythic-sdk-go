#!/usr/bin/env python3
"""Fix NewClient to detect JWT tokens in APIToken config and reclassify them as AccessToken."""

with open("pkg/mythic/client.go", "r") as f:
    content = f.read()

# 1. Add "strings" to imports
old_imports = '''import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"sync"

	"github.com/hasura/go-graphql-client"
)'''

new_imports = '''import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"sync"

	"github.com/hasura/go-graphql-client"
)'''

if old_imports not in content:
    print("ERROR: Could not find imports block")
    exit(1)
content = content.replace(old_imports, new_imports)

# 2. Add JWT detection after config validation, before cookie jar creation
old_section = '''\tif err := config.Validate(); err != nil {
		return nil, WrapError("NewClient", err, "invalid configuration")
	}

	// Create cookie jar for session management'''

new_section = '''\tif err := config.Validate(); err != nil {
		return nil, WrapError("NewClient", err, "invalid configuration")
	}

	// Detect if the APIToken is actually a JWT (starts with "eyJ" — the
	// base64 encoding of '{"'). JWTs must be sent via "Authorization:
	// Bearer" header, not the "apitoken" header which is reserved for
	// Mythic's long-lived API tokens stored in the database.
	if config.APIToken != "" && strings.HasPrefix(config.APIToken, "eyJ") {
		config.AccessToken = config.APIToken
		config.APIToken = ""
	}

	// Create cookie jar for session management'''

if old_section not in content:
    print("ERROR: Could not find config.Validate section")
    exit(1)
content = content.replace(old_section, new_section)

with open("pkg/mythic/client.go", "w") as f:
    f.write(content)

print("OK: NewClient now detects JWT-as-APIToken and reclassifies it")
