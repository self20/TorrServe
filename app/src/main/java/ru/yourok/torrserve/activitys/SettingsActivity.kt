package ru.yourok.torrserve.activitys

import android.content.Intent
import android.net.Uri
import android.os.Build
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.provider.Settings
import android.support.v7.app.AppCompatActivity
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_settings.*
import ru.yourok.torrserve.App
import ru.yourok.torrserve.BuildConfig
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.ServerSettings
import kotlin.concurrent.thread


class SettingsActivity : AppCompatActivity() {

    private val serverAddr = "http://localhost:8090"

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_settings)
        loadSettings()

        buttonOk.setOnClickListener {
            saveSettings()
            finish()
        }

        buttonCancel.setOnClickListener {
            finish()
        }

        buttonRetrieveSettings.setOnClickListener {
            loadSettings()
        }

        textViewVersion.setText("YouROK " + getText(R.string.app_name) + " ${BuildConfig.FLAVOR} ${BuildConfig.VERSION_NAME}")

        checkBoxShowWndInfo.setOnCheckedChangeListener { compoundButton, b ->
            if (b)
                checkPermission()
        }
    }

    override fun onResume() {
        super.onResume()
        checkServer()
    }

    override fun onPause() {
        super.onPause()
        isCheck = false
    }

    var isCheck = false
    fun checkServer() {
        synchronized(isCheck) {
            if (isCheck)
                return
            isCheck = true
        }
        thread {
            while (isCheck) {
                val echo = ServerApi.echo()
                runOnUiThread {
                    if (echo)
                        textViewStatus.setText(R.string.stat_server_is_running)
                    else
                        textViewStatus.setText(R.string.stat_server_is_not_running)
                }
                Thread.sleep(1000)
            }
        }
    }

    fun loadSettings(): Boolean {
        val addr = Preferences.getServerAddress()
        editTextServerAddr.setText(addr)

        val autoStart = Preferences.isAutoStart()
        checkBoxStartOnBoot.isChecked = autoStart

        val showWnd = Preferences.isShowState()
        checkBoxShowWndInfo.isChecked = showWnd

        val sets = ServerApi.readSettings()
        if (sets == null) {
            Toast.makeText(this, R.string.error_retrieving_settings, Toast.LENGTH_SHORT).show()
            return false
        }

        editTextCacheSize.setText(sets.CacheSize.toString())
        editTextPreloadBufferSize.setText(sets.PreloadBufferSize.toString())

        checkBoxDisableTCP.setChecked(sets.DisableTCP)
        checkBoxDisableUTP.setChecked(sets.DisableUTP)
        checkBoxDisableUPNP.setChecked(sets.DisableUPNP)
        checkBoxDisableDHT.setChecked(sets.DisableDHT)
        checkBoxDisableUpload.setChecked(sets.DisableUpload)

        editTextEncryption.setText(sets.Encryption.toString())
        editTextConnectionsLimit.setText(sets.ConnectionsLimit.toString())
        editTextDownloadRateLimit.setText(sets.DownloadRateLimit.toString())
        editTextUploadRateLimit.setText(sets.UploadRateLimit.toString())

        return true
    }

    fun saveSettings() {
        val addr = editTextServerAddr.text.toString()
        Preferences.setServerAddress(addr)
        val autoStart = checkBoxStartOnBoot.isChecked
        Preferences.setAutoStart(autoStart)
        val showWnd = checkBoxShowWndInfo.isChecked
        Preferences.setShowState(showWnd)

        val sets = ServerSettings(
                editTextCacheSize.text.toString().toInt(),
                editTextPreloadBufferSize.text.toString().toInt(),
                checkBoxDisableTCP.isChecked,
                checkBoxDisableUTP.isChecked,
                checkBoxDisableUPNP.isChecked,
                checkBoxDisableDHT.isChecked,
                checkBoxDisableUpload.isChecked,
                editTextEncryption.text.toString().toInt(),
                editTextDownloadRateLimit.text.toString().toInt(),
                editTextUploadRateLimit.text.toString().toInt(),
                editTextConnectionsLimit.text.toString().toInt())
        thread {
            val err = ServerApi.writeSettings(sets)
            if (err.isNotEmpty())
                Handler(Looper.getMainLooper()).post {
                    Toast.makeText(App.getContext(), R.string.error_sending_settings, Toast.LENGTH_SHORT).show()
                }
        }
    }

    fun checkPermission() {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.M) {
            if (!Settings.canDrawOverlays(this)) {
                val intent = Intent(Settings.ACTION_MANAGE_OVERLAY_PERMISSION,
                        Uri.parse("package:$packageName"))
                startActivity(intent)
            }
        }
    }
}
