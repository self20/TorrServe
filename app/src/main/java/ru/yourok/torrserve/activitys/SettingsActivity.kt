package ru.yourok.torrserve.activitys

import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.support.v7.app.AppCompatActivity
import android.widget.ArrayAdapter
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_settings.*
import ru.yourok.torrserve.App
import ru.yourok.torrserve.BuildConfig
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.ServerSettings
import ru.yourok.torrserve.utils.Player
import ru.yourok.torrserve.utils.Players
import kotlin.concurrent.thread


class SettingsActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_settings)

        buttonOk.setOnClickListener {
            saveSettings()
            finish()
        }

        buttonCancel.setOnClickListener {
            finish()
        }

        buttonRetrieveSettings.setOnClickListener {
            loadSettings(true)
        }

        val plist = Players.getList()
        plist.add(0, Player(getString(R.string.default_player), "0"))
        plist.add(1, Player(getString(R.string.choose_player), "1"))

        val adp1 = ArrayAdapter<Player>(this, android.R.layout.simple_list_item_1, plist)
        adp1.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        spinnerPlayer.setAdapter(adp1)

        textViewVersion.setText("YouROK " + getText(R.string.app_name) + " ${BuildConfig.FLAVOR} ${BuildConfig.VERSION_NAME}")

        loadSettings(false)
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

    fun loadSettings(load: Boolean): Boolean {
        if (!load) {
            val addr = Preferences.getServerAddress()
            editTextServerAddr.setText(addr)

            val autoStart = Preferences.isAutoStart()
            checkBoxStartOnBoot.isChecked = autoStart

            val showWnd = Preferences.isShowPreloadWnd()
            checkBoxShowPreload.isChecked = showWnd

            val player = Preferences.getPlayer()
            if (player.isEmpty() || player == "0")
                spinnerPlayer.setSelection(0)
            if (player == "1")
                spinnerPlayer.setSelection(1)
            else {
                val ind = Players.getList().indexOfFirst { it.Package == player }
                spinnerPlayer.setSelection(ind + 2)
            }
        }

        val sets = ServerApi.readSettings()
        if (sets == null) {
            Toast.makeText(this, R.string.error_retrieving_settings, Toast.LENGTH_SHORT).show()
            return false
        }

        editTextCacheSize.setText(sets.CacheSize.toString())
        editTextPreloadBufferSize.setText(sets.PreloadBufferSize.toString())

        checkBoxIsElementumCache.setChecked(sets.IsElementumCache)

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
        val showWnd = checkBoxShowPreload.isChecked
        Preferences.setShowPreloadWnd(showWnd)
        val player = spinnerPlayer.selectedItem as Player
        Preferences.setPlayer(player.Package)

        try {
            val sets = ServerSettings(
                    editTextCacheSize.text.toString().toInt(),
                    editTextPreloadBufferSize.text.toString().toInt(),
                    checkBoxIsElementumCache.isChecked,
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
        } catch (e: Exception) {
        }
    }
}
