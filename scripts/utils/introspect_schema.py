#!/usr/bin/env python3
"""
GraphQL Schema Introspection Tool for Mythic

This script introspects the Mythic GraphQL schema to discover available types and fields.
Useful for debugging and ensuring SDK queries match the actual Mythic API.

Usage:
    python introspect_schema.py [--type TYPE_NAME] [--url URL] [--username USER] [--password PASS]

Examples:
    # Introspect callback type
    python introspect_schema.py --type callback

    # Introspect operator type
    python introspect_schema.py --type operator

    # Introspect with custom credentials
    python introspect_schema.py --type callback --url https://mythic.example.com:7443
"""

import argparse
import json
import os
import requests
import sys
from urllib3.exceptions import InsecureRequestWarning

# Suppress SSL warnings for self-signed certificates
requests.packages.urllib3.disable_warnings(category=InsecureRequestWarning)


def login(url, username, password):
    """Authenticate with Mythic and return access token."""
    auth_url = f"{url}/auth"
    payload = {"username": username, "password": password}

    try:
        response = requests.post(auth_url, json=payload, verify=False)
        response.raise_for_status()
        data = response.json()
        return data.get("access_token")
    except Exception as e:
        print(f"Error: Failed to login: {e}", file=sys.stderr)
        sys.exit(1)


def introspect_type(graphql_url, token, type_name):
    """Introspect a specific GraphQL type."""
    query = """
    query IntrospectType($typeName: String!) {
        __type(name: $typeName) {
            name
            kind
            fields {
                name
                type {
                    name
                    kind
                    ofType {
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
    """

    headers = {
        "Content-Type": "application/json",
        "Authorization": f"Bearer {token}"
    }

    payload = {
        "query": query,
        "variables": {"typeName": type_name}
    }

    try:
        response = requests.post(graphql_url, json=payload, headers=headers, verify=False)
        response.raise_for_status()
        return response.json()
    except Exception as e:
        print(f"Error: Failed to introspect schema: {e}", file=sys.stderr)
        sys.exit(1)


def format_type_info(type_info):
    """Format type information for display."""
    if not type_info:
        return "unknown"

    if type_info.get("kind") == "NON_NULL":
        inner = format_type_info(type_info.get("ofType"))
        return f"{inner}!"
    elif type_info.get("kind") == "LIST":
        inner = format_type_info(type_info.get("ofType"))
        return f"[{inner}]"
    else:
        return type_info.get("name", "unknown")


def print_type_info(data):
    """Print formatted type information."""
    if "errors" in data:
        print("GraphQL Errors:", file=sys.stderr)
        for error in data["errors"]:
            print(f"  - {error.get('message')}", file=sys.stderr)
        sys.exit(1)

    type_data = data.get("data", {}).get("__type")
    if not type_data:
        print("Error: Type not found in schema", file=sys.stderr)
        sys.exit(1)

    print(f"\nType: {type_data['name']}")
    print(f"Kind: {type_data['kind']}")
    print(f"\nFields ({len(type_data.get('fields', []))}):")
    print("=" * 80)

    for field in sorted(type_data.get("fields", []), key=lambda x: x["name"]):
        field_name = field["name"]
        field_type = format_type_info(field["type"])
        print(f"  {field_name:<30} {field_type}")


def main():
    parser = argparse.ArgumentParser(
        description="Introspect Mythic GraphQL schema",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog=__doc__
    )
    parser.add_argument(
        "--type", "-t",
        default="callback",
        help="GraphQL type to introspect (default: callback)"
    )
    parser.add_argument(
        "--url",
        default=os.getenv("MYTHIC_URL", "https://127.0.0.1:7443"),
        help="Mythic server URL (default: $MYTHIC_URL or https://127.0.0.1:7443)"
    )
    parser.add_argument(
        "--username", "-u",
        default=os.getenv("MYTHIC_USERNAME", "mythic_admin"),
        help="Mythic username (default: $MYTHIC_USERNAME or mythic_admin)"
    )
    parser.add_argument(
        "--password", "-p",
        default=os.getenv("MYTHIC_PASSWORD"),
        help="Mythic password (default: $MYTHIC_PASSWORD)"
    )
    parser.add_argument(
        "--graphql-url",
        help="Direct GraphQL URL (default: inferred from --url)"
    )

    args = parser.parse_args()

    if not args.password:
        print("Error: Password required (use --password or set MYTHIC_PASSWORD)", file=sys.stderr)
        sys.exit(1)

    # Determine GraphQL endpoint
    if args.graphql_url:
        graphql_url = args.graphql_url
    else:
        # Use Hasura direct endpoint (no auth required from localhost in default config)
        base_url = args.url.replace("https://", "http://").replace(":7443", ":8080")
        graphql_url = f"{base_url}/v1/graphql"

    print(f"Connecting to: {args.url}")
    print(f"GraphQL endpoint: {graphql_url}")
    print(f"Authenticating as: {args.username}")

    # Authenticate
    token = login(args.url, args.username, args.password)
    print(f"✓ Authentication successful\n")

    # Introspect type
    data = introspect_type(graphql_url, token, args.type)
    print_type_info(data)


if __name__ == "__main__":
    main()
