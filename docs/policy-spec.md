# Regional, Age & Legal Policy Specifications

## 1. 개요 및 관할권별 준수 프레임워크
본 문서는 글로벌 소셜 네트워크 플랫폼 구축 시 적용되는 한국(KR), 미국(US), 유럽연합(EU)의 법률 및 규제 요구사항과 시스템적 격리(Feature Flagging) 정책을 정의합니다.

---

## 2. 지역별 규제 매트릭스

| 관할권 | 주요 법령 | 시스템 요구사항 | 기술적 구현 방안 |
|---|---|---|---|
| **대한민국 (KR)** | 개인정보보호법 (PIPA)<br>정보통신망법<br>전자상거래법<br>청소년보호법 | 1. 14세 미만 아동 법정대리인 동의 필수<br>2. 마케팅/선택정보 별도 동의<br>3. 통신판매중개자 신원표시 의무<br>4. 본인확인/실명인증 (휴대폰 OTP/KMC) | - `UserIdentity.Age < 14` 시 가입 차단 또는 대리인 인증 게이트<br>- Marketplace C2C 분쟁 해결 기구 연동<br>- 1년 이상 미이용자 분리보관 |
| **미국 (US)** | COPPA (아동 온라인 사생활 보호법)<br>CCPA/CPRA (캘리포니아 프라이버시법)<br>Section 230 | 1. 13세 미만 개인정보 수집 금지<br>2. Do Not Sell/Share My Info 옵션 제공<br>3. 투명성 리포트 및 게시물 면책 조항 준수 | - 만 13세 미만 계정 등록 원천 차단<br>- 타겟 광고 옵트아웃 토글 (`ad_personalization_opt_out = true`)<br>- Ad Library 공개 보관 7년 보장 |
| **유럽연합 (EU)** | GDPR (일반 데이터 보호 규정)<br>DSA (디지털 서비스법)<br>DMA (디지털 시장법) | 1. 망각권(Right to be Forgotten) 및 데이터 이동권<br>2. 미성년자 프로파일링/타겟 광고 전면 금지<br>3. 추천 알고리즘의 투명성 및 비개인화 피드 선택권<br>4. 불법 유해 콘텐츠 24시간 내 신속 조치 | - 원클릭 데이터 내보내기 (ZIP / JSON)<br>- 30일 유예 기간 후 완전 파기 Worker<br>- Chronological (시간 역순) 피드 강제 토글 버튼 제공<br>- DSA 인증 신뢰할 수 있는 신고자(Trusted Flaggers) 전용 큐 |

---

## 3. 연령별 접근 통제 정책 (ABAC)

1. **미성년자 보호 (Minor Protection):**
   - 만 18세 미만 사용자의 경우 Dating 기능(`FEAT-047`) 접근 불가.
   - 검색 인덱스에서 미성년자 프로필의 "전체 공개" 기본값 비활성화 ("친구의 친구" 까지만 노출).
   - 미성년자에게는 행동 기반 맞춤형 광고(Behavioral Retargeting) 서빙 차단.

2. **Dating 서비스 격리 (`FEAT-047`):**
   - 별도 약관 동의 및 만 18세 이상 성인 인증 통과 시에만 진입 허용.
   - 페이스북 기존 친구 목록과 데이팅 프로필의 상호 검색 및 노출 원천 차단.

---

## 4. Feature Flag 매핑 (`system_feature_flags`)

```yaml
features:
  feature_dating:
    enabled: true
    min_age: 18
    allowed_countries: ["KR", "US", "CA", "EU"]
  feature_marketplace_escrow:
    enabled: true
    allowed_countries: ["KR", "US"]
    kyc_required: true
  feature_ads_behavioral_targeting:
    enabled: true
    forbidden_roles: ["ROLE_MINOR_USER"]
    gdpr_consent_required: true
```
