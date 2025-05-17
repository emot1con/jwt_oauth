// Dashboard page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    console.log('Dashboard page loaded');
    
    // Check if user is logged in
    let token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    if (!token) {
        console.log('No token found, redirecting to login');
        // Redirect to login if not logged in
        window.location.replace('login.html');
        return;
    }
    
    console.log('Token found, user is logged in');
    
    const userName = document.getElementById('user-name');
    const userEmail = document.getElementById('user-email');
    const userId = document.getElementById('user-id');
    const logoutBtn = document.getElementById('logout-btn');
    const refreshTokenBtn = document.getElementById('refresh-token-btn');
    const deleteAccountBtn = document.getElementById('delete-account-btn');

    // Function to fetch user profile
    function fetchUserProfile() {
        // Show loading spinner
        loadingSpinner.show();

        fetch(`${config.API_URL}${config.ENDPOINTS.PROFILE}`, {
            method: 'GET',
            headers: {
                'Authorization': token
            }
        })
        .then(response => {
            if (!response.ok) {
                if (response.status === 401) {
                    // Token expired, try to refresh
                    console.log('Token expired, attempting refresh');
                    return refreshToken()
                        .then(() => {
                            // Retry with new token
                            console.log('Token refreshed, retrying profile fetch');
                            return fetchUserProfile();
                        })
                        .catch(error => {
                            throw new Error('Failed to refresh token: ' + error.message);
                        });
                }
                throw new Error('Failed to fetch profile');
            }
            return response.json();
        })
        .then(data => {
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
    }

    // Function to refresh token
    function refreshToken() {
        console.log('Attempting to refresh token');
        
        // Show loading spinner
        loadingSpinner.show();
        
        const refreshToken = localStorage.getItem(config.STORAGE_KEYS.REFRESH_TOKEN);
        if (!refreshToken) {
            console.error('No refresh token available');
            loadingSpinner.hide();
            logout();
            return Promise.reject(new Error('No refresh token available'));
        }
        
        return fetch(`${config.API_URL}${config.ENDPOINTS.REFRESH_TOKEN}`, {
            method: 'POST',
            headers: {
                'Authorization': refreshToken
            }
        })
        .then(response => {
            if (!response.ok) {
                console.error('Refresh token is invalid or expired');
                // If refresh token is invalid, logout
                logout();
                throw new Error('Session expired. Please login again.');
            }
            return response.json();
        })
        .then(data => {
            console.log('Token refreshed successfully:', data);
            
            // Update tokens in storage
            localStorage.setItem(config.STORAGE_KEYS.TOKEN, data.token);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN, data.refresh_token);
            localStorage.setItem(config.STORAGE_KEYS.TOKEN_EXPIRY, data.token_expired_at);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY, data.refresh_token_expired_at);
            
            // Update token variable for current session
            token = data.token;
            
            // Hide loading spinner
            loadingSpinner.hide();
            
            return data.token; // Return the new token
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            console.error('Error refreshing token:', error);
            throw error; // Re-throw the error for the caller to handle
        });
    }

    // Function to logout
    function logout() {
        // Show loading spinner
        loadingSpinner.show();
        
        fetch(`${config.API_URL}${config.ENDPOINTS.LOGOUT}`, {
            method: 'POST',
            headers: {
                'Authorization': token
            }
        })
        .then(response => {
            // Clear local storage and redirect to login page
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.replace('login.html');
        })
        .catch(error => {
            console.error('Error:', error);
            // Even if logout fails on server, clear local storage and redirect
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.replace('login.html');
        });
    }

    // Function to delete account
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
            })
            .then(data => {
                // Hide loading spinner
                loadingSpinner.hide();
                
                alert(data.message || 'Account deleted successfully');
                // Clear local storage and redirect to login page
                localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
                localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
                localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
                localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
                window.location.replace('index.html');
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
    refreshTokenBtn.addEventListener('click', function() {
        refreshToken()
            .then(() => {
                alert('Token refreshed successfully');
                // Optionally reload user profile with new token
                fetchUserProfile();
            })
            .catch(error => {
                alert('Failed to refresh token: ' + error.message);
            });
    });
    deleteAccountBtn.addEventListener('click', deleteAccount);

    // Fetch user profile on page load
    fetchUserProfile();
});
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
                    console.log('Token expired, attempting refresh');
                    return refreshToken()
                        .then(() => {
                            // Retry with new token
                            console.log('Token refreshed, retrying profile fetch');
                            return fetchUserProfile();
                        })
                        .catch(error => {
                            throw new Error('Failed to refresh token: ' + error.message);
                        });
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
        console.log('Attempting to refresh token');
        
        // Show loading spinner
        loadingSpinner.show();
        
        const refreshToken = localStorage.getItem(config.STORAGE_KEYS.REFRESH_TOKEN);
        if (!refreshToken) {
            console.error('No refresh token available');
            loadingSpinner.hide();
            logout();
            return Promise.reject(new Error('No refresh token available'));
        }
        
        return fetch(`${config.API_URL}${config.ENDPOINTS.REFRESH_TOKEN}`, {
            method: 'POST',
            headers: {
                'Authorization': refreshToken
            }
        })
        .then(response => {
            if (!response.ok) {
                console.error('Refresh token is invalid or expired');
                // If refresh token is invalid, logout
                logout();
                throw new Error('Session expired. Please login again.');
            }
            return response.json();
        })
        .then(data => {
            console.log('Token refreshed successfully:', data);
            
            // Update tokens in storage
            localStorage.setItem(config.STORAGE_KEYS.TOKEN, data.token);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN, data.refresh_token);
            localStorage.setItem(config.STORAGE_KEYS.TOKEN_EXPIRY, data.token_expired_at);
            localStorage.setItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY, data.refresh_token_expired_at);
            
            // Update token variable for current session
            token = data.token;
            
            // Hide loading spinner
            loadingSpinner.hide();
            
            return data.token; // Return the new token
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            console.error('Error refreshing token:', error);
            throw error; // Re-throw the error for the caller to handle
        });
    }
    function logout() {
        // Show loading spinner
        loadingSpinner.show();
        
        fetch(`${config.API_URL}${config.ENDPOINTS.LOGOUT}`, {
            method: 'POST',
            headers: {
                'Authorization': token
            }
        })
        .then(response => {
            // Clear local storage and redirect to login page
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.replace('login.html');
        })
        .catch(error => {
            console.error('Error:', error);
            // Even if logout fails on server, clear local storage and redirect
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.replace('login.html');
        });
    }
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
                window.location.replace('index.html');
            })
            .catch(error => {
                // Hide loading spinner
                loadingSpinner.hide();
                
                console.error('Error:', error);
                alert('Error deleting account: ' + error.message);
            });
        }
    }    // Add event listeners
    logoutBtn.addEventListener('click', logout);
    refreshTokenBtn.addEventListener('click', function() {
        refreshToken()
            .then(() => {
                alert('Token refreshed successfully');
                // Optionally reload user profile with new token
                fetchUserProfile();
            })
            .catch(error => {
                alert('Failed to refresh token: ' + error.message);
            });
    });
    deleteAccountBtn.addEventListener('click', deleteAccount);

    // Fetch user profile on page load
    fetchUserProfile();
