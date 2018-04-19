package ru.yourok.torrserve.services

import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.content.Context
import android.content.Context.NOTIFICATION_SERVICE
import android.content.Intent
import android.os.Build
import android.support.v4.app.NotificationCompat
import ru.yourok.torrserve.R
import ru.yourok.torrserve.activitys.MainActivity


object NotificationServer {

    private val channelId = "ru.yourok.torrserve"
    private val channelName = "ru.yourok.torrserve"

    fun Show(context: Context) {

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


        val builder = NotificationCompat.Builder(context, channelId)
                .setSmallIcon(R.mipmap.ic_launcher)
                .setContentTitle(context.getString(R.string.app_name))
                .setContentText(context.getText(R.string.stat_server_is_running))
                .setAutoCancel(false)
                .setOngoing(true)
                .setContentIntent(pendingIntent)
                .addAction(android.R.drawable.stat_notify_sync, context.getText(R.string.restart_server), restartPendingIntent)
                .addAction(android.R.drawable.ic_delete, context.getText(R.string.exit), exitPendingIntent)
        val notification = builder.build()

        val notificationManager = context.getSystemService(NOTIFICATION_SERVICE) as NotificationManager
        notificationManager.notify(0, notification)
    }

    fun Close(context: Context) {
        val notificationManager = context.getSystemService(NOTIFICATION_SERVICE) as NotificationManager
        notificationManager.cancel(0)
    }
}