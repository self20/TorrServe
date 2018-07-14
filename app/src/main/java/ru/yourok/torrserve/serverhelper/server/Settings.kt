package ru.yourok.torrserve.serverhelper.server

import org.json.JSONObject

class Settings(val json: JSONObject) {
    var cacheSize: Long
        get() {
            if (json.has("CacheSize"))
                return json.getLong("CacheSize")
            return 0
        }
        set(value) {
            json.put("CacheSize", value)
        }

    var preloadBufferSize: Long
        get() {
            if (json.has("PreloadBufferSize"))
                return json.getLong("PreloadBufferSize")
            return 0
        }
        set(value) {
            json.put("PreloadBufferSize", value)
        }

    var retrackersMode: Int
        get() {
            if (json.has("RetrackersMode"))
                return json.getInt("RetrackersMode")
            return 0
        }
        set(value) {
            json.put("RetrackersMode", value)
        }

    var disableTCP: Boolean
        get() {
            if (json.has("DisableTCP"))
                return json.getBoolean("DisableTCP")
            return false
        }
        set(value) {
            json.put("DisableTCP", value)
        }

    var disableUTP: Boolean
        get() {
            if (json.has("DisableUTP"))
                return json.getBoolean("DisableUTP")
            return false
        }
        set(value) {
            json.put("DisableUTP", value)
        }

    var disableUPNP: Boolean
        get() {
            if (json.has("DisableUPNP"))
                return json.getBoolean("DisableUPNP")
            return false
        }
        set(value) {
            json.put("DisableUPNP", value)
        }

    var disableDHT: Boolean
        get() {
            if (json.has("DisableDHT"))
                return json.getBoolean("DisableDHT")
            return false
        }
        set(value) {
            json.put("DisableDHT", value)
        }

    var disableUpload: Boolean
        get() {
            if (json.has("DisableUpload"))
                return json.getBoolean("DisableUpload")
            return false
        }
        set(value) {
            json.put("DisableUpload", value)
        }

    var encryption: Int
        get() {
            if (json.has("Encryption"))
                return json.getInt("Encryption")
            return 0
        }
        set(value) {
            json.put("Encryption", value)
        }

    var downloadRateLimit: Int
        get() {
            if (json.has("DownloadRateLimit"))
                return json.getInt("DownloadRateLimit")
            return 0
        }
        set(value) {
            json.put("DownloadRateLimit", value)
        }

    var uploadRateLimit: Int
        get() {
            if (json.has("UploadRateLimit"))
                return json.getInt("UploadRateLimit")
            return 0
        }
        set(value) {
            json.put("UploadRateLimit", value)
        }

    var connectionsLimit: Int
        get() {
            if (json.has("ConnectionsLimit"))
                return json.getInt("ConnectionsLimit")
            return 0
        }
        set(value) {
            json.put("ConnectionsLimit", value)
        }
}