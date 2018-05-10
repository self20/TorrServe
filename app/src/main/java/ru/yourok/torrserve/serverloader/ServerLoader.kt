package ru.yourok.torrserve.serverloader

import android.os.Environment
import android.util.Log
import ru.yourok.torrserve.App
import java.io.File
import java.io.FileInputStream
import java.io.FileOutputStream
import java.io.InputStream

object ServerLoader {

    private val servPath = File(App.getContext().cacheDir, "serv")

    fun copy() {
        if (File("/sdcard/serv").exists()) {
            val input = FileInputStream("/sdcard/serv")

            copy(input, servPath)
            servPath.setExecutable(true)
        }
    }

    private fun copy(input: InputStream, out: File) {
        try {
            val out = FileOutputStream(out)
            try {
                val buf = ByteArray(1024)
                while (true) {
                    val len = input.read(buf)
                    if (len > 0)
                        out.write(buf, 0, len)
                    else
                        break
                }
            } finally {
                out.close()
            }
        } finally {
            input.close()
        }
    }

    fun run() {
        var data = File(Environment.getExternalStorageDirectory().path, "TorrServe")
        data.mkdir()
        val process = Process(servPath.path, data.path)
        process.onOutput {
            Log.i("GoLog", it)
        }

        process.onError {
            Log.i("GoLogErr", it)
        }
        process.exec()

    }
}