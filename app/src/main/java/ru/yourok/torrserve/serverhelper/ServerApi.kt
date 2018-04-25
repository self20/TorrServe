package ru.yourok.torrserve.serverhelper

import android.content.Context
import android.content.Intent
import android.net.Uri
import org.json.JSONArray
import org.json.JSONObject
import ru.yourok.torrserve.App
import ru.yourok.torrserve.services.TorrService
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
        val ret = torrentserver.Torrentserver.torrentServerAdd(addr, link)
        val tor = Torrent()
        val js = JSONObject(ret)
        tor.Name = js.getString("Name")
        tor.Magnet = js.getString("Magnet")
        tor.Hash = js.getString("Hash")
        if (isRemove)
            File(link).delete()
        return tor
    }

    fun rem(id: String): String {
        var addr = Preferences.getServerAddress()
        try {
            torrentserver.Torrentserver.torrentServerRem(addr, id)
            return ""
        } catch (e: Exception) {
            e.printStackTrace()
            return e.message ?: "Error remove torrent"
        }
    }

    fun get(id: String): Torrent? {
        val addr = Preferences.getServerAddress()
        try {
            val js = torrentserver.Torrentserver.torrentServerGet(addr, id)
            return js2Torrent(js)
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return null
    }

    fun list(): List<Torrent> {
        val retArr = mutableListOf<Torrent>()
        try {
            val addr = Preferences.getServerAddress()
            val ret = torrentserver.Torrentserver.torrentServerList(addr)
            val jsArr = JSONArray(ret)

            for (i in 0 until jsArr.length()) {
                val js = jsArr.getJSONObject(i)
                val tor = js2Torrent(js)
                retArr.add(tor)
            }
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return retArr
    }

    fun info(hash: String): Info? {
        val addr = Preferences.getServerAddress()
        try {
            val js = torrentserver.Torrentserver.torrentServerInfo(addr, hash)
            return js2Info(js)
        } catch (e: Exception) {
            e.printStackTrace()
        }
        return null
    }

    fun cleanCache() {
        try {
            val addr = Preferences.getServerAddress()
            torrentserver.Torrentserver.torrentServerCleanCache(addr)
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
        TorrService.showInfoWindow(hash)
    }

    fun echo(): Boolean {
        var addr = Preferences.getServerAddress()
        var echo = false
        thread {
            try {
                torrentserver.Torrentserver.torrentServerEcho(addr)
                echo = true
            } catch (e: Exception) {
            }
        }.join(1000)
        return echo
    }

    fun readSettings(): ServerSettings? {
        try {
            val addr = Preferences.getServerAddress()
            val strSets = torrentserver.Torrentserver.torrentServerReadSets(addr)

            val js = JSONObject(strSets)

            return ServerSettings(
                    js.getInt("CacheSize") / (1024 * 1024),
                    js.getInt("PreloadBufferSize") / (1024 * 1024),
                    js.getBoolean("DisableTCP"),
                    js.getBoolean("DisableUTP"),
                    js.getBoolean("DisableUPNP"),
                    js.getBoolean("DisableDHT"),
                    js.getBoolean("DisableUpload"),
                    js.getInt("Encryption"),
                    js.getInt("DownloadRateLimit"),
                    js.getInt("UploadRateLimit"),
                    js.getInt("ConnectionsLimit"))

        } catch (e: Exception) {
            e.printStackTrace()
            return null
        }
    }

    fun writeSettings(sets: ServerSettings): String {
        try {
            val addr = Preferences.getServerAddress()
            val js = JSONObject()
            js.put("CacheSize", sets.CacheSize * (1024 * 1024))
            js.put("PreloadBufferSize", sets.PreloadBufferSize * (1024 * 1024))
            js.put("DisableTCP", sets.DisableTCP)
            js.put("DisableUTP", sets.DisableUTP)
            js.put("DisableUPNP", sets.DisableUPNP)
            js.put("DisableDHT", sets.DisableDHT)
            js.put("DisableUpload", sets.DisableUpload)
            js.put("Encryption", sets.Encryption)
            js.put("DownloadRateLimit", sets.DownloadRateLimit)
            js.put("UploadRateLimit", sets.UploadRateLimit)
            js.put("ConnectionsLimit", sets.ConnectionsLimit)
            torrentserver.Torrentserver.torrentServerWriteSets(addr, js.toString())
            return ""
        } catch (e: Exception) {
            e.printStackTrace()
            return e.message ?: "error parse settings"
        }
    }
}

data class ServerSettings(var CacheSize: Int,
                          var PreloadBufferSize: Int,
                          var DisableTCP: Boolean,
                          var DisableUTP: Boolean,
                          var DisableUPNP: Boolean,
                          var DisableDHT: Boolean,
                          var DisableUpload: Boolean,
                          var Encryption: Int,
                          var DownloadRateLimit: Int,
                          var UploadRateLimit: Int,
                          var ConnectionsLimit: Int)

