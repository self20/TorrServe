package ru.yourok.torrserve.activitys

import android.content.ClipData
import android.content.ClipboardManager
import android.content.Context
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.support.v7.widget.PopupMenu
import android.view.Gravity
import android.view.MenuItem
import android.view.View
import android.widget.ListView
import android.widget.TextView
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_files.*
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListFileAdapter
import ru.yourok.torrserve.serverhelper.File
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread


/**
 * Created by yourok on 20.02.18.
 */
class FilesActivity : AppCompatActivity() {
    var torrId = ""
    var torrent: Torrent? = null

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_files)
        if (!intent.hasExtra("Hash")) {
            finish()
            return
        }
        torrId = intent.getStringExtra("Hash")
        progressBar.visibility = View.VISIBLE

        thread {
            try {
                torrent = ServerApi.get(torrId)
                if (torrent == null) {
                    Toast.makeText(this, R.string.stat_server_is_not_running, Toast.LENGTH_SHORT).show()
                    finish()
                    return@thread
                }
                torrent?.let { torr ->
                    runOnUiThread {
                        findViewById<TextView>(R.id.textViewTorrFileName).setText(torr.Name)
                        textViewTorrSize.visibility = View.VISIBLE
                        textViewTorrSize.setText(Utils.byteFmt(torr.Length))
                    }

                    val adapter = TorrentListFileAdapter(this, torrId)
                    val listViewFiles = findViewById<ListView>(R.id.listViewTorrentFiles)
                    runOnUiThread {
                        listViewFiles.adapter = adapter
                        listViewFiles.setOnItemClickListener { _, _, i, _ ->
                            val file = torr.Files[i]
                            adapter.torrent?.let {
                                it.Files[i].Viewed = true
                                adapter.notifyDataSetChanged()
                            }
                            thread {
                                ServerApi.view(this, torr.Hash, file.Name, file.Link)
                            }
                        }
                        listViewFiles.setOnItemLongClickListener { _, view, i, _ ->
                            showPopupMenu(view, torr.Files[i])
                            true
                        }
                        progressBar.visibility = View.GONE
                    }
                }
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
    }

    private fun showPopupMenu(v: View, file: File) {
        val popupMenu = PopupMenu(this, v, Gravity.CENTER)
        popupMenu.inflate(R.menu.file_list_menu)
        popupMenu.setOnMenuItemClickListener(object : PopupMenu.OnMenuItemClickListener {
            override fun onMenuItemClick(item: MenuItem): Boolean {
                var addr = Preferences.getServerAddress()
                addr += file.Link

                when (item.getItemId()) {
                    R.id.itemCopyUrl -> {
                        val clipboard = getSystemService(Context.CLIPBOARD_SERVICE) as ClipboardManager
                        val clip = ClipData.newPlainText(file.Name, addr)
                        clipboard.primaryClip = clip
                        return true
                    }
                    R.id.itemOpen -> {
                        ServerApi.view(this@FilesActivity, torrent!!.Hash, file.Name, file.Link)
                        return true
                    }
                    else -> return false
                }
            }
        })
        popupMenu.show()
    }
}