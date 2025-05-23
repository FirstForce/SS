import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

interface ProtectedRouteProps {
  authRequired: boolean;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ authRequired }) => {
  const { isLoggedIn, loading } = useAuth();
  
  // While auth state is loading, show nothing (or could add a loading spinner here)
  if (loading) {
    return <div className="flex justify-center items-center min-h-[60vh]">
      <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-sky-500"></div>
    </div>;
  }
  
  // If auth is required and user is not logged in, redirect to login
  if (authRequired && !isLoggedIn) {
    return <Navigate to="/login" replace />;
  }
  
  // If auth is not required and user is logged in, redirect to root
  if (!authRequired && isLoggedIn) {
    return <Navigate to="/" replace />;
  }
  
  // Otherwise render the children
  return <Outlet />;
};

export default ProtectedRoute; 