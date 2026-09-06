# Session security is a set of behaviours, not a login form

The workspace in my portfolio is a private administrative application. It is not a public identity product with registration, account recovery or multiple tenants. That scope matters when explaining its security.

## Credentials and sessions have different jobs

Passwords are checked with bcrypt. The implementation validates length, including bcrypt's byte limit, rather than silently accepting a different password. On successful login, a cryptographically random token is sent in an HTTP-only cookie and a SHA-256 hash of that token is stored server-side.

A random bearer token is not a human-chosen password; its storage requirements are different. Production cookies use Secure and SameSite protections, but those settings do not make an unsafe deployment or cross-site scripting harmless.

## Expiry and logout must be real

The session service checks the persisted session, administrator status, idle timeout and absolute expiry. The browser cookie is not the sole authority for whether access continues.

Logout must revoke the stored token. Clearing a cookie while silently ignoring a database failure would claim more than the application accomplished. The reviewed implementation returns an error when revocation fails, and a regression test covers that path.

## Defensive controls have boundaries too

Request limits bound body reads; Journal limits bound stored content; cross-origin and token checks protect relevant requests. Login throttling has bounded entries so arbitrary credentials cannot grow its maps indefinitely.

That limiter is process-local. Multiple instances and reverse proxies require an explicit trusted-client-identity and shared-limiting strategy. Trusting arbitrary forwarded headers would not be an appropriate shortcut.

## Test the failures

Expired sessions, inactive accounts, missing CSRF tokens, oversized inputs and failed revocation are useful tests alongside successful login. Vulnerability scanning finds known affected code paths; a clean reachable-code result is not proof of universal security. Module advisories and deployed dependencies still need maintenance.

Explore the [application source](https://github.com/danieljmanningdev/danieljmanningdev-portfolio), [case study](https://danieljmanningdev.com/work/portfolio), and [software development services](https://danieljmanningdev.com/software-development/).

References: [OWASP Authentication](https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html) and [Session Management](https://cheatsheetseries.owasp.org/cheatsheets/Session_Management_Cheat_Sheet.html).
