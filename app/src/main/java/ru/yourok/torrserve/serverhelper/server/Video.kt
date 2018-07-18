package ru.yourok.torrserve.serverhelper.server

import org.json.JSONArray
import org.json.JSONObject

class Video(val json: JSONObject) {
    var ID: Int
        get() {
            if (json.has("ID"))
                return json.getInt("ID")
            return 0
        }
        set(value) {
            json.put("ID", value)
        }

    var Name: String
        get() {
            if (json.has("Name"))
                return json.getString("Name")
            return ""
        }
        set(value) {
            json.put("Name", value)
        }

    var OriginalName: String
        get() {
            if (json.has("OriginalName"))
                return json.getString("OriginalName")
            return ""
        }
        set(value) {
            json.put("OriginalName", value)
        }

    var Overview: String
        get() {
            if (json.has("Overview"))
                return json.getString("Overview")
            return ""
        }
        set(value) {
            json.put("Overview", value)
        }

    var Genres: List<Genre>
        get() {
            if (json.has("Genres")) {
                val arr = json.getJSONArray("Genres")
                val ret = mutableListOf<Genre>()
                for (i in 0 until arr.length()) {
                    ret.add(Genre(arr.getJSONObject(i)))
                }
                return ret
            }
            return mutableListOf()
        }
        set(value) {
            val arr = JSONArray()
            value.forEach {
                arr.put(it.json)
            }
            json.put("Genres", arr)
        }
    var Year: String
        get() {
            if (json.has("Year"))
                return json.getString("Year")
            return ""
        }
        set(value) {
            json.put("Year", value)
        }

    var Tagline: String
        get() {
            if (json.has("Tagline"))
                return json.getString("Tagline")
            return ""
        }
        set(value) {
            json.put("Tagline", value)
        }

    var Poster: String
        get() {
            if (json.has("Poster"))
                return json.getString("Poster")
            return ""
        }
        set(value) {
            json.put("Poster", value)
        }

    var Backdrop: String
        get() {
            if (json.has("Backdrop"))
                return json.getString("Backdrop")
            return ""
        }
        set(value) {
            json.put("Backdrop", value)
        }

    var AllArts: List<String>
        get() {
            if (json.has("AllArts")) {
                val arr = json.getJSONArray("AllArts")

            }
            return listOf()
        }
        set(value) {
            json.put("AllArts", value)
        }

    var Seasons: Int
        get() {
            if (json.has("Seasons"))
                return json.getInt("Seasons")
            return 0
        }
        set(value) {
            json.put("Seasons", value)
        }

    var Episodes: Int
        get() {
            if (json.has("Episodes"))
                return json.getInt("Episodes")
            return 0
        }
        set(value) {
            json.put("Episodes", value)
        }
}

class Genre(val json: JSONObject) {
    var ID: Int
        get() {
            if (json.has("ID"))
                return json.getInt("ID")
            return 0
        }
        set(value) {
            json.put("ID", value)
        }
    var Name: String
        get() {
            if (json.has("Name"))
                return json.getString("Name")
            return ""
        }
        set(value) {
            json.put("Name", value)
        }
}