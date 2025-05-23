import React, { useState } from 'react';
import fallbackImage from '../../assets/photo-fallback.svg';

interface PhotoCardProps {
  imageUrl: string;
  altText?: string;
  extractedText?: string;
}

const PhotoCard: React.FC<PhotoCardProps> = ({ imageUrl, altText = 'Photo', extractedText = '' }) => {
  const [isZoomed, setIsZoomed] = useState(false);
  const [imageError, setImageError] = useState(false);

  const handleImageError = () => {
    setImageError(true);
  };

  const toggleZoom = () => {
    setIsZoomed(!isZoomed);
  };

  // Handle click outside the zoomed image to close it
  const handleModalClick = (e: React.MouseEvent<HTMLDivElement>) => {
    if (e.target === e.currentTarget) {
      setIsZoomed(false);
    }
  };

  return (
    <>
      <div className="bg-white rounded-lg shadow-md overflow-hidden transition-all hover:shadow-lg">
        <div className="relative h-48 cursor-pointer" onClick={toggleZoom}>
          <img 
            src={imageError ? fallbackImage : imageUrl} 
            alt={altText}
            onError={handleImageError}
            className="w-full h-full object-cover"
          />
        </div>

        {extractedText && (
          <div className="p-3 border-t border-gray-100">
            <p className="text-sm text-gray-600 truncate">{extractedText}</p>
          </div>
        )}
      </div>

      {/* Improved zoom modal overlay with animations and better styling */}
      {isZoomed && (
        <div 
          className="fixed inset-0 bg-black bg-opacity-60 backdrop-blur-sm flex items-center justify-center z-50 transition-opacity duration-300 ease-in-out"
          onClick={handleModalClick}
        >
          <div 
            className="relative bg-white rounded-xl shadow-2xl max-w-4xl max-h-[90vh] overflow-hidden transform transition-all duration-300 ease-in-out animate-scaleIn"
          >
            <div className="absolute top-0 right-0 left-0 bg-gradient-to-b from-black/50 to-transparent h-20 z-10 flex justify-between items-start p-4">
              <div className="text-white text-lg font-medium truncate pr-10">{altText}</div>
              <button 
                className="bg-white/20 hover:bg-white/40 text-white rounded-full p-2 backdrop-blur-sm transition-all duration-200"
                onClick={toggleZoom}
              >
                <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            
            <div className="p-4 pt-20">
              <img 
                src={imageError ? fallbackImage : imageUrl} 
                alt={altText}
                className="max-w-full max-h-[65vh] object-contain mx-auto rounded-md"
              />
            </div>
            
            {extractedText && (
              <div className="bg-gray-50 p-6 border-t border-gray-100">
                <h3 className="text-sm font-medium text-gray-500 mb-2">Extracted Text</h3>
                <p className="text-gray-800 text-base">{extractedText}</p>
              </div>
            )}
          </div>
        </div>
      )}
    </>
  );
};

export default PhotoCard; 