package ru.yourok.torrserve.services

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Context.NOTIFICATION_SERVICE
import android.content.Intent
import android.os.Build
import android.support.v4.app.NotificationCompat
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R
import ru.yourok.torrserve.activitys.MainActivity
import ru.yourok.torrserve.activitys.ViewActivity
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.utils.Utils
import kotlin.concurrent.thread


object NotificationServer {

    private val channelId = "ru.yourok.torrserve"
    private val channelName = "ru.yourok.torrserve"
    private var update = false
    private var hash = ""
    private var builder: NotificationCompat.Builder? = null
    private val lock = Any()

    private fun build(context: Context, msg: String) {
        synchronized(lock) {
            val restartIntent = Intent(context, TorrService::class.java)
            restartIntent.setAction("ru.yourok.torrserve.notifications.action_restart")
            val restartPendingIntent = PendingIntent.getService(context, 0, restartIntent, 0)

            val exitIntent = Intent(context, TorrService::class.java)
            exitIntent.setAction("ru.yourok.torrserve.notifications.action_exit")
            val exitPendingIntent = PendingIntent.getService(context, 0, exitIntent, 0)

            val intent = Intent(context, MainActivity::class.java)
            intent.addFlags(Intent.FLAG_ACTIVITY_CLEAR_TOP)
            val pendingIntent = PendingIntent.getActivity(context, 0, intent, PendingIntent.FLAG_ONE_SHOT)

            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O) {
                val channel = NotificationChannel(channelId, channelName, NotificationManager.IMPORTANCE_LOW)
                context.getSystemService<NotificationManager>(NotificationManager::class.java)!!.createNotificationChannel(channel)
            }

            if (builder == null)
                builder = NotificationCompat.Builder(context, channelId)
                        .setSmallIcon(R.drawable.ic_launcher)
                        .setContentTitle(context.getString(R.string.app_name))
//                .setContentText(msg)
                        .setAutoCancel(false)
                        .setOngoing(true)
                        .setContentIntent(pendingIntent)
                        .setStyle(NotificationCompat.BigTextStyle().bigText(msg))
                        .addAction(android.R.drawable.stat_notify_sync, context.getText(R.string.restart_server), restartPendingIntent)
                        .addAction(android.R.drawable.ic_delete, context.getText(R.string.exit), exitPendingIntent)
            else
                builder?.setStyle(NotificationCompat.BigTextStyle().bigText(msg))

            builder?.let {
                val notificationManager = context.getSystemService(NOTIFICATION_SERVICE) as NotificationManager
                notificationManager.notify(0, it.build())
            }
        }
    }

    fun Show(context: Context, hash: String) {
        this.hash = hash
        synchronized(update) {
            if (update)
                return
            update = true
        }

        thread {
            var isShow = false
            while (update) {
                val info = ServerApi.info(this.hash)
                info?.let {
                    var msg = context.getString(R.string.peers) + ": " + it.ConnectedSeeders.toString() + " / " + it.TotalPeers.toString() + "\n" +
                            context.getString(R.string.download_speed) + ": " + Utils.byteFmt(it.DownloadSpeed)
                    if (it.UploadSpeed > 0)
                        msg += "\n" + context.getString(R.string.upload_speed) + ": " + Utils.byteFmt(it.UploadSpeed)

                    if (info.IsPreload && !isShow && info.PreloadLength > 0) {
                        msg += "\n" + context.getString(R.string.buffer) + ": " + (info.PreloadOffset * 100 / info.PreloadLength).toString() + "% " + Utils.byteFmt(info.PreloadOffset) + "/" + Utils.byteFmt(info.PreloadLength)
                        if (Preferences.isShowPreloadWnd()) {
                            val intent = Intent(App.getContext(), ViewActivity::class.java)
                            intent.putExtra("Preload", this.hash)
                            intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
                            App.getContext().startActivity(intent)
                            isShow = true
                        }
                    }
                    build(context, msg)

                } ?: let {
                    build(context, context.getText(R.string.stat_server_is_running).toString())
                }
                Thread.sleep(1000)
            }
            build(context, context.getText(R.string.stat_server_is_running).toString())
            update = false
        }
    }

    fun Close(context: Context) {
        synchronized(lock) {
            update = false
            builder = null
            val notificationManager = context.getSystemService(NOTIFICATION_SERVICE) as NotificationManager
            notificationManager.cancel(0)
        }
    }
}