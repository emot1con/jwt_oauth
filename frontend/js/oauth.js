// OAuth authentication functions
class OAuthHandler {
    constructor() {
        // OAuth endpoints - these would need to be implemented on the backend
        this.endpoints = {
            google: `${config.API_URL}${config.ENDPOINTS.OAUTH_GOOGLE}`,
            github: `${config.API_URL}${config.ENDPOINTS.OAUTH_GITHUB}`,
            facebook: `${config.API_URL}${config.ENDPOINTS.OAUTH_FACEBOOK}`
        };
        
        // OAuth redirect URL
        this.redirectUri = config.OAUTH.REDIRECT_URI;
        
        // OAuth window settings
        this.windowSettings = 'width=500,height=600,menubar=no,toolbar=no,location=no,resizable=yes,scrollbars=yes';
        
        // Generate a random state for CSRF protection
        this.state = this.generateRandomState();
    }
    
    // Generate a random state for CSRF protection
    generateRandomState() {
        return Math.random().toString(36).substring(2, 15) + 
               Math.random().toString(36).substring(2, 15);
    }
      // Initiate OAuth login flow
    initiateOAuth(provider) {
        // Show loading spinner
        loadingSpinner.show();
        
        // Store state in local storage for validation on callback
        localStorage.setItem('oauth_state', this.state);
        
        // Get correct authorization URL based on provider
        let authUrl;
        switch (provider) {
            case 'google':
                authUrl = this.buildGoogleAuthUrl();
                break;
            case 'github':
                authUrl = this.buildGithubAuthUrl();
                break;
            case 'facebook':
                authUrl = this.buildFacebookAuthUrl();
                break;
            default:
                loadingSpinner.hide();
                alert(`Unknown provider: ${provider}`);
                return;
        }
        
        // For demonstration, show a message before continuing
        setTimeout(() => {
            loadingSpinner.hide();
            
            // In a real implementation with backend support, open OAuth window
            if (confirm(`This will redirect you to ${provider} for authentication.\n\nNote: This is just a frontend demonstration. The OAuth flow won't complete without backend implementation.`)) {
                // When backend is implemented, this would open the auth window
                // const authWindow = window.open(authUrl, `${provider}Auth`, this.windowSettings);
                
                // For now, just demonstrate the auth URL format
                console.log(`Auth URL for ${provider}:`, authUrl);
            }
        }, 500);
    }
    
    // Build Google OAuth URL
    buildGoogleAuthUrl() {
        const googleAuthUrl = 'https://accounts.google.com/o/oauth2/v2/auth';
        const params = new URLSearchParams({
            client_id: config.OAUTH.CLIENT_IDS.GOOGLE,
            redirect_uri: this.redirectUri,
            response_type: 'code',
            scope: config.OAUTH.SCOPES.GOOGLE,
            state: this.state,
            access_type: 'offline',
            prompt: 'consent'
        });
        
        return `${googleAuthUrl}?${params.toString()}`;
    }
    
    // Build GitHub OAuth URL
    buildGithubAuthUrl() {
        const githubAuthUrl = 'https://github.com/login/oauth/authorize';
        const params = new URLSearchParams({
            client_id: config.OAUTH.CLIENT_IDS.GITHUB,
            redirect_uri: this.redirectUri,
            scope: config.OAUTH.SCOPES.GITHUB,
            state: this.state
        });
        
        return `${githubAuthUrl}?${params.toString()}`;
    }
    
    // Build Facebook OAuth URL
    buildFacebookAuthUrl() {
        const facebookAuthUrl = 'https://www.facebook.com/v12.0/dialog/oauth';
        const params = new URLSearchParams({
            client_id: config.OAUTH.CLIENT_IDS.FACEBOOK,
            redirect_uri: this.redirectUri,
            scope: config.OAUTH.SCOPES.FACEBOOK,
            state: this.state
        });
        
        return `${facebookAuthUrl}?${params.toString()}`;
    }
      // Process OAuth callback from the backend
    processCallback(params) {
        // Verify state to prevent CSRF attacks
        const savedState = localStorage.getItem('oauth_state');
        if (params.state !== savedState) {
            alert('Invalid state parameter. Authentication rejected.');
            return false;
        }
        
        // Remove state from storage after verification
        localStorage.removeItem('oauth_state');
        
        // Show loading spinner while exchanging code for token
        loadingSpinner.show();
        
        // In a real implementation, we would exchange the code for a token via backend
        // For now, let's simulate the process
        setTimeout(() => {
            // Simulate token exchange success
            const demoToken = "demo_oauth_token_" + Math.random().toString(36).substring(2);
            
            // Store token in local storage
            localStorage.setItem(config.STORAGE_KEYS.TOKEN, demoToken);
            
            // Hide spinner and redirect
            loadingSpinner.hide();
            window.location.href = 'dashboard.html';
        }, 1000);
        
        return true;
    }
    
    // Exchange authorization code for access token (would be handled by backend)
    exchangeCodeForToken(provider, code) {
        // In a real implementation, this would make a request to the backend to exchange code
        return fetch(`${config.API_URL}${config.ENDPOINTS.OAUTH_CALLBACK}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                provider,
                code
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to exchange code for token');
            }
            return response.json();
        });
    }
    
    // Fetch user info after OAuth login (would use the token from previous step)
    fetchUserInfo(token) {
        return fetch(`${config.API_URL}${config.ENDPOINTS.PROFILE}`, {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to fetch user profile');
            }
            return response.json();
        });
    }
}

// Create global OAuth handler instance
const oauthHandler = new OAuthHandler();

// Add event listeners when DOM is ready
document.addEventListener('DOMContentLoaded', () => {
    // Login buttons
    const googleLoginBtn = document.getElementById('google-login');
    const githubLoginBtn = document.getElementById('github-login');
    const facebookLoginBtn = document.getElementById('facebook-login');
    
    // Register buttons
    const googleRegisterBtn = document.getElementById('google-register');
    const githubRegisterBtn = document.getElementById('github-register');
    const facebookRegisterBtn = document.getElementById('facebook-register');
    
    // Add event listeners for login
    if (googleLoginBtn) {
        googleLoginBtn.addEventListener('click', () => oauthHandler.initiateOAuth('google'));
    }
    
    if (githubLoginBtn) {
        githubLoginBtn.addEventListener('click', () => oauthHandler.initiateOAuth('github'));
    }
    
    if (facebookLoginBtn) {
        facebookLoginBtn.addEventListener('click', () => oauthHandler.initiateOAuth('facebook'));
    }
    
    // Add event listeners for register
    if (googleRegisterBtn) {
        googleRegisterBtn.addEventListener('click', () => oauthHandler.initiateOAuth('google'));
    }
    
    if (githubRegisterBtn) {
        githubRegisterBtn.addEventListener('click', () => oauthHandler.initiateOAuth('github'));
    }
    
    if (facebookRegisterBtn) {
        facebookRegisterBtn.addEventListener('click', () => oauthHandler.initiateOAuth('facebook'));
    }
});
