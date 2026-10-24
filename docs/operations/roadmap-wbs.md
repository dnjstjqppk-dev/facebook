# Master 100-Step Work Breakdown Structure (WBS) & Engineering Roadmap

## 1. Roadmap Architecture: 10 Phases x 10 Steps = 100 Increments

```text
PHASE 1: Product Scope, Threat Modeling & Architecture Blueprint (Steps 01-10) [COMPLETED]
PHASE 2: Monorepo, Infrastructure & Foundational Scaffolding (Steps 11-20)
PHASE 3: Database Schemas, Transactional Outbox & Core Protocols (Steps 21-30)
PHASE 4: Identity, Profile, Privacy & Social Graph Services (Steps 31-40)
PHASE 5: Content Creation, Transcoding & Media Delivery Pipeline (Steps 41-50)
PHASE 6: News Feed Ranking, Realtime Notifications & Search (Steps 51-60)
PHASE 7: Messenger, Realtime Chat, WebRTC Calling & Channels (Steps 61-70)
PHASE 8: Communities, Groups, Pages & Event Management (Steps 71-80)
PHASE 9: Marketplace, Commerce, Ads Engine & Creator Economy (Steps 81-90)
PHASE 10: Trust & Safety, Admin Console, E2E Verification & Release Gates (Steps 91-100)
```

---

## 2. Detailed Work Breakdown Structure (Steps 01 - 30)

| Step ID | Task Name | Responsible Subsystem | Dependencies | Key Deliverable | Priority |
|---|---|---|---|---|---|
| `TSK-01` | Feature Inventory Catalog | `docs/` | None | `docs/facebook-feature-inventory.csv` | P0 |
| `TSK-02` | RBAC & ABAC Matrix Definition | `docs/security/` | TSK-01 | `docs/security/rbac-abac-matrix.md` | P0 |
| `TSK-03` | Screen Registry & Route Directory | `docs/product/` | TSK-01 | `docs/product/screen-registry.yaml` | P0 |
| `TSK-04` | User Action & Event Processing Map | `docs/architecture/` | TSK-01, 02 | `docs/architecture/user-actions-events.md` | P0 |
| `TSK-05` | Feature Parity Tracking Matrix | `docs/product/` | TSK-01 | `docs/product/feature-parity-matrix.md` | P0 |
| `TSK-06` | Regional & Age Compliance Policies | `docs/product/` | TSK-01 | `docs/product/compliance-policies.md` | P1 |
| `TSK-07` | Nova UI Design System Specification| `docs/architecture/`| TSK-03 | `docs/architecture/design-system-spec.md` | P1 |
| `TSK-08` | Privacy & GDPR Directive | `docs/security/` | TSK-02, 06 | `docs/security/privacy-gdpr-policy.md` | P0 |
| `TSK-09` | Security Threat Model & Mitigations| `docs/security/` | TSK-02, 04 | `docs/security/threat-model.md` | P0 |
| `TSK-10` | Engineering Roadmap & WBS Spec | `docs/operations/` | TSK-01-09 | `docs/operations/roadmap-wbs.md` | P0 |
| `TSK-11` | Monorepo Setup (pnpm, Turborepo, Go) | Root config | TSK-10 | `package.json`, `pnpm-workspace.yaml`, `go.work`| P0 |
| `TSK-12` | User Web Vue 3 / Nuxt Scaffolding | `apps/user-web` | TSK-07, 11 | Nuxt 3 app skeleton, Pinia, Tailwind | P0 |
| `TSK-13` | Admin Web Console Foundation | `apps/admin-console`| TSK-02, 11 | Vue 3 / Vite admin portal skeleton | P1 |
| `TSK-14` | Business & Ads Center Scaffolding | `apps/business-center`| TSK-11 | Business portal skeleton | P1 |
| `TSK-15` | Flutter Cross-Platform Baseline | `apps/mobile-flutter`| TSK-03, 11 | Flutter project baseline with clean architecture | P1 |
| `TSK-16` | Go API Gateway Foundation | `services/api-gateway`| TSK-04, 11 | Go HTTP proxy, CORS, JWT middleware, rate limiter| P0 |
| `TSK-17` | Go Identity & Auth Microservice | `services/identity` | TSK-02, 16 | User registration, login, bcrypt, JWT rotation | P0 |
| `TSK-18` | Shared Domain Microservices Proto | `proto/` | TSK-04, 17 | gRPC definitions for Identity, Social, Feed | P0 |
| `TSK-19` | Local Development Docker Compose | `infrastructure/docker`| TSK-16, 17 | PostgreSQL, Redis, Kafka, OpenSearch compose | P0 |
| `TSK-20` | CI/CD GitHub Actions Pipeline | `.github/workflows`| TSK-11 | Automated build & test workflow | P0 |
| `TSK-21` | PostgreSQL Master Schema & Migrations| `database/` | TSK-01, 04 | `001_initial_schema.sql` (Tables & Indices) | P0 |
| `TSK-22` | Transactional Outbox Pattern Pipeline| `services/common` | TSK-21 | Outbox table + Go publisher worker | P0 |

---

## 3. Risk Register & Mitigation Strategy

1. **Risk R-01 (Fan-out Write Amplification on High-Follower Celebrities)**:
   - *Impact*: Redis memory saturation and worker lag during high-profile posts.
   - *Mitigation*: Hybrid Fan-out model. Standard users (<5k friends) use write fan-out; VIP/celebrity creators use pull fan-out dynamically merged in the read path.
2. **Risk R-02 (Distributed Event Ordering & Race Conditions in Chat)**:
   - *Impact*: Out-of-order message delivery in multi-device messenger sync.
   - *Mitigation*: Snowflake / UUIDv7 millisecond time-ordered IDs + Kafka partition key pegged to `conversation_id`.
3. **Risk R-03 (Asset Copyright & Trademark Invalidation)**:
   - *Impact*: Legal infringement claims.
   - *Mitigation*: 100% bespoke design tokens, custom SVG icons, distinct brand color palette (Nova UI), and independent source code implementation.
