# RBAC & ABAC Security and Permission Matrix

## 1. Overview
This document defines the Role-Based Access Control (RBAC) and Attribute-Based Access Control (ABAC) matrix for the Facebook Functional Parity Platform. Every operation across Web, Mobile, API Gateway, and Microservices must enforce these access boundaries.

---

## 2. Defined Roles and Entities

| Role ID | Role Name | Scope & Description | Default Trust Tier |
|---|---|---|---|
| `ROLE_ANONYMOUS` | Anonymous Visitor | Unauthenticated visitors; restricted to public landing, auth, and public post links. | Tier 0 |
| `ROLE_USER_STANDARD` | Standard Verified User | Authenticated regular user with verified email or phone; standard social interactions. | Tier 1 |
| `ROLE_USER_MINOR` | Minor User (Age 13-17) | Authenticated minor; additional safeguards on privacy, messaging from strangers, ads targeting. | Tier 1 (Protected) |
| `ROLE_PROFILE_OWNER` | Profile Owner | User acting as administrator of their personal timeline and security settings. | Tier 1 (Contextual) |
| `ROLE_CREATOR_PRO` | Professional Creator | Users with Professional Mode enabled; access to monetized content and analytics. | Tier 2 |
| `ROLE_PAGE_OWNER` | Page Owner / Admin | Primary manager of Facebook Page asset; full financial, post, and staff management. | Tier 2 (Contextual) |
| `ROLE_PAGE_EDITOR` | Page Editor / Moderator | Publishing posts, answering messages, moderating comments; no ownership/financial control. | Tier 2 (Contextual) |
| `ROLE_GROUP_ADMIN` | Group Admin | Full authority over group settings, member approvals, rules, and moderation. | Tier 2 (Contextual) |
| `ROLE_GROUP_MOD` | Group Moderator | Can approve posts, remove comments, mute or kick members. | Tier 2 (Contextual) |
| `ROLE_MARKET_SELLER` | Marketplace Seller | Verified identity to list goods, interact with buyers, manage escrow transactions. | Tier 2 |
| `ROLE_ADVERTISER` | Direct Advertiser | Manages ad campaigns, billing methods, and creative assets. | Tier 2 |
| `ROLE_AGENCY_OPERATOR` | Ad Agency Partner | Manages assigned ad accounts with delegate permissions on behalf of clients. | Tier 2 (Contextual) |
| `ROLE_BUSINESS_OWNER` | Business Account Admin | Manages corporate assets, payment methods, domains, and partner permissions. | Tier 3 (Contextual) |
| `ROLE_DEV_PARTNER` | External Developer | OAuth2 API client developer with app credentials, webhooks, and Graph API access. | Tier 2 |
| `ROLE_SUPPORT_AGENT` | Tier 1/2 Customer Support | Read-only audit access to user tickets and account state upon explicit consent. | Tier 3 (Staff) |
| `ROLE_MODERATOR_SAFETY` | Trust & Safety Moderator | Reviews reported content, applies moderation tags, executes strikes/bans. | Tier 3 (Staff) |
| `ROLE_BILLING_SPECIALIST` | Finance & Payments Admin | Manages chargebacks, payout approvals, refund disputes, and reconciliation. | Tier 4 (Staff) |
| `ROLE_PLATFORM_ADMIN` | Platform Superadmin | System-level configuration, feature flags, global incident response. | Tier 4 (Restricted) |
| `ROLE_SRE_OPERATOR` | Infrastructure / SRE | Production deployment, secret rotation, system telemetry. No user content access. | Tier 4 (Restricted) |

---

## 3. RBAC Operation Permission Matrix

| Resource / Operation | Anonymous | Standard User | Minor User | Group/Page Admin | Seller/Adv | Support/Mod Staff | Platform Admin |
|---|---|---|---|---|---|---|---|
| `auth:register` | ALLOW | DENY | DENY | DENY | DENY | DENY | DENY |
| `auth:login` | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW |
| `auth:mfa:verify` | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW |
| `profile:read:public` | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW | ALLOW |
| `profile:read:private` | DENY | OWNER_ONLY | OWNER_ONLY | OWNER_ONLY | OWNER_ONLY | AUDIT_ONLY | AUDIT_ONLY |
| `post:create` | DENY | ALLOW | ALLOW | ALLOW | ALLOW | DENY | ALLOW |
| `post:delete` | DENY | OWNER_ONLY | OWNER_ONLY | ADMIN_IN_GROUP | OWNER_ONLY | MOD_ONLY | ALLOW |
| `comment:create` | DENY | ALLOW | RESTRICTED | ALLOW | ALLOW | DENY | ALLOW |
| `message:send:direct` | DENY | ALLOW | RESTRICTED | ALLOW | ALLOW | DENY | DENY |
| `media:upload` | DENY | ALLOW | ALLOW | ALLOW | ALLOW | DENY | ALLOW |
| `group:create` | DENY | ALLOW | RESTRICTED | ALLOW | ALLOW | DENY | ALLOW |
| `marketplace:create_listing`| DENY | ALLOW (18+) | DENY | ALLOW | ALLOW | DENY | ALLOW |
| `ads:create_campaign` | DENY | ALLOW (18+) | DENY | ALLOW | ALLOW | DENY | ALLOW |
| `business:manage_assets` | DENY | ALLOW (18+) | DENY | ALLOW | ALLOW | DENY | ALLOW |
| `moderation:take_action` | DENY | DENY | DENY | GROUP_SCOPE | DENY | ALLOW | ALLOW |
| `admin:access_console` | DENY | DENY | DENY | DENY | DENY | ALLOW (Staff) | ALLOW |
| `system:manage_flags` | DENY | DENY | DENY | DENY | DENY | DENY | ALLOW |

---

## 4. Attribute-Based Access Control (ABAC) Rules

In addition to role-based privileges, every state mutation and sensitive read evaluates fine-grained attributes:

```yaml
ABAC_Rules:
  - id: ABAC-RULE-001-FRIENDSHIP-VISIBILITY
    description: Posts marked 'FRIENDS_ONLY' require an active friendship relationship between author and requester.
    attributes:
      requester.id: "!author.id"
      post.privacy: "FRIENDS_ONLY"
      relationship(requester.id, author.id): "ACCEPTED"
      block_status(requester.id, author.id): "NONE"
    evaluation: ALLOW

  - id: ABAC-RULE-002-BLOCK-MASKING
    description: If either user has blocked the other, all queries return 404 Not Found (bidirectional invisibility).
    attributes:
      blocked_pair: "is_blocked(requester.id, target.id) == true"
    evaluation: DENY_IMMEDIATE_404

  - id: ABAC-RULE-003-MINOR-PROTECTION
    description: Adults cannot initiate unsolicited direct messages with unrelated minors (<18).
    attributes:
      requester.age: ">= 18"
      target.age: "< 18"
      relationship(requester.id, target.id): "!= ACCEPTED"
    evaluation: DENY_FORBIDDEN

  - id: ABAC-RULE-004-PAGE-DELEGATION
    description: Access to Page actions is evaluated based on the business asset assignment record.
    attributes:
      asset_type: "PAGE"
      requester.role_in_asset: "IN ['ADMIN', 'EDITOR', 'MODERATOR']"
      business.status: "ACTIVE"
    evaluation: ALLOW
```

---

## 5. Security Enforcement Architecture
1. **API Gateway (Edge)**: Verifies JWT signature, role claims, and IP blocklist. Rejects unauthenticated requests to protected endpoints.
2. **Go Middleware (RPC/REST)**: Evaluates role scopes and decodes context claims (`user_id`, `roles`, `tenant_id`, `device_id`).
3. **Domain Service (ABAC Engine)**: Queries graph relations in Redis cache or PostgreSQL read replicas to execute relationship-bound rules.
4. **Database (Row-Level Security / Tenant Isolation)**: Tenant and asset ID isolation enforced via partitioned schemas and query parameters.
