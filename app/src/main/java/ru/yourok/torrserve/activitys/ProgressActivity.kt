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
import ru.yourok.torrserve.serverhelper.File
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
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
                NotificationServer.Show(this, it.Hash)
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
                    ServerApi.drop(it.Hash)
                }
            }
    }


    private fun waitPreload(): String {
        var errMsg = ""
        torrent?.let { torrent ->
            file?.let { file ->
                isPreload = true
                val th = thread {
                    errMsg = ServerApi.preload(torrent.Hash, file.Link)
                    isPreload = false
                }
                while (isPreload) {
                    Thread.sleep(1000)
                    val info = ServerApi.info(torrent.Hash)
                    if (info == null) {
                        Thread.sleep(1000)
                        continue
                    }

                    if (!info.IsGettingInfo && (!info.IsPreload || info.PreloadSize >= info.PreloadLength))
                        return errMsg

                    var msg = ""
                    var prc = 0
                    if (info.PreloadLength > 0) {
                        prc = (info.PreloadSize * 100 / info.PreloadLength).toInt()
                        msg += getString(R.string.buffer) + ": " + (prc).toString() + "% " + Utils.byteFmt(info.PreloadSize) + "/" + Utils.byteFmt(info.PreloadLength) + "\n"
                    }
                    msg += getString(R.string.peers) + ": [" + info.ConnectedSeeders.toString() + "] " + info.ActivePeers.toString() + "/" + info.TotalPeers.toString() + "\n"
                    msg += getString(R.string.download_speed) + ": " + Utils.byteFmt(info.DownloadSpeed) + "/Sec"
                    setMessage(msg, prc)
                }
                th.join(15000)
                if (ServerApi.get(torrent.Hash) == null) {
                    if (errMsg.isNotEmpty())
                        return errMsg
                    else
                        return getString(R.string.error_open_torrent)
                }
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
