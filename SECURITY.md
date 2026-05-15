# Security Policy

## Reporting a Vulnerability

If you discover a security vulnerability in this project, please report it responsibly.

**Do NOT open a public GitHub issue for security vulnerabilities.**

### How to Report

Send an email to **security@azion.com** with:

1. A description of the vulnerability
2. Steps to reproduce the issue
3. Potential impact assessment
4. Any suggested remediation (optional)

### What to Expect

- **Acknowledgment** within 2 business days
- **Initial assessment** within 5 business days
- **Resolution timeline** communicated after assessment
- **Credit** given to reporters (unless anonymity is requested)

### Scope

This policy applies to the latest release of the `azion` CLI and its direct dependencies.

## Supported Versions

| Version | Supported |
|---------|-----------|
| Latest  | Yes       |
| < Latest | No       |

## Security Measures

- Dependencies are monitored via Dependabot and govulncheck
- SBOM (Software Bill of Materials) is generated for each release
- Secrets scanning is enforced via gitleaks in CI
