// Register page JavaScript
document.addEventListener('DOMContentLoaded', function() {
    // Check if user is already logged in
    const token = localStorage.getItem(config.STORAGE_KEYS.TOKEN);
    if (token) {
        // Redirect to dashboard if already logged in
        window.location.href = 'dashboard.html';
    }

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
        })        .then(data => {
            // Hide loading spinner
            loadingSpinner.hide();
            
            // Show success message and redirect to login page
            alert(data.message || 'Registration successful! Please login.');
            window.location.href = 'login.html';
        })
        .catch(error => {
            // Hide loading spinner
            loadingSpinner.hide();
            errorMessage.textContent = error.message;
        });
    });
});
