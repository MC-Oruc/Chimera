// API service for making authenticated requests
import { getCurrentUserToken } from './authService';

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Helper function to make authenticated API requests
export const makeAuthenticatedRequest = async (endpoint, options = {}) => {
  try {
    const token = await getCurrentUserToken();
    
    if (!token) {
      throw new Error('No authentication token available');
    }

    const url = `${API_BASE_URL}${endpoint}`;
    const headers = {
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json',
      ...options.headers
    };

    const response = await fetch(url, {
      ...options,
      headers
    });

    // Handle authentication errors
    if (response.status === 401) {
      // Token might be expired, try to refresh
      localStorage.removeItem('token');
      localStorage.removeItem('userData');
      window.location.href = '/login';
      throw new Error('Authentication failed');
    }

    return response;
  } catch (error) {
    console.error('API request failed:', error);
    throw error;
  }
};

// Helper function for GET requests
export const apiGet = async (endpoint) => {
  const response = await makeAuthenticatedRequest(endpoint, {
    method: 'GET'
  });
  return response.json();
};

// Helper function for POST requests
export const apiPost = async (endpoint, data) => {
  const response = await makeAuthenticatedRequest(endpoint, {
    method: 'POST',
    body: JSON.stringify(data)
  });
  return response.json();
};

// Helper function for PUT requests
export const apiPut = async (endpoint, data) => {
  const response = await makeAuthenticatedRequest(endpoint, {
    method: 'PUT',
    body: JSON.stringify(data)
  });
  return response.json();
};

// Helper function for DELETE requests
export const apiDelete = async (endpoint) => {
  const response = await makeAuthenticatedRequest(endpoint, {
    method: 'DELETE'
  });
  return response.json();
};

// Public API requests (no authentication required)
export const makePublicRequest = async (endpoint, options = {}) => {
  const url = `${API_BASE_URL}${endpoint}`;
  const headers = {
    'Content-Type': 'application/json',
    ...options.headers
  };

  return fetch(url, {
    ...options,
    headers
  });
};

// Helper function for public GET requests
export const publicGet = async (endpoint) => {
  const response = await makePublicRequest(endpoint, {
    method: 'GET'
  });
  return response.json();
};

// Check API health and get service info
export const getServiceInfo = async () => {
  try {
    return await publicGet('/api/info');
  } catch (error) {
    console.error('Failed to get service info:', error);
    return { mode: 'unknown', firebaseEnabled: false };
  }
};

export default {
  makeAuthenticatedRequest,
  makePublicRequest,
  apiGet,
  apiPost,
  apiPut,
  apiDelete,
  publicGet,
  getServiceInfo
}; 