package ru.yourok.torrserve.serverloader

import java.io.IOException
import java.nio.charset.Charset
import kotlin.concurrent.thread

class Process(vararg val args: String) {

    private var callbackOutput: ((line: String) -> Unit)? = null
    private var callbackError: ((line: String) -> Unit)? = null
    private var proc: java.lang.Process? = null

    fun onOutput(callback: ((line: String) -> Unit)?) {
        callbackOutput = callback
    }

    fun onError(callback: ((line: String) -> Unit)?) {
        callbackError = callback
    }

    fun exec() {
        try {
            proc = Runtime.getRuntime().exec(args)
        } catch (e: IOException) {
            throw RuntimeException(e)
        } catch (e: InterruptedException) {
            throw RuntimeException(e)
        }

        thread {
            proc?.waitFor()
            proc = null
        }

        thread {
            while (true) {
                proc?.let { p ->
                    val line = p.inputStream.bufferedReader(Charset.defaultCharset()).readLine()
                    callbackOutput?.let {
                        line?.let { ln ->
                            it.invoke(ln)
                        }
                    }
                }
                if (proc == null)
                    break
            }
        }

        thread {
            while (true) {
                proc?.let { p ->
                    val line = p.errorStream.bufferedReader(Charset.defaultCharset()).readLine()
                    callbackError?.let {
                        line?.let { ln ->
                            it.invoke(ln)
                        }
                    }
                }
                if (proc == null)
                    break
            }
        }
    }

    fun isRunning(): Boolean {
        proc?.let {
            return try {
                it.exitValue()
                false
            } catch (e: Exception) {
                true
            }
        }
        return false
    }

    fun stop() {
        proc?.destroy()
        proc = null
    }

}