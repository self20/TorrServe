package ru.yourok.torrserve

import android.app.Activity
import android.app.AlertDialog
import android.content.Context
import android.content.Intent
import android.net.Uri
import android.os.Handler
import android.os.Looper
import android.support.design.widget.Snackbar
import android.support.v4.content.ContextCompat.startActivity
import ru.yourok.torrserve.serverhelper.Preferences
import java.util.*
import kotlin.concurrent.thread


object Donate {

    fun donateDialog(context: Context) {
        AlertDialog.Builder(context)
                .setTitle(R.string.donate)
                .setMessage(R.string.donate_msg)
                .setPositiveButton(R.string.paypal, { _, _ ->
                    val cur = Currency.getInstance(Locale.getDefault())
                    val mon = cur.toString()
                    val link = "https://www.paypal.me/yourok/0$mon"
                    val browserIntent = Intent(Intent.ACTION_VIEW, Uri.parse(link))
                    startActivity(context, browserIntent, null)
                    Preferences.setLastViewDonate(-1L)
                })
                .setNegativeButton(R.string.yandex_money, { _, _ ->
                    val browserIntent = Intent(Intent.ACTION_VIEW, Uri.parse("https://money.yandex.ru/to/410013733697114/100"))
                    startActivity(context, browserIntent, null)
                    Preferences.setLastViewDonate(-1L)
                })
                .setNeutralButton(R.string.google_play, { _, _ ->
                    val appPackageName = "ru.yourok.torrserve"
                    try {
                        context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse("market://details?id=$appPackageName")))
                    } catch (anfe: android.content.ActivityNotFoundException) {
                        context.startActivity(Intent(Intent.ACTION_VIEW, Uri.parse("https://play.google.com/store/apps/details?id=$appPackageName")))
                    }
                    Preferences.setLastViewDonate(-1L)
                })
                .show()
    }

    @Volatile
    private var showDonate = false

    fun showDonate(activity: Activity) {
        if (BuildConfig.FLAVOR != "pay")
            thread {
                synchronized(showDonate) {
                    val last: Long = Preferences.getLastViewDonate()
                    if (last == -1L || System.currentTimeMillis() < last || showDonate)
                        return@thread
                    showDonate = true
                    Preferences.setLastViewDonate(System.currentTimeMillis() + 5 * 60 * 1000)
                }

                val snackbar = Snackbar.make(activity.findViewById(android.R.id.content), R.string.donate, Snackbar.LENGTH_INDEFINITE)
                Handler(Looper.getMainLooper()).postDelayed(Runnable {
                    snackbar
                            .setAction(android.R.string.ok) {
                                Preferences.setLastViewDonate(System.currentTimeMillis())
                                donateDialog(activity)
                            }
                            .show()
                }, 5000)
                Handler(Looper.getMainLooper()).postDelayed(Runnable {
                    if (snackbar.isShown)
                        snackbar.dismiss()
                    showDonate = false
                }, 15000)
            }
    }
}