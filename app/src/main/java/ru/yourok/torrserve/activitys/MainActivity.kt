package ru.yourok.torrserve.activitys

import android.Manifest
import android.content.Intent
import android.os.Bundle
import android.support.design.widget.Snackbar
import android.support.v4.app.ActivityCompat
import android.support.v7.app.AppCompatActivity
import android.view.KeyEvent
import android.view.KeyEvent.KEYCODE_DPAD_LEFT
import android.view.KeyEvent.KEYCODE_DPAD_RIGHT
import android.view.View
import android.widget.ListView
import android.widget.ProgressBar
import com.mikepenz.materialdrawer.Drawer
import kotlinx.android.synthetic.main.activity_main.*
import ru.yourok.torrserve.Donate
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.TorrentListAdapter
import ru.yourok.torrserve.menu.TorrentListSelectionMenu
import ru.yourok.torrserve.navigationBar.NavigationBar
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.services.TorrService
import java.util.*
import kotlin.concurrent.schedule
import kotlin.concurrent.thread


class MainActivity : AppCompatActivity() {
    val torrAdapter = TorrentListAdapter(this)
    lateinit var drawer: Drawer

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        drawer = NavigationBar.setup(this)

        requestPermissionWithRationale()

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
    }

    override fun onResume() {
        super.onResume()
        startServer()
        torrAdapter.updateList()
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
                showMenuHelp()
                Donate.showDonate(this)
            }
        }
    }

    private fun showMenuHelp() {
        if (listViewTorrent.count == 0)
            Timer().schedule(2000) {
                if (listViewTorrent.count == 0)
                    runOnUiThread { drawer.openDrawer() }
            }
    }

    override fun onKeyUp(keyCode: Int, event: KeyEvent?): Boolean {
        event?.let {
            when (keyCode) {
                KEYCODE_DPAD_RIGHT -> {
                    drawer.openDrawer()
                    return true
                }
                KEYCODE_DPAD_LEFT -> {
                    drawer.closeDrawer()
                    return true
                }
            //Add
                KeyEvent.KEYCODE_1, KeyEvent.KEYCODE_NUMPAD_1, KeyEvent.KEYCODE_BUTTON_1 -> {
                    startActivity(Intent(this, AddActivity::class.java))
                }
            //Remove All
                KeyEvent.KEYCODE_2, KeyEvent.KEYCODE_NUMPAD_2, KeyEvent.KEYCODE_BUTTON_2 -> {
                    thread {
                        val torrList = ServerApi.list()
                        torrList.forEach {
                            ServerApi.rem(it.Hash)
                        }
                        runOnUiThread {
                            torrAdapter.updateList()
                        }
                    }
                }
            //Donate
                KeyEvent.KEYCODE_3, KeyEvent.KEYCODE_NUMPAD_3, KeyEvent.KEYCODE_BUTTON_3 -> {
                    ServerApi.cleanCache()
                }
            //Clear cache
                KeyEvent.KEYCODE_4, KeyEvent.KEYCODE_NUMPAD_4, KeyEvent.KEYCODE_BUTTON_4 -> {
                    ServerApi.cleanCache()
                }
            //Exit
                KeyEvent.KEYCODE_5, KeyEvent.KEYCODE_NUMPAD_5, KeyEvent.KEYCODE_BUTTON_5 -> {
                    TorrService.stopAndExit()
                }
            //Settings
                KeyEvent.KEYCODE_6, KeyEvent.KEYCODE_NUMPAD_6, KeyEvent.KEYCODE_BUTTON_6 -> {
                    startActivity(Intent(this, SettingsActivity::class.java))
                }
                else -> return super.onKeyUp(keyCode, event)
            }
        }
        return super.onKeyUp(keyCode, event)
    }

    private fun requestPermissionWithRationale() {
        thread {
            if (ActivityCompat.shouldShowRequestPermissionRationale(this, Manifest.permission.WRITE_EXTERNAL_STORAGE)) {
                Snackbar.make(findViewById<View>(R.id.main_layout), R.string.permission_storage_msg, Snackbar.LENGTH_INDEFINITE)
                        .setAction(R.string.permission_btn, {
                            ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.WRITE_EXTERNAL_STORAGE), 1)
                        })
                        .show()
            } else {
                ActivityCompat.requestPermissions(this, arrayOf(Manifest.permission.WRITE_EXTERNAL_STORAGE), 1)
            }
        }
    }

    override fun onBackPressed() {
        if (drawer.isDrawerOpen) {
            drawer.closeDrawer()
            return
        }

        super.onBackPressed()
    }
}

