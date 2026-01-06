# Compatibility Matrix

This document tracks compatibility between mythic-sdk-go versions and Mythic C2 server versions.

## Supported Versions

| SDK Version | Mythic Version | Status        | Release Date | End of Support |
|-------------|----------------|---------------|--------------|----------------|
| v0.1.0      | 3.3.0+         | In Development | TBD         | -              |

## Version Support Policy

- **Active Support**: Bug fixes, security updates, and new features
- **Maintenance Support**: Critical security fixes only
- **End of Life (EOL)**: No updates

We recommend always using the latest SDK version with the latest compatible Mythic version.

## Feature Compatibility

### SDK v0.1.0 (Target: Mythic 3.3.0+)

**Supported Features**:
- ✅ GraphQL API (Mythic 3.0+)
- ✅ API Token Authentication
- ✅ Username/Password Authentication
- ✅ Callback Management
- ✅ Task Operations
- ✅ File Operations
- ✅ Payload Generation
- ✅ Operator Management
- ✅ WebSocket Subscriptions
- ✅ Real-time Updates
- ✅ Process Enumeration
- ✅ File Browser
- ✅ Credential Management
- ✅ Analytics

**Breaking Changes from Mythic 2.x**:
- ❌ REST API removed (Mythic 3.0+)
- ✅ GraphQL only

## Testing Against Mythic Versions

Our CI/CD pipeline automatically tests against:
- **Latest Stable**: Most recent Mythic release
- **Master Branch**: Bleeding edge Mythic development

## Reporting Compatibility Issues

If you encounter compatibility issues:
1. Check this document for known issues
2. Verify you're using supported versions
3. Open an issue with:
   - SDK version
   - Mythic version
   - Error message
   - Steps to reproduce

## Mythic API Changes

We track Mythic API changes and update the SDK accordingly:

### Mythic 3.3.0 (Current Target)
- Full GraphQL API
- Hasura integration
- WebSocket subscriptions
- Enhanced OPSEC checks

### Mythic 3.0.0
- Complete REST to GraphQL migration
- Breaking change: All REST endpoints removed

### Future Changes

We monitor the [Mythic repository](https://github.com/its-a-feature/Mythic) for upcoming changes and aim to provide SDK updates within 1 week of Mythic releases.

---

**Last Updated**: 2026-01-06
**Next Review**: Upon each SDK or Mythic release
