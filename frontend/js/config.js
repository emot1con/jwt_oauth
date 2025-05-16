// Configuration file for the frontend
const config = {
    // API base URL - use the same hostname as the browser is using, but different port
    API_URL: `http://${window.location.hostname}:8080`,
    
    // Token Storage Keys
    STORAGE_KEYS: {
        TOKEN: 'token',
        REFRESH_TOKEN: 'refreshToken',
        TOKEN_EXPIRY: 'tokenExpiry',
        REFRESH_TOKEN_EXPIRY: 'refreshTokenExpiry'
    },
    
    // API Endpoints
    ENDPOINTS: {
        REGISTER: '/auth/register',
        LOGIN: '/auth/login',
        LOGOUT: '/user/logout',
        PROFILE: '/user/profile',
        REFRESH_TOKEN: '/user/refresh',
        DELETE_ACCOUNT: '/user/delete',
          // OAuth endpoints
        OAUTH_GOOGLE: '/oauth/google',
        OAUTH_GITHUB: '/oauth/github',
        OAUTH_FACEBOOK: '/oauth/facebook',
        OAUTH_CALLBACK: '/oauth/callback'
    },
      // OAuth configuration
    OAUTH: {
        REDIRECT_URI: `${window.location.origin}/oauth-callback.html`,
        SCOPES: {
            GOOGLE: 'email profile',
            GITHUB: 'user:email',
            FACEBOOK: 'email,public_profile'
        },
        CLIENT_IDS: {
            // These are placeholder values - replace with your actual OAuth app client IDs
            GOOGLE: 'YOUR_GOOGLE_CLIENT_ID',
            GITHUB: 'YOUR_GITHUB_CLIENT_ID',
            FACEBOOK: 'YOUR_FACEBOOK_CLIENT_ID'
        }
    }
};
