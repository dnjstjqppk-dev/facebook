# Comprehensive Threat Model & Defense Architecture

## 1. STRIDE Threat Analysis Matrix

| Threat Category | Specific Attack Vector | Potential Impact | Automated Preventive Control | Verification / Test Case |
|---|---|---|---|---|
| **Spoofing** | Credential stuffing / Brute force login | Account takeover | Distributed rate limiting in Redis (5 failed attempts / IP / min); mandatory MFA challenge. | `tests/security/brute_force_test.go` |
| **Tampering** | Parameter manipulation / IDOR (BOLA) | Unauthorized post deletion, viewing private messages | ABAC policy engine checks contextual ownership on every entity mutation. Direct object IDs are unpredictable UUIDv7. | `tests/security/idor_test.go` |
| **Repudiation** | Denying payment / admin action | Financial fraud, rogue staff changes | Immutable Append-Only Audit Log in PostgreSQL with HMAC signature. | `tests/security/audit_trail_test.go` |
| **Information Disclosure** | Exposing email / phone via profile API | Privacy breach, mass scraping | Strict GraphQL / REST DTO projection. Field-level privacy filters remove unshared attributes. | `tests/security/privacy_projection_test.go` |
| **Denial of Service** | Feed query bomb / Media upload exhaustion | Platform outage, storage exhaustion | Leaky bucket rate limiter in API Gateway; client-side chunking; strict max file size (100MB video, 10MB image). | `tests/security/rate_limit_test.go` |
| **Elevation of Privilege** | Normal user calling `/admin/*` routes | Full platform compromise | API Gateway enforces separate JWT audience claim (`aud: admin_console`) + mTLS on admin routes. | `tests/security/admin_privilege_test.go` |

---

## 2. Specific Attack Mitigations

### 2.1 Cross-Site Scripting (XSS) & Content Injection
- Frontend utilizes Vue 3 reactive text interpolation with automatic HTML escaping.
- HTML rendering (`v-html`) strictly prohibited except for sanitized rich text via DOMPurify with strict whitelist tags (`<b>`, `<i>`, `<a>`, `<p>`).
- Strict HTTP Headers enforced at Gateway:
  ```http
  Content-Security-Policy: default-src 'self'; img-src 'self' data: https://cdn.social.internal; script-src 'self'; frame-ancestors 'none';
  X-Content-Type-Options: nosniff
  X-Frame-Options: DENY
  Referrer-Policy: strict-origin-when-cross-origin
  ```

### 2.2 Broken Object-Level Authorization (BOLA / IDOR)
Every controller follows a strict verification pattern:
```go
// Enforce contextual ownership
post, err := repo.GetPostByID(ctx, postID)
if err != nil {
    return nil, ErrNotFound
}
if post.AuthorID != session.UserID && !session.HasRole("ROLE_PLATFORM_ADMIN") {
    // If inside a group, verify if user is group admin
    if post.GroupID != nil && !groupService.IsGroupAdmin(ctx, *post.GroupID, session.UserID) {
        return nil, ErrForbidden
    }
}
```

### 2.3 Malicious File Upload Attacks
1. Content-Type sniffing rejected (`http.DetectContentType` matches against declared MIME).
2. Files run through FFmpeg / libvips to sanitize EXIF GPS data and re-encode buffers before CDN placement.
3. Media served exclusively from isolated subdomains (`https://cdn.social.internal`) without execute permissions (`Content-Disposition: inline`).

### 2.4 CSRF & Webhook Replay Mitigation
- SameSite=Lax/Strict cookie configuration with HTTPS `Secure` flag.
- Custom HTTP header requirement (`X-Requested-With` or `Authorization: Bearer`) for all state-changing API endpoints.
- Webhook endpoints enforce HMAC-SHA256 signature verification with a 5-minute timestamp tolerance window.
