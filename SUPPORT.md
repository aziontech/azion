# Azion CLI Feature Support & Breaking Changes

This document provides a comprehensive overview of feature availability, breaking changes, and deprecated functionality across Azion CLI versions based on the official GitHub releases.

## :fire: Major Version Breaking Changes

### Version 4.0.0 - **MAJOR BREAKING RELEASE**

**Breaking Changes:**
- Complete migration to Azion API v4
- New command structure and organization
- Changed error handling standards across all commands
- V3 API endpoints moved to legacy support via `backwards compatibility` (v3 account is necessary to use v3 commands)

**Migration Impact:**
- Existing scripts using v3 API may need updates
- Manifest files require structure updates
- Error handling in integrations needs review



**:warning: Deprecated:**
- Direct V3 API usage (moved to backwards compatibility for v3 accounts)

## :arrows_counterclockwise: Migration Guidelines

### Upgrading to v4.x
1. **API Migration**: Update integrations from v3 to v4 API endpoints
2. **Azion.json Updates**: Convert azion.json to new v4 structure (a script is available at 'scripts/v3_to_v4_converter.sh')
3. **Azion.config Updates**: Convert azion.config to new v4 structure 
3. **Command Updates**: Review and update any automated scripts
4. **Error Handling**: Update error parsing for new error standards

## :bar_chart: Support Matrix

| Feature Category | v1.x | v2.x | v3.x | v4.x | Notes |
|------------------|------|------|------|------|-------|
| **Core CLI Commands** | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Stable across versions |
| **V3 API Support** | :white_check_mark: | :white_check_mark: | :white_check_mark: | :arrows_counterclockwise: | :arrows_counterclockwise:: Legacy support via backwards compatibility |
| **V4 API Support** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Application** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Functions** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Rule Engine** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Cache Settings** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Workload** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Workload Deployment** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Edge Connector** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **V4 Function Instance** | :x: | :x: | :x: | :white_check_mark: | New in v4.4.0 |
| **V4 Cache Warming** | :x: | :x: | :x: | :white_check_mark: | New in v4.0.0 |
| **Profile Management** | :x: | :x: | :x: | :white_check_mark: | v4.12.0+ |
| **Bundler 5.0.0** | :x: | :x: | :white_check_mark: | :white_check_mark: | Required from v3.0.0 |
| **Concurrent Uploads** | :x: | :x: | :x: | :white_check_mark: | Optimized in v4.x |

## :label: Version Recommendations

- **Production Use**: v4.11.0+ (Latest stable with all features)
- **Legacy Projects**: v3.6.0 (If v4 migration not feasible) | there's also the possibility of using the latest version with a v3 account, thus making use of the backwards compatibility
- **Minimum Supported**: v2.x.x (Consider upgrading for security and features - we cannot guarantee support for v2.x)
- **1.x and below**: Not supported

## :memo: Notes

- **Breaking Change Pattern**: Major versions (x.0.0) introduce breaking changes
- **Backward Compatibility**: Minor versions maintain backward compatibility
- **API Evolution**: V3 API support maintained for transition period

---

*Last Updated: Based on releases through v4.11.0*
*For the most current information, check the [official releases page](https://github.com/aziontech/azion/releases)*