package ru.yourok.torrserve.activitys

import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.*
import ru.yourok.torrserve.BuildConfig
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerRequests
import ru.yourok.torrserve.serverloader.ServerLoader
import kotlin.concurrent.thread

class ServerLoaderActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_server_loader)

        updateVersion()

        findViewById<Button>(R.id.update_remotly_button).setOnClickListener {
            thread {
                updateRemotly()
            }
        }

        findViewById<Button>(R.id.update_local_button).setOnClickListener {
            thread {
                updateLocaly()
            }
        }

        findViewById<Button>(R.id.delete_server_button).setOnClickListener {
            ServerLoader.deleteServer()
            updateVersion()
        }

        findViewById<ImageButton>(R.id.buttonCheckUpdate).setOnClickListener {
            updateVersion()
        }
    }

    fun updateVersion() {
        var servVer = ""
        if (ServerLoader.serverExists()) {
            thread {
                try {
                    servVer = ServerRequests.echo()
                } catch (e: Exception) {
                    e.printStackTrace()
                }
            }.join()
        }

        var version = "${getString(R.string.version)} Android: ${getString(R.string.app_name)} ${BuildConfig.VERSION_NAME}"
        if (servVer.isNotEmpty())
            version += "\n${getString(R.string.version)} Server: ${servVer}"
        else if (!ServerLoader.serverExists())
            version += "\n${getString(R.string.version)} Server: " + getString(R.string.warn_server_not_exists)
        else
            version += "\n${getString(R.string.version)} Server: not connected"

        findViewById<TextView>(R.id.current_info).setText(version)
    }

    fun updateRemotly() {
        Handler(Looper.getMainLooper()).post {
            findViewById<ProgressBar>(R.id.progress_bar).visibility = View.VISIBLE
        }
        ServerLoader.stop()
        val err = ServerLoader.download()
        if (err.isNotEmpty()) {
            Handler(Looper.getMainLooper()).post {
                val msg = getString(R.string.warn_error_download_server) + ": " + err
                findViewById<TextView>(R.id.update_info).setText(msg)
                Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
                findViewById<ProgressBar>(R.id.progress_bar).visibility = View.GONE
            }
            ServerLoader.run()
            return
        }
        ServerLoader.run()
        Handler(Looper.getMainLooper()).post {
            findViewById<TextView>(R.id.update_info).setText(R.string.stat_server_is_running)
            updateVersion()
            findViewById<ProgressBar>(R.id.progress_bar).visibility = View.GONE
        }
    }

    fun updateLocaly() {
        Handler(Looper.getMainLooper()).post {
            findViewById<ProgressBar>(R.id.progress_bar).visibility = View.VISIBLE
        }
        if (ServerLoader.checkLocal() == null) {
            Handler(Looper.getMainLooper()).post {
                val msg = getString(R.string.warn_no_localy_updates) + ": Download/TorrServer-android-${ServerLoader.getArch()}"
                findViewById<TextView>(R.id.update_info).setText(msg)
                Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
                updateVersion()
                findViewById<ProgressBar>(R.id.progress_bar).visibility = View.GONE
            }
            return
        }
        ServerLoader.stop()
        if (!ServerLoader.copyLocal()) {
            Handler(Looper.getMainLooper()).post {
                val msg = "Error copy server: Download/TorrServer-android-${ServerLoader.getArch()}"
                findViewById<TextView>(R.id.update_info).setText(msg)
                Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
                updateVersion()
                findViewById<ProgressBar>(R.id.progress_bar).visibility = View.GONE
            }
            return
        }
        ServerLoader.run()
        Handler(Looper.getMainLooper()).post {
            findViewById<TextView>(R.id.update_info).setText(R.string.stat_server_is_running)
            updateVersion()
            findViewById<ProgressBar>(R.id.progress_bar).visibility = View.GONE
        }
    }
}
