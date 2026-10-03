# Agent 2: Frontend Settings UI Implementation

**Task:** Phase 13b Part 3 & 4 — Frontend Settings Route & Theme System  
**Worktree:** `/home/denny/Project/headscale-settings-frontend`  
**Branch:** `feature/settings-frontend`  
**Base Commit:** `d95c5eff`

---

## Objective

Implement Headplane settings UI:
- Settings page accessible after password login
- API key management (save/update, masked display)
- Password change form
- Theme selector (light/dark with persistence)
- Profile name input
- All with proper validation, error handling, and feedback

---

## Context

Phase 13a (password login) works, but users must re-enter API key at each login. This task adds a Settings page where users can:
- Save their Headscale API key (stored securely on backend)
- Change their Headplane password
- Choose theme (light/dark)
- Set profile name

Backend Agent 1 is implementing the API endpoints in parallel. You can mock them initially.

---

## API Contract (Backend Endpoints)

### GET /api/v1/headplane/settings
**Request:**
```
GET /api/v1/headplane/settings
Authorization: Bearer <session-token>
```

**Response:**
```json
{
  "apiKey": "hs_abc123...",
  "theme": "light",
  "profileName": "John Doe"
}
```

### POST /api/v1/headplane/settings
**Request:**
```
POST /api/v1/headplane/settings
Authorization: Bearer <session-token>
Content-Type: application/json

{
  "apiKey": "hs_new_key",
  "theme": "dark",
  "profileName": "Jane Smith"
}
```

**Response:**
```json
{"success": true}
```

### POST /api/v1/headplane/change-password
**Request:**
```
POST /api/v1/headplane/change-password
Authorization: Bearer <session-token>
Content-Type: application/json

{
  "currentPassword": "old-password",
  "newPassword": "new-password"
}
```

**Response:**
```json
{"success": true}
```

**Errors:** 401 (invalid current password), 400 (validation error)

---

## Implementation Tasks

### Task 1: Add Backend API Methods

**File:** `headplane/app/server/headscale/api/index.ts`

Add methods to Headscale API interface:
```typescript
export interface HeadscaleAPI {
  // ... existing methods
  
  getSettings(): Promise<{
    apiKey: string;
    theme: string;
    profileName: string;
  }>;
  
  updateSettings(data: {
    apiKey?: string;
    theme?: string;
    profileName?: string;
  }): Promise<{success: boolean}>;
  
  changePassword(data: {
    currentPassword: string;
    newPassword: string;
  }): Promise<{success: boolean}>;
}
```

Implement in transport class.

### Task 2: Create Settings Route

**File:** Create `headplane/app/routes/settings/profile/page.tsx`

Structure:
```tsx
export default function SettingsProfilePage() {
  return (
    <div className="settings-container">
      <h1>Settings</h1>
      
      <section className="settings-section">
        <h2>Account</h2>
        <PasswordChangeForm />
        <SessionInfo />
      </section>
      
      <section className="settings-section">
        <h2>Integration</h2>
        <APIKeyForm />
      </section>
      
      <section className="settings-section">
        <h2>Preferences</h2>
        <ThemeSelector />
      </section>
      
      <section className="settings-section">
        <h2>Profile</h2>
        <ProfileNameForm />
      </section>
    </div>
  );
}
```

### Task 3: API Key Management Component

**File:** Create `headplane/app/routes/settings/profile/components/api-key-form.tsx`

Features:
- Input field for API key
- Masked display by default (e.g., `hs_****...****abc`)
- "Show/Hide" toggle button
- "Save" button
- Success/error feedback
- Load current value on mount

Validation:
- API key format: starts with `hs_`
- Minimum length check

### Task 4: Password Change Component

**File:** Create `headplane/app/routes/settings/profile/components/password-change-form.tsx`

Features:
- Current password field
- New password field
- Confirm new password field
- Password strength indicator
- "Change Password" button
- Success/error feedback

Validation:
- Current password required
- New password minimum 8 characters
- New password != current password
- Confirm matches new

### Task 5: Theme Selector Component

**File:** Create `headplane/app/routes/settings/profile/components/theme-selector.tsx`

Features:
- Radio buttons or toggle for light/dark
- Immediate preview (apply theme on select)
- Auto-save on change
- Visual feedback

Implementation:
```tsx
function ThemeSelector() {
  const [theme, setTheme] = useState<'light' | 'dark'>('light');
  
  const handleChange = async (newTheme: 'light' | 'dark') => {
    setTheme(newTheme);
    applyTheme(newTheme);
    await api.updateSettings({theme: newTheme});
  };
  
  return (
    <div>
      <label>
        <input type="radio" checked={theme === 'light'} 
               onChange={() => handleChange('light')} />
        Light
      </label>
      <label>
        <input type="radio" checked={theme === 'dark'} 
               onChange={() => handleChange('dark')} />
        Dark
      </label>
    </div>
  );
}
```

### Task 6: Profile Name Component

**File:** Create `headplane/app/routes/settings/profile/components/profile-name-form.tsx`

Features:
- Text input for profile name
- "Save" button
- Success/error feedback
- Load current value on mount

Validation:
- Optional (can be empty)
- Max 50 characters

### Task 7: Theme System Implementation

**File:** `headplane/app/root.tsx` or similar

Implement theme application:
```typescript
function applyTheme(theme: 'light' | 'dark') {
  document.documentElement.setAttribute('data-theme', theme);
}

// On app load
useEffect(() => {
  api.getSettings().then(settings => {
    applyTheme(settings.theme);
  });
}, []);
```

**File:** `headplane/app/styles/themes.css` or similar

Define CSS variables:
```css
:root {
  --bg-primary: #ffffff;
  --text-primary: #000000;
  /* ... more variables */
}

[data-theme="dark"] {
  --bg-primary: #1a1a1a;
  --text-primary: #ffffff;
  /* ... more variables */
}
```

Use variables in components:
```css
.container {
  background: var(--bg-primary);
  color: var(--text-primary);
}
```

### Task 8: Add Settings Link to Navigation

**File:** Modify navigation component (check existing nav structure)

Add "Settings" link visible only after password authentication.

Conditional rendering:
```tsx
{isPasswordAuthenticated && (
  <Link to="/settings/profile">Settings</Link>
)}
```

### Task 9: Session Info Component

**File:** Create `headplane/app/routes/settings/profile/components/session-info.tsx`

Display:
- Session expiry time
- "Logout" button
- Current authentication method (password vs API key)

### Task 10: Form Validation & Error Handling

For all forms:
- Client-side validation before submit
- Display field-level errors
- Display API errors (401, 429, 500, etc.)
- Loading states during submission
- Success feedback (toast/banner)

Use existing Headplane form patterns.

---

## Testing

### Unit Tests

**Files:**
- `headplane/tests/unit/settings/api-key-form.test.ts`
- `headplane/tests/unit/settings/password-change-form.test.ts`
- `headplane/tests/unit/settings/theme-selector.test.ts`
- `headplane/tests/unit/settings/profile-name-form.test.ts`

Test:
- Component rendering
- Form validation
- Submit handlers
- Error display
- Success feedback

### E2E Tests

**File:** `headplane/tests/e2e/settings.spec.ts`

Test flows:
1. Login with password → navigate to settings → see settings page
2. Save API key → logout → login → API key loaded
3. Change theme → logout → login → theme persisted
4. Change password → logout → login with new password
5. Update profile name → see name displayed

### Run Tests

```bash
cd /home/denny/Project/headscale-settings-frontend/headplane

pnpm install
pnpm test:unit
pnpm test:e2e
pnpm typecheck
pnpm lint
```

---

## Mocking Backend (Initial Development)

If backend endpoints aren't ready yet, mock them:

**File:** `headplane/app/server/headscale/api/mock.ts` (or similar)

```typescript
export const mockSettings = {
  apiKey: 'hs_mock_key_12345',
  theme: 'light',
  profileName: 'Test User'
};

export async function mockGetSettings() {
  return mockSettings;
}

export async function mockUpdateSettings(data: any) {
  Object.assign(mockSettings, data);
  return {success: true};
}

export async function mockChangePassword(data: any) {
  if (data.currentPassword === 'current') {
    return {success: true};
  }
  throw new Error('Invalid current password');
}
```

Switch to real API once backend is ready.

---

## Files to Create/Modify

### Create:
- `headplane/app/routes/settings/profile/page.tsx`
- `headplane/app/routes/settings/profile/components/api-key-form.tsx`
- `headplane/app/routes/settings/profile/components/password-change-form.tsx`
- `headplane/app/routes/settings/profile/components/theme-selector.tsx`
- `headplane/app/routes/settings/profile/components/profile-name-form.tsx`
- `headplane/app/routes/settings/profile/components/session-info.tsx`
- `headplane/app/styles/themes.css` (or modify existing)
- `headplane/tests/unit/settings/` (test files)
- `headplane/tests/e2e/settings.spec.ts`

### Modify:
- `headplane/app/server/headscale/api/index.ts` (add methods)
- `headplane/app/server/headscale/api/transport.ts` (implement methods)
- `headplane/app/root.tsx` (theme application)
- `headplane/app/routes/settings/overview.tsx` (add Profile link)
- Navigation component (add Settings link)

---

## Design Guidelines

- **Match existing Headplane UI:** Use same components, colors, spacing
- **Responsive:** Works on mobile and desktop
- **Accessible:** Proper labels, ARIA attributes, keyboard navigation
- **Consistent:** Follow patterns from other settings routes (auth-keys, restrictions)
- **Clear feedback:** Success messages, error messages, loading states

---

## Acceptance Criteria

- [ ] `/settings/profile` route exists and accessible after login
- [ ] Settings page has 4 sections: Account, Integration, Preferences, Profile
- [ ] API key input with masked display and save
- [ ] Password change form with validation
- [ ] Theme selector with immediate application
- [ ] Profile name input with save
- [ ] "Settings" link in navigation (password auth only)
- [ ] All forms have validation and error handling
- [ ] Theme persists across page refreshes
- [ ] All unit tests pass
- [ ] All E2E tests pass

---

## Coordination with Other Agents

### With Backend Agent:
- **API Contract:** Confirm endpoint paths and request/response formats
- **Mock → Real:** Start with mocks, switch to real endpoints when ready
- **Integration Testing:** Test together after both complete

### With Docs Agent:
- **Screenshots:** Provide screenshots of settings page for docs
- **User Flows:** Explain how to use settings features

---

## Completion

When done:
1. Run all tests and verify they pass
2. Test manually in browser (npm run dev)
3. Commit your changes
4. Push branch: `git push -u origin feature/settings-frontend`
5. Create summary: files changed, tests added, commit SHA
6. Report completion with merge readiness

**Do not merge yet.** Integration testing happens after all agents complete.

---

## References

- **Implementation Plan:** `docs/phase13b-implementation-plan.md`
- **Existing Settings:** `headplane/app/routes/settings/`
- **Auth UI:** `headplane/app/routes/auth/login/`
- **ROADMAP:** `ROADMAP.md` Phase 13b
