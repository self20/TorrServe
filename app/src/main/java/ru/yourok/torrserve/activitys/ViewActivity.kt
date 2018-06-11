package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.*
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListFileAdapter
import ru.yourok.torrserve.serverhelper.File
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.services.TorrService
import ru.yourok.torrserve.utils.Mime
import java.net.URLDecoder
import kotlin.concurrent.thread

class ViewActivity : AppCompatActivity() {
    var torrentLink = ""
    var saveInDB = true

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
            intent.data?.let {
                torrentLink = URLDecoder.decode(it.toString(), "UTF-8")
            }
        }

        if (intent.hasExtra("DontSave"))
            saveInDB = false

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

    private fun prepareTorrent() {
        setMessage(R.string.starting_server)
        val run = TorrService.waitServer()
        if (!run) {
            showToast(R.string.error_server_start)
            finish()
            return
        }
        setMessage(R.string.preparing_torrent)
        val tors = addTorrent()

        if (tors.isEmpty()) {
            showToast(R.string.error_add_torrent)
            finish()
            return
        }
        play(tors[0])
        return
    }

    private fun setMessage(msg: Int) {
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

    private fun showToast(msg: Int) {
        runOnUiThread {
            Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
        }
    }

    private fun play(tor: Torrent) {
        val fpList = findPlayableFiles(tor)
        if (fpList.size == 1) {
            finish()
            thread {
                ProgressActivity.show(tor, fpList[0]!!)
            }
        } else if (fpList.size > 1) {
            runOnUiThread {
                findViewById<TextView>(R.id.textViewStatus).visibility = View.GONE
                findViewById<ProgressBar>(R.id.progressBar).visibility = View.GONE
                findViewById<Button>(R.id.buttonPlaylist).visibility = View.VISIBLE
                findViewById<Button>(R.id.buttonPlaylist).setOnClickListener {
                    ServerApi.openPlayList(tor)
                    finish()
                }
                val adapter = TorrentListFileAdapter(this, tor.Hash)
                val listViewFiles = findViewById<ListView>(R.id.listViewTorrentFiles)
                listViewFiles.adapter = adapter
                listViewFiles.setOnItemClickListener { _, _, i, _ ->
                    val file = adapter.getItem(i) as File?
                    file?.let {
                        finish()
                        thread {
                            ProgressActivity.show(tor, file)
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

    private fun addTorrent(): List<Torrent> {
        try {
            return ServerApi.add(torrentLink, saveInDB)
        } catch (e: Exception) {
            val msg = e.message ?: getString(R.string.error_add_torrent)
            runOnUiThread {
                Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
            }
            return emptyList()
        }
    }

    private fun findPlayableFiles(tor: Torrent): Map<Int, File> {
        val retList = mutableMapOf<Int, File>()
        tor.Files.forEachIndexed { index, it ->
            if (Mime.getMimeType(it.Name) != "*/*")
                retList[index] = it
        }
        return retList
    }
}
