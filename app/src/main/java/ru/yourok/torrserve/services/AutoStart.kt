package ru.yourok.torrserve.services

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.widget.Toast
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerApi


class BootCompletedReceiver : BroadcastReceiver() {
    override fun onReceive(context: Context, intent: Intent) {
        if (Preferences.isAutoStart()) {
            val intent = Intent(context, TorrService::class.java)
            intent.putExtra("Cmd", "Start")
//            if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.O)
//                context.startForegroundService(intent)
//            else
            context.startService(intent)

            if (!waitServer()) {
                Toast.makeText(App.getContext(), context.getResources().getString(R.string.error_server_start), Toast.LENGTH_LONG).show()
            } else
                Toast.makeText(App.getContext(), context.getResources().getString(R.string.stat_server_is_running), Toast.LENGTH_LONG).show()
        }
    }

    private fun waitServer(): Boolean {
        var count = 0
        while (!ServerApi.echo()) {
            Thread.sleep(1000)
            count++
            if (count % 10 == 0)
                TorrService.start()
            if (count > 60)
                return false
        }
        return true
    }
}