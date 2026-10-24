# Master Database Entity-Relationship Diagram (ERD)

## 1. Core Architectural Strategy
The database architecture employs PostgreSQL 16+ as the single source of truth (ACID transactional consistency), augmented with:
- **UUIDv7**: Time-ordered, cluster-friendly 128-bit unique identifiers for all primary keys.
- **Partitioning**: Range partitioning by month for high-throughput append-only tables (`posts`, `comments`, `reactions`, `messages`, `audit_logs`).
- **Transactional Outbox Table**: For zero-loss, asynchronous CDC/Kafka event publishing.

---

## 2. Mermaid Entity-Relationship Diagram

```mermaid
erDiagram
    USERS ||--o{ IDENTITIES : has
    USERS ||--|| PROFILES : has
    USERS ||--o{ SESSIONS : establishes
    USERS ||--o{ RELATIONSHIPS : initiates
    USERS ||--o{ POSTS : authors
    USERS ||--o{ COMMENTS : writes
    USERS ||--o{ REACTIONS : emits
    USERS ||--o{ CONVERSATION_MEMBERS : participates
    USERS ||--o{ GROUP_MEMBERS : belongs_to
    USERS ||--o{ PAGE_MEMBERS : manages
    USERS ||--o{ MARKETPLACE_ITEMS : sells

    POSTS ||--o{ POST_ATTACHMENTS : contains
    POSTS ||--o{ COMMENTS : receives
    POSTS ||--o{ REACTIONS : aggregates
    POSTS }o--o| GROUPS : posted_in
    POSTS }o--o| PAGES : published_by

    CONVERSATIONS ||--o{ CONVERSATION_MEMBERS : groups
    CONVERSATIONS ||--o{ MESSAGES : contains
    MESSAGES ||--o{ MESSAGE_ATTACHMENTS : includes

    GROUPS ||--o{ GROUP_MEMBERS : contains
    PAGES ||--o{ PAGE_MEMBERS : contains

    USERS {
        uuid id PK
        varchar email UK
        varchar phone UK
        varchar password_hash
        varchar status
        boolean email_verified
        boolean phone_verified
        boolean mfa_enabled
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    PROFILES {
        uuid user_id PK,FK
        varchar username UK
        varchar first_name
        varchar last_name
        text bio
        text avatar_url
        text cover_url
        date date_of_birth
        varchar gender
        varchar privacy_level
        jsonb preferences
        timestamp updated_at
    }

    RELATIONSHIPS {
        uuid id PK
        uuid user_id FK
        uuid target_user_id FK
        varchar relationship_type
        varchar status
        timestamp created_at
        timestamp updated_at
    }

    POSTS {
        uuid id PK
        uuid author_id FK
        uuid group_id FK
        uuid page_id FK
        text content
        varchar privacy
        varchar post_type
        integer reaction_count
        integer comment_count
        integer share_count
        timestamp created_at
        timestamp updated_at
        timestamp deleted_at
    }

    COMMENTS {
        uuid id PK
        uuid post_id FK
        uuid author_id FK
        uuid parent_comment_id FK
        text content
        integer reaction_count
        timestamp created_at
        timestamp updated_at
    }

    REACTIONS {
        uuid id PK
        uuid target_id
        varchar target_type
        uuid user_id FK
        varchar reaction_type
        timestamp created_at
    }

    CONVERSATIONS {
        uuid id PK
        varchar conversation_type
        varchar title
        uuid created_by FK
        timestamp created_at
        timestamp updated_at
    }

    MESSAGES {
        uuid id PK
        uuid conversation_id FK
        uuid sender_id FK
        text body
        varchar message_type
        boolean is_deleted
        timestamp created_at
    }

    OUTBOX_EVENTS {
        uuid id PK
        varchar aggregate_type
        uuid aggregate_id
        varchar topic
        text payload_json
        varchar status
        timestamp created_at
        timestamp published_at
    }
```
