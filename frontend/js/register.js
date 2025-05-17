// Register page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log('Register page loaded, checking for existing login...');
    
    // Check if user is already logged in
    const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    if (token) {
        console.log('User already logged in, redirecting to dashboard');
        // Redirect to dashboard if already logged in
        window.location.replace('dashboard.html');
        return; // Stop further execution
    }
    console.log('No existing login found, showing register page');

    const registerForm = document.getElementById('register-form');
    const errorMessage = document.getElementById('error-message');

    registerForm.addEventListener('submit', function(e) {
        e.preventDefault();
        
        const name = document.getElementById('name').value;
        const email = document.getElementById('email').value;
        const password = document.getElementById('password').value;
        errorMessage.textContent = '';
        
        // Show loading spinner
        loadingSpinner.show();
          
        // Send registration request to API
        fetch(`${config.API_URL}${config.ENDPOINTS.REGISTER}`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                name: name,
                email: email,
                password: password
            })
        })
        .then(response => {
            if (!response.ok) {
                return response.json().then(data => {
                    throw new Error(data.error || 'Registration failed');
                });
            }
            return response.json();
        })
        .then(data => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            // Show success message and redirect to login page
            alert(data.message || 'Registration successful! Please login.');
            window.location.replace('login.html');
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            errorMessage.textContent = error.message;
        });
    });

    // Initialize OAuth handler for registration
    function initOAuthFlow(provider) {
        console.log(`Initiating ${provider} OAuth registration flow`);
        
        // First check if user is already logged in
        const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
        if (token) {
            console.log('User already logged in, redirecting to dashboard');
            window.location.replace('dashboard.html');
            return;
        }
        
        // Store provider in local storage for callback identification
        localStorage.setItem('oauth_provider', provider);
        localStorage.setItem('oauth_action', 'register');
        
        // Show loading spinner
        loadingSpinner.show();
        
        // Redirect to the backend's OAuth endpoints
        const endpoint = `${config.API_URL}${config.ENDPOINTS['OAUTH_' + provider.toUpperCase()]}`;
        console.log(`Redirecting to ${provider} auth:`, endpoint);
        window.location.href = endpoint;
    }
    
    // Add event listeners for OAuth registration buttons
    const googleBtn = document.getElementById('google-register');
    const githubBtn = document.getElementById('github-register');
    const facebookBtn = document.getElementById('facebook-register');
    
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
