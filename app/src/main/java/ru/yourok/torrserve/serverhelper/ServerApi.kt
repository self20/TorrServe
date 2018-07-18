package ru.yourok.torrserve.serverhelper

import android.content.Intent
import android.net.Uri
import ru.yourok.torrserve.App
import ru.yourok.torrserve.serverhelper.server.Settings
import ru.yourok.torrserve.serverhelper.server.Torrent
import ru.yourok.torrserve.serverhelper.server.TorrentStats
import java.io.File
import java.io.FileInputStream
import java.io.FileOutputStream
import kotlin.concurrent.thread


/**
 * Created by yourok on 20.02.18.
 */


object ServerApi {
    fun add(link: String, save: Boolean): String {
        if (link.startsWith("magnet:", true))
            return addLink(ServerRequests.fixLink(link), save)
        else if (
                link.startsWith("http:", true) ||
                link.startsWith("https:", true))
            return addLink(link, save)
        else
            return addFile(link, save)
    }

    private fun addLink(link: String, save: Boolean): String {
        return ServerRequests.addTorrent(link, save)
    }

    private fun addFile(path: String, save: Boolean): String {
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
            link = Uri.parse(link).path

        val hash = ServerRequests.uploadFile(link, save)
        if (isRemove)
            File(link).delete()
        return hash
    }

    fun rem(hash: String) {
        ServerRequests.removeTorrent(hash)
    }

    fun get(hash: String): Torrent {
        return ServerRequests.getTorrent(hash)
    }

    fun list(): List<Torrent> {
        return ServerRequests.listTorrent()
    }

    fun stat(hash: String): TorrentStats {
        return ServerRequests.statTorrent(hash)
    }

    fun drop(hash: String) {
        if (hash.isEmpty())
            return
        ServerRequests.dropTorrent(hash)
    }

    fun preload(fileLink: String) {
        if (fileLink.isEmpty())
            return
        ServerRequests.preloadTorrent(fileLink)
    }

    fun echo(): Boolean {
        var echo = false
        thread {
            try {
                ServerRequests.echo()
                echo = true
            } catch (e: Exception) {
            }
        }.join(1000)
        return echo
    }

    fun readSettings(): Settings {
        return ServerRequests.readSettings()
    }

    fun writeSettings(sets: Settings) {
        ServerRequests.writeSettings(sets)
    }

    fun restartTorrentClient() {
        ServerRequests.restartTorrentClient()
    }

    fun shutdownServer() {
        ServerRequests.serverShutdown()
    }

    fun openPlayList(torrent: Torrent) {
        if (torrent.Hash().isEmpty())
            return

        val intent = Intent(Intent.ACTION_VIEW)
        val url = Uri.parse(ServerRequests.getHostUrl(torrent.Hash()))
        intent.setDataAndType(url, "audio/x-mpegurl")
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        if (torrent.Name().isNotEmpty()) {
            intent.putExtra("title", torrent.Name())
            intent.putExtra("name", torrent.Name())
        }
        App.getContext().startActivity(intent)
    }

    fun openPlayList() {
        val intent = Intent(Intent.ACTION_VIEW)
        val url = Uri.parse(ServerRequests.getHostUrl("/torrent/playlist.m3u"))
        intent.setDataAndType(url, "audio/x-mpegurl")
        intent.setFlags(Intent.FLAG_ACTIVITY_NEW_TASK)
        intent.putExtra("title", "TorrServePlaylist")
        intent.putExtra("name", "TorrServePlaylist")
        App.getContext().startActivity(intent)
    }
}


