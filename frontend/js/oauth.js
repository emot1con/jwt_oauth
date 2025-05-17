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
    }    // Generate a random state for CSRF protection
    generateRandomState() {
        return Math.random().toString(36).substring(2, 15) + 
               Math.random().toString(36).substring(2, 15);
    }    // Initiate OAuth login flow
    initiateOAuth(provider) {
        console.log('initiateOAuth called for provider:', provider);
        
        // Double-check if user is already logged in
        const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
        if (token) {
            console.log('User already has token, redirecting to dashboard');
            window.location.replace('dashboard.html');
            return;
        }
        
        // Show loading spinner
        loadingSpinner.show();
        
        // Store provider in local storage for callback identification
        localStorage.setItem('oauth_provider', provider);
        
        // Redirect to the backend's OAuth endpoints which handle OAuth flow
        console.log(`Initiating ${provider} OAuth flow, redirecting to:`, this.endpoints[provider]);
        
        // Redirect to the appropriate OAuth provider endpoint
        window.location.href = this.endpoints[provider];
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
    }    // Process OAuth callback from the provider
    processCallback(params) {
        // If there's an error in the params, show it
        if (params.error) {
            loadingSpinner.hide();
            return false;
        }
        
        // Get the code from the URL params
        const code = params.code;
        const provider = localStorage.getItem('oauth_provider') || 'google';
        localStorage.removeItem('oauth_provider');
        
        console.log(`Processing ${provider} OAuth callback with code: ${code}`);
        
        // Exchange the code for a token with the backend
        return this.exchangeCodeForToken(provider, code)
            .then(data => {
                console.log('OAuth token response:', data);
                
                // Store tokens in local storage
                localStorage.setItem(config.STORAGE_KEYS.TOKEN, data.token);
                localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN, data.refresh_token);
                localStorage.setItem(config.STORAGE_KEYS.TOKEN_EXPIRY, data.token_expired_at);
                localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY, data.refresh_token_expired_at);
                
                // Hide spinner and redirect to dashboard
                loadingSpinner.hide();
                window.location.href = 'dashboard.html';
                return true;
            })            .catch(error => {
                console.error('OAuth error:', error);
                loadingSpinner.hide();
                alert('Authentication failed: ' + error.message);
                return false;
            });
    }    // Exchange authorization code for access token
    exchangeCodeForToken(provider, code) {
        // Provider-specific callback endpoints
        const endpoints = {
            google: `${config.API_URL}/oauth/google/callback`,
            github: `${config.API_URL}/oauth/github/callback`,
            facebook: `${config.API_URL}/oauth/facebook/callback`
        };

        // Show loading spinner
        loadingSpinner.show();

        // Send the code to the backend to exchange for a token
        return fetch(endpoints[provider], {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ code })
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Failed to exchange code for token');
                });
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
