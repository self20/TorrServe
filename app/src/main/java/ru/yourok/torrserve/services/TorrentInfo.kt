package ru.yourok.torrserve.services

import android.os.Handler
import android.os.Looper
import android.view.View
import android.widget.TextView
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.utils.Utils
import ru.yourok.torrserve.views.FloatingView
import kotlin.concurrent.thread


object TorrentInfo {

    private var Hash: String = ""
    private var isWatching: Boolean = false

    fun showWindow(hash: String?) {
        if (hash == null) {
            closeWindow()
            return
        }
        Hash = hash
        watching()
    }

    fun closeWindow() {
        isWatching = false
    }

    private fun watching() {

        synchronized(isWatching) {
            if (isWatching)
                return
            isWatching = true
        }

        val view = FloatingView()

        Handler(Looper.getMainLooper()).post {
            try {
                view?.create()
            } catch (e: Exception) {
                isWatching = false
                e.printStackTrace()
            }
        }
        if (!isWatching)
            return

        thread {
            var isShow = false
            while (isWatching && Preferences.isShowState()) {
                view?.getView()?.let { view ->
                    val info = ServerApi.info(Hash)
                    info?.let {
                        isShow = true
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
                        if (isShow)
                            isWatching = false
                    }
                    Thread.sleep(1000)
                }
            }
            isWatching = false
            view?.remove()
        }
    }
}