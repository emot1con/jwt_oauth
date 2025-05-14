# Authentication System Frontend

This is a simple frontend for the authentication system backend. It provides a user interface for registration, login, profile viewing, and account management.

## Features

- User registration
- User login with email/password
- OAuth login with Google, GitHub, and Facebook
- View profile information
- Token refresh
- Logout
- Account deletion

## Setup and Running

### Prerequisites

- Go installed on your system (to run the simple frontend server)
- Backend server running (configured in the JS files)

### Running the Frontend

#### Option 1: Running with Docker (Recommended)

The frontend can be run in a Docker container alongside the backend:

```
cd ..
make up_build
```

This will build and start both the frontend and backend services together.

Open your browser and navigate to: `http://localhost:3000`

#### Option 2: Running Locally

1. To start the frontend server locally, run:
   ```
   cd frontend
   go run server.go
   ```

2. Open your browser and navigate to: `http://localhost:3000`

The frontend is automatically configured to connect to the backend API using the same hostname as the browser is using, but on port 8080.

## Directory Structure

- `index.html` - Landing page
- `login.html` - Login page
- `register.html` - Registration page
- `dashboard.html` - User dashboard (protected, requires login)
- `oauth-callback.html` - OAuth authentication callback handler
- `css/` - CSS styles
  - `styles.css` - Main stylesheet
  - `oauth.css` - OAuth button styles
- `img/` - Image assets
  - `oauth/` - OAuth provider icons
- `js/` - JavaScript files
  - `main.js` - Main JavaScript file
  - `config.js` - Configuration settings
  - `login.js` - Login page functionality
  - `register.js` - Registration page functionality
  - `dashboard.js` - Dashboard functionality (profile, refresh token, logout, delete account)
  - `oauth.js` - OAuth authentication functionality
  - `spinner.js` - Loading indicator component

## Authentication Flow

1. **Registration**: 
   - Option 1: User provides name, email, and password to create an account
   - Option 2: User registers via OAuth provider (Google, GitHub, or Facebook)

2. **Login**:
   - Option 1: User provides email and password to authenticate
   - Option 2: User authenticates via OAuth provider (Google, GitHub, or Facebook)

3. **OAuth Flow** (when using OAuth):
   - User clicks OAuth provider button
   - User is redirected to the provider's login page
   - Upon successful authentication, provider redirects to our callback page
   - Callback page exchanges the auth code for a token via the backend API
   - User is authenticated and redirected to the dashboard

4. **Token Management**: 
   - Access token is used for API requests
   - Refresh token is used to obtain a new access token when it expires

5. **Logout**: Invalidates the current token

6. **Account Deletion**: Permanently removes the user account

## Security Features

- JWT token-based authentication
- OAuth 2.0 authentication flow
- CSRF protection via state parameter in OAuth
- Token refresh mechanism
- Rate limiting protection (handled by backend)
- Secure password handling (never stored in browser)

## OAuth Configuration

To enable OAuth login functionality:

1. Update the client IDs in `js/config.js`:

```javascript
OAUTH: {
    // ...
    CLIENT_IDS: {
        GOOGLE: 'YOUR_GOOGLE_CLIENT_ID',
        GITHUB: 'YOUR_GITHUB_CLIENT_ID',
        FACEBOOK: 'YOUR_FACEBOOK_CLIENT_ID'
    }
}
```

2. Create OAuth applications with each provider:

- **Google**: https://console.developers.google.com/
- **GitHub**: https://github.com/settings/developers
- **Facebook**: https://developers.facebook.com/

3. Set the redirect URI in each OAuth application to:
   `http://your-domain/oauth-callback.html`

4. Implement the backend API endpoints for OAuth authentication (see the backend README).
