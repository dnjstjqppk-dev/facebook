# Global Social Network System Architecture Overview

## 1. 하이레벨 아키텍처 다이어그램

```mermaid
flowchart TD
    subgraph Clients["Clients"]
        Web["Web Client (Vue 3 / Nuxt 3)"]
        Mobile["Mobile Client (Flutter iOS/Android)"]
        AdminWeb["Admin & Business Suite (Vue 3)"]
    end

    subgraph Edge["Edge & Gateway"]
        CF["Cloudflare CDN / WAF"]
        Gateway["Go API Gateway (HTTP/2, gRPC-Web, WebSocket)"]
    end

    subgraph CoreServices["Go Microservices Core"]
        IdentitySvc["Identity & Auth Service"]
        ProfileSvc["Profile & Social Graph Service"]
        FeedSvc["Feed & Ranking Service"]
        ContentSvc["Content & Post Service"]
        MediaSvc["Media & Streaming Service"]
        MessengerSvc["Messenger & Realtime Service"]
        MarketplaceSvc["Marketplace & Commerce Service"]
        AdsSvc["Ads & Auction Engine"]
        AdminSvc["Admin & Moderation Service"]
    end

    subgraph DataBackbone["Data & Messaging Infrastructure"]
        Postgres[(PostgreSQL 16 Multi-Master/Replica)]
        RedisCache[(Redis Cluster: Cache, RateLimit, PubSub)]
        KafkaCluster[(Apache Kafka: Event Streaming)]
        OpenSearchCluster[(OpenSearch: Universal Search Index)]
        ClickHouseDB[(ClickHouse: OLAP Metrics & Ad Telemetry)]
        ObjectStorage[(MinIO / S3: Photos, Videos, HLS Chunks)]
    end

    Web --> CF
    Mobile --> CF
    AdminWeb --> CF
    CF --> Gateway

    Gateway --> IdentitySvc
    Gateway --> ProfileSvc
    Gateway --> FeedSvc
    Gateway --> ContentSvc
    Gateway --> MediaSvc
    Gateway --> MessengerSvc
    Gateway --> MarketplaceSvc
    Gateway --> AdsSvc
    Gateway --> AdminSvc

    CoreServices --> Postgres
    CoreServices --> RedisCache
    CoreServices --> KafkaCluster
    CoreServices --> OpenSearchCluster
    CoreServices --> ClickHouseDB
    MediaSvc --> ObjectStorage
```

---

## 2. 서비스 통신 및 프로토콜 규격
- **클라이언트 <-> Gateway:**
  - RESTful HTTP/2 (JSON/Protobuf) for standard queries & mutations
  - Secure WebSocket (`wss://`) for Messenger, Presence, and Notifications
  - WebRTC (STUN/TURN) for 1:1 and SFU Video Calls
- **Gateway <-> 마이크로서비스:**
  - gRPC over HTTP/2 with mutual TLS (mTLS)
- **비동기 이벤트 메시징:**
  - Apache Kafka topics with schema registry for domain lifecycle events
