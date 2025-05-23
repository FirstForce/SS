import React from 'react';

interface NavbarProps {
  title: string;
}

const Navbar: React.FC<NavbarProps> = ({ title }) => {
  return (
    <nav className="fixed top-0 left-0 right-0 bg-sky-50 shadow-sm z-50">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div className="w-24">
          {/* Placeholder for left side content/logo */}
        </div>
        
        <h1 className="text-xl font-semibold text-sky-700">{title}</h1>
        
        <div className="w-24 flex justify-end">
          {/* Placeholder for right side actions */}
        </div>
      </div>
    </nav>
  );
};

export default Navbar; 