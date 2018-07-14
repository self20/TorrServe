package ru.yourok.torrserve.serverhelper.server

import org.json.JSONObject

class Torrent(val json: JSONObject) {
    private val files: List<File>

    init {
        if (json.has("Files")) {
            val arr = json.getJSONArray("Files")
            val flist = mutableListOf<File>()
            for (i in 0 until arr.length())
                flist.add(File(arr.getJSONObject(i)))
            files = flist.toList()
        } else
            files = listOf<File>()
    }

    fun Name(): String {
        if (json.has("Name"))
            return json.getString("Name")
        return Hash()
    }

    fun Magnet(): String {
        if (json.has("Magnet"))
            return json.getString("Magnet")
        return ""
    }

    fun Hash(): String {
        if (json.has("Hash"))
            return json.getString("Hash")
        return ""
    }

    fun AddTime(): Long {
        if (json.has("AddTime"))
            return json.getLong("AddTime")
        return 0
    }

    fun Length(): Long {
        if (json.has("Length"))
            return json.getLong("Length")
        return 0
    }

    fun Status(): Int {
        if (json.has("Status"))
            return json.getInt("Status")
        return 0
    }

    fun Files() = files

}

class File(val json: JSONObject) {

    var Name: String
        get() {
            if (json.has("Name"))
                return json.getString("Name")
            return ""
        }
        set(value) {
            json.put("Name", value)
        }

    var Link: String
        get() {
            if (json.has("Link"))
                return json.getString("Link")
            return ""
        }
        set(value) {
            json.put("Link", value)
        }

    var Preload: String
        get() {
            if (json.has("Preload"))
                return json.getString("Preload")
            return ""
        }
        set(value) {
            json.put("Preload", value)
        }

    var Size: Long
        get() {
            if (json.has("Size"))
                return json.getLong("Size")
            return 0
        }
        set(value) {
            json.put("Size", value)
        }

    var Viewed: Boolean
        get() {
            if (json.has("Viewed"))
                return json.getBoolean("Viewed")
            return false
        }
        set(value) {
            json.put("Viewed", value)
        }
}