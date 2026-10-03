# Headplane Settings Frontend

React-based settings interface for Headplane with authentication and profile management.

## Features

- **Password Authentication**: Secure login with session management
- **Settings Management**:
  - Account: Password change, session information
  - Integration: Headscale API key management (encrypted storage)
  - Preferences: Light/dark theme selection with persistence
  - Profile: Display name customization
- **Beautiful UI**: Built with Chakra UI for a modern, accessible interface
- **Responsive**: Works on mobile and desktop devices
- **Type-Safe**: Written in TypeScript

## Tech Stack

- **React 18** with TypeScript
- **Vite** for fast development and building
- **Chakra UI** for beautiful, accessible components
- **React Router** for client-side routing
- **Vitest** for unit testing
- **Playwright** for E2E testing

## Getting Started

### Prerequisites

- Node.js 18+ and npm

### Installation

```bash
cd headplane
npm install
```

### Development

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser.

**Demo Credentials:**
- Password: `password123`

### Building

```bash
npm run build
```

### Testing

```bash
# Unit tests
npm test

# E2E tests
npm run test:e2e

# Type checking
npm run type-check
```

## Project Structure

```
headplane/
├── app/
│   ├── main.tsx              # Application entry point
│   ├── App.tsx               # Root component with routing
│   ├── theme.ts              # Chakra UI theme configuration
│   ├── components/           # Shared components
│   │   └── Layout.tsx        # Main layout with header
│   ├── contexts/             # React contexts
│   │   └── AuthContext.tsx   # Authentication state
│   ├── routes/               # Page routes
│   │   ├── auth/login/       # Login page
│   │   └── settings/profile/ # Settings page with components
│   ├── server/               # API layer
│   │   └── headscale/api/    # Mock API (replace with real backend)
│   └── styles/               # Global styles
├── tests/
│   ├── unit/                 # Component unit tests
│   └── e2e/                  # End-to-end tests
├── index.html                # HTML entry point
├── package.json              # Dependencies and scripts
├── tsconfig.json             # TypeScript configuration
├── vite.config.ts            # Vite configuration
└── playwright.config.ts      # Playwright configuration
```

## API Integration

The app currently uses mock API endpoints in `app/server/headscale/api/index.ts`. To integrate with the real backend:

1. Replace mock implementations with real HTTP calls
2. Update the `baseUrl` in the API client
3. Ensure session tokens are properly passed in Authorization headers

### API Contract

The following endpoints are expected:

- `POST /api/v1/headplane/login` - Authenticate user
- `GET /api/v1/headplane/settings` - Retrieve user settings
- `POST /api/v1/headplane/settings` - Update user settings
- `POST /api/v1/headplane/change-password` - Change password

See `AGENT_INSTRUCTIONS.md` for detailed API specifications.

## Features by Section

### Account Section
- **Password Change**: Change Headplane password with validation
- **Session Info**: Display current session status and age

### Integration Section
- **API Key Management**: Save Headscale API key (masked display, validation)

### Preferences Section
- **Theme Selector**: Switch between light and dark themes (persisted)

### Profile Section
- **Display Name**: Set optional profile name

## Accessibility

- All forms have proper labels and ARIA attributes
- Keyboard navigation support
- Focus indicators on all interactive elements
- Color contrast meets WCAG AA standards
- Screen reader friendly

## Browser Support

- Chrome/Edge (latest)
- Firefox (latest)
- Safari (latest)

## License

Same as parent Headscale project.
