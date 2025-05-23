import React from 'react';
import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

interface ProtectedRouteProps {
  authRequired: boolean;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ authRequired }) => {
  const { isLoggedIn } = useAuth();
  
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