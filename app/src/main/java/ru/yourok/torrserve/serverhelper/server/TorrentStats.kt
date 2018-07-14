package ru.yourok.torrserve.serverhelper.server

import org.json.JSONObject

const val TorrentAdded = 0
const val TorrentGettingInfo = 1
const val TorrentPreload = 2
const val TorrentWorking = 3
const val TorrentClosed = 4

class TorrentStats(val json: JSONObject) {

    private val filesStats: List<TorrentFileStat>

    init {
        if (json.has("FileStats")) {
            val arr = json.getJSONArray("FileStats")
            val flist = mutableListOf<TorrentFileStat>()
            for (i in 0 until arr.length())
                flist.add(TorrentFileStat(arr.getJSONObject(i)))
            filesStats = flist.toList()
        } else
            filesStats = listOf<TorrentFileStat>()
    }

    fun Name(): String {
        if (json.has("Name"))
            return json.getString("Name")
        return ""
    }

    fun Hash(): String {
        if (json.has("Hash"))
            return json.getString("Hash")
        return ""
    }

    fun TorrentStatus(): Int {
        if (json.has("TorrentStatus"))
            return json.getInt("TorrentStatus")
        return 0
    }

    fun TorrentStatusString(): String {
        if (json.has("TorrentStatusString"))
            return json.getString("TorrentStatusString")
        return ""
    }

    fun LoadedSize(): Long {
        if (json.has("LoadedSize"))
            return json.getLong("LoadedSize")
        return 0
    }

    fun TorrentSize(): Long {
        if (json.has("TorrentSize"))
            return json.getLong("TorrentSize")
        return 0
    }

    fun PreloadedBytes(): Long {
        if (json.has("PreloadedBytes"))
            return json.getLong("PreloadedBytes")
        return 0
    }

    fun PreloadSize(): Long {
        if (json.has("PreloadSize"))
            return json.getLong("PreloadSize")
        return 0
    }

    fun DownloadSpeed(): Double {
        if (json.has("DownloadSpeed"))
            return json.getDouble("DownloadSpeed")
        return 0.0
    }

    fun UploadSpeed(): Double {
        if (json.has("UploadSpeed"))
            return json.getDouble("UploadSpeed")
        return 0.0
    }

    fun TotalPeers(): Int {
        if (json.has("TotalPeers"))
            return json.getInt("TotalPeers")
        return 0
    }

    fun PendingPeers(): Int {
        if (json.has("PendingPeers"))
            return json.getInt("PendingPeers")
        return 0
    }

    fun ActivePeers(): Int {
        if (json.has("ActivePeers"))
            return json.getInt("ActivePeers")
        return 0
    }

    fun ConnectedSeeders(): Int {
        if (json.has("ConnectedSeeders"))
            return json.getInt("ConnectedSeeders")
        return 0
    }

    fun HalfOpenPeers(): Int {
        if (json.has("HalfOpenPeers"))
            return json.getInt("HalfOpenPeers")
        return 0
    }

    fun BytesWritten(): Long {
        if (json.has("BytesWritten"))
            return json.getLong("BytesWritten")
        return 0
    }

    fun BytesWrittenData(): Long {
        if (json.has("BytesWrittenData"))
            return json.getLong("BytesWrittenData")
        return 0
    }

    fun BytesRead(): Long {
        if (json.has("BytesRead"))
            return json.getLong("BytesRead")
        return 0
    }

    fun BytesReadData(): Long {
        if (json.has("BytesReadData"))
            return json.getLong("BytesReadData")
        return 0
    }

    fun BytesReadUsefulData(): Long {
        if (json.has("BytesReadUsefulData"))
            return json.getLong("BytesReadUsefulData")
        return 0
    }

    fun ChunksWritten(): Long {
        if (json.has("ChunksWritten"))
            return json.getLong("ChunksWritten")
        return 0
    }

    fun ChunksRead(): Long {
        if (json.has("ChunksRead"))
            return json.getLong("ChunksRead")
        return 0
    }

    fun ChunksReadUseful(): Long {
        if (json.has("ChunksReadUseful"))
            return json.getLong("ChunksReadUseful")
        return 0
    }

    fun ChunksReadWasted(): Long {
        if (json.has("ChunksReadWasted"))
            return json.getLong("ChunksReadWasted")
        return 0
    }

    fun PiecesDirtiedGood(): Long {
        if (json.has("PiecesDirtiedGood"))
            return json.getLong("PiecesDirtiedGood")
        return 0
    }

    fun PiecesDirtiedBad(): Long {
        if (json.has("PiecesDirtiedBad"))
            return json.getLong("PiecesDirtiedBad")
        return 0
    }

    fun FileStats() = filesStats
}

class TorrentFileStat(val json: JSONObject) {
    fun Id(): Int {
        if (json.has("Id"))
            return json.getInt("Id")
        return 0
    }

    fun Path(): String {
        if (json.has("Path"))
            return json.getString("Path")
        return ""
    }

    fun Length(): Long {
        if (json.has("Length"))
            return json.getLong("Length")
        return 0
    }
}