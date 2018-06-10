package ru.yourok.torrserve.services

import android.app.Service
import android.content.Intent
import android.os.IBinder
import ru.yourok.torrserve.App
import ru.yourok.torrserve.serverhelper.ServerApi
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
        if (!ServerApi.echo()) {
            NotificationServer.Show(this, "")
        }
    }

    private fun stopServer() {
        NotificationServer.Close(this)
        stopSelf()
    }

    private fun restartServer() {
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