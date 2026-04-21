# Security

## Reporting a vulnerability

Do **not** open a public issue for anything that could compromise users or deployments. Issues are public right away, which makes coordinated fixes harder.

Use **[private vulnerability reporting](https://github.com/trelay-dev/trelay/security)** on this repository instead: open the **Security** tab, then **Report a vulnerability**. That starts a private thread with maintainers so details stay off the public tracker until there is a fix.

We do not publish a separate security email. GitHub’s flow above is the right channel for this project.

If the **Report a vulnerability** option is missing, a maintainer may need to turn on private reporting under the repo’s **Settings → Security**.

## Scope

Reports we care about include authentication and session handling bugs, privilege escalation, unsafe redirects, injection or path issues in the API or server, and anything that leaks secrets or data between tenants on a shared instance. If you are unsure, report it anyway.

## After you report

We will acknowledge when we can, work on a fix or mitigation, and coordinate disclosure (for example via a security advisory) once a release is ready. Please give us a reasonable window before going public with exploit details.

## Supported versions

Security fixes are applied to the latest release line on `main`. Run a current build or image from this repository for production.
