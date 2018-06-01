package ru.yourok.torrserve.adapters

import android.app.Activity
import android.content.Context
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.TextView
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverhelper.Torrent
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread

/**
 * Created by yourok on 20.02.18.
 */

class TorrentListAdapter(val activity: Activity) : BaseAdapter() {
    var torrList: List<Torrent> = listOf()

    fun updateList() {
        thread {
            while (!ServerApi.echo()) {
                Thread.sleep(1000)
            }
            torrList = ServerApi.list()
            activity.runOnUiThread {
                notifyDataSetChanged()
            }
        }
    }

    fun checkList() {
        val tmpList = ServerApi.list()
        if (tmpList != torrList) {
            torrList = tmpList
            activity.runOnUiThread {
                notifyDataSetChanged()
            }
        }
    }

    override fun getView(index: Int, convertView: View?, viewGroup: ViewGroup): View {
        val vi: View = convertView ?: (activity.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater).inflate(R.layout.adapter_torrent_item, null)
        if (index in 0 until torrList.size) {
            vi.findViewById<TextView>(R.id.textViewTorrName).text = torrList[index].Name
            vi.findViewById<TextView>(R.id.textViewTorrSize).text = Utils.byteFmt(torrList[index].Length)
            vi.findViewById<TextView>(R.id.textViewTorrMagnet).text = torrList[index].Magnet
        }
        return vi
    }

    override fun getItem(p0: Int): Any {
        return torrList[p0]
    }

    override fun getItemId(p0: Int): Long {
        return p0.toLong()
    }

    override fun getCount(): Int {
        return torrList.size
    }

}