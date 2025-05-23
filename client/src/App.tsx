import { BrowserRouter, Routes, Route, Navigate, Outlet } from 'react-router-dom';
import Navbar from './components/navbar';
import HomePage from './pages/homePage';

const Layout = () => {
  return (
    <>
      <Navbar title="Security of Systems - First Force" />
      <div className="pt-16 px-4">
        <Outlet />
      </div>
    </>
  );
};

const App = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route path="/" element={<HomePage />} />
          <Route path="*" element={<Navigate to="/" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
};

export default App;
