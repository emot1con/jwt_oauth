// Login page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log('Login page loaded, checking for existing login...');
    
    // Check if user is already logged in
    const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    if (token) {
        console.log('User already logged in, redirecting to dashboard');
        // Redirect to dashboard if already logged in
        window.location.replace('dashboard.html');
        return; // Stop further execution
    }
    console.log('No existing login found, showing login page');

    const loginForm = document.getElementById('login-form');
    const errorMessage = document.getElementById('error-message');

    loginForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
          errorMessage.textContent = '';
        
        // Show loading spinner
        loadingSpinner.show();
        
        // Send login request to API
        fetch(`${config.API_URL}${config.ENDPOINTS.LOGIN}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                password: password
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Login failed');
                });
            }
            return response.json();
        })        .then(data => {
            // Save token to local storage
            localStorage.setItem(config.STORAGE_KEYS.TOKEN, data.token);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN, data.refresh_token);
            localStorage.setItem(config.STORAGE_KEYS.TOKEN_EXPIRY, data.token_expired_at);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY, data.refresh_token_expired_at);
              // Hide loading spinner
            loadingSpinner.hide();
            
            // Redirect to dashboard
            window.location.replace('dashboard.html');
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            errorMessage.textContent = error.message;
        });
    });    // Initialize OAuth handler
    const oauthHandler = new OAuthHandler();
    
    // Add event listeners for OAuth login buttons
    const googleBtn = document.getElementById('google-login');
    const githubBtn = document.getElementById('github-login');
    const facebookBtn = document.getElementById('facebook-login');
    
    function initOAuthFlow(provider) {
        console.log(`Initiating ${provider} OAuth flow`);
        
        // First check if user is already logged in
        const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
        if (token) {
            console.log('User already logged in, redirecting to dashboard');
            window.location.replace('dashboard.html');
            return;
        }
        
        // Store provider in local storage for callback identification
        localStorage.setItem('oauth_provider', provider);
        
        // Show loading spinner
        loadingSpinner.show();
        
        // Redirect to the backend's OAuth endpoints
        const endpoint = `${config.API_URL}${config.ENDPOINTS['OAUTH_' + provider.toUpperCase()]}`;
        console.log(`Redirecting to ${provider} auth:`, endpoint);
        window.location.href = endpoint;
    }
    
    if (googleBtn) {
        googleBtn.addEventListener('click', function(e) {
            e.preventDefault();
            initOAuthFlow('google');
        });
    }
    
    if (githubBtn) {
        githubBtn.addEventListener('click', function(e) {
            e.preventDefault();
            initOAuthFlow('github');
        });
    }
    
    if (facebookBtn) {
        facebookBtn.addEventListener('click', function(e) {
            e.preventDefault();
            initOAuthFlow('facebook');
        });
    }
});
