package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nbaertsch/mythic-sdk-go/pkg/mythic"
)

func main() {
	// Create client
	client, err := mythic.NewClient(&mythic.Config{
		ServerURL:     os.Getenv("MYTHIC_URL"),
		Username:      os.Getenv("MYTHIC_USERNAME"),
		Password:      os.Getenv("MYTHIC_PASSWORD"),
		SSL:           true,
		SkipTLSVerify: true,
		Timeout:       30 * time.Second,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create client: %v\n", err)
		os.Exit(1)
	}

	// Authenticate
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := client.Login(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to authenticate: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("=== Authenticated successfully ===\n")

	// Introspect types
	types := []string{
		"c2profile",
		"taskartifact",
		"operation",
		"payload",
		"operator",
		"tag",
		"tagtype",
	}

	for _, typeName := range types {
		fmt.Printf("\n=== Type: %s ===\n", typeName)

		query := fmt.Sprintf(`{
			__type(name: "%s") {
				name
				fields {
					name
					type {
						name
						kind
					}
				}
			}
		}`, typeName)

		var result map[string]interface{}
		if err := client.RawGraphQLQuery(ctx, query, &result); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to introspect %s: %v\n", typeName, err)
			continue
		}

		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	}

	// Introspect mutations
	fmt.Println("\n\n=== Mutations ===")
	mutationsQuery := `{
		__schema {
			mutationType {
				fields {
					name
					args {
						name
						type {
							name
							kind
							ofType {
								name
								kind
							}
						}
					}
				}
			}
		}
	}`

	var mutationsResult map[string]interface{}
	if err := client.RawGraphQLQuery(ctx, mutationsQuery, &mutationsResult); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to introspect mutations: %v\n", err)
	} else {
		data, _ := json.MarshalIndent(mutationsResult, "", "  ")
		fmt.Println(string(data))
	}
}
