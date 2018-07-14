package ru.yourok.torrserve.activitys

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.ProgressBar
import android.widget.TextView
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.server.File
import ru.yourok.torrserve.serverhelper.server.Torrent
import ru.yourok.torrserve.serverhelper.server.TorrentWorking
import ru.yourok.torrserve.services.NotificationServer
import ru.yourok.torrserve.utils.Mime
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread

class ProgressActivity : AppCompatActivity() {

    private var isClosed = false
    private var isPreload = false

    companion object {
        private var torrent: Torrent? = null
        private var file: File? = null

        fun show(torrent: Torrent, file: File) {
            this.torrent = torrent
            this.file = file
            val intent = Intent(App.getContext(), ProgressActivity::class.java)
            intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
            App.getContext().startActivity(intent)
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_progress)

        if (torrent == null || file == null) {
            finish()
            return
        }
        setMessage(getText(R.string.buffering_torrent).toString(), 0)
        thread {
            isClosed = false
            torrent?.let {
                NotificationServer.Show(this, it.Hash())
            }

            if (!isClosed && Preferences.isShowPreloadWnd()) {
                val errMsg = waitPreload()
                if (!isClosed && errMsg.isNotEmpty()) {
                    try {
                        App.Toast(errMsg)
                    } catch (e: Exception) {
                    }
                    isClosed = true
                }
            }
            if (!isClosed)
                play()
            finish()
        }
    }

    override fun onBackPressed() {
        super.onBackPressed()
        isClosed = true
        if (isPreload)
            torrent?.let {
                thread {
                    ServerApi.drop(it.Hash())
                }
            }
    }


    private fun waitPreload(): String {
        torrent?.let { torrent ->
            file?.let { file ->
                isPreload = true
                val th = thread {
                    ServerApi.preload(file.Link)
                }
                Thread.sleep(1000)
                while (isPreload) {
                    try {
                        val stat = ServerApi.stat(torrent.Hash())
                        if (stat.TorrentStatus() == TorrentWorking || stat.PreloadedBytes() > stat.PreloadSize())
                            break

                        var msg = ""
                        var prc = 0
                        if (stat.PreloadSize() > 0) {
                            prc = (stat.PreloadedBytes() * 100 / stat.PreloadSize()).toInt()
                            msg += getString(R.string.buffer) + ": " + (prc).toString() + "% " + Utils.byteFmt(stat.PreloadedBytes()) + "/" + Utils.byteFmt(stat.PreloadSize()) + "\n"
                        }
                        msg += getString(R.string.peers) + ": [" + stat.ConnectedSeeders().toString() + "] " + stat.ActivePeers().toString() + "/" + stat.TotalPeers().toString() + "\n"
                        msg += getString(R.string.download_speed) + ": " + Utils.byteFmt(stat.DownloadSpeed()) + "/Sec"
                        setMessage(msg, prc)
                        Thread.sleep(100)
                    } catch (e: Exception) {
                        Thread.sleep(1000)
                    }
                }
                th.join(15000)
                return ""
            }
        }
        return getString(R.string.error_open_torrent)
    }

    private fun setMessage(msg: String, progress: Int) {
        runOnUiThread {
            if (msg.isNotEmpty()) {
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.VISIBLE
                findViewById<ProgressBar>(R.id.progressBar).isIndeterminate = progress == 0
                if (progress > 0) {
                    findViewById<ProgressBar>(R.id.progressBar).progress = progress
                }

                findViewById<TextView>(R.id.textViewStatus).visibility = View.VISIBLE
                findViewById<TextView>(R.id.textViewStatus).setText(msg)
            }
        }
    }

    private fun play() {
        torrent?.let { torrent ->
            file?.let { file ->
                thread {
                    val addr = Preferences.getServerAddress() + file.Link
                    val pkg = Preferences.getPlayer()

                    val intent = Intent(Intent.ACTION_VIEW, Uri.parse(addr))
                    val mime = Mime.getMimeType(file.Name)
                    intent.setDataAndType(Uri.parse(addr), mime)
                    intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                    intent.putExtra("title", file.Name)

                    if (pkg.isEmpty() or pkg.equals("0")) {
                        if (intent.resolveActivity(packageManager) != null) {
                            startActivity(intent)
                            return@thread
                        }
                    }
                    if (pkg.isNotEmpty() and !pkg.equals("0") and !pkg.equals("1")) {
                        intent.`package` = pkg
                        if (intent.resolveActivity(packageManager) != null) {
                            startActivity(intent)
                            return@thread
                        }
                    }

                    val intentC = Intent.createChooser(intent, "")
                    startActivity(intentC)
                }
            }
        }
    }
}
