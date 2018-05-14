package ru.yourok.torrserve.adapters

import android.app.Activity
import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.ImageView
import android.widget.TextView
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread

/**
 * Created by yourok on 20.02.18.
 */

class TorrentListFileAdapter(val activity: Activity, val torrId: String) : BaseAdapter() {
    var torrent: Torrent? = null

    init {
        updateList()
    }

    fun updateList() {
        if (torrId.isNotEmpty())
            thread {
                while (!ServerApi.echo()) {
                    Thread.sleep(1000)
                }
                torrent = ServerApi.get(torrId)
                activity.runOnUiThread {
                    notifyDataSetChanged()
                }
            }
    }

    override fun getView(index: Int, convertView: View?, viewGroup: ViewGroup): View {
        val vi: View = convertView ?: (activity.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater).inflate(R.layout.adapter_torrent_list_files, null)
        torrent?.let {
            val file = it.Files[index]
            vi.findViewById<TextView>(R.id.textViewFileName).text = file.Name
            vi.findViewById<TextView>(R.id.textViewFileSize).text = Utils.byteFmt(file.Size)
            if (file.Viewed)
                vi.findViewById<ImageView>(R.id.imageViewPlayed).visibility = View.VISIBLE
            else
                vi.findViewById<ImageView>(R.id.imageViewPlayed).visibility = View.GONE
        }
        return vi
    }

    override fun getItem(p0: Int): Any? {
        torrent?.let {
            return it.Files[p0]
        }
        return null
    }

    override fun getItemId(p0: Int): Long {
        return p0.toLong()
    }

    override fun getCount(): Int {
        return torrent?.Files?.size ?: 0
    }
}