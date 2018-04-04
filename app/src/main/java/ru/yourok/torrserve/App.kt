package ru.yourok.torrserve

import android.app.Application
import android.content.Context
import android.os.PowerManager


class App : Application() {
    companion object {
        private lateinit var contextApp: Context
        private lateinit var wakeLock: PowerManager.WakeLock

        fun getContext(): Context {
            return contextApp
        }

        fun wakeLock(timeout: Long) {
            wakeLock.acquire(timeout)
        }
    }

    override fun onCreate() {
        super.onCreate()

        val powerManager = getSystemService(Context.POWER_SERVICE) as PowerManager
        wakeLock = powerManager.newWakeLock(PowerManager.PARTIAL_WAKE_LOCK, "TorrServeWakeLock")
        contextApp = applicationContext

        ACR.get(this)
                .setEmailAddresses("8yourok8@gmail.com")
                .setEmailSubject(getString(R.string.app_name) + " Crash Report")
                .start()
    }
}
