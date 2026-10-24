# Security Threat Model & Defense-in-Depth Specification

## 1. STRIDE 보안 위협 분석

| 위협 범주 | 잠재적 공격 시나리오 | 위험도 | 시스템 방어 대책 (Mitigation) |
|---|---|---|---|
| **Spoofing (신원 도용)** | - 무차별 대입(Brute-force) 로그인<br>- 세션 토큰 탈취 및 재전송 | HIGH | - IP 및 계정 기반 Exponential Rate Limiter (Redis)<br>- 비정상 지역 로그인 시 2FA 강제<br>- 단기 수명 JWT (15분) + Refresh Token Family Rotation |
| **Tampering (데이터 변조)** | - 게시물 작성자 ID 위조 (BOLA/IDOR)<br>- 결제 금액 조작 | CRITICAL | - 모든 백엔드 엔드포인트에서 컨텍스트 토큰의 `subject_id`와 리소스 소유자 검증<br>- 결제 금액은 클라이언트 파라미터가 아닌 서버 DB 카탈로그 기준으로만 결제 세션 생성 |
| **Repudiation (부인 방지)** | - 관리자/광고비 환불 승인 부인<br>- 악성 콘텐츠 게시 부인 | MEDIUM | - Append-only 불변 감사 로그(Audit Log) 체계 구축<br>- 관리자 작업 시 2FA 재인증 및 서명 기록 보관 |
| **Information Disclosure (정보 유출)** | - 비공개 게시물 URL 직접 열람<br>- API 응답에 민감 PII 노출 | HIGH | - S3 서명된 URL 만료시간(TTL) 60분 제한<br>- DTO 직렬화 시 `@Exclude` 어노테이션 적용하여 비밀번호/2FA 키 노출 원천 방지 |
| **Denial of Service (서비스 거부)** | - 대용량 피드 요청 폭주<br>- 초대형 동영상 반복 업로드 | HIGH | - API Gateway 레벨 분산 Rate Limit (100 req/sec per IP)<br>- Chunked Resumable Upload 시 단일 청크 5MB 및 총 파일 용량 제한 (동영상 4GB, 사진 25MB) |
| **Elevation of Privilege (권한 상승)** | - 일반 유저가 Super Admin API 호출<br>- 비즈니스 대행사의 타 고객사 자산 접근 | CRITICAL | - API Gateway 및 각 마이크로서비스 RPC 레벨에서 RBAC/ABAC 이중 인가 검사<br>- 테넌트 ID 격리 필터 적용 |

---

## 2. 웹 보안 헤더 및 콘텐츠 보안 정책 (CSP)

```http
Content-Security-Policy: default-src 'self'; script-src 'self' 'nonce-rAnd0m'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https://*.storage.ourdomain.com; media-src 'self' https://*.storage.ourdomain.com; frame-ancestors 'none';
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
Strict-Transport-Security: max-age=63072000; includeSubDomains; preload
Referrer-Policy: strict-origin-when-cross-origin
Cross-Origin-Opener-Policy: same-origin
Cross-Origin-Embedder-Policy: require-corp
```

---

## 3. 웹소켓 및 실시간 통신 보안

1. **핸드셰이크 인증:** WebSocket 연결 수립 시 Query string이 아닌 `Sec-WebSocket-Protocol` 또는 Cookie 기반의 짧은 일회용 인증 티켓(Ticket)을 발급받아 검증.
2. **룸 격리 (Room Boundary Enforcement):** 채팅방 ID 진입 시 RDBMS의 대화방 멤버십(`conversation_members`) 활성 여부를 캐시(Redis Set)에서 확인 후 불일치 시 즉시 연결 차단(Disconnect).
