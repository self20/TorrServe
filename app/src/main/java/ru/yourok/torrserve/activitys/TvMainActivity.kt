package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.ListView
import android.widget.ProgressBar
import kotlinx.android.synthetic.main.activity_main_tv.*
import ru.yourok.torrserve.Donate
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListAdapter
import ru.yourok.torrserve.menu.TorrentListSelectionMenu
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.server.Torrent
import ru.yourok.torrserve.serverloader.ServerLoader
import ru.yourok.torrserve.services.TorrService
import kotlin.concurrent.thread

class TvMainActivity : AppCompatActivity() {
    val torrAdapter = TorrentListAdapter(this)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main_tv)

        val listViewTorrent = findViewById<ListView>(R.id.listViewTorrent)
        listViewTorrent.adapter = torrAdapter
        listViewTorrent.setOnItemClickListener { adapterView, view, i, l ->
            if (!ServerApi.echo()) {
                startServer()
                return@setOnItemClickListener
            }

            if ((torrAdapter.getItem(i) as Torrent).Files().isEmpty()) {
                return@setOnItemClickListener
            }

            val name = (torrAdapter.getItem(i) as Torrent).Name()
            val hash = (torrAdapter.getItem(i) as Torrent).Hash()
            val progressBar = view.findViewById<ProgressBar>(R.id.progressBar)
            progressBar.visibility = View.VISIBLE
            view.isEnabled = false
            thread {
                ServerApi.get(hash)
                runOnUiThread {
                    progressBar.visibility = View.GONE
                    view.isEnabled = true
                }
                val intent = Intent(this, FilesActivity::class.java)
                intent.putExtra("Name", name)
                intent.putExtra("Hash", hash)
                startActivity(intent)
            }
        }
        listViewTorrent.setOnItemLongClickListener { adapterView, view, i, l ->
            listViewTorrent.choiceMode = ListView.CHOICE_MODE_MULTIPLE_MODAL
            listViewTorrent.setItemChecked(i, true)
            true
        }
        listViewTorrent.setMultiChoiceModeListener(TorrentListSelectionMenu(this, torrAdapter))

        ///Button listeners
        buttonAdd.setOnClickListener {
            startActivity(Intent(this, AddActivity::class.java))
        }

        buttonRemoveAll.setOnClickListener {
            thread {
                val torrList = ServerApi.list()
                torrList.forEach {
                    ServerApi.rem(it.Hash())
                }
                runOnUiThread {
                    val ada = this.findViewById<ListView>(R.id.listViewTorrent).adapter
                    (ada as TorrentListAdapter).updateList()
                }
            }
        }

        buttonPlaylist.setOnClickListener {
            ServerApi.openPlayList()
        }

        buttonDonate.setOnClickListener {
            Donate.donateDialog(this)
        }

        buttonUpdate.setOnClickListener {
            startActivity(Intent(this, ServerLoaderActivity::class.java))
        }

        buttonSearch.setOnClickListener {
            startActivity(Intent(this, SearchActivity::class.java))
        }

        buttonExit.setOnClickListener {
            TorrService.exit()
        }

        buttonSettings.setOnClickListener {
            startActivity(Intent(this, SettingsActivity::class.java))
        }
    }

    override fun onResume() {
        super.onResume()
        startServer()
        torrAdapter.updateList()
        autoUpdateList()
    }

    override fun onPause() {
        super.onPause()
        isUpdate = false
    }

    private var isUpdate = false
    private fun autoUpdateList() {
        thread {
            synchronized(isUpdate) {
                if (isUpdate)
                    return@thread
                isUpdate = true
            }

            while (isUpdate) {
                torrAdapter?.checkList()
                Thread.sleep(1000)
            }
        }
    }

    private fun startServer() {
        progressBar.visibility = View.GONE
        textViewStatus.visibility = View.GONE
        if (!ServerApi.echo()) {
            textViewStatus.visibility = View.VISIBLE
            if (!ServerLoader.serverExists()) {
                textViewStatus.setText(R.string.warn_server_not_exists)
                return
            }
            progressBar.visibility = View.VISIBLE
            textViewStatus.setText(R.string.starting_server)
            thread {
                if (TorrService.waitServer()) {
                    runOnUiThread { textViewStatus.setText(R.string.loading_torrent_list) }
                    ServerApi.list()
                    runOnUiThread {
                        progressBar.visibility = View.GONE
                        textViewStatus.visibility = View.GONE
                        torrAdapter.updateList()
                    }
                } else {
                    runOnUiThread {
                        progressBar.visibility = View.GONE
                        textViewStatus.setText(R.string.error_server_start)
                    }
                }
                Donate.showDonate(this)
            }
        }
    }
}