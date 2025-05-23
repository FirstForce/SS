import React from 'react';

interface DeviceCardProps {
  deviceId: string;
  deviceName: string;
  onNormalModeClick?: () => void;
  onLiveModeClick?: () => void;
}

const DeviceCard: React.FC<DeviceCardProps> = ({
  deviceId,
  deviceName,
  onNormalModeClick = () => console.log(`Normal mode clicked for device ${deviceId}`),
  onLiveModeClick = () => console.log(`Live mode clicked for device ${deviceId}`),
}) => {
  return (
    <div className="bg-white rounded-lg shadow-md overflow-hidden transition-all hover:shadow-lg p-4">
      <h3 className="text-lg font-medium text-gray-800 mb-4">{deviceId} - {deviceName}</h3>
      <div className="flex gap-2">
        <button 
          onClick={onNormalModeClick}
          className="flex-1 px-3 py-2 bg-sky-600 text-white rounded-md hover:bg-sky-700 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:ring-offset-2 transition-colors text-sm"
        >
          Normal Mode
        </button>
        <button 
          onClick={onLiveModeClick}
          className="flex-1 px-3 py-2 bg-emerald-600 text-white rounded-md hover:bg-emerald-700 focus:outline-none focus:ring-2 focus:ring-emerald-500 focus:ring-offset-2 transition-colors text-sm"
        >
          Live Mode
        </button>
      </div>
    </div>
  );
};

export default DeviceCard; 