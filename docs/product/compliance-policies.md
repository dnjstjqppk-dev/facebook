# Regional, Age & Regulatory Compliance Policy Matrix

## 1. Regulatory Contexts & Jurisdiction Scope
The platform implements jurisdiction-aware dynamic policy enforcement across three primary regional frameworks:
1. **Republic of Korea (KR)**: PIPA (개인정보보호법), Information and Communications Network Act, Juvenile Protection Act (청소년 보호법), Electronic Commerce Act (전자상거래법).
2. **United States (US)**: COPPA (Children's Online Privacy Protection Act), Section 230 CDA, CCPA/CPRA (California), state-level age verification laws.
3. **European Union (EU)**: GDPR (General Data Protection Regulation), DSA (Digital Services Act), DMA (Digital Markets Act).

---

## 2. Dynamic Feature Flag & Restriction Rules

| Policy Domain | Region / Criterion | Regulatory Requirement | System Enforcement / Feature Flag | Fallback / Block State |
|---|---|---|---|---|
| **Age Eligibility** | Global | Minimum age: 13 years | `FLAG_MIN_AGE_13_GATE` | Registration blocked if DOB < 13y |
| **Real-Name Verification**| KR (South Korea) | Telecom or PASS identity verification required for adult content and payment features | `FLAG_KR_REALNAME_VERIFY` | Sensitive feeds and checkout locked until CI (Connecting Information) confirmed |
| **Direct Messaging Safety**| Global (Minors <18) | Prevent unsolicited contact from unknown adults (>18) | `FLAG_MINOR_STRANGER_DM_SHIELD` | Message requests automatically rejected/hidden unless already mutual friends |
| **Ad Targeting Restrictions**| EU (GDPR) & Minor accounts | Zero behavioral / profiling targeting for minors; strict consent for adults | `FLAG_GDPR_CONSENT_STRICT`, `FLAG_NO_BEHAVIORAL_ADS_UNDER_18` | Only contextual ads (category, language) served to under-18 users |
| **Right to Erasure / Export**| EU (GDPR Art. 17/20) & CA (CCPA) | One-click machine-readable data export (JSON/ZIP) and full deletion within 30 days | `FLAG_GDPR_DATA_PORTABILITY`, `FLAG_PERMANENT_DELETION_WORKFLOW` | Async export pipeline packages media + text; tombstone record retained for audit |
| **Marketplace Escrow & ID**| KR & EU | Seller identification, tax registration number verification, escrow guarantee | `FLAG_MARKETPLACE_SELLER_KYC` | Selling restricted to buyer-pickup / offline or verified escrow transactions |
| **Dating & Romantic Matches**| Global | Strict 18+ enforcement, isolated profile from standard feed, geolocation opt-in | `FLAG_DATING_ENABLED` | Feature completely hidden from minor accounts and jurisdictions where prohibited |
| **Creator Monetization & Payouts**| US / KR / Global | W-9 / W-8BEN tax form collection, AML / OFAC sanction sanctions check | `FLAG_CREATOR_PAYOUTS_KYC` | Balances accumulate but disbursement locked pending tax & banking verification |

---

## 3. Dynamic Policy Resolver Architecture
The API Gateway and frontends query the **Policy Service**:
```json
{
  "user_id": "usr_99182374",
  "geo_country": "KR",
  "account_age": 16,
  "identity_verified": false,
  "active_flags": {
    "FLAG_MIN_AGE_13_GATE": true,
    "FLAG_KR_REALNAME_VERIFY": true,
    "FLAG_MINOR_STRANGER_DM_SHIELD": true,
    "FLAG_NO_BEHAVIORAL_ADS_UNDER_18": true,
    "FLAG_DATING_ENABLED": false,
    "FLAG_MARKETPLACE_SELLER_KYC": true
  }
}
```
All frontend buttons and backend RPC handlers inspect `active_flags` prior to execution.
