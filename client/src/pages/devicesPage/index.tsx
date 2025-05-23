import React, { useState, useEffect } from 'react';
import DeviceCard from '../../components/devicesCards';
import { useAuth } from '../../contexts/AuthContext';

// Interface for device data
interface Device {
  id: string;
  device_id: string;
  device_name: string;
  device_status: string;
}

// Interface for tracking device action states
interface DeviceActionState {
  [deviceId: string]: {
    loading: boolean;
    success: boolean;
    error: string | null;
    lastMode: string | null;
  };
}

const DevicesPage: React.FC = () => {
  const [devices, setDevices] = useState<Device[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  
  // Track status of actions for individual devices
  const [deviceActionStates, setDeviceActionStates] = useState<DeviceActionState>({});
  
  const { token } = useAuth();
  
  // Fetch devices from API
  useEffect(() => {
    const fetchDevices = async () => {
      setLoading(true);
      setError(null);
      
      try {
        const response = await fetch('https://api.ss.stefaniordache.com/devices', {
          method: 'GET',
          headers: {
            'Authorization': `Bearer ${token}`,
            'Content-Type': 'application/json',
          },
        });
        
        if (!response.ok) {
          throw new Error(`Failed to fetch devices: ${response.status} ${response.statusText}`);
        }
        
        const data = await response.json();
        
        // Filter devices to only include active ones
        const activeDevices = Array.isArray(data) 
          ? data.filter(device => device.device_status === "active") 
          : [];
        
        setDevices(activeDevices);
        
        // Initialize action states for all active devices
        const initialActionStates: DeviceActionState = {};
        activeDevices.forEach(device => {
          initialActionStates[device.device_id] = {
            loading: false,
            success: false,
            error: null,
            lastMode: null,
          };
        });
        setDeviceActionStates(initialActionStates);
        
      } catch (error) {
        console.error('Error fetching devices:', error);
        setError((error as Error).message || 'Failed to load devices');
        setDevices([]);
      } finally {
        setLoading(false);
      }
    };
    
    fetchDevices();
  }, [token]);

  // Clear success/error messages after delay
  useEffect(() => {
    const timeout = setTimeout(() => {
      // Clear success/error messages but keep the last mode for displaying status
      setDeviceActionStates(prevStates => {
        const newStates = { ...prevStates };
        Object.keys(newStates).forEach(deviceId => {
          if (newStates[deviceId].success || newStates[deviceId].error) {
            newStates[deviceId] = {
              ...newStates[deviceId],
              success: false,
              error: null,
            };
          }
        });
        return newStates;
      });
    }, 3000);
    
    return () => clearTimeout(timeout);
  }, [deviceActionStates]);

  const switchDeviceMode = async (deviceId: string, mode: 'manual' | 'live') => {
    // Find device object with matching device_id to get its id
    const device = devices.find(d => d.device_id === deviceId);
    if (!device) return;
    
    // Set loading state
    setDeviceActionStates(prev => ({
      ...prev,
      [deviceId]: {
        ...prev[deviceId],
        loading: true,
        success: false,
        error: null,
      }
    }));
    
    try {
      const response = await fetch('https://api.ss.stefaniordache.com/devices/switch', {
        method: 'POST',
        headers: {
          'Authorization': `Bearer ${token}`,
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          id: deviceId,
          mode: mode
        })
      });
      
      if (!response.ok) {
        throw new Error(`Failed to switch mode: ${response.status} ${response.statusText}`);
      }
      
      // Set success state
      setDeviceActionStates(prev => ({
        ...prev,
        [deviceId]: {
          loading: false,
          success: true,
          error: null,
          lastMode: mode,
        }
      }));
      
    } catch (error) {
      console.error(`Error switching device ${deviceId} to ${mode} mode:`, error);
      
      // Set error state
      setDeviceActionStates(prev => ({
        ...prev,
        [deviceId]: {
          loading: false,
          success: false,
          error: (error as Error).message || `Failed to switch to ${mode} mode`,
          lastMode: prev[deviceId].lastMode,
        }
      }));
    }
  };

  const handleNormalMode = (deviceId: string) => {
    switchDeviceMode(deviceId, 'manual');
  };
  
  const handleLiveMode = (deviceId: string) => {
    switchDeviceMode(deviceId, 'live');
  };
  
  return (
    <div className="container mx-auto">
      <h1 className="text-2xl font-semibold text-sky-700 mb-6">Devices</h1>
      
      {/* Devices grid with fixed height and scroll */}
      <div className="bg-gray-50 p-4 rounded-lg shadow-sm overflow-y-auto max-h-[60vh]">
        {/* Loading state */}
        {loading && (
          <div className="flex justify-center items-center h-40">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-sky-500"></div>
          </div>
        )}
        
        {/* Error state */}
        {!loading && error && (
          <div className="bg-red-50 border border-red-200 text-red-700 p-4 rounded-md">
            <p className="font-medium">You don't have permission to view connected devices</p>
            <p className="mt-1 text-sm">Please contact your administrator for access</p>
          </div>
        )}
        
        {/* Success state - display devices */}
        {!loading && !error && (
          <>
            {devices.length === 0 ? (
              <div className="text-center text-gray-500 py-10">
                No devices found
              </div>
            ) : (
              <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
                {devices.map(device => {
                  const actionState = deviceActionStates[device.device_id] || {
                    loading: false,
                    success: false,
                    error: null,
                    lastMode: null,
                  };
                  
                  return (
                    <div key={device.id} className="relative">
                      <DeviceCard 
                        deviceId={device.device_id}
                        deviceName={device.device_name}
                        onNormalModeClick={() => handleNormalMode(device.device_id)}
                        onLiveModeClick={() => handleLiveMode(device.device_id)}
                      />
                      
                      {/* Loading overlay */}
                      {actionState.loading && (
                        <div className="absolute inset-0 bg-white bg-opacity-70 flex items-center justify-center rounded-lg">
                          <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-sky-500"></div>
                        </div>
                      )}
                      
                      {/* Success notification */}
                      {actionState.success && (
                        <div className="absolute top-0 right-0 left-0 bg-green-100 text-green-800 text-sm p-2 rounded-t-lg text-center">
                          Successfully switched to {actionState.lastMode === 'manual' ? 'Normal' : 'Live'} mode
                        </div>
                      )}
                      
                      {/* Error notification */}
                      {actionState.error && (
                        <div className="absolute top-0 right-0 left-0 bg-red-100 text-red-800 text-sm p-2 rounded-t-lg text-center">
                          {actionState.error}
                        </div>
                      )}
                      
                      {/* Mode indicator */}
                      {!actionState.loading && !actionState.error && !actionState.success && actionState.lastMode && (
                        <div className={`absolute bottom-0 right-0 left-0 ${
                          actionState.lastMode === 'manual' ? 'bg-sky-100 text-sky-800' : 'bg-emerald-100 text-emerald-800'
                        } text-xs p-1 text-center`}>
                          Current mode: {actionState.lastMode === 'manual' ? 'Normal' : 'Live'}
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default DevicesPage; 