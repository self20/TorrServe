package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.ListView
import android.widget.ProgressBar
import android.widget.TextView
import android.widget.Toast
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListFileAdapter
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.services.TorrService
import ru.yourok.torrserve.utils.Mime
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread

class ViewActivity : AppCompatActivity() {
    var torrentLink = ""

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_view)

        if (intent == null) {
            finish()
            return
        }

        setFinishOnTouchOutside(false)
        ///Intent open
        if (intent.action != null && intent.action.equals(Intent.ACTION_VIEW)) {
            torrentLink = intent.data?.toString() ?: ""
        }

        ///Intent send
        if (intent.action != null && intent.action.equals(Intent.ACTION_SEND)) {
            if (intent.getStringExtra(Intent.EXTRA_TEXT) != null)
                torrentLink = intent.getStringExtra(Intent.EXTRA_TEXT)
            if (intent.extras.get(Intent.EXTRA_STREAM) != null)
                torrentLink = intent.extras.get(Intent.EXTRA_STREAM).toString()
        }

        if (torrentLink.isEmpty() && !intent.hasExtra("Preload")) {
            finish()
            return
        }

        thread {
            ///Intent preload
            if (intent.hasExtra("Preload")) {
                if (Preferences.isShowPreloadWnd()) {
                    val tor = ServerApi.get(intent.getStringExtra("Preload"))
                    tor?.let {
                        waitPreload(tor)
                    }
                }
                finish()
            } else {
                prepareTorrent()
            }
        }
    }

    fun prepareTorrent() {
        setMessage(R.string.starting_server)
        val run = TorrService.waitServer()
        if (!run) {
            showToast(R.string.error_server_start)
            finish()
            return
        }
        setMessage(R.string.preparing_torrent)
        //TODO не возвращает файлы из сериала (полицейский с рублевки)
        val tor = addTorrent()

        if (tor == null) {
            showToast(R.string.error_add_torrent)
            finish()
            return
        }
        play(tor)
        return
    }

    fun setMessage(msg: Int) {
        runOnUiThread {
            if (msg != -1) {
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.VISIBLE
                findViewById<TextView>(R.id.textViewStatus).visibility = View.VISIBLE
                findViewById<TextView>(R.id.textViewStatus).setText(msg)
            } else {
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.GONE
                findViewById<TextView>(R.id.textViewStatus).visibility = View.GONE
            }
        }
    }

    fun setMessage(msg: String, progress: Int) {
        runOnUiThread {
            if (msg.isNotEmpty()) {
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.VISIBLE
                findViewById<ProgressBar>(R.id.progressBar).isIndeterminate = progress == 0
                if (progress > 0) {
                    findViewById<ProgressBar>(R.id.progressBar).progress = progress
                }

                findViewById<TextView>(R.id.textViewStatus).visibility = View.VISIBLE
                findViewById<TextView>(R.id.textViewStatus).setText(msg)
            } else {
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.GONE
                findViewById<TextView>(R.id.textViewStatus).visibility = View.GONE
            }
        }
    }


    fun showToast(msg: Int) {
        runOnUiThread {
            Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
        }
    }

    fun play(tor: Torrent) {
        val fpList = findPlayableFiles(tor)
        if (fpList.size == 1) {
            finish()
            thread {
                Thread.sleep(500)
                ServerApi.view(this, tor.Hash, tor.Name, fpList.values.first())
            }
        } else if (fpList.size > 1) {
            runOnUiThread {
                findViewById<TextView>(R.id.textViewStatus).visibility = View.GONE
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.GONE
                val adapter = TorrentListFileAdapter(this, tor.Hash)
                val listViewFiles = findViewById<ListView>(R.id.listViewTorrentFiles)
                listViewFiles.adapter = adapter
                listViewFiles.setOnItemClickListener { _, _, i, _ ->
                    val link = fpList[i]
                    link?.let {
                        finish()
                        thread {
                            Thread.sleep(500)
                            ServerApi.view(this, tor.Hash, tor.Name, it)
                        }
                    }
                }
            }
        } else {
            val intent = Intent(this, FilesActivity::class.java)
            intent.putExtra("Hash", tor.Hash)
            startActivity(intent)
            finish()
        }
    }

    fun waitPreload(tor: Torrent) {
        var err = 0
        setMessage(R.string.buffering_torrent)
        while (true) {
            if (err > 15) {
                return
            }
            Thread.sleep(1000)

            val info = ServerApi.info(tor.Hash)
            if (info == null) {
                err++
                continue
            }

            if (!info.IsPreload)
                return
            if (info.PreloadLength > 0) {
                var msg = ""
                val prc = (info.PreloadOffset * 100 / info.PreloadLength).toInt()
                msg += getString(R.string.buffer) + ": " + (prc).toString() + "% " + Utils.byteFmt(info.PreloadOffset) + "/" + Utils.byteFmt(info.PreloadLength) + "\n"
                msg += getString(R.string.peers) + ": " + info.ConnectedSeeders.toString() + "/" + info.TotalPeers.toString() + "\n"
                msg += getString(R.string.download_speed) + ": " + Utils.byteFmt(info.DownloadSpeed) + "/Sec"
                setMessage(msg, prc)
            }
        }
    }

    fun addTorrent(): Torrent? {
        try {
            return ServerApi.add(torrentLink)
        } catch (e: Exception) {
            val msg = e.message ?: getString(R.string.error_add_torrent)
            runOnUiThread {
                Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
            }
            return null
        }
    }

    fun findPlayableFiles(tor: Torrent): Map<Int, String> {
        val retList = mutableMapOf<Int, String>()
        tor.Files.forEachIndexed { index, it ->
            if (Mime.getMimeType(it.Name) != "*/*")
                retList[index] = it.Link
        }
        return retList
    }
}
