// Login page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Check if user is already logged in
    const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    if (token) {
        // Redirect to dashboard if already logged in
        window.location.href = 'dashboard.html';
    }

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
            window.location.href = 'dashboard.html';
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            errorMessage.textContent = error.message;
        });
    });
});
