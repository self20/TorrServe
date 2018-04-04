package ru.yourok.torrserve.serverhelper

import org.json.JSONObject

/**
 * Created by yourok on 11.03.18.
 */

class Torrent {
    var Name: String = ""
    var Magnet: String = ""
    var Hash: String = ""
    var Length: Long = 0
    var Files: MutableList<File> = mutableListOf()
}

class File {
    var Name: String = ""
    var Link: String = ""
    var Size: Long = 0
    var Viewed: Boolean = false
}

class Info {
    var BytesWritten: Long = 0
    var BytesWrittenData: Long = 0

    var BytesRead: Long = 0
    var BytesReadData: Long = 0
    var BytesReadUsefulData: Long = 0

    var ChunksWritten: Long = 0

    var ChunksRead: Long = 0
    var ChunksReadUseful: Long = 0
    var ChunksReadUnwanted: Long = 0

    var PiecesDirtiedGood: Long = 0
    var PiecesDirtiedBad: Long = 0

    var TotalPeers: Int = 0
    var PendingPeers: Int = 0
    var ActivePeers: Int = 0
    var ConnectedSeeders: Int = 0
    var HalfOpenPeers: Int = 0
}

fun js2Torrent(jsStr: String): Torrent {
    val js = JSONObject(jsStr)
    return js2Torrent(js)
}

fun js2Info(jsStr: String): Info {
    val js = JSONObject(jsStr)
    return js2Info(js)
}

fun js2Info(js: JSONObject): Info {
    val info = Info()
    info.BytesWritten = js.getLong("BytesWritten")
    info.BytesWrittenData = js.getLong("BytesWrittenData")

    info.BytesRead = js.getLong("BytesRead")
    info.BytesReadData = js.getLong("BytesReadData")
    info.BytesReadUsefulData = js.getLong("BytesReadUsefulData")

    info.ChunksWritten = js.getLong("ChunksWritten")

    info.ChunksRead = js.getLong("ChunksRead")
    info.ChunksReadUseful = js.getLong("ChunksReadUseful")
    info.ChunksReadUnwanted = js.getLong("ChunksReadUnwanted")

    info.PiecesDirtiedGood = js.getLong("PiecesDirtiedGood")
    info.PiecesDirtiedBad = js.getLong("PiecesDirtiedBad")

    info.TotalPeers = js.getInt("TotalPeers")
    info.PendingPeers = js.getInt("PendingPeers")
    info.ActivePeers = js.getInt("ActivePeers")
    info.ConnectedSeeders = js.getInt("ConnectedSeeders")
    info.HalfOpenPeers = js.getInt("HalfOpenPeers")
    return info
}

fun js2Torrent(js: JSONObject): Torrent {
    val tor = Torrent()
    tor.Name = js.getString("Name")
    tor.Magnet = js.getString("Magnet")
    tor.Hash = js.getString("Hash")
    if (js.has("Length"))
        tor.Length = js.getLong("Length")
    if (js.has("Files")) {
        var jsArr = js.getJSONArray("Files")
        for (i in 0 until jsArr.length()) {
            val file = File()
            val jsFile = jsArr.getJSONObject(i)
            file.Name = jsFile.getString("Name")
            file.Link = jsFile.getString("Link")
            file.Size = jsFile.getLong("Size")
            file.Viewed = jsFile.getBoolean("Viewed")
            tor.Files.add(file)
        }
    }
    return tor
}