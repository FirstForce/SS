import React, { useState, useEffect } from 'react';
import PhotoCard from '../../components/photosCards';
import { useAuth } from '../../contexts/AuthContext';

// Interface for device data
interface Device {
  id: string;
  device_id: string;
  device_name: string;
  device_status: string;
}

// Interface for photo data
interface Photo {
  id: string;
  timestamp: string;
  image_type: string;
  presigned_url: string;
  device_id: string;
  text: string;
}

// Interface for search parameters to store in localStorage
interface SearchParams {
  searchText: string;
  startDate: string;
  endDate: string;
  selectedDevice: string;
}

const STORAGE_KEY = 'photoSearchParams';

const PhotosPage: React.FC = () => {
  // Initialize state with values from localStorage if available
  const getStoredSearchParams = (): SearchParams => {
    const storedParams = localStorage.getItem(STORAGE_KEY);
    if (storedParams) {
      return JSON.parse(storedParams);
    }
    
    // Default values if nothing is stored
    const today = new Date();
    return {
      searchText: '',
      startDate: `${today.getFullYear()}-01-01`,
      endDate: today.toISOString().slice(0, 10),
      selectedDevice: 'all'
    };
  };
  
  const storedParams = getStoredSearchParams();
  
  const [searchText, setSearchText] = useState(storedParams.searchText);
  const [startDate, setStartDate] = useState(storedParams.startDate);
  const [endDate, setEndDate] = useState(storedParams.endDate);
  const [selectedDevice, setSelectedDevice] = useState(storedParams.selectedDevice);
  
  const [devices, setDevices] = useState<Device[]>([]);
  const [deviceError, setDeviceError] = useState<boolean>(false);
  const [deviceLoading, setDeviceLoading] = useState<boolean>(true);
  
  // States for photos
  const [photos, setPhotos] = useState<Photo[]>([]);
  const [photosLoading, setPhotosLoading] = useState<boolean>(false);
  const [photosError, setPhotosError] = useState<string | null>(null);
  
  const { token } = useAuth();

  // Save search parameters to localStorage whenever they change
  useEffect(() => {
    const searchParams: SearchParams = {
      searchText,
      startDate,
      endDate,
      selectedDevice
    };
    
    localStorage.setItem(STORAGE_KEY, JSON.stringify(searchParams));
  }, [searchText, startDate, endDate, selectedDevice]);

  // Fetch devices from API
  useEffect(() => {
    const fetchDevices = async () => {
      setDeviceLoading(true);
      setDeviceError(false);
      
      try {
        const response = await fetch('https://api.ss.stefaniordache.com/devices', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });
        
        if (!response.ok) {
          throw new Error('Failed to fetch devices');
        }
        
        const data = await response.json();
        setDevices(data);
        
      } catch (error) {
        console.error('Error fetching devices:', error);
        setDeviceError(true);
      } finally {
        setDeviceLoading(false);
      }
    };
    
    fetchDevices();
  }, [token]);

  // Initial search on page load
  useEffect(() => {
    handleSearch();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const handleSearch = async () => {
    setPhotosLoading(true);
    setPhotosError(null);
    
    try {
      // Convert dates to Unix timestamps (seconds)
      const startTimestamp = Math.floor(new Date(startDate).getTime() / 1000);
      const endTimestamp = Math.floor(new Date(endDate).getTime() / 1000) + 86399; // End of day

      // Build query parameters
      const queryParams = new URLSearchParams();
      queryParams.append('start', startTimestamp.toString());
      queryParams.append('end', endTimestamp.toString());
      
      if (searchText.trim()) {
        queryParams.append('text', searchText.trim());
      }
      
      // Only add device_id if a specific device (not "all") is selected
      if (selectedDevice !== 'all') {
        queryParams.append('device_id', selectedDevice);
      }
      
      // Make API request
      const response = await fetch(`https://api.ss.stefaniordache.com/photos?${queryParams.toString()}`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
      });
      
      if (!response.ok) {
        throw new Error(`Failed to fetch photos: ${response.status} ${response.statusText}`);
      }
      
      const data = await response.json();
      // Ensure data is an array, otherwise use empty array
      setPhotos(Array.isArray(data) ? data : []);
      
    } catch (error) {
      console.error('Error fetching photos:', error);
      setPhotosError((error as Error).message || 'Failed to load photos');
      setPhotos([]);
    } finally {
      setPhotosLoading(false);
    }
  };

  return (
    <div className="container mx-auto">
      <h1 className="text-2xl font-semibold text-sky-700 mb-6">Photos</h1>
      
      {/* Search and filter section */}
      <div className="bg-white p-4 rounded-lg shadow-sm mb-6">
        <div className="flex flex-wrap items-end gap-4">
          {/* Text search */}
          <div className="flex-1 min-w-[200px]">
            <label htmlFor="search" className="block text-sm font-medium text-gray-700 mb-1">
              Search Text
            </label>
            <input
              id="search"
              type="text"
              value={searchText}
              onChange={(e) => setSearchText(e.target.value)}
              placeholder="Search text in photos..."
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
            />
          </div>
          
          {/* Start date */}
          <div>
            <label htmlFor="start-date" className="block text-sm font-medium text-gray-700 mb-1">
              Start Date
            </label>
            <input
              id="start-date"
              type="date"
              value={startDate}
              onChange={(e) => setStartDate(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
            />
          </div>
          
          {/* End date */}
          <div>
            <label htmlFor="end-date" className="block text-sm font-medium text-gray-700 mb-1">
              End Date
            </label>
            <input
              id="end-date"
              type="date"
              value={endDate}
              onChange={(e) => setEndDate(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
            />
          </div>
          
          {/* Device dropdown - always shown, with or without device data */}
          {!deviceLoading && (
            <div>
              <label htmlFor="device" className="block text-sm font-medium text-gray-700 mb-1">
                Device
              </label>
              <select
                id="device"
                value={selectedDevice}
                onChange={(e) => setSelectedDevice(e.target.value)}
                className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent"
              >
                <option value="all">All</option>
                {!deviceError && devices.map(device => (
                  <option key={device.id} value={device.device_id}>
                    {device.device_id} - {device.device_name}
                  </option>
                ))}
              </select>
            </div>
          )}
          
          {/* Loading indicator for devices */}
          {deviceLoading && (
            <div className="flex items-end">
              <div className="h-10 flex items-center">
                <div className="animate-spin rounded-full h-5 w-5 border-t-2 border-b-2 border-sky-500 mr-2"></div>
                <span className="text-sm text-gray-500">Loading devices...</span>
              </div>
            </div>
          )}
          
          {/* Search button */}
          <div>
            <button
              onClick={handleSearch}
              disabled={photosLoading}
              className="px-4 py-2 bg-sky-600 text-white rounded-md hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {photosLoading ? 'Searching...' : 'Search'}
            </button>
          </div>
        </div>
      </div>
      
      {/* Photos section with fixed height and scroll */}
      <div className="bg-gray-50 p-4 rounded-lg shadow-sm overflow-y-auto max-h-[60vh]">
        {/* Loading state */}
        {photosLoading && (
          <div className="flex justify-center items-center h-40">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-sky-500"></div>
          </div>
        )}
        
        {/* Error state */}
        {!photosLoading && photosError && (
          <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded-md">
            <p className="font-medium">Error loading photos</p>
            <p className="mt-1">{photosError}</p>
          </div>
        )}
        
        {/* Results grid */}
        {!photosLoading && !photosError && (
          <>
            {(photos || []).length === 0 ? (
              <div className="text-center text-gray-500 py-10">
                No photos found matching your search criteria
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {(photos || []).map(photo => (
                  <PhotoCard 
                    key={photo.id}
                    imageUrl={photo.presigned_url}
                    extractedText={photo.text}
                    altText={`Photo from ${new Date(photo.timestamp).toLocaleDateString()}`}
                  />
                ))}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default PhotosPage; 