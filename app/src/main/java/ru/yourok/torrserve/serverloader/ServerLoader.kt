package ru.yourok.torrserve.serverloader

import android.os.Build
import android.os.Environment
import android.util.Log
import cz.msebera.android.httpclient.client.methods.HttpGet
import cz.msebera.android.httpclient.impl.client.HttpClients
import ru.yourok.torrserve.App
import java.io.File
import java.io.FileInputStream
import java.io.FileOutputStream
import java.io.InputStream

object ServerLoader {
    private val servPath = File(App.getContext().filesDir, "torrserver")
    private var process: Process? = null

    fun serverExists(): Boolean {
        return servPath.exists()
    }

    fun deleteServer(): Boolean {
        if (!serverExists())
            return true
        return servPath.delete()
    }

    fun checkLocal(): File? {
        val dir = Environment.getExternalStoragePublicDirectory(Environment.DIRECTORY_DOWNLOADS)
        val name = getArch()
        if (name.isEmpty())
            return null

        val file = File(dir, "TorrServer-android-$name")
        if (!file.exists())
            return null
        return file
    }

    fun copyLocal(): Boolean {
        val file = checkLocal()
        file?.let {
            servPath.delete()
            val input = FileInputStream(file)
            copy(input, servPath)
            servPath.setExecutable(true)
            return true
        }
        return false
    }

    fun download(): String {
        val name = getArch()
        if (name.isEmpty())
            return "error get arch"

        val url = "https://raw.githubusercontent.com/YouROK/TorrServe/master/TorrServer/dist/TorrServer-android-$name"

        val httpclient = HttpClients.custom().disableRedirectHandling().build()
        val httpreq = HttpGet(url)
        val response = httpclient.execute(httpreq)
        val status = response.statusLine?.statusCode ?: -1
        if (status == 200) {
            servPath.delete()
            val output = FileOutputStream(servPath)
            response.entity.writeTo(output)
            servPath.setExecutable(true)
            return ""
        } else {
            return response.statusLine?.reasonPhrase ?: "error load file"
        }
    }

    fun getArch(): String {
        val arch = if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.LOLLIPOP)
            Build.SUPPORTED_ABIS[0]
        else
            Build.CPU_ABI

        when (arch) {
            "arm64-v8a" -> return "arm64"
            "armeabi-v7a" -> return "arm7"
            "x86_64" -> return "amd64"
            "x86" -> return "386"
        }
        return ""
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
        if (!ServerLoader.serverExists())
            return
        if (process == null || !process!!.isRunning()) {
            val process = Process(servPath.path, "-d", servPath.parent)
            process.onOutput {
                Log.i("GoLog", it)
            }

            process.onError {
                Log.i("GoLogErr", it)
            }
            try {
                process.exec()
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }
    }

    fun stop() {
        process?.stop()
    }
}