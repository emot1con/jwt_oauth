// Immediately check if user is logged in and redirect to dashboard if token exists
(function() {
    console.log('Index page loaded, checking for token...');
    
    // Check if user is already logged in
    const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    const tokenExpiry = localStorage.getItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
    
    if (token) {
        console.log('Token found, checking expiry...');
        
        // If token exists and is not expired, redirect to dashboard
        const now = new Date();
        const expiryDate = tokenExpiry ? new Date(tokenExpiry) : null;
        
        if (!expiryDate || expiryDate > now) {
            console.log('Token is valid, redirecting to dashboard');
            window.location.replace('dashboard.html');
        } else {
            console.log('Token is expired, attempting to use refresh token');
            
            // Check if refresh token exists and try to use it
            const refreshToken = localStorage.getItem(config.STORAGE_KEYS.REFRESH_TOKEN);
            if (refreshToken) {
                // Will handle this in the dashboard page
                window.location.replace('dashboard.html');
            } else {
                console.log('No refresh token available, staying on index page');
                // Clear any invalid tokens
                localStorage.removeItem(config.STORAGE_KEYS.TOKEN);
                localStorage.removeItem(config.STORAGE_KEYS.TOKEN_EXPIRY);
            }
        }
    } else {
        console.log('No token found, staying on index page');
    }
})();
