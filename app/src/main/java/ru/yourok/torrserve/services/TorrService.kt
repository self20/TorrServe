package ru.yourok.torrserve.services

import android.app.Service
import android.content.Intent
import android.os.Handler
import android.os.IBinder
import android.widget.Toast
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.serverloader.ServerLoader
import kotlin.concurrent.thread


/**
 * Created by yourok on 20.02.18.
 */

class TorrService : Service() {
    override fun onBind(p0: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        intent?.let {
            if (it.action != null) {
                when (it.action) {
                    "ru.yourok.torrserve.notifications.action_exit" -> {
                        thread {
                            Thread.sleep(1000)
                            stopServer()
                            System.exit(0)
                        }
                        return START_NOT_STICKY
                    }
                    "ru.yourok.torrserve.notifications.action_restart" -> {
                        restartServer()
                        return START_STICKY
                    }
                }
            }

            if (it.hasExtra("Cmd")) {
                val cmd = it.getStringExtra("Cmd")
                when (cmd) {
                    "Stop" -> stopServer()
                    "Restart" -> restartServer()
                    else -> startServer()
                }
            }
        }
        startServer()
        return START_STICKY
    }

    private fun startServer() {
        thread {
            if (!ServerApi.echo()) {
                if (ServerLoader.serverExists())
                    ServerLoader.run()
                NotificationServer.Show(this, "")
            }
        }
    }

    private fun stopServer() {
        NotificationServer.Close(this)
        if (ServerApi.echo()) {
            ServerApi.shutdownServer()
            ServerLoader.stop()
            Handler(this.getMainLooper()).post(Runnable {
                Toast.makeText(this, R.string.server_stoped, Toast.LENGTH_LONG).show()
            })
        }
        stopSelf()
    }

    private fun restartServer() {
        thread {
            ServerLoader.stop()
            ServerLoader.run()
            Handler(this.getMainLooper()).post(Runnable {
                Toast.makeText(this, R.string.stat_server_is_running, Toast.LENGTH_SHORT).show()
            })
        }
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

        fun exit() {
            try {
                val intent = Intent(App.getContext(), TorrService::class.java)
                intent.action = "ru.yourok.torrserve.notifications.action_exit"
                App.getContext().startService(intent)
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
    }
}