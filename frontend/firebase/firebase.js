// src/firebase/firebase.js
import { initializeApp } from "firebase/app";
import { getAuth } from "firebase/auth";
import { getStorage } from "firebase/storage";

// Check if Firebase should be enabled
const isFirebaseEnabled = () => {
  const firebaseEnabled = process.env.NEXT_PUBLIC_FIREBASE_ENABLE !== 'false';
  const hasFirebaseConfig = !!(
    process.env.NEXT_PUBLIC_FIREBASE_API_KEY &&
    process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN &&
    process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID
  );
  
  return firebaseEnabled && hasFirebaseConfig;
};

let app = null;
let auth = null;
let storage = null;

// Only initialize Firebase if enabled and configured
if (isFirebaseEnabled()) {
  try {
    const firebaseConfig = {
      apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
      authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
      projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
      storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
      messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
      appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
      measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID,
    };

    app = initializeApp(firebaseConfig);
    auth = getAuth(app);
    storage = getStorage(app);
    
    console.log("🔥 Firebase initialized successfully");
  } catch (error) {
    console.warn("⚠️ Firebase initialization failed:", error.message);
    console.log("📁 Falling back to local auth mode");
  }
} else {
  console.log("📁 Firebase disabled, using local auth mode");
}

export { auth, storage };
export default app;
