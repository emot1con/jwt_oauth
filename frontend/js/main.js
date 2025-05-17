// Main JavaScript file
document.addEventListener('DOMContentLoaded', function() {
    console.log('Main page loaded');
    
    // Function to check if the token is still valid
    function isTokenValid() {
        const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
        if (!token) {
            console.log('No token found');
            return false;
        }
        
        const tokenExpiry = localStorage.getItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
        if (!tokenExpiry) {
            console.log('No token expiry found');
            return false;
        }
        
        // Check if token has expired
        const expiryDate = new Date(tokenExpiry);
        const now = new Date();
        
        console.log('Token expires at:', expiryDate);
        console.log('Current time:', now);
        console.log('Token valid?', expiryDate > now);
        
        return expiryDate > now;
    }
    
    // Check if user is already logged in with a valid token
    if (isTokenValid()) {
        console.log('User has valid token, redirecting to dashboard');
        // Redirect to dashboard if already logged in
        window.location.replace('dashboard.html');
    } else {
        console.log('Token is invalid or not present, staying on index page');
        // Clear any invalid tokens
        localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
        localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
        localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
        localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
    }
    
    // Global utility to check token and redirect if needed
    window.authUtils = {
        // Redirect to login page if not authenticated
        requireAuth: function() {
            if (!isTokenValid()) {
                console.log('Authentication required, redirecting to login');
                window.location.replace('login.html');
                return false;
            }
            return true;
        },
        
        // Redirect to dashboard if already authenticated
        redirectIfAuthenticated: function() {
            if (isTokenValid()) {
                console.log('User already authenticated, redirecting to dashboard');
                window.location.replace('dashboard.html');
                return true;
            }
            return false;
        },
        
        // Logout user by clearing storage and redirecting
        logout: function() {
            console.log('Logging out user');
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            localStorage.removeItem(config.STORAGE_KEYS.REFRESH_TOKEN_EXPIRY);
            window.location.replace('login.html');
        }
    };
});
