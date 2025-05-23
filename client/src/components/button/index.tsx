import React from 'react';

interface ButtonProps {
  text: string;
  onClick?: () => void;
  variant?: 'primary' | 'secondary' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  type?: 'button' | 'submit' | 'reset';
}

const Button: React.FC<ButtonProps> = ({
  text,
  onClick,
  variant = 'primary',
  size = 'md',
  type = 'button',
}) => {
  // Base styles
  const baseStyles = 'inline-flex items-center justify-center rounded-md font-medium transition-colors focus:outline-none focus:ring-2 focus:ring-sky-400 focus:ring-offset-2';
  
  // Size variations
  const sizeStyles = {
    sm: 'px-3 py-1 text-sm',
    md: 'px-4 py-2 text-base',
    lg: 'px-6 py-3 text-lg',
  };
  
  // Variant styles
  const variantStyles = {
    primary: 'bg-sky-600 text-white hover:bg-sky-700',
    secondary: 'bg-sky-100 text-sky-800 hover:bg-sky-200',
    outline: 'bg-transparent border border-sky-600 text-sky-600 hover:bg-sky-50',
  };
  
  return (
    <button
      className={`${baseStyles} ${sizeStyles[size]} ${variantStyles[variant]}`}
      onClick={onClick}
      type={type}
    >
      {text}
    </button>
  );
};

export default Button; 