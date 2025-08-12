# Security System - Android Application

A native Kotlin mobile client that extends monitoring capabilities to on-the-go scenarios, allowing remote system control from the UI. The Android application serves as a security device, capturing photos and sending them to the server via MQTT protocol with secure SSL/TLS communication.

## Overview

This Android application implements a mobile security device that can capture photos and communicate with the central security system. It provides both manual and automatic photo capture modes, real-time communication with the server, and secure data transmission using TLS encryption.

## Key Features

### Camera Integration
- **CameraX API**: Uses the Android CameraX API to preview and capture images
- **Image Processing**: Converts ImageProxy to Bitmap, compresses to JPEG, and encodes in Base64
- **Real-time Preview**: Live camera preview with capture functionality

### MQTT Communication
- **Eclipse Paho Client**: Integrates with the Eclipse Paho MQTT client to publish images
- **Topic-based Messaging**: Publishes photos to `photos/ID` topic
- **Device Registration**: Sends registration messages to notify server of device presence
- **Mode Control**: Receives mode changes via `device/setup/ID` topic

### Secure Communication
- **SSL/TLS Support**: Supports SSL/TLS connections using CA, certificate, and private key files
- **Certificate Storage**: Security certificates stored in the `res/raw` directory
- **Encrypted Transmission**: All data transmitted over encrypted channels

### User Interaction
- **Manual Capture**: Button to capture and send image manually
- **Mode Switching**: Buttons to change application mode (live/manual)
- **Transmission Control**: Start or stop transmission functionality
- **Device Identification**: Stores unique device identifier (UUID) persistently

## Application Architecture

### MainActivity
- **Camera Setup**: Handles camera initialization and configuration
- **MQTT Connection**: Manages MQTT client connection and lifecycle
- **User Interactions**: Processes UI button clicks and user input
- **UI Updates**: Updates interface based on application state

### Image Processing
- **ImageProxy to Bitmap**: Converts camera ImageProxy to Bitmap format
- **JPEG Compression**: Compresses images to reduce transmission size
- **Base64 Encoding**: Encodes images for MQTT transmission

### MQTT Logic
- **Asynchronous Publishing**: Publishes on separate thread to avoid blocking UI
- **Topic Subscription**: Subscribes to topics with message callbacks
- **Message Handling**: Processes incoming messages for mode changes

### Threading
- **Handler and Runnables**: Uses Handler and Runnables for periodic UI updates
- **Background Processing**: Image processing and MQTT operations run in background

## MQTT Message Types

### Register Message
- **Purpose**: Notify server that device is connected
- **Content**: Device ID and connection status
- **Topic**: Device registration topic

### Photos
- **Purpose**: Send captured images to server
- **Content**: Base64 encoded JPEG image
- **Topic**: `photos/ID` where ID is the device identifier

### Manual/Live Mode
- **Purpose**: Receive mode configuration from server
- **Content**: Mode setting (Live or Manual)
- **Topic**: `device/setup/ID`

### Disconnect Message
- **Purpose**: Notify server when application is closed
- **Content**: Device ID and disconnect status
- **Topic**: Device disconnect topic

## Technology Stack

- **Kotlin**: Primary programming language
- **Android CameraX**: Camera functionality and image capture
- **Eclipse Paho MQTT**: Real-time messaging protocol
- **Jetpack Compose**: Modern UI framework
- **SharedPreferences**: Persistent storage for device settings
- **BouncyCastle**: SSL/TLS security implementation

## Development Setup

### Prerequisites
- Android Studio Arctic Fox or later
- Android SDK API level 21 or higher
- Target SDK 35
- Kotlin 1.9.0 or later

### Dependencies
```kotlin
// Camera dependencies
implementation("androidx.camera:camera-core:1.3.0")
implementation("androidx.camera:camera-extensions:1.3.0")
implementation("androidx.camera:camera-camera2:1.3.0")
implementation("androidx.camera:camera-lifecycle:1.3.0")
implementation("androidx.camera:camera-view:1.3.0")

// MQTT dependencies
implementation("org.eclipse.paho:org.eclipse.paho.client.mqttv3:1.2.5")
implementation("org.eclipse.paho:org.eclipse.paho.android.service:1.1.1")

// Security dependencies
implementation("org.bouncycastle:bcpkix-jdk15on:1.59")
```

### Installation
1. Clone the repository and navigate to android directory
2. Open the project in Android Studio
3. Configure MQTT broker settings in the application
4. Add SSL certificates to `res/raw` directory
5. Build and install on Android device

## Configuration

### MQTT Settings
- **Broker URL**: Configure MQTT broker address
- **Port**: SSL/TLS port (typically 8883)
- **Client ID**: Unique device identifier
- **Topics**: Configure publish/subscribe topics

### Security Certificates
Place the following files in `res/raw` directory:
- `ca.crt`: Certificate Authority certificate
- `client.crt`: Client certificate
- `client.key`: Client private key

### Device Settings
- **Device ID**: Unique identifier for the device
- **Capture Mode**: Manual or Live mode
- **Transmission Interval**: For Live mode operation

## Usage

### Manual Mode
1. Launch the application
2. Grant camera permissions when prompted
3. Use "Capture" button to take photos manually
4. Photos are automatically sent to the server

### Live Mode
1. Switch to Live mode using the mode button
2. Application automatically captures photos at intervals
3. Photos are continuously transmitted to the server
4. Use "Stop" button to halt automatic capture

### Device Control
1. Server can send mode change commands
2. Application automatically switches between modes
3. Device status is reported to server
4. Disconnect messages sent when app closes

## Security Considerations

- **Certificate Validation**: All SSL/TLS certificates are validated
- **Secure Storage**: Sensitive data stored in Android Keystore
- **Permission Management**: Minimal required permissions
- **Data Encryption**: All transmitted data is encrypted

## Troubleshooting

### Common Issues
- **Camera Permissions**: Ensure camera permissions are granted
- **MQTT Connection**: Check broker URL and network connectivity
- **Certificate Issues**: Verify SSL certificates are properly configured
- **Image Transmission**: Check image size and compression settings

### Debug Information
- Check Android logs for MQTT connection status
- Verify camera initialization in logs
- Monitor image processing performance
- Check network connectivity for server communication

## Testing

The application includes comprehensive testing:
- **Unit Tests**: Core functionality testing
- **Integration Tests**: MQTT communication testing
- **UI Tests**: User interface testing
- **Security Tests**: SSL/TLS configuration testing

## Build and Deployment

### Debug Build
```bash
./gradlew assembleDebug
```

### Release Build
```bash
./gradlew assembleRelease
```

### APK Location
Built APK files are located in `app/build/outputs/apk/` directory.
