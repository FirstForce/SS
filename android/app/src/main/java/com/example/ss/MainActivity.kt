package com.example.ss

import android.Manifest
import android.content.pm.PackageManager
import android.os.Bundle
//import android.util.Base64
import android.util.Log
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.camera.core.CameraSelector
import androidx.camera.core.ImageCapture
import androidx.camera.core.ImageCaptureException
import androidx.camera.core.ImageProxy
import androidx.camera.core.Preview
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.app.ActivityCompat
import androidx.core.content.ContextCompat
import com.example.ss.databinding.ActivityMainBinding
//import kotlinx.coroutines.flow.internal.NoOpContinuation.context
//import org.bouncycastle.jce.provider.BouncyCastleProvider
//import org.bouncycastle.openssl.PEMKeyPair
//import org.bouncycastle.openssl.PEMParser
//import org.bouncycastle.openssl.jcajce.JcaPEMKeyConverter
import org.eclipse.paho.client.mqttv3.IMqttDeliveryToken
import org.eclipse.paho.client.mqttv3.MqttCallback
import org.eclipse.paho.client.mqttv3.MqttClient
import org.eclipse.paho.client.mqttv3.MqttConnectOptions
import org.eclipse.paho.client.mqttv3.MqttException
import org.eclipse.paho.client.mqttv3.MqttMessage
import java.io.BufferedInputStream
import java.io.ByteArrayOutputStream
import java.io.InputStream
import java.io.InputStreamReader
import java.nio.charset.StandardCharsets
import java.security.KeyPair
import java.security.KeyStore
import java.security.Security
import java.security.cert.X509Certificate
import java.util.concurrent.Executors
import kotlin.concurrent.fixedRateTimer
import java.security.cert.CertificateFactory
import javax.net.ssl.KeyManagerFactory
import javax.net.ssl.SSLContext
import javax.net.ssl.SSLSocketFactory
import javax.net.ssl.TrustManagerFactory
//import kotlin.coroutines.jvm.internal.CompletedContinuation.context
import java.security.KeyFactory
import java.security.spec.PKCS8EncodedKeySpec
import java.util.Base64



class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding
    private lateinit var imageCapture: ImageCapture
    private val executor = Executors.newSingleThreadExecutor()
    private var stopTransmission = true
    private var manualMode = false
    private var sendManual = false
    private var deviceID = (1..0xFFFF).random()

    private lateinit var mqttClient: MqttClient

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        binding.stateButton.setOnClickListener {
            stopTransmission = !stopTransmission
            //updateText()
            //Toast.makeText(this, "Stop transmission flag set to true", Toast.LENGTH_SHORT).show()
        }

        binding.captureButton.setOnClickListener {
            sendManual = true
            captureImageAndSendManual()
            Log.d("SS", "Sent Manual")
            //Toast.makeText(this, "Send manual flag set to true", Toast.LENGTH_SHORT).show()
        }

        binding.modeButton.setOnClickListener {
            manualMode = !manualMode
        }

        val handler = android.os.Handler(mainLooper)
        val updateRunnable = object : Runnable {
            override fun run() {
                runOnUiThread {
                    Log.d("SS", if (manualMode) "MANUAL" else "LIVE")
                    val status = if (manualMode) "Manual Mode" else "Live Mode"
                    binding.stateIndicator.text = status
                    binding.stateIndicator.setTextColor(
                        ContextCompat.getColor(this@MainActivity,
                            if (manualMode) android.R.color.holo_blue_dark
                            else if (stopTransmission) android.R.color.holo_red_dark
                            else android.R.color.holo_green_dark)
                    )
                }
                handler.postDelayed(this, 5000) // Repeat every 5 seconds
            }
        }
        handler.post(updateRunnable)

        if (allPermissionsGranted()) {
            startCamera()
            setupMQTT()
        } else {
            ActivityCompat.requestPermissions(
                this, REQUIRED_PERMISSIONS, REQUEST_CODE_PERMISSIONS
            )
        }
    }

    override fun onPause() {
        super.onPause()
        publish("device/id/$deviceID", "Device Disconnected")
    }

    override fun onStop() {
        super.onStop()
        publish("device/id/$deviceID", "Device Disconnected")
    }

    private fun setupMQTT() {
        val serverURI = "ssl://31.97.52.8:8883" // Or your broker
        val clientId = MqttClient.generateClientId()

        mqttClient = MqttClient(serverURI, clientId, null)

        mqttClient.setCallback(object : MqttCallback {
            override fun connectionLost(cause: Throwable?) {
                Log.w("SS", "Connection lost", cause)
            }
            override fun messageArrived(topic: String?, message: MqttMessage?) {
                val payload = message?.toString()

                Log.d("ss", topic + payload)

                if (payload == "start manual") {
                    manualMode = true
                    Log.d("SS", "Start Manual")
                } else if (payload == "start live") {
                    manualMode = false
                    Log.d("SS", "Start Live")
                }

                //updateText()

            }
            override fun deliveryComplete(token: IMqttDeliveryToken?) {
                Log.d("SS", "Delivery complete")
            }
        })

        val caCrtFile = resources.openRawResource(R.raw.ca)
        val keyFile = resources.openRawResource(R.raw.clientkey)
        val crtFile = resources.openRawResource(R.raw.clientcrt)

        Log.d("SS", "Getting ssl socket")
        val sslSocketFactory = getSocketFactory(caCrtFile, crtFile, keyFile, "")
        Log.d("SS", "Got ssl socket")

        val options = MqttConnectOptions()
        options.isCleanSession = false
        options.socketFactory = sslSocketFactory
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

        val devname = android.os.Build.MANUFACTURER + " " + android.os.Build.MODEL
        publish("register/$deviceID", devname)
        Toast.makeText(applicationContext, "MQTT Connected", Toast.LENGTH_SHORT).show()
        Log.d("SS", "CONNECTED!!")

        mqttClient.subscribe("setup/$deviceID")

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
        if (::imageCapture.isInitialized && !stopTransmission && !manualMode) {
            imageCapture.takePicture(ContextCompat.getMainExecutor(this), object :
                ImageCapture.OnImageCapturedCallback() {
                override fun onCaptureSuccess(imageProxy: ImageProxy) {
                    val buffer = imageProxy.planes[0].buffer
                    val bytes = ByteArray(buffer.remaining())
                    buffer.get(bytes)

                    //val base64Image = Base64.encodeToString(bytes, Base64.DEFAULT)
                    //sendToMQTT(base64Image.toByteArray(StandardCharsets.UTF_8))
                    sendToMQTT(bytes)
                    imageProxy.close()
                }

                override fun onError(exception: ImageCaptureException) {
                    Log.e("CameraX", "Capture failed: ${exception.message}", exception)
                }
            })
        }
    }

    private fun captureImageAndSendManual() {
        val outputOptions = ImageCapture.OutputFileOptions.Builder(ByteArrayOutputStream()).build()
        if (::imageCapture.isInitialized) {
            imageCapture.takePicture(ContextCompat.getMainExecutor(this), object :
                ImageCapture.OnImageCapturedCallback() {
                override fun onCaptureSuccess(imageProxy: ImageProxy) {
                    val buffer = imageProxy.planes[0].buffer
                    val bytes = ByteArray(buffer.remaining())
                    buffer.get(bytes)

                    //val base64Image = Base64.encodeToString(bytes, Base64.DEFAULT)
                    //sendToMQTT(base64Image.toByteArray(StandardCharsets.UTF_8))
                    sendToMQTT(bytes)
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
            Thread {
                val message = MqttMessage(data)
                message.qos = 0
                mqttClient.publish("photos/$deviceID", message)
            }.start()
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
            Thread {
                mqttClient.publish(topic, mqttMessage.payload, qos, false)
            }.start()
            Log.d("SS", "Message published to topic `$topic`: $msg")
        } catch (e: MqttException) {
            Log.w("SS", "Error publishing to $topic: " + e.message, e)
            // Handle publishing failure
        }
    }


    private fun getSocketFactory(
        caInput: InputStream,
        certInput: InputStream,
        keyInput: InputStream,
        password: String = ""
    ): SSLSocketFactory {
        // Load CA certificate
        val caCert = CertificateFactory.getInstance("X.509")
            .generateCertificate(caInput) as X509Certificate

        val trustStore = KeyStore.getInstance(KeyStore.getDefaultType())
        trustStore.load(null)
        trustStore.setCertificateEntry("ca-cert", caCert)

        val tmf = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm())
        tmf.init(trustStore)

        // Load client certificate
        val clientCert = CertificateFactory.getInstance("X.509")
            .generateCertificate(certInput) as X509Certificate

        // Load private key from PEM
        val keyBytes = keyInput.bufferedReader().useLines { lines ->
            lines
                .filter { !it.startsWith("-----") }
                .joinToString("") { it.trim() }
        }
        val decodedKey = Base64.getDecoder().decode(keyBytes)
        val keySpec = PKCS8EncodedKeySpec(decodedKey)
        val privateKey = KeyFactory.getInstance("RSA").generatePrivate(keySpec)

        val keyStore = KeyStore.getInstance(KeyStore.getDefaultType())
        keyStore.load(null)
        keyStore.setKeyEntry("client-key", privateKey, password.toCharArray(), arrayOf(clientCert))

        val kmf = KeyManagerFactory.getInstance(KeyManagerFactory.getDefaultAlgorithm())
        kmf.init(keyStore, password.toCharArray())

        // Create SSL context
        val context = SSLContext.getInstance("TLS")
        context.init(kmf.keyManagers, tmf.trustManagers, null)

        return context.socketFactory
    }

}