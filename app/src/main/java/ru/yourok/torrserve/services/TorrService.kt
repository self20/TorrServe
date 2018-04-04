package ru.yourok.torrserve.services

import android.app.Service
import android.content.Intent
import android.os.Handler
import android.os.IBinder
import android.widget.Toast
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread


/**
 * Created by yourok on 20.02.18.
 */

class TorrService : Service() {
    override fun onBind(p0: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        thread {
            intent?.let {
                if (it.hasExtra("Cmd")) {
                    val cmd = it.getStringExtra("Cmd")
                    when (cmd) {
                        "Stop" -> stopServer()
                        "Restart" -> restartServer()
                        else -> startServer()
                    }
                    return@thread
                }

            }
            startServer()
        }
        return START_STICKY
    }

    private fun startServer() {
        if (!ServerApi.echo())
            torrentserver.Torrentserver.start(Utils.getAppPath())
    }

    private fun stopServer() {
        thread {
            if (ServerApi.echo()) {
                torrentserver.Torrentserver.stop()
                torrentserver.Torrentserver.waitServer()
                Handler(this.getMainLooper()).post(Runnable {
                    Toast.makeText(this, R.string.server_stoped, Toast.LENGTH_LONG).show()
                })
            }
            stopSelf()
        }
    }

    private fun restartServer() {
        torrentserver.Torrentserver.stop()
        torrentserver.Torrentserver.waitServer()
        torrentserver.Torrentserver.start(Utils.getAppPath())
    }

    companion object {
        fun start() {
            try {
                val intent = Intent(App.getContext(), TorrService::class.java)
                intent.putExtra("Cmd", "Start")
                App.getContext().startService(intent)
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }

        fun stop() {
            try {
                val intent = Intent(App.getContext(), TorrService::class.java)
                intent.putExtra("Cmd", "Stop")
                App.getContext().startService(intent)
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }

        fun waitServer(): Boolean {
            start()
            var count = 0
            while (!ServerApi.echo()) {
                Thread.sleep(1000)
                count++
                if (count % 10 == 0)
                    start()
                if (count > 60)
                    return false
            }
            return true
        }
    }
}