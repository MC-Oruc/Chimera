// Auth service abstraction for both Firebase and Local auth
import { 
  signInWithEmailAndPassword as firebaseSignIn,
  signInWithPopup,
  GoogleAuthProvider,
  signOut as firebaseSignOut,
  onAuthStateChanged as firebaseOnAuthStateChanged,
  createUserWithEmailAndPassword as firebaseCreateUser
} from "firebase/auth";
import { auth } from "@/firebase/firebase";

// Auth mode configuration
const getAuthMode = () => {
  // Check if Firebase is enabled via environment variable
  const firebaseEnabled = process.env.NEXT_PUBLIC_FIREBASE_ENABLE !== 'false';
  
  // Check if Firebase config is available
  const hasFirebaseConfig = !!(
    process.env.NEXT_PUBLIC_FIREBASE_API_KEY &&
    process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN &&
    process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID
  );
  
  return firebaseEnabled && hasFirebaseConfig ? 'firebase' : 'local';
};

// Local auth implementation
class LocalAuth {
  constructor() {
    this.currentUser = null;
    this.listeners = new Set();
    this.baseURL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';
    
    // Initialize from localStorage on startup
    this.initializeFromStorage();
  }

  initializeFromStorage() {
    const token = localStorage.getItem('token');
    const userData = localStorage.getItem('userData');
    
    if (token && userData) {
      try {
        this.currentUser = JSON.parse(userData);
        this.notifyListeners(this.currentUser);
      } catch (error) {
        console.error('Error parsing stored user data:', error);
        this.signOut();
      }
    } else {
      this.notifyListeners(null);
    }
  }

  async signInWithEmailAndPassword(email, password) {
    try {
      // For development, create a simple token based on email
      const token = `local_${email.split('@')[0]}_${Date.now()}`;
      const userData = {
        uid: email.split('@')[0],
        email: email,
        displayName: email.split('@')[0],
        emailVerified: true
      };

      // Store token and user data
      localStorage.setItem('token', token);
      localStorage.setItem('userData', JSON.stringify(userData));
      
      this.currentUser = userData;
      this.notifyListeners(this.currentUser);

      return { user: userData };
    } catch (error) {
      throw new Error(`Local auth error: ${error.message}`);
    }
  }

  async createUserWithEmailAndPassword(email, password) {
    // For local development, treat signup same as signin
    return this.signInWithEmailAndPassword(email, password);
  }

  async signInWithPopup(provider) {
    // For local development, simulate Google login
    const email = 'demo@localhost.com';
    return this.signInWithEmailAndPassword(email, 'password');
  }

  async signOut() {
    localStorage.removeItem('token');
    localStorage.removeItem('userData');
    this.currentUser = null;
    this.notifyListeners(null);
  }

  onAuthStateChanged(callback) {
    this.listeners.add(callback);
    
    // Immediately call with current state
    callback(this.currentUser);
    
    return () => {
      this.listeners.delete(callback);
    };
  }

  async getIdToken(forceRefresh = false) {
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('No token available');
    }
    return token;
  }

  notifyListeners(user) {
    this.listeners.forEach(callback => callback(user));
  }
}

// Create auth service instance based on mode
const createAuthService = () => {
  const mode = getAuthMode();
  console.log(`🔐 Auth mode: ${mode.toUpperCase()}`);
  
  if (mode === 'local') {
    return new LocalAuth();
  } else {
    // Firebase auth wrapper
    return {
      signInWithEmailAndPassword: firebaseSignIn,
      createUserWithEmailAndPassword: firebaseCreateUser,
      signInWithPopup,
      signOut: firebaseSignOut,
      onAuthStateChanged: firebaseOnAuthStateChanged,
      getIdToken: async (forceRefresh = false) => {
        if (auth.currentUser) {
          return auth.currentUser.getIdToken(forceRefresh);
        }
        throw new Error('No authenticated user');
      },
      get currentUser() {
        return auth.currentUser;
      }
    };
  }
};

// Export the auth service
export const authService = createAuthService();
export const authMode = getAuthMode();

// Google provider for Firebase
export const googleProvider = new GoogleAuthProvider();

// Helper function to get current user token
export const getCurrentUserToken = async () => {
  try {
    if (authMode === 'local') {
      return localStorage.getItem('token');
    } else {
      return await authService.getIdToken();
    }
  } catch (error) {
    console.error('Error getting user token:', error);
    return null;
  }
};

// Helper function to check if user is authenticated
export const isAuthenticated = () => {
  if (authMode === 'local') {
    return !!(localStorage.getItem('token') && localStorage.getItem('userData'));
  } else {
    return !!authService.currentUser;
  }
};

// API helper with authentication
export const makeAuthenticatedRequest = async (url, options = {}) => {
  const token = await getCurrentUserToken();
  
  if (!token) {
    throw new Error('No authentication token available');
  }

  const headers = {
    'Authorization': `Bearer ${token}`,
    'Content-Type': 'application/json',
    ...options.headers
  };

  return fetch(url, {
    ...options,
    headers
  });
};

export default authService; 