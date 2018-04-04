package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.ListView
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_view.*
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListFileAdapter
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.services.TorrService
import ru.yourok.torrserve.utils.Mime
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

        if (torrentLink.isEmpty()) {
            finish()
            return
        }
        thread {
            prepareTorrent()
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
                progressBar.visibility = View.VISIBLE
                textViewStatus.visibility = View.VISIBLE
                textViewStatus.setText(msg)
            } else {
                progressBar.visibility = View.GONE
                textViewStatus.visibility = View.GONE
            }
        }
    }

    fun showToast(msg: Int) {
        runOnUiThread {
            Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
        }
    }

    fun play(torr: Torrent) {
        var tor: Torrent? = torr
        if (torr.Files.size == 0)
            tor = ServerApi.get(torr.Hash)

        if (tor == null) {
            showToast(R.string.error_open_torrent)
            finish()
            return
        }
        tor.let {
            val fpList = findPlayableFiles(it)
            if (fpList.size == 1) {
                ServerApi.view(this, torr.Name, fpList.values.first())
                finish()
            } else if (fpList.size > 1) {
                runOnUiThread {
                    textViewStatus.visibility = View.GONE
                    progressBar.visibility = View.GONE
                    val adapter = TorrentListFileAdapter(this, it.Hash)
                    val listViewFiles = findViewById<ListView>(R.id.listViewTorrentFiles)
                    listViewFiles.adapter = adapter
                    listViewFiles.setOnItemClickListener { _, _, i, _ ->
                        val link = fpList[i]
                        link?.let {
                            ServerApi.view(this, torr.Name, it)
                        }
                        finish()
                    }
                }
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

    fun findTorrent(): Torrent? {
        val torrList = ServerApi.list()
        torrList.forEach {
            if (it.Magnet == torrentLink)
                return ServerApi.get(it.Hash)
        }
        return null
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
