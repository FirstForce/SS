package com.example.ss

import android.Manifest
import android.content.pm.PackageManager
import android.os.Bundle
import android.util.Base64
import android.util.Log
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.camera.core.*
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import com.example.ss.databinding.ActivityMainBinding
import org.eclipse.paho.android.service.MqttAndroidClient
import org.eclipse.paho.client.mqttv3.*
import java.io.ByteArrayOutputStream
import java.nio.charset.StandardCharsets
import java.util.concurrent.Executors
import kotlin.concurrent.fixedRateTimer

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var imageCapture: ImageCapture
    private val executor = Executors.newSingleThreadExecutor()

    private lateinit var mqttClient: MqttClient

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        if (allPermissionsGranted()) {
            startCamera()
            setupMQTT()
        } else {
            ActivityCompat.requestPermissions(
                this, REQUIRED_PERMISSIONS, REQUEST_CODE_PERMISSIONS
            )
        }
    }

    private fun setupMQTT() {
        val serverURI = "tcp://10.0.2.2:1883" // Or your broker
        val clientId = MqttClient.generateClientId()

        mqttClient = MqttClient(serverURI, clientId, null)

        mqttClient.setCallback(object : MqttCallback {
            override fun connectionLost(cause: Throwable?) {}
            override fun messageArrived(topic: String?, message: MqttMessage?) {}
            override fun deliveryComplete(token: IMqttDeliveryToken?) {}
        })

        val options = MqttConnectOptions()
        options.isCleanSession = true
        //options.isAutomaticReconnect = CONNECTION_RECONNECT
        //options.isCleanSession = CONNECTION_CLEAN_SESSION
        //options.userName = CLIENT_USER_NAME
        //options.password = CLIENT_PASSWORD.toCharArray()
        //options.connectionTimeout = CONNECTION_TIMEOUT
        //options.keepAliveInterval = CONNECTION_KEEP_ALIVE_INTERVA

        /*
        mqttClient.connect(options, null, object : IMqttActionListener {
            override fun onSuccess(asyncActionToken: IMqttToken?) {
                Toast.makeText(applicationContext, "MQTT Connected", Toast.LENGTH_SHORT).show()
                startImageCaptureLoop()
            }

            override fun onFailure(asyncActionToken: IMqttToken?, exception: Throwable?) {
                Toast.makeText(applicationContext, "MQTT Connection Failed", Toast.LENGTH_SHORT).show()
            }
        })
         */

        try {
            Log.d("SS", "Try Connect")
            mqttClient.connect(options)
        } catch (ex: MqttException) {
            Toast.makeText(applicationContext, "MQTT Connection Failed", Toast.LENGTH_SHORT).show()
        }
        publish("aa", "bbb")
        Toast.makeText(applicationContext, "MQTT Connected", Toast.LENGTH_SHORT).show()
        Log.d("SS", "CONNECTED!!")
        startImageCaptureLoop()
    }

    private fun startCamera() {
        val cameraProviderFuture = ProcessCameraProvider.getInstance(this)

        cameraProviderFuture.addListener({
            val cameraProvider: ProcessCameraProvider = cameraProviderFuture.get()
            val preview = Preview.Builder().build().also {
                it.setSurfaceProvider(binding.viewFinder.surfaceProvider)
            }

            imageCapture = ImageCapture.Builder()
                .setCaptureMode(ImageCapture.CAPTURE_MODE_MINIMIZE_LATENCY)
                .build()

            val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA

            cameraProvider.unbindAll()
            cameraProvider.bindToLifecycle(
                this, cameraSelector, preview, imageCapture
            )
        }, ContextCompat.getMainExecutor(this))
    }

    private fun startImageCaptureLoop() {
        fixedRateTimer("cameraTimer", false, 0L, 5000) {
            captureImageAndSend()
        }
    }

    private fun captureImageAndSend() {
        val outputOptions = ImageCapture.OutputFileOptions.Builder(ByteArrayOutputStream()).build()
        if (::imageCapture.isInitialized) {
            imageCapture.takePicture(ContextCompat.getMainExecutor(this), object :
                ImageCapture.OnImageCapturedCallback() {
                override fun onCaptureSuccess(imageProxy: ImageProxy) {
                    val buffer = imageProxy.planes[0].buffer
                    val bytes = ByteArray(buffer.remaining())
                    buffer.get(bytes)

                    val base64Image = Base64.encodeToString(bytes, Base64.DEFAULT)
                    sendToMQTT(base64Image.toByteArray(StandardCharsets.UTF_8))
                    imageProxy.close()
                }

                override fun onError(exception: ImageCaptureException) {
                    Log.e("CameraX", "Capture failed: ${exception.message}", exception)
                }
            })
        }
    }

    private fun sendToMQTT(data: ByteArray) {
        if (mqttClient.isConnected) {
            val message = MqttMessage(data)
            message.qos = 0
            mqttClient.publish("camera/image", message)
        }
    }

    private fun allPermissionsGranted() = REQUIRED_PERMISSIONS.all {
        ContextCompat.checkSelfPermission(baseContext, it) == PackageManager.PERMISSION_GRANTED
    }

    override fun onRequestPermissionsResult(
        requestCode: Int, permissions: Array<String>, grantResults: IntArray
    ) {
        super.onRequestPermissionsResult(requestCode, permissions, grantResults)
        if (requestCode == REQUEST_CODE_PERMISSIONS) {
            if (allPermissionsGranted()) {
                startCamera()
                setupMQTT()
            } else {
                Toast.makeText(this, "Permissions not granted", Toast.LENGTH_SHORT).show()
                finish()
            }
        }
    }

    companion object {
        private const val REQUEST_CODE_PERMISSIONS = 10
        private val REQUIRED_PERMISSIONS =
            arrayOf(Manifest.permission.CAMERA, Manifest.permission.INTERNET)
    }

    fun publish(topic: String, msg: String, qos: Int = 0) {
        try {
            val mqttMessage = MqttMessage(msg.toByteArray())
            mqttClient.publish(topic, mqttMessage.payload, qos, false)
            Log.d("SS", "Message published to topic `$topic`: $msg")
        } catch (e: MqttException) {
            Log.w("SS", "Error publishing to $topic: " + e.message, e)
            // Handle publishing failure
        }
    }
}