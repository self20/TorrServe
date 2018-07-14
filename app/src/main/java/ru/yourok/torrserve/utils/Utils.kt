package ru.yourok.torrserve.utils

import android.os.Environment
import ru.yourok.torrserve.App
import java.io.File

/**
 * Created by yourok on 23.02.18.
 */
object Utils {
    fun byteFmt(bytes: Double): String {
        if (bytes < 1024)
            return bytes.toString() + " B"
        val exp = (Math.log(bytes) / Math.log(1024.0)).toInt()
        val pre = "KMGTPE"[exp - 1].toString()
        return "%.1f %sB".format(bytes / Math.pow(1024.0, exp.toDouble()), pre)
    }

    fun byteFmt(bytes: Float): String {
        return byteFmt(bytes.toDouble())
    }

    fun byteFmt(bytes: Long): String {
        return byteFmt(bytes.toDouble())
    }

    fun byteFmt(bytes: Int): String {
        return byteFmt(bytes.toDouble())
    }

    fun getAppPath(): String {
        return App.getContext().filesDir.path

//        val state = Environment.getExternalStorageState()
//        var filesDir: File?
//        if (Environment.MEDIA_MOUNTED.equals(state))
//            filesDir = App.getContext().getExternalFilesDir(null)
//        else
//            filesDir = App.getContext().getFilesDir()
//
//        if (filesDir == null)
//            filesDir = App.getContext().cacheDir
//
//        if (!filesDir.exists())
//            filesDir.mkdirs()
//        return filesDir.path
    }
}