package ru.yourok.torrserve.serverhelper

import android.net.Uri
import cz.msebera.android.httpclient.client.methods.HttpEntityEnclosingRequestBase
import cz.msebera.android.httpclient.client.methods.HttpGet
import cz.msebera.android.httpclient.client.methods.HttpPost
import cz.msebera.android.httpclient.entity.ContentType
import cz.msebera.android.httpclient.entity.StringEntity
import cz.msebera.android.httpclient.entity.mime.HttpMultipartMode
import cz.msebera.android.httpclient.entity.mime.MultipartEntityBuilder
import cz.msebera.android.httpclient.entity.mime.content.FileBody
import cz.msebera.android.httpclient.entity.mime.content.StringBody
import cz.msebera.android.httpclient.impl.client.HttpClients
import cz.msebera.android.httpclient.util.EntityUtils
import org.json.JSONArray
import org.json.JSONObject
import ru.yourok.torrserve.serverhelper.server.Settings
import ru.yourok.torrserve.serverhelper.server.Torrent
import ru.yourok.torrserve.serverhelper.server.TorrentStats
import java.io.IOException
import java.net.URI

object ServerRequests {

    fun addTorrent(link: String, save: Boolean): String {
        val url = getHostUrl("/torrent/add")
        val js = JSONObject()
        js.put("Link", link)
        js.put("DontSave", !save)
        val req = js.toString(0)
        return Post(url, req)
    }

    fun uploadFile(path: String, save: Boolean): String {
        val url = getHostUrl("/torrent/upload")
        val hashes = Upload(url, path, save)
        return hashes[0]
    }

    fun getTorrent(hash: String): Torrent {
        val url = getHostUrl("/torrent/get")
        val js = JSONObject()
        js.put("Hash", hash)
        val req = js.toString(0)
        return Torrent(JSONObject(Post(url, req)))
    }

    fun removeTorrent(hash: String) {
        val url = getHostUrl("/torrent/rem")
        val js = JSONObject()
        js.put("Hash", hash)
        val req = js.toString(0)
        Post(url, req)
    }

    fun listTorrent(): List<Torrent> {
        val url = getHostUrl("/torrent/list")
        val torjs = JSONArray(Post(url, ""))
        val list = mutableListOf<Torrent>()
        for (i in 0 until torjs.length())
            list.add(Torrent(torjs.getJSONObject(i)))
        return list.toList()
    }

    fun statTorrent(hash: String): TorrentStats {
        val url = getHostUrl("/torrent/stat")
        val js = JSONObject()
        js.put("Hash", hash)
        val req = js.toString(0)
        return TorrentStats(JSONObject(Post(url, req)))
    }

    fun cacheTorrent(hash: String): JSONObject {
        val url = getHostUrl("/torrent/cache")
        val js = JSONObject()
        js.put("Hash", hash)
        val req = js.toString(0)
        return JSONObject(Post(url, req))
    }

    fun dropTorrent(hash: String) {
        val url = getHostUrl("/torrent/drop")
        val js = JSONObject()
        js.put("Hash", hash)
        val req = js.toString(0)
        Post(url, req)
    }

    fun preloadTorrent(fileLink: String) {
        val link = fileLink.replace("/torrent/view/", "/torrent/preload/")
        val url = getHostUrl(link)
        Get(url)
    }

    fun restartTorrentClient() {
        val url = getHostUrl("/torrent/restart")
        Get(url)
    }

    fun readSettings(): Settings {
        val url = getHostUrl("/settings/read")
        return Settings(JSONObject(Post(url, "")))
    }

    fun writeSettings(sets: Settings) {
        val url = getHostUrl("/settings/write")
        Post(url, sets.json.toString(0))
    }

    fun echo(): String {
        val url = getHostUrl("/echo")
        return Get(url)
    }

    fun serverShutdown() {
        val url = getHostUrl("/shutdown")
        Post(url, "")
    }

    ///////////////////////////////////////////////////////////

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

    fun getHostUrl(path: String): String {
        val url = Preferences.getServerAddress()
        if (url.last() == '/')
            return url + path.substring(1)
        else
            return url + path
    }

    private fun Upload(url: String, path: String, save: Boolean): List<String> {
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

    private fun Post(url: String, req: String): String {
        val httpclient = HttpClients.custom().disableRedirectHandling().build()
        val httpreq = HttpPost(url)
        if (req.isNotEmpty())
            (httpreq as HttpEntityEnclosingRequestBase).setEntity(StringEntity(req))

        val response = httpclient.execute(httpreq)
        val status = response.statusLine?.statusCode ?: -1
        if (status == 200) {
            val entity = response.entity ?: return ""
            return EntityUtils.toString(entity)
        } else if (status == 302) {
            return ""
        } else {
            val resp = EntityUtils.toString(response.entity)
            resp?.let {
                if (it.isNotEmpty()) {
                    var errMsg = response.statusLine.reasonPhrase
                    try {
                        errMsg = JSONObject(it).getString("Message")
                    } catch (e: Exception) {
                        try {
                            errMsg = JSONObject(it).getString("message")
                        } catch (e: Exception) {
                        }
                    }
                    throw IOException(errMsg)
                }
            }
            throw IOException(response.statusLine.reasonPhrase)
        }
    }

    private fun Get(url: String): String {
        val httpclient = HttpClients.custom().disableRedirectHandling().build()

        val httpreq = HttpGet(url)
        val response = httpclient.execute(httpreq)
        val status = response.statusLine?.statusCode ?: -1
        if (status == 200) {
            val entity = response.entity ?: return ""
            return EntityUtils.toString(entity)
        } else if (status == 302) {
            return ""
        } else {
            val resp = EntityUtils.toString(response.entity)
            resp?.let {
                if (it.isNotEmpty()) {
                    var errMsg = response.statusLine.reasonPhrase
                    try {
                        errMsg = JSONObject(it).getString("Message")
                    } catch (e: Exception) {
                        try {
                            errMsg = JSONObject(it).getString("message")
                        } catch (e: Exception) {
                        }
                    }
                    throw IOException(errMsg)
                }
            }
            throw IOException(response.statusLine.reasonPhrase)
        }
    }
}