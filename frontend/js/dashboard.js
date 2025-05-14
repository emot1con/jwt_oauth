// Dashboard page JavaScript
document.addEventListener('DOMContentLoaded', function() {    // Check if user is logged in
    let token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    if (!token) {
        // Redirect to login if not logged in
        window.location.href = 'login.html';
        return;
    }
    
    const userName = document.getElementById('user-name');
    const userEmail = document.getElementById('user-email');
    const userId = document.getElementById('user-id');
    const logoutBtn = document.getElementById('logout-btn');
    const refreshTokenBtn = document.getElementById('refresh-token-btn');
    const deleteAccountBtn = document.getElementById('delete-account-btn');    // Function to fetch user profile
    function fetchUserProfile() {
        // Show loading spinner
        loadingSpinner.show();

        fetch(`${config.API_URL}${config.ENDPOINTS.PROFILE}`, {
            method: 'GET',
            headers: {
                'Authorization': token
            }
        })        .then(response => {
            if (!response.ok) {
                if (response.status === 401) {
                    // Token expired, try to refresh
                    loadingSpinner.hide(); // Hide current spinner before refreshing
                    refreshToken();
                    return Promise.reject('Token expired'); // Properly exit this promise chain
                }
                throw new Error('Failed to fetch profile');
            }
            return response.json();
        }).then(data => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            if (data) {
                userName.textContent = data.name;
                userEmail.textContent = data.email;
                userId.textContent = data.id;
            }
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            console.error('Error:', error);
            alert('Error fetching profile: ' + error.message);
        });
    }    // Function to refresh token
    function refreshToken() {
        // Show loading spinner
        loadingSpinner.show();
        
        const refreshToken = localStorage.getItem(config.STORAGE_KEYS.REFRESH_TOKEN);
        
        fetch(`${config.API_URL}${config.ENDPOINTS.REFRESH_TOKEN}`, {
            method: 'POST',
            headers: {
                'Authorization': refreshToken
            }
        })
        .then(response => {
            if (!response.ok) {
                // If refresh token is invalid, logout
                logout();
                throw new Error('Session expired. Please login again.');
            }
            return response.json();
        })        .then(data => {
            // Update tokens in storage
            localStorage.setItem(config.STORAGE_KEYS.TOKEN, data.token);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN, data.refresh_token);
            localStorage.setItem(config.STORAGE_KEYS.TOKEN_EXPIRY, data.token_expired_at);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY, data.refresh_token_expired_at);
            
            // Update token variable for current session
            token = data.token;
            
            // Hide loading spinner
            loadingSpinner.hide();
            
            // Display success message
            alert('Token refreshed successfully');
            
            // Optionally reload user profile with new token (uncomment if needed)
            // fetchUserProfile();
        }).catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            console.error('Error:', error);
            alert(error.message);
        });
    }    // Function to logout
    function logout() {
        // Show loading spinner
        loadingSpinner.show();
        
        fetch(`${config.API_URL}${config.ENDPOINTS.LOGOUT}`, {
            method: 'POST',
            headers: {
                'Authorization': token
            }
        })        .then(response => {
            // Clear local storage and redirect to login page
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.href = 'login.html';
        })
        .catch(error => {
            console.error('Error:', error);            // Even if logout fails on server, clear local storage and redirect
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.href = 'login.html';
        });
    }    // Function to delete account
    function deleteAccount() {
        if (confirm('Are you sure you want to delete your account? This action cannot be undone.')) {
            // Show loading spinner
            loadingSpinner.show();
            
            fetch(`${config.API_URL}${config.ENDPOINTS.DELETE_ACCOUNT}`, {
                method: 'DELETE',
                headers: {
                    'Authorization': token
                }
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Failed to delete account');
                }
                return response.json();
            })            .then(data => {
                // Hide loading spinner
                loadingSpinner.hide();
                
                alert(data.message || 'Account deleted successfully');
                // Clear local storage and redirect to login page
                localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
                localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
                localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
                localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
                window.location.href = 'index.html';
            })
            .catch(error => {
                // Hide loading spinner
                loadingSpinner.hide();
                
                console.error('Error:', error);
                alert('Error deleting account: ' + error.message);
            });
        }
    }

    // Add event listeners
    logoutBtn.addEventListener('click', logout);
    refreshTokenBtn.addEventListener('click', refreshToken);
    deleteAccountBtn.addEventListener('click', deleteAccount);

    // Fetch user profile on page load
    fetchUserProfile();
});
