package ru.yourok.torrserve.adapters

import android.app.Activity
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Handler
import android.os.Looper
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.Button
import org.json.JSONArray
import ru.yourok.torrserve.R
import ru.yourok.torrserve.activitys.AddDSActivity
import ru.yourok.torrserve.serverhelper.ServerApi
import kotlin.concurrent.thread

class SearchTorrentAdapter(val activity: Activity) : BaseAdapter() {
    var list: JSONArray = JSONArray()

    fun setResult(list: JSONArray) {
        this.list = list
        notifyDataSetChanged()
    }

    override fun getView(index: Int, convertView: View?, viewGroup: ViewGroup): View {
        val vi: View = convertView ?: (activity.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater).inflate(R.layout.adapter_search_torrent_item, null)
        val info = list.getJSONObject(index)
        var name = ""
        if (info.has("Name")) {
            name = info.getString("Name")
            if (info.has("Size") && info.getString("Size").isNotEmpty())
                name += "\n" + info.getString("Size")
            if (info.has("PeersUl") && info.getInt("PeersUl") >= 0) {
                name += " | ▲ " + info.getInt("PeersUl")
                name += " | ▼ " + info.getInt("PeersDl")
            }
        }
        if (info.has("Magnet")) {
            val plButton = vi.findViewById<Button>(R.id.buttonPlayTorrent)
            if (info.getString("Magnet").isNotEmpty()) {
                plButton.text = name
                val mag = info.getString("Magnet")
                plButton.setOnClickListener {
                    val intent = Intent(activity, AddDSActivity::class.java)
                    intent.data = Uri.parse(mag)
                    intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
                    intent.action = Intent.ACTION_VIEW
                    activity.startActivity(intent)
                }
                vi.findViewById<Button>(R.id.buttonAddTorrent).setOnClickListener {
                    thread {
                        try {
                            ServerApi.add(mag, true)
                            Handler(Looper.getMainLooper()).post {
                                it.isEnabled = false
                            }
                        } catch (e: Exception) {

                        }
                    }
                }
            } else
                plButton.text = "Magnet not found"
        }
        return vi
    }

    override fun getItem(p0: Int): Any? {
        try {
            return list.getJSONObject(p0)
        } catch (e: Exception) {
            return 0
        }
    }

    override fun getItemId(p0: Int): Long {
        return p0.toLong()
    }

    override fun getCount(): Int {
        if (list.length() > 0) {
            if (list.getJSONObject(0).has("Magnet"))
                return list.length()
        }
        return 0
    }
}