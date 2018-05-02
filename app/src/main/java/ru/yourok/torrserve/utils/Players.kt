package ru.yourok.torrserve.utils

import android.content.Intent
import android.net.Uri
import android.provider.MediaStore
import ru.yourok.torrserve.App

class Player(var Name: String, var Package: String) {
    override fun toString(): String {
        return Name
    }
}

object Players {
    fun getList(): MutableList<Player> {
        val list = mutableListOf<Player>()
        list.addAll(getList("video/*"))
        list.addAll(getList("audio/*"))
        return list.distinctBy { it.Package }.toMutableList()
    }

    private fun getList(mime: String): List<Player> {
        val intent = Intent(Intent.ACTION_VIEW)
        val uri = Uri.withAppendedPath(MediaStore.Audio.Media.INTERNAL_CONTENT_URI, "1")
        intent.data = uri
        intent.type = mime
        val apps = App.getContext().getPackageManager().queryIntentActivities(intent, 0)
        val list = mutableListOf<Player>()
        for (a in apps) {
            var name = a.loadLabel(App.getContext().packageManager)?.toString() ?: a.activityInfo.packageName
            list.add(Player(name, a.activityInfo.packageName))
        }
        return list
    }
}