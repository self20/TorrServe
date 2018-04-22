package ru.yourok.torrserve.services

import android.app.Service
import android.content.Intent
import android.os.Handler
import android.os.IBinder
import android.os.Looper
import android.view.View
import android.widget.TextView
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.utils.Utils
import ru.yourok.torrserve.views.FloatingView
import kotlin.concurrent.thread


class TorrentInfo : Service() {
    override fun onBind(p0: Intent?): IBinder? = null

    override fun onCreate() {
        super.onCreate()
        watching()
    }

    private fun watching() {

        synchronized(isWatching) {
            if (isWatching)
                return
            isWatching = true
        }

        val view = FloatingView(this)
        view.create()

        thread {
            view?.getView()?.let { view ->
                while (isWatching && Preferences.isShowState()) {
                    val info = ServerApi.info(Hash)
                    info?.let {
                        Handler(Looper.getMainLooper()).post {
                            view.visibility = View.VISIBLE
                            (view.findViewById(R.id.textViewPeers) as TextView?)?.text = "Peers: " + it.ConnectedSeeders.toString() + " / " + it.TotalPeers.toString()
                            (view.findViewById(R.id.textViewSpeedDL) as TextView?)?.text = "Download speed: " + Utils.byteFmt(it.DownloadSpeed)
                            (view.findViewById(R.id.textViewSpeedUL) as TextView?)?.text = "Upload speed: " + Utils.byteFmt(it.UploadSpeed)
                        }
                    } ?: let {
                        Handler(Looper.getMainLooper()).post {
                            view.visibility = View.GONE
                        }
                    }
                    Thread.sleep(1000)
                }
            }
            isWatching = false

        }
    }

    companion object {
        private var Hash: String = ""
        private var isWatching: Boolean = false


        fun showWindow(hash: String) {
            Hash = hash
            if (!isWatching)
                try {
                    val intent = Intent(App.getContext(), TorrentInfo::class.java)
                    App.getContext().startService(intent)
                } catch (e: Exception) {
                    e.printStackTrace()
                }
        }

        fun closeWindow() {
            isWatching = false
        }
    }
}