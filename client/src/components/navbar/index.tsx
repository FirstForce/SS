import React from 'react';
import Button from '../button';

interface ButtonProps {
  text: string;
  onClick?: () => void;
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'sm' | 'md' | 'lg';
}

interface NavbarProps {
  title: string;
  buttons?: ButtonProps[];
}

const Navbar: React.FC<NavbarProps> = ({ title, buttons = [] }) => {
  return (
    <nav className="fixed top-0 left-0 right-0 bg-sky-50 shadow-sm z-50">
      <div className="container mx-auto px-4 py-3 flex items-center justify-between">
        <div className="w-24">
          {/* Placeholder for left side content/logo */}
        </div>
        
        <h1 className="text-xl font-semibold text-sky-700">{title}</h1>
        
        <div className="flex justify-end space-x-2">
          {buttons.map((button, index) => (
            <Button 
              key={index}
              text={button.text}
              onClick={button.onClick}
              variant={button.variant || 'outline'}
              size="sm"
            />
          ))}
        </div>
      </div>
    </nav>
  );
};

export default Navbar; 