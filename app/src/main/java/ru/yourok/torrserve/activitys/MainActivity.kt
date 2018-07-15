package ru.yourok.torrserve.activitys

import android.Manifest
import android.app.UiModeManager
import android.content.Intent
import android.content.res.Configuration
import android.os.Bundle
import android.support.design.widget.Snackbar
import android.support.v4.app.ActivityCompat
import android.support.v7.app.AppCompatActivity
import android.util.DisplayMetrics
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
import ru.yourok.torrserve.serverhelper.server.Torrent
import ru.yourok.torrserve.serverloader.ServerLoader
import ru.yourok.torrserve.services.TorrService
import java.util.*
import kotlin.concurrent.schedule
import kotlin.concurrent.thread


class MainActivity : AppCompatActivity() {
    val torrAdapter = TorrentListAdapter(this)
    lateinit var drawer: Drawer

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        var isAndroidTV = false
        val uiModeManager = getSystemService(UI_MODE_SERVICE) as UiModeManager
        if (uiModeManager.getCurrentModeType() == Configuration.UI_MODE_TYPE_TELEVISION)
            isAndroidTV = true
        val dm = DisplayMetrics()
        windowManager.defaultDisplay.getMetrics(dm)
        val wi = dm.widthPixels.toDouble() / dm.xdpi.toDouble()
        val hi = dm.heightPixels.toDouble() / dm.ydpi.toDouble()
        val screenInches = Math.sqrt(wi * wi + hi * hi)
        if (screenInches >= 8 || isAndroidTV) {
            startActivity(Intent(this, TvMainActivity::class.java))
            finish()
            return
        }

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

            if ((torrAdapter.getItem(i) as Torrent).Files().isEmpty())
                return@setOnItemClickListener

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
                torrAdapter.checkList()
                Thread.sleep(1000)
            }
        }
    }

    private fun startServer() {
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

