"use client";

import { createContext, useContext, useState, useEffect } from "react";
import { authService, authMode, getCurrentUserToken } from "@/services/authService";

export const AuthContext = createContext({});

export function AuthProvider({ children }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    console.log(`🔐 AuthContext initializing in ${authMode.toUpperCase()} mode`);
    
    const unsubscribe = authService.onAuthStateChanged(async (user) => {
      if (user) {
        setUser(user);
        
        // Store/update token for API calls
        try {
          const token = await getCurrentUserToken();
          if (token) {
            localStorage.setItem("token", token);
            console.log(`✅ Token stored for user: ${user.email || user.uid}`);
          }
        } catch (error) {
          console.error("Error getting token:", error);
        }
      } else {
        setUser(null);
        // Clear token when user logs out
        localStorage.removeItem("token");
        localStorage.removeItem("userData");
      }
      setLoading(false);
    });

    return () => unsubscribe();
  }, []);

  const logout = async () => {
    try {
      await authService.signOut();
      setUser(null);
      return true;
    } catch (error) {
      console.error("Error signing out:", error);
      throw error;
    }
  };

  const value = {
    user,
    loading,
    logout,
    authMode, // Expose auth mode to components
  };

  return (
    <AuthContext.Provider value={value}>
      {!loading && children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
