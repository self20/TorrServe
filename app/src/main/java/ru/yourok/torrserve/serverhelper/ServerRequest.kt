package ru.yourok.torrserve.serverhelper

import android.net.Uri
import cz.msebera.android.httpclient.client.methods.HttpPost
import cz.msebera.android.httpclient.entity.ContentType
import cz.msebera.android.httpclient.entity.mime.HttpMultipartMode
import cz.msebera.android.httpclient.entity.mime.MultipartEntityBuilder
import cz.msebera.android.httpclient.entity.mime.content.FileBody
import cz.msebera.android.httpclient.entity.mime.content.StringBody
import cz.msebera.android.httpclient.impl.client.HttpClients
import cz.msebera.android.httpclient.util.EntityUtils
import org.json.JSONArray
import org.json.JSONObject
import java.io.DataOutputStream
import java.io.IOException
import java.net.HttpURLConnection
import java.net.URI
import java.net.URL
import java.nio.charset.Charset


data class Torrent(
        var Name: String,
        var Magnet: String,
        var Hash: String,
        var Length: Long,
        var AddTime: Long,
        var Size: Long,
        var IsGettingInfo: Boolean,
        var Playlist: String,
        var Files: List<File>
)

data class File(
        var Name: String,
        var Link: String,
        var Size: Long,
        var Viewed: Boolean
)

data class Info(
        var Name: String,
        var Hash: String,

        var BytesWritten: Long,
        var BytesWrittenData: Long,

        var BytesRead: Long,
        var BytesReadData: Long,
        var BytesReadUsefulData: Long,

        var ChunksWritten: Long,

        var ChunksRead: Long,
        var ChunksReadUseful: Long,
        var ChunksReadUnwanted: Long,

        var PiecesDirtiedGood: Long,
        var PiecesDirtiedBad: Long,

        var DownloadSpeed: Double,
        var UploadSpeed: Double,

        var TotalPeers: Int,
        var PendingPeers: Int,
        var ActivePeers: Int,
        var ConnectedSeeders: Int,
        var HalfOpenPeers: Int,

        var IsGettingInfo: Boolean,
        var IsPreload: Boolean,
        var PreloadSize: Long,
        var PreloadLength: Long
)


data class ServerSettings(var CacheSize: Int,
                          var PreloadBufferSize: Int,
                          var RetrackersMode: Int,
                          var DisableTCP: Boolean,
                          var DisableUTP: Boolean,
                          var DisableUPNP: Boolean,
                          var DisableDHT: Boolean,
                          var DisableUpload: Boolean,
                          var Encryption: Int,
                          var DownloadRateLimit: Int,
                          var UploadRateLimit: Int,
                          var ConnectionsLimit: Int)

fun fixLink(link: String): String {
    try {
        if (link.isNotEmpty()) {
            val url = Uri.parse(link)
            val uri = URI(url.scheme, url.userInfo, url.host, url.port, url.path, url.query, url.fragment)
            return uri.toASCIIString()
        }
    } catch (e: Exception) {
    }
    return link
}

fun getRequest(link: String, hash: String): String {
    val js = JSONObject()
    js.put("Link", link)
    js.put("Hash", hash)
    return js.toString(0)
}

fun getRequest(link: String, dontSave: Boolean): String {
    val js = JSONObject()
    js.put("Link", link)
    js.put("DontSave", dontSave)
    return js.toString(0)
}

fun getTorrent(jsStr: String): Torrent {
    return getTorrent(JSONObject(jsStr))
}

fun getTorrent(js: JSONObject): Torrent {
    val fileList = mutableListOf<File>()
    if (js.has("Files")) {
        val jsFiles = js.getJSONArray("Files")

        for (i in 0 until jsFiles.length()) {
            val jsf = jsFiles.getJSONObject(i)
            val tf = File(
                    jsf.getString("Name"),
                    jsf.getString("Link"),
                    jsf.getLong("Size"),
                    jsf.getBoolean("Viewed"))
            fileList.add(tf)
        }
    }

    val ret = Torrent(
            js.getString("Name"),
            js.getString("Magnet"),
            js.getString("Hash"),
            js.getLong("Length"),
            js.getLong("AddTime"),
            js.getLong("Size"),
            js.getBoolean("IsGettingInfo"),
            js.getString("Playlist"),
            fileList.toList()
    )
    return ret
}

fun js2Info(jsStr: String): Info {
    val js = JSONObject(jsStr)
    return js2Info(js)
}

fun js2Info(js: JSONObject): Info {
    return Info(
            js.getString("Name"),
            js.getString("Hash"),

            js.getLong("BytesWritten"),
            js.getLong("BytesWrittenData"),

            js.getLong("BytesRead"),
            js.getLong("BytesReadData"),
            js.getLong("BytesReadUsefulData"),

            js.getLong("ChunksWritten"),

            js.getLong("ChunksRead"),
            js.getLong("ChunksReadUseful"),
            js.getLong("ChunksReadUnwanted"),

            js.getLong("PiecesDirtiedGood"),
            js.getLong("PiecesDirtiedBad"),

            js.getDouble("DownloadSpeed"),
            js.getDouble("UploadSpeed"),

            js.getInt("TotalPeers"),
            js.getInt("PendingPeers"),
            js.getInt("ActivePeers"),
            js.getInt("ConnectedSeeders"),
            js.getInt("HalfOpenPeers"),

            js.getBoolean("IsGettingInfo"),
            js.getBoolean("IsPreload"),
            js.getLong("PreloadSize"),
            js.getLong("PreloadLength")
    )
}

object ServerRequest {
    fun joinUrl(url: String, path: String): String {
        if (url.last() == '/')
            return url + path.substring(1)
        else
            return url + path
    }

    private fun requestTorr(post: Boolean, url: String, req: String): Torrent {
        return getTorrent(requestStr(post, url, req))
    }

    private fun requestFile(url: String, path: String, save: Boolean): List<String> {
        val file = java.io.File(path)

        val httpclient = HttpClients.custom().build()
        val httppost = HttpPost(url)

        val mpEntity = MultipartEntityBuilder.create()
        mpEntity.setMode(HttpMultipartMode.BROWSER_COMPATIBLE)
        mpEntity.addPart(file.name, FileBody(file))
        if (!save)
            mpEntity.addPart("DontSave", StringBody("true", ContentType.DEFAULT_TEXT))

        val entity = mpEntity.build()
        httppost.setEntity(entity)
        val response = httpclient.execute(httppost)
        val str = EntityUtils.toString(response.getEntity())
        val arr = JSONArray(str)

        val hashList = mutableListOf<String>()
        for (i in 0 until arr.length()) {
            val str = arr.getString(i)
            hashList.add(str)
        }
        return hashList
    }

    private fun requestStr(post: Boolean, url: String, req: String): String {
        val url = URL(url)
        val conn = url.openConnection() as HttpURLConnection

        if (post) {
            conn.requestMethod = "POST"
            conn.setRequestProperty("Content-Type", "application/json")
        } else
            conn.requestMethod = "GET"
        conn.connect()

        if (req.isNotEmpty()) {
            val os = DataOutputStream(conn.outputStream)
            os.writeBytes(req)
            os.flush()
            os.close()
        }

        val input = conn.inputStream
        val str = input.bufferedReader(Charset.defaultCharset()).readText()

        if (conn.responseCode != 200) {
            val input = conn.errorStream
            var err = "Error connect to server"
            val str = input.bufferedReader(Charset.defaultCharset()).readText()
            if (str.isNotEmpty()) {
                err = try {
                    JSONObject(str).getString("Message")
                } catch (e: Exception) {
                    str
                }
            } else
                conn.responseMessage?.let {
                    err = it
                }
            input.close()
            throw IOException(err)
        }

        input.close()
        conn.disconnect()
        return str
    }

    fun serverAdd(host: String, link: String, save: Boolean): Torrent {
        val url = joinUrl(host, "/torrent/add")
        val req = getRequest(link, !save)
        val hash = requestStr(true, url, req)
        Thread.sleep(1000)
        return serverGet(host, hash)
    }

    fun serverAddFile(host: String, link: String, save: Boolean): List<Torrent> {
        val url = joinUrl(host, "/torrent/upload")
        val hashes = requestFile(url, link, save)
        Thread.sleep(1000)
        val torrs = serverList(host)
        return torrs.filter { tor ->
            val list = hashes.find {
                it == tor.Hash
            }
            list != null && list.isNotEmpty()
        }
    }

    fun serverGet(host: String, hash: String): Torrent {
        val url = joinUrl(host, "/torrent/get")
        val req = getRequest("", hash)
        return requestTorr(true, url, req)
    }

    fun serverRem(host: String, hash: String) {
        val url = joinUrl(host, "/torrent/rem")
        val req = getRequest("", hash)
        requestStr(true, url, req)
    }

    fun serverList(host: String): List<Torrent> {
        val url = joinUrl(host, "/torrent/list")
        val str = requestStr(true, url, "")
        val arr = JSONArray(str)

        val torrList = mutableListOf<Torrent>()
        for (i in 0 until arr.length()) {
            val js = arr.getJSONObject(i)
            val tor = getTorrent(js)
            torrList.add(tor)
        }
        return torrList
    }

    fun serverInfo(host: String, hash: String): Info {
        val url = joinUrl(host, "/torrent/stat")
        val req = getRequest("", hash)
        val str = requestStr(true, url, req)
        return js2Info(str)
    }

    fun serverEcho(host: String) {
        val url = joinUrl(host, "/echo")
        requestStr(false, url, "")
    }

    fun readSettings(): ServerSettings? {
        val host = Preferences.getServerAddress()
        val url = joinUrl(host, "/settings/read")
        val str = requestStr(true, url, "")

        val js = JSONObject(str)

        return ServerSettings(
                js.getInt("CacheSize") / (1024 * 1024),
                js.getInt("PreloadBufferSize") / (1024 * 1024),
                js.getInt("RetrackersMode"),
                js.getBoolean("DisableTCP"),
                js.getBoolean("DisableUTP"),
                js.getBoolean("DisableUPNP"),
                js.getBoolean("DisableDHT"),
                js.getBoolean("DisableUpload"),
                js.getInt("Encryption"),
                js.getInt("DownloadRateLimit"),
                js.getInt("UploadRateLimit"),
                js.getInt("ConnectionsLimit"))
    }

    fun writeSettings(sets: ServerSettings) {
        val host = Preferences.getServerAddress()
        val url = joinUrl(host, "/settings/write")
        val js = JSONObject()
        js.put("CacheSize", sets.CacheSize * (1024 * 1024))
        js.put("PreloadBufferSize", sets.PreloadBufferSize * (1024 * 1024))
        js.put("RetrackersMode", sets.RetrackersMode)
        js.put("DisableTCP", sets.DisableTCP)
        js.put("DisableUTP", sets.DisableUTP)
        js.put("DisableUPNP", sets.DisableUPNP)
        js.put("DisableDHT", sets.DisableDHT)
        js.put("DisableUpload", sets.DisableUpload)
        js.put("Encryption", sets.Encryption)
        js.put("DownloadRateLimit", sets.DownloadRateLimit)
        js.put("UploadRateLimit", sets.UploadRateLimit)
        js.put("ConnectionsLimit", sets.ConnectionsLimit)

        requestStr(true, url, js.toString(0))
    }

    fun serverPreload(host: String, fileLink: String) {
        val link = fileLink.replace("/torrent/view/", "/torrent/preload/")
        val url = joinUrl(host, link)
        requestStr(false, url, "")
    }
}