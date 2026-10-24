# Independent Design System Specification (Nova UI)

## 1. Brand Philosophy & Identity
The platform utilizes an independent modern design system named **Nova UI**. It ensures enterprise-tier ergonomic usability, accessible contrast ratios (WCAG 2.1 AA compliant), and responsive behavior across web and mobile while maintaining complete independence from third-party trademarked assets.

---

## 2. Design Tokens

### 2.1 Color Palette
```css
:root {
  /* Brand Primary */
  --nova-primary-50:  #eef4ff;
  --nova-primary-100: #d9e6ff;
  --nova-primary-300: #84b0ff;
  --nova-primary-500: #1b64f2; /* Core Accent */
  --nova-primary-600: #134ecc;
  --nova-primary-700: #0f3ba1;

  /* Neutrals & Surfaces */
  --nova-surface-bg:       #f0f2f5; /* Light feed background */
  --nova-surface-card:     #ffffff; /* Card / dialog surface */
  --nova-surface-hover:    #f2f4f7;
  --nova-surface-border:   #e4e6eb;
  
  /* Text & Hierarchy */
  --nova-text-primary:     #050505;
  --nova-text-secondary:   #65676b;
  --nova-text-tertiary:    #8a8d91;
  --nova-text-inverse:     #ffffff;

  /* Semantic Feedback */
  --nova-semantic-success: #31a24c;
  --nova-semantic-warning: #f7b125;
  --nova-semantic-danger:  #fa383e;
  --nova-semantic-info:    #1877f2;
}

/* Dark Mode Tokens */
[data-theme="dark"] {
  --nova-surface-bg:       #18191a;
  --nova-surface-card:     #242526;
  --nova-surface-hover:    #3a3b3c;
  --nova-surface-border:   #393a3b;
  --nova-text-primary:     #e4e6eb;
  --nova-text-secondary:   #b0b3b8;
  --nova-text-tertiary:    #8a8d91;
}
```

### 2.2 Typography Scale
- Font Family: `Inter, -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif`
- Display Large: `32px / line-height: 38px / 700 bold`
- Title Main: `20px / line-height: 24px / 600 semibold`
- Body Standard: `15px / line-height: 20px / 400 normal`
- Caption / Subtext: `13px / line-height: 16px / 400 normal`
- Micro Label: `11px / line-height: 14px / 500 medium`

### 2.3 Elevation & Shadows
- `shadow-sm`: `0 1px 2px rgba(0, 0, 0, 0.08)`
- `shadow-card`: `0 1px 2px rgba(0, 0, 0, 0.1), 0 2px 4px rgba(0, 0, 0, 0.05)`
- `shadow-modal`: `0 12px 28px 0 rgba(0, 0, 0, 0.2), 0 2px 4px 0 rgba(0, 0, 0, 0.1)`

---

## 3. Core Component Library

1. **Button Component (`<NovaButton>`)**:
   - Variants: `primary`, `secondary`, `tertiary`, `danger`, `ghost`
   - Sizes: `sm` (32px), `md` (38px), `lg` (44px)
   - Features: Loading spinner overlay, disabled accessibility aria-disabled, icon slots.

2. **Card Component (`<NovaCard>`)**:
   - Default container for Feed Posts, Profile bio, Marketplace listings.
   - Built-in header (avatar, author name, timestamp, privacy icon, meatball menu), content slot, action footer.

3. **Reactions Pill Bar (`<NovaReactionsBar>`)**:
   - 6 standard emotional reactions (Like, Love, Care, Haha, Wow, Sad, Angry) with smooth CSS spring transitions.
   - Live counter popover displaying users who reacted by category.

4. **Media Lightbox Modal (`<NovaMediaViewer>`)**:
   - High-performance full-screen image/video viewer.
   - Keyboard navigation (Left/Right arrows, Escape to close).
   - Side rail for comments, reactions, and post metadata.

5. **Responsive Breakpoints**:
   - Mobile: `< 640px` (Single column, bottom navigation bar)
   - Tablet: `640px - 1024px` (Feed + collapsible sidebar)
   - Desktop: `> 1024px` (3-column layout: Left navigation, Center feed, Right contacts/chat)
