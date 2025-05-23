import { BrowserRouter, Routes, Route, Navigate, Outlet, useNavigate } from 'react-router-dom';
import Navbar from './components/navbar';
import HomePage from './pages/homePage';
import LoginPage from './pages/loginPage';
import { AuthProvider, useAuth } from './contexts/AuthContext';

const Layout = () => {
  const navigate = useNavigate();
  const { isLoggedIn, logout } = useAuth();
  
  const navButtons = isLoggedIn 
    ? [
        {
          text: 'Logout',
          variant: 'outline' as const,
          onClick: () => {
            logout();
            navigate('/');
          }
        }
      ]
    : [
        {
          text: 'Login',
          variant: 'outline' as const,
          onClick: () => navigate('/login')
        },
        {
          text: 'Register',
          variant: 'primary' as const,
          onClick: () => console.log('Register clicked')
        }
      ];

  return (
    <>
      <Navbar 
        title="Security of Systems - First Force"
        buttons={navButtons}
      />
      <div className="pt-16 px-4">
        <Outlet />
      </div>
    </>
  );
};

const App = () => {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          <Route element={<Layout />}>
            <Route path="/" element={<HomePage />} />
            <Route path="/login" element={<LoginPage />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Route>
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
};

export default App;
