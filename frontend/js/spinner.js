// Loading spinner functionality
class LoadingSpinner {
    constructor() {
        this.spinnerHTML = `
            <div class="spinner-overlay">
                <div class="spinner"></div>
            </div>
        `;
        this.styleAdded = false;
    }

    // Add spinner CSS to the document
    addSpinnerStyles() {
        if (this.styleAdded) return;
        
        const style = document.createElement('style');
        style.textContent = `
            .spinner-overlay {
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background-color: rgba(0, 0, 0, 0.5);
                display: flex;
                justify-content: center;
                align-items: center;
                z-index: 1000;
            }
            
            .spinner {
                width: 50px;
                height: 50px;
                border: 5px solid #f3f3f3;
                border-top: 5px solid #3498db;
                border-radius: 50%;
                animation: spin 1s linear infinite;
            }
            
            @keyframes spin {
                0% { transform: rotate(0deg); }
                100% { transform: rotate(360deg); }
            }
        `;
        document.head.appendChild(style);
        this.styleAdded = true;
    }

    // Show the spinner
    show() {
        this.addSpinnerStyles();
        
        const spinnerDiv = document.createElement('div');
        spinnerDiv.id = 'loading-spinner';
        spinnerDiv.innerHTML = this.spinnerHTML;
        document.body.appendChild(spinnerDiv);
    }

    // Hide the spinner
    hide() {
        const spinner = document.getElementById('loading-spinner');
        if (spinner) {
            spinner.remove();
        }
    }
}

// Create a global spinner instance
const loadingSpinner = new LoadingSpinner();
