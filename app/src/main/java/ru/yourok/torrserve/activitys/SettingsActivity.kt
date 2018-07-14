package ru.yourok.torrserve.activitys

import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.support.v7.app.AlertDialog
import android.support.v7.app.AppCompatActivity
import android.view.inputmethod.EditorInfo
import android.widget.ArrayAdapter
import android.widget.EditText
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_settings.*
import ru.yourok.torrserve.App
import ru.yourok.torrserve.BuildConfig
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.utils.Player
import ru.yourok.torrserve.utils.Players
import kotlin.concurrent.thread


class SettingsActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_settings)

        val autocomplete = Preferences.getSaveHosts()
        val editTextServAddr = findViewById<EditText>(R.id.editTextServerAddr)
        editTextServAddr.setOnLongClickListener {
            val builder = AlertDialog.Builder(this)
            builder.setItems(autocomplete.toTypedArray()) { _, i ->
                val host = autocomplete[i]
                editTextServAddr.setText(host)
                Preferences.setServerAddress(host)
                loadSettings(true)
            }
            builder.create().show()
            false
        }

        editTextServAddr.setOnEditorActionListener { v, actionId, event ->
            if (actionId == EditorInfo.IME_ACTION_SEND) {
                val addr = editTextServerAddr.text.toString()
                Preferences.setServerAddress(addr)
                loadSettings(true)
            }
            false
        }

        buttonOk.setOnClickListener {
            saveSettings()
            finish()
        }

        buttonCancel.setOnClickListener {
            finish()
        }

        buttonRetrieveSettings.setOnClickListener {
            val addr = editTextServerAddr.text.toString()
            Preferences.setServerAddress(addr)
            loadSettings(true)
        }

        val plist = Players.getList()
        plist.add(0, Player(getString(R.string.default_player), "0"))
        plist.add(1, Player(getString(R.string.choose_player), "1"))

        val adp1 = ArrayAdapter<Player>(this, android.R.layout.simple_list_item_1, plist)
        adp1.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        spinnerPlayer.setAdapter(adp1)

        val adp2 = ArrayAdapter<String>(this, android.R.layout.simple_list_item_1, getResources().getStringArray(R.array.retracker_mode))
        adp2.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        spinnerRetracker.setAdapter(adp2)

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

    fun loadSettings(fromServer: Boolean) {
        if (!fromServer) {
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
        thread {
            try {
                val sets = ServerApi.readSettings()
                Handler(Looper.getMainLooper()).post {
                    editTextCacheSize.setText((sets.cacheSize / (1024 * 1024)).toString())
                    editTextPreloadBufferSize.setText((sets.preloadBufferSize / (1024 * 1024)).toString())
                    spinnerRetracker.setSelection(sets.retrackersMode)

                    checkBoxDisableTCP.setChecked(sets.disableTCP)
                    checkBoxDisableUTP.setChecked(sets.disableUTP)
                    checkBoxDisableUPNP.setChecked(sets.disableUPNP)
                    checkBoxDisableDHT.setChecked(sets.disableDHT)
                    checkBoxDisableUpload.setChecked(sets.disableUpload)

                    editTextEncryption.setText(sets.encryption.toString())
                    editTextConnectionsLimit.setText(sets.connectionsLimit.toString())
                    editTextDownloadRateLimit.setText(sets.downloadRateLimit.toString())
                    editTextUploadRateLimit.setText(sets.uploadRateLimit.toString())
                }
            } catch (e: Exception) {
                Handler(Looper.getMainLooper()).post {
                    Toast.makeText(this, R.string.error_retrieving_settings, Toast.LENGTH_SHORT).show()
                }
            }
        }
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
        thread {
            try {
                val sets = ServerApi.readSettings()

                sets.cacheSize = editTextCacheSize.text.toString().toLong() * (1024 * 1024)
                sets.preloadBufferSize = editTextPreloadBufferSize.text.toString().toLong() * (1024 * 1024)
                sets.retrackersMode = spinnerRetracker.selectedItemPosition
                sets.disableTCP = checkBoxDisableTCP.isChecked
                sets.disableUTP = checkBoxDisableUTP.isChecked
                sets.disableUPNP = checkBoxDisableUPNP.isChecked
                sets.disableDHT = checkBoxDisableDHT.isChecked
                sets.disableUpload = checkBoxDisableUpload.isChecked
                sets.encryption = editTextEncryption.text.toString().toInt()
                sets.downloadRateLimit = editTextDownloadRateLimit.text.toString().toInt()
                sets.uploadRateLimit = editTextUploadRateLimit.text.toString().toInt()
                sets.connectionsLimit = editTextConnectionsLimit.text.toString().toInt()

                thread {
                    try {
                        ServerApi.writeSettings(sets)
                        Preferences.addSaveHost(addr)
                        ServerApi.restartTorrentClient()
                    } catch (e: Exception) {
                        Handler(Looper.getMainLooper()).post {
                            Toast.makeText(App.getContext(), R.string.error_sending_settings, Toast.LENGTH_SHORT).show()
                        }
                    }
                }
            } catch (e: Exception) {
                Handler(Looper.getMainLooper()).post {
                    Toast.makeText(App.getContext(), R.string.error_sending_settings, Toast.LENGTH_SHORT).show()
                }
            }
        }.join()
    }
}
