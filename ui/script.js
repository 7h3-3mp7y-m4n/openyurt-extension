const API_BASE_URL = 'http://localhost:8080';
const POLL_INTERVAL = 2000;

const statusDot = document.getElementById('statusDot');
const statusText = document.getElementById('statusText');
const statusDetails = document.getElementById('statusDetails');
const installBtn = document.getElementById('installBtn');
const uninstallBtn = document.getElementById('uninstallBtn');
const dashboardBtn = document.getElementById('dashboardBtn');

let statusPolling = null;
let dashboardUrl = null;

document.addEventListener('DOMContentLoaded', function() {
    checkStatus();
    startStatusPolling();
});

async function apiCall(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE_URL}${endpoint}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API call failed:', error);
        throw error;
    }
}

async function checkStatus() {
    try {
        const result = await apiCall('/status');
        if (result.success && result.data) {
            updateUIFromStatus(result.data);
        }
    } catch (error) {
        updateStatus('Connection Error', 'Unable to connect to backend service');
        console.error('Status check failed:', error);
    }
}

function updateUIFromStatus(status) {
    const { installed, status: currentStatus, message } = status;
    updateStatus(getDisplayStatus(currentStatus), message);
    updateButtonStates(installed, currentStatus);
}

function getDisplayStatus(status) {
    const statusMap = {
        'not_installed': 'Not Installed',
        'installing': 'Installing',
        'running': 'Running',
        'uninstalling': 'Uninstalling',
        'failed': 'Failed'
    };
    return statusMap[status] || status;
}

function updateButtonStates(installed, status) {
    const isInstalling = status === 'installing';
    const isUninstalling = status === 'uninstalling';
    const isProcessing = isInstalling || isUninstalling;
    installBtn.disabled = installed || isProcessing;
    document.getElementById('installText').textContent = isInstalling ? 'Installing...' : 'Install OpenYurt';
    uninstallBtn.disabled = !installed || isProcessing;
    document.getElementById('uninstallText').textContent = isUninstalling ? 'Uninstalling...' : 'Uninstall OpenYurt';
    dashboardBtn.disabled = !installed || status !== 'running';
}

function updateStatus(status, details) {
    statusText.textContent = `OpenYurt Status: ${status}`;
    statusDetails.textContent = details;
    statusDot.className = 'status-dot';
    switch(status.toLowerCase()) {
        case 'running':
            statusDot.classList.add('connected');
            break;
        case 'installing':
        case 'uninstalling':
            statusDot.classList.add('installing');
            break;
        case 'failed':
        case 'connection error':
            statusDot.classList.add('error');
            break;
    }
}

async function installOpenYurt() {
    try {
        setLoading(installBtn, document.getElementById('installText'));
        
        const result = await apiCall('/install', {
            method: 'POST'
        });
        
        if (result.success) {
            updateStatus('Installing', result.message);
        } else {
            throw new Error(result.message || 'Installation failed');
        }
    } catch (error) {
        updateStatus('Failed', `Installation failed: ${error.message}`);
        installBtn.disabled = false;
        document.getElementById('installText').textContent = 'Install OpenYurt';
    }
}

async function uninstallOpenYurt() {
    if (!confirm('Are you sure you want to uninstall OpenYurt? This will remove all components.')) {
        return;
    }
    
    try {
        setLoading(uninstallBtn, document.getElementById('uninstallText'));
        
        const result = await apiCall('/uninstall', {
            method: 'POST'
        });
        
        if (result.success) {
            updateStatus('Uninstalling', result.message);
            dashboardUrl = null;
        } else {
            throw new Error(result.message || 'Uninstallation failed');
        }
    } catch (error) {
        updateStatus('Failed', `Uninstallation failed: ${error.message}`);
        uninstallBtn.disabled = false;
        document.getElementById('uninstallText').textContent = 'Uninstall OpenYurt';
    }
}

async function openDashboard() {
    try {
        const result = await apiCall('/dashboard');
        
        if (result.success && result.data && result.data.url) {
            const url = result.data.url.startsWith('http') ? result.data.url : `http://${result.data.url}`;
            window.open(url, '_blank');
        } else {
            throw new Error(result.message || 'Dashboard URL not available');
        }
    } catch (error) {
        alert(`Unable to open dashboard: ${error.message}`);
        console.error('Dashboard open failed:', error);
    }
}

function setLoading(button, textElement, loading = true) {
    if (loading) {
        textElement.innerHTML = '<div class="loading"></div>Processing...';
        button.disabled = true;
    }
}


function startStatusPolling() {
    if (statusPolling) {
        clearInterval(statusPolling);
    }
    
    statusPolling = setInterval(checkStatus, POLL_INTERVAL);
}

function stopStatusPolling() {
    if (statusPolling) {
        clearInterval(statusPolling);
        statusPolling = null;
    }
}

document.addEventListener('visibilitychange', function() {
    if (document.hidden) {
        stopStatusPolling();
    } else {
        startStatusPolling();
    }
});

window.addEventListener('beforeunload', function() {
    stopStatusPolling();
});