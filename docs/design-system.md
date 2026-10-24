# Global Social Network Design System (Codename: Aurora Design)

## 1. 개요 및 디자인 철학
Aurora Design System은 고밀도 소셜 정보 소비와 비즈니스/커머스 워크플로우를 모두 지원하는 모던 반응형 디자인 시스템입니다.
Facebook의 친숙한 사용성을 존중하되, 독자적인 브랜드 아이덴티티와 최신 웹/모바일 트렌드(Fluent + Glassmorphism accents, Dark Mode first-class support)를 결합합니다.

---

## 2. 디자인 토큰 (Design Tokens)

### 2.1 색상 팔레트 (Color Palette)

```css
:root {
  /* Primary Brand (Electric Cobalt) */
  --color-primary-50:  #eef2ff;
  --color-primary-100: #e0e7ff;
  --color-primary-500: #3b82f6;
  --color-primary-600: #1877f2; /* Classic Trust Blue */
  --color-primary-700: #1d4ed8;

  /* Neutrals & Surfaces (Light Mode) */
  --surface-bg:        #f0f2f5;
  --surface-card:      #ffffff;
  --surface-elevated:  #ffffff;
  --text-primary:      #050505;
  --text-secondary:    #65676b;
  --border-subtle:     #ced0d4;
  --divider:           #e4e6eb;

  /* Reaction Emotives */
  --reaction-like:     #1877f2;
  --reaction-love:     #f02849;
  --reaction-haha:     #f7b125;
  --reaction-wow:      #f7b125;
  --reaction-sad:      #f7b125;
  --reaction-angry:    #e41e3f;

  /* Semantic Alerts */
  --status-success:    #31a24c;
  --status-warning:    #f59e0b;
  --status-danger:     #ef4444;
}

[data-theme="dark"] {
  --surface-bg:        #18191a;
  --surface-card:      #242526;
  --surface-elevated:  #3a3b3c;
  --text-primary:      #e4e6eb;
  --text-secondary:    #b0b3b8;
  --border-subtle:     #3e4042;
  --divider:           #393a3b;
}
```

### 2.2 반응형 중단점 (Breakpoints)
- **Mobile (`< 640px`)**: 단일 컬럼 피드, 하단 네비게이션 바 (Bottom Navigation Bar)
- **Tablet (`640px ~ 1024px`)**: 슬림 사이드바 + 메인 피드 + 플로팅 메신저
- **Desktop (`1024px ~ 1440px`)**: 좌측 바로가기 사이드바 (280px) + 중앙 메인 피드 (590px) + 우측 연락처/스폰서 (280px)
- **Ultra-Wide (`> 1440px`)**: 좌우 고정 여백 및 3단 분할 그리드 정렬

---

## 3. 원자 컴포넌트 규격 (Atomic Components)

1. **Button (`AppButton`)**:
   - Variants: `primary`, `secondary`, `tertiary`, `danger`, `ghost`
   - Sizes: `sm` (32px), `md` (36px), `lg` (44px)
   - States: `default`, `hover`, `active`, `disabled`, `loading` (Spinner overlay)

2. **Avatar (`AppAvatar`)**:
   - Sizes: `xs` (24px), `sm` (32px), `md` (40px), `lg` (56px), `xl` (96px), `2xl` (168px)
   - Indicators: Active status dot (Green online), Story ring (Gradient border)

3. **Card (`AppCard`)**:
   - Rounded corners (`rounded-xl` / 12px), Subtle shadow, Border separator

4. **Input & Composer (`PostComposerModal`, `CommentInput`)**:
   - Auto-expanding textarea, Emoji picker, Rich media attachment preview, Mention dropdown
