package ru.yourok.torrserve.serverhelper

import android.content.Context
import android.content.Intent
import android.net.Uri
import ru.yourok.torrserve.App
import ru.yourok.torrserve.services.NotificationServer
import ru.yourok.torrserve.utils.Mime
import java.io.File
import java.io.FileInputStream
import java.io.FileOutputStream
import kotlin.concurrent.thread


/**
 * Created by yourok on 20.02.18.
 */


object ServerApi {
    fun add(magnet: String): Torrent {
        var link = magnet
        var isRemove = false
        if (magnet.startsWith("content://", true)) {
            val outputDir = App.getContext().getCacheDir() // context being the Activity pointer
            val outputFile = File.createTempFile("tmp", ".torr", outputDir)
            val fd = App.getContext().contentResolver.openFileDescriptor(Uri.parse(magnet), "r")
            val inStream = FileInputStream(fd.fileDescriptor)
            val outStream = FileOutputStream(outputFile)
            val inChannel = inStream.getChannel()
            val outChannel = outStream.getChannel()
            inChannel.transferTo(0, inChannel.size(), outChannel)
            inStream.close()
            outStream.close()
            link = outputFile.path
            isRemove = true
        }

        val addr = Preferences.getServerAddress()
        val tor = ServerRequest.serverAdd(addr, link)
        if (isRemove)
            File(link).delete()
        return tor
    }

    fun rem(hash: String): String {
        var addr = Preferences.getServerAddress()
        try {
            ServerRequest.serverRem(addr, hash)
            return ""
        } catch (e: Exception) {
            e.printStackTrace()
            return e.message ?: "Error remove torrent"
        }
    }

    fun get(hash: String): Torrent? {
        val addr = Preferences.getServerAddress()
        try {
            return ServerRequest.serverGet(addr, hash)
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return null
    }

    fun list(): List<Torrent> {
        val retArr = mutableListOf<Torrent>()
        try {
            val addr = Preferences.getServerAddress()
            return ServerRequest.serverList(addr)
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return retArr
    }

    fun info(hash: String): Info? {
        if (hash.isEmpty())
            return null
        try {
            val addr = Preferences.getServerAddress()
            return ServerRequest.serverInfo(addr, hash)
        } catch (e: Exception) {
        }
        return null
    }

    fun cleanCache() {
        try {
            val addr = Preferences.getServerAddress()
            ServerRequest.serverCleanCache(addr)
        } catch (e: Exception) {
            e.printStackTrace()
        }
    }

    fun view(context: Context, hash: String, name: String, link: String) {
        var addr = Preferences.getServerAddress()
        addr += link
        val browserIntent = Intent(Intent.ACTION_VIEW)
        browserIntent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        browserIntent.putExtra("title", name)
        browserIntent.setDataAndType(Uri.parse(addr), Mime.getMimeType(link))
        if (browserIntent.resolveActivity(context.packageManager) != null) {
            context.startActivity(browserIntent)
        } else {
            val intent = Intent.createChooser(browserIntent, "")
            context.startActivity(intent)
        }
        NotificationServer.Show(context, hash)
    }

    fun echo(): Boolean {
        var addr = Preferences.getServerAddress()
        var echo = false
        thread {
            try {
                ServerRequest.serverEcho(addr)
                echo = true
            } catch (e: Exception) {
            }
        }.join(1000)
        return echo
    }

    fun readSettings(): ServerSettings? {
        try {
            var sets: ServerSettings? = null
            thread {
                sets = ServerRequest.readSettings()
            }.join()
            return sets
        } catch (e: Exception) {
            e.printStackTrace()
            return null
        }
    }

    fun writeSettings(sets: ServerSettings): String {
        var err = ""
        thread {
            try {
                ServerRequest.writeSettings(sets)
            } catch (e: Exception) {
                e.printStackTrace()
                err = e.message ?: "error parse settings"
            }
        }.join()
        return err
    }
}


