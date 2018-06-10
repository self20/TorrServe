package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.KeyEvent
import android.view.View
import android.widget.ListView
import android.widget.ProgressBar
import kotlinx.android.synthetic.main.activity_main_tv.*
import ru.yourok.torrserve.Donate
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListAdapter
import ru.yourok.torrserve.menu.TorrentListSelectionMenu
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
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

            if ((torrAdapter.getItem(i) as Torrent).Files.size == 0) {
                return@setOnItemClickListener
            }

            val name = (torrAdapter.getItem(i) as Torrent).Name
            val hash = (torrAdapter.getItem(i) as Torrent).Hash
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
                    ServerApi.rem(it.Hash)
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

        buttonClearCache.setOnClickListener {
            ServerApi.cleanCache("")
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

    override fun onKeyUp(keyCode: Int, event: KeyEvent?): Boolean {
        event?.let {
            when (keyCode) {
            //Add
                KeyEvent.KEYCODE_1, KeyEvent.KEYCODE_NUMPAD_1, KeyEvent.KEYCODE_BUTTON_1 -> {
                    buttonAdd.performClick()
                }
            //Remove All
                KeyEvent.KEYCODE_2, KeyEvent.KEYCODE_NUMPAD_2, KeyEvent.KEYCODE_BUTTON_2 -> {
                    buttonRemoveAll.performClick()
                }
            //Playlist
                KeyEvent.KEYCODE_3, KeyEvent.KEYCODE_NUMPAD_3, KeyEvent.KEYCODE_BUTTON_3 -> {
                    buttonPlaylist.performClick()
                }
            //Donate
                KeyEvent.KEYCODE_4, KeyEvent.KEYCODE_NUMPAD_4, KeyEvent.KEYCODE_BUTTON_4 -> {
                    buttonDonate.performClick()
                }
            //Clear Cache
                KeyEvent.KEYCODE_5, KeyEvent.KEYCODE_NUMPAD_5, KeyEvent.KEYCODE_BUTTON_5 -> {
                    buttonClearCache.performClick()
                }
            //Exit
                KeyEvent.KEYCODE_6, KeyEvent.KEYCODE_NUMPAD_6, KeyEvent.KEYCODE_BUTTON_6 -> {
                    buttonExit.performClick()
                }
            //Settings
                KeyEvent.KEYCODE_7, KeyEvent.KEYCODE_NUMPAD_7, KeyEvent.KEYCODE_BUTTON_7 -> {
                    buttonSettings.performClick()
                }
                else -> return super.onKeyUp(keyCode, event)
            }
        }
        return super.onKeyUp(keyCode, event)
    }

    private fun startServer() {
        if (!ServerApi.echo()) {
            progressBar.visibility = View.VISIBLE
            textViewStatus.visibility = View.VISIBLE
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