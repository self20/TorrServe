package ru.yourok.torrserve.menu

import android.app.Activity
import android.content.Intent
import android.view.ActionMode
import android.view.Menu
import android.view.MenuItem
import android.widget.AbsListView
import android.widget.ListView
import android.widget.Toast
import ru.yourok.torrserve.R
import ru.yourok.torrserve.activitys.InfoActivity
import ru.yourok.torrserve.adapters.TorrentListAdapter
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.server.Torrent
import kotlin.concurrent.thread

/**
 * Created by yourok on 19.11.17.
 */
class TorrentListSelectionMenu(val activity: Activity, val adapter: TorrentListAdapter) : AbsListView.MultiChoiceModeListener {

    private val selected: MutableSet<Int> = mutableSetOf()

    override fun onCreateActionMode(mode: ActionMode?, menu: Menu?): Boolean {
        mode?.menuInflater?.inflate(R.menu.torrent_list_menu, menu)
        selected.clear()
        return true
    }

    override fun onPrepareActionMode(mode: ActionMode?, menu: Menu?): Boolean {
        return false
    }

    override fun onActionItemClicked(mode: ActionMode, item: MenuItem): Boolean {
        when (item.itemId) {
            R.id.itemShareMagnet -> {
                var msg = ""
                selected.forEach {
                    val torrent = (adapter.getItem(it) as Torrent)
                    msg += "${torrent.Name()}:\n${torrent.Magnet()}\n\n"
                }
                if (msg.isNotEmpty()) {
                    val share = Intent(Intent.ACTION_SEND)
                    share.setType("text/plain")
                    share.putExtra(Intent.EXTRA_TEXT, msg)
                    val intent = Intent.createChooser(share, "")
                    activity.startActivity(intent)
                }
            }
            R.id.itemInfoTorrent -> {
                val hashs = mutableListOf<String>()
                selected.forEach {
                    val Id = (adapter.getItem(it) as Torrent).Hash()
                    hashs.add(Id)
                }
                val intent = Intent(activity, InfoActivity::class.java)
                intent.apply {
                    putExtra("Hashs", hashs.toTypedArray())
                }
                activity.startActivity(intent)
            }
            R.id.itemRemove -> {
                selected.forEach {
                    val Id = (adapter.getItem(it) as Torrent).Hash()
                    thread {
                        try {
                            ServerApi.rem(Id)
                            adapter.updateList()
                        } catch (e: Exception) {
                            activity.runOnUiThread {
                                val msg = e.message.run { if (this != null) ": " + this else "" }
                                Toast.makeText(activity, activity.getText(R.string.error_remove_torrent).toString() + msg, Toast.LENGTH_SHORT).show()
                            }
                        }
                    }
                }
            }
        }
        mode.finish()
        return false
    }

    override fun onDestroyActionMode(mode: ActionMode?) {
    }

    override fun onItemCheckedStateChanged(mode: ActionMode?, position: Int, id: Long, checked: Boolean) {
        if (checked)
            selected.add(position)
        else
            selected.remove(position)
    }

    private fun clearChoice() {
        val lv = activity.findViewById<ListView>(R.id.listViewTorrent)
        lv.clearChoices()
        for (i in 0 until lv.getCount())
            lv.setItemChecked(i, false)
        lv.post(Runnable { lv.choiceMode = ListView.CHOICE_MODE_NONE })
        adapter.notifyDataSetChanged()
        lv.requestLayout()
    }
}