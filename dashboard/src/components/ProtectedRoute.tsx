import { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';
import { UserRole } from '../lib/api';

interface ProtectedRouteProps {
  readonly children: ReactNode;
  readonly requiredRole?: UserRole;
  readonly allowedRoles?: UserRole[];
}

// Role hierarchy: admin > rssi > operator > analyst > viewer
const roleHierarchy: Record<UserRole, number> = {
  admin: 5,
  rssi: 4,
  operator: 3,
  analyst: 2,
  viewer: 1,
};

function hasRequiredRole(userRole: UserRole | undefined, requiredRole?: UserRole, allowedRoles?: UserRole[]): boolean {
  if (!userRole) return false;

  // If allowedRoles is specified, check if user's role is in the list
  if (allowedRoles && allowedRoles.length > 0) {
    return allowedRoles.includes(userRole);
  }

  // If requiredRole is specified, check role hierarchy
  if (requiredRole) {
    return roleHierarchy[userRole] >= roleHierarchy[requiredRole];
  }

  // No role requirement - any authenticated user is allowed
  return true;
}

export function ProtectedRoute({ children, requiredRole, allowedRoles }: ProtectedRouteProps) {
  const { isAuthenticated, isLoading, user } = useAuth();
  const location = useLocation();

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-100">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading...</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    // Redirect to login, preserving the intended destination
    return <Navigate to="/login" state={{ from: location }} replace />;
  }

  // Check role-based access
  if ((requiredRole || allowedRoles) && !hasRequiredRole(user?.role as UserRole, requiredRole, allowedRoles)) {
    // User doesn't have required role - show unauthorized
    return (
      <div className="min-h-screen flex items-center justify-center bg-gray-900">
        <div className="text-center">
          <h1 className="text-4xl font-bold text-red-500">403</h1>
          <p className="mt-4 text-gray-400">Access Denied</p>
          <p className="mt-2 text-sm text-gray-500">You don't have permission to access this page.</p>
          <a href="/dashboard" className="mt-4 inline-block text-primary-400 hover:text-primary-300">
            Return to Dashboard
          </a>
        </div>
      </div>
    );
  }

  return <>{children}</>;
}
