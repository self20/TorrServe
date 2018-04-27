package ru.yourok.torrserve.serverhelper

import org.json.JSONArray
import org.json.JSONObject
import java.io.DataOutputStream
import java.io.IOException
import java.net.HttpURLConnection
import java.net.URL
import java.nio.charset.Charset


data class Torrent(
        var Name: String,
        var Magnet: String,
        var Hash: String,
        var Length: Long,
        var AddTime: Long,
        var Size: Long,
        var Files: List<File>
)

data class File(
        var Name: String,
        var Link: String,
        var Size: Long,
        var Viewed: Boolean
)

data class Info(
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

        var IsPreload: Boolean,
        var PreloadOffset: Long,
        var PreloadLength: Long
)


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

fun getRequest(link: String, hash: String): String {
    val js = JSONObject()
    js.put("Link", link)
    js.put("Hash", hash)
    return js.toString(0)
}

fun getTorrent(jsStr: String): Torrent {
    return getTorrent(JSONObject(jsStr))
}

fun getTorrent(js: JSONObject): Torrent {
    val jsFiles = js.getJSONArray("Files")
    val fileList = mutableListOf<File>()

    for (i in 0 until jsFiles.length()) {
        val jsf = jsFiles.getJSONObject(i)
        val tf = File(
                jsf.getString("Name"),
                jsf.getString("Link"),
                jsf.getLong("Size"),
                jsf.getBoolean("Viewed"))
        fileList.add(tf)
    }

    val ret = Torrent(
            js.getString("Name"),
            js.getString("Magnet"),
            js.getString("Hash"),
            js.getLong("Length"),
            js.getLong("AddTime"),
            js.getLong("Size"),
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

            js.getBoolean("IsPreload"),
            js.getLong("PreloadOffset"),
            js.getLong("PreloadLength")
    )
}

object ServerRequest {
    private fun joinUrl(url: String, path: String): String {
        if (url.last() == '/')
            return url + path.substring(1)
        else
            return url + path
    }

    private fun requestTorr(post: Boolean, url: String, req: String): Torrent {
        return getTorrent(requestStr(post, url, req))
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

    fun serverAdd(host: String, link: String): Torrent {
        val url = joinUrl(host, "/torrent/add")
        val req = getRequest(link, "")
        return requestTorr(true, url, req)
    }


    fun serverGet(host: String, hash: String): Torrent {
        val url = joinUrl(host, "/torrent/get")
        val req = getRequest("", hash)
        return requestTorr(true, url, req)
    }

    fun serverRem(host: String, hash: String) {
        val url = joinUrl(host, "/torrent/rem")
        val req = getRequest("", hash)
        requestTorr(true, url, req)
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

    fun serverCleanCache(host: String) {
        val url = joinUrl(host, "/torrent/cleancache")
        requestStr(true, url, "")
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
}