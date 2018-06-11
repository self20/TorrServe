package ru.yourok.torrserve.serverhelper

import android.content.Intent
import android.net.Uri
import ru.yourok.torrserve.App
import java.io.File
import java.io.FileInputStream
import java.io.FileOutputStream
import kotlin.concurrent.thread


/**
 * Created by yourok on 20.02.18.
 */


object ServerApi {
    fun add(link: String, save: Boolean): List<Torrent> {
        if (
                link.startsWith("magnet:", true) ||
                link.startsWith("http:", true) ||
                link.startsWith("https:", true))
            return listOf(addLink(link, save))
        else
            return addFile(link, save)
    }

    private fun addLink(link: String, save: Boolean): Torrent {
        val addr = Preferences.getServerAddress()
        val tor = ServerRequest.serverAdd(addr, link, save)
        return tor
    }

    private fun addFile(path: String, save: Boolean): List<Torrent> {
        var link = path
        var isRemove = false
        if (link.startsWith("content://", true)) {
            val outputDir = App.getContext().getCacheDir() // context being the Activity pointer
            val outputFile = File.createTempFile("tmp", ".torr", outputDir)
            val fd = App.getContext().contentResolver.openFileDescriptor(Uri.parse(link), "r")
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

        if (link.startsWith("file://", true))
            link = link.substring(7)

        val addr = Preferences.getServerAddress()
        val tor = ServerRequest.serverAddFile(addr, link, save)
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

    fun cleanCache(hash: String) {
        thread {
            try {
                val addr = Preferences.getServerAddress()
                ServerRequest.serverCleanCache(addr, hash)
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
    }

    fun preload(hash: String, fileLink: String) {
        if (hash.isEmpty() || fileLink.isEmpty())
            return
        try {
            val addr = Preferences.getServerAddress()
            ServerRequest.serverPreload(addr, fileLink)
        } catch (e: Exception) {
            e.printStackTrace()
        }
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
        var sets: ServerSettings? = null
        thread {
            try {
                sets = ServerRequest.readSettings()
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }.join()
        return sets
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

    fun openPlayList(torrent: Torrent) {
        if (torrent.Playlist.isEmpty()) {
            return
        }
        val addr = Preferences.getServerAddress()
        val intent = Intent(Intent.ACTION_VIEW)
        val url = Uri.parse(ServerRequest.joinUrl(addr, torrent.Playlist))
        intent.setDataAndType(url, "audio/x-mpegurl")
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        intent.putExtra("title", torrent.Name)
        intent.putExtra("name", torrent.Name)
        App.getContext().startActivity(intent)
    }

    fun openPlayList() {
        val addr = Preferences.getServerAddress()
        val intent = Intent(Intent.ACTION_VIEW)
        val url = Uri.parse(ServerRequest.joinUrl(addr, "/torrent/playlist.m3u"))
        intent.setDataAndType(url, "audio/x-mpegurl")
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        intent.putExtra("title", "TorrServePlaylist")
        intent.putExtra("name", "TorrServePlaylist")
        App.getContext().startActivity(intent)
    }
}


