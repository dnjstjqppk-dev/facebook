# Privacy & Data Protection Architecture

## 1. 데이터 분류 체계 (Data Classification Tier)

| 등급 | 유형 | 대상 필드 | 암호화 및 보관 정책 |
|---|---|---|---|
| **Tier 1 (극비/인증)** | 비밀번호, 2FA 시크릿, 결제 토큰 | `password_hash` (Argon2id/Bcrypt), `totp_secret`, `stripe_customer_id` | 필드 레벨 AES-256-GCM 암호화, HSM/KMS 키 관리, 평문 로깅 절대 불가 |
| **Tier 2 (민감 PII)** | 실명, 생년월일, 전화번호, 이메일 | `real_name`, `birth_date`, `phone_number`, `email` | DB 컬럼 암호화, API 마스킹 전송 (예: `u***@example.com`) |
| **Tier 3 (소셜 데이터)** | 비공개 게시물, 1:1 메시지, 친구 목록 | `posts.visibility != PUBLIC`, `chat_messages` | 저장 시 암호화(At-Rest), RBAC/ABAC 인가 없이는 DB 조회 불가 |
| **Tier 4 (공개 데이터)** | 프로필 사진, 닉네임, 공개 게시글 | `avatar_url`, `display_name`, `public_posts` | CDN 캐싱 허용, OpenSearch 공개 검색 색인 |

---

## 2. 세부 프라이버시 가시성 모델 (Granular Visibility Matrix)

사용자는 모든 게시물, 프로필 정보, 친구 목록에 대해 5단계 가시성을 지정할 수 있습니다.
- `PUBLIC`: 전 세계 인터넷 방문자 및 검색 엔진 공개
- `FRIENDS`: 상호 친구 승인된 사용자에게만 노출
- `FRIENDS_EXCEPT`: 특정 사용자를 제외한 친구에게 노출
- `CUSTOM`: 지정된 특정 친구 또는 친구 목록에만 한정 노출
- `ONLY_ME`: 작성자 본인만 열람 가능

---

## 3. 개인정보 다운로드 (GDPR/PIPA DYI Flow)

1. 사용자가 웹/앱 설정에서 **[내 정보 다운로드]** 요청.
2. `services/identity`에서 비동기 내보내기 작업(`gdpr_export_tasks`) 생성 및 Kafka 발행.
3. 백그라운드 Worker가 분산 데이터베이스에서 해당 사용자의 모든 데이터(프로필, 게시물, 댓글, 미디어 링크, 메시지 기록)를 JSON 및 정적 HTML 형태로 패키징.
4. 패키지 ZIP 파일을 S3 호환 Object Storage에 업로드하고 단기 만료(48시간) 서명된 URL(Presigned URL) 생성.
5. 사용자 이메일 및 인앱 알림으로 보안 다운로드 링크 발송 (다운로드 시 2FA 재인증 요구).

---

## 4. 계정 탈퇴 및 단계별 영구 삭제 파기 파이프라인

1. **탈퇴 신청 (D-Day):** 즉시 세션 강제 종료, 프로필 상태를 `DEACTIVATED`로 전환. 검색 인덱스 및 피드 노출에서 즉시 제외.
2. **30일 유예 기간 (Grace Period):** 30일 이내 로그인 시 탈퇴 철회 및 계정 복구 지원.
3. **영구 파기 (D+30일 자정):**
   - RDBMS: 사용자 식별자 익명화 처리 (`users.deleted_at = NOW()`, PII 완전 마스킹 및 `is_purged = true`).
   - S3 Media: 사용자가 업로드한 원본 사진 및 동영상 삭제 태스크 실행.
   - Elastic/OpenSearch: 사용자 관련 도큐먼트 Hard Delete.
   - Redis: 캐시 세션 및 팔로우 관계 캐시 키 삭제.
   - 감사 로그(Audit Log)는 법정 의무 보관 기간(통신비밀보호법 3개월, 전자상거래법 5년) 동안 안전하게 분리 보관.
