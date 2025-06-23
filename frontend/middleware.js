import { NextResponse } from "next/server";
import { authService, authMode, isAuthenticated } from "@/services/authService";

// Token refresh function to be called client-side
export const refreshToken = async () => {
  try {
    if (authMode === 'local') {
      // For local mode, just return the stored token
      const token = localStorage.getItem('token');
      return token;
    } else {
      // Firebase mode
      const currentUser = authService.currentUser;
      if (currentUser) {
        const token = await currentUser.getIdToken(true); // Force token refresh
        localStorage.setItem("token", token);
        return token;
      }
    }
  } catch (error) {
    console.error("Error refreshing token:", error);
  }
  return null;
};

export async function middleware(request) {
  // List of protected routes that require authentication
  const protectedPaths = [
    "/dashboard",
    "/chat",
    "/create-avatar",
    "/canvas",
    "/profile",
  ];

  // List of public routes that should always be accessible
  const publicPaths = [
    "/",
    "/about",
    "/login",
    "/signup",
    "/marketplace", // Making marketplace public
  ];

  const path = request.nextUrl.pathname;

  // Allow access to public paths without authentication
  if (publicPaths.some((publicPath) => path.startsWith(publicPath))) {
    return NextResponse.next();
  }

  // Check if the current path is protected
  const isProtectedPath = protectedPaths.some((protectedPath) =>
    path.startsWith(protectedPath)
  );

  if (isProtectedPath) {
    // For client-side authentication check, we'll rely on the auth context
    // This middleware primarily handles static route protection
    
    // Check for token in cookies or headers for server-side validation
    const token = request.headers.get('authorization') || 
                  request.cookies.get('token')?.value ||
                  request.nextUrl.searchParams.get('token');

    // If no token is found, redirect to login
    if (!token) {
      console.log(`🔒 No token found for protected path: ${path}, redirecting to login`);
      return NextResponse.redirect(new URL("/login", request.url));
    }

    // For local mode, we can do basic token validation
    if (token.startsWith('local_')) {
      // Token appears to be a local token, allow access
      console.log(`🔓 Local token found for path: ${path}`);
      return NextResponse.next();
    }

    // For Firebase tokens or other cases, let the client-side auth handle validation
    console.log(`🔓 Token found for path: ${path}`);
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/dashboard/:path*",
    "/chat/:path*",
    "/marketplace/:path*",
    "/create-avatar/:path*",
    "/canvas/:path*",
    "/profile/:path*",
  ],
};
