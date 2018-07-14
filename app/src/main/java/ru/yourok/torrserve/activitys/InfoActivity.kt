package ru.yourok.torrserve.activitys

import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import kotlinx.android.synthetic.main.activity_info.*
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerApi
import kotlin.concurrent.thread

class InfoActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_info)

        val hashs = intent.getStringArrayExtra("Hashs")
        if (hashs == null || hashs!!.isEmpty()) {
            finish()
            return
        }

        thread {
            hashs?.let {
                while (!this@InfoActivity.isFinishing) {
                    var msg = ""
                    it.forEach {
                        val info = ServerApi.stat(it)

                        val torr = ServerApi.get(it)
                        msg += "Name: ${torr?.Name()}\n"
                        msg += "PreloadLength: ${torr?.Length()}\n"
                        info?.let { info ->
                            msg += "BytesWritten: ${info.BytesWritten()}\n"
                            msg += "BytesWrittenData: ${info.BytesWrittenData()}\n"
                            msg += "BytesRead: ${info.BytesRead()}\n"
                            msg += "BytesReadData: ${info.BytesReadData()}\n"
                            msg += "BytesReadUsefulData: ${info.BytesReadUsefulData()}\n"
                            msg += "ChunksWritten: ${info.ChunksWritten()}\n"
                            msg += "ChunksRead: ${info.ChunksRead()}\n"
                            msg += "ChunksReadUseful: ${info.ChunksReadUseful()}\n"
                            msg += "ChunksReadWasted: ${info.ChunksReadWasted()}\n"
                            msg += "PiecesDirtiedGood: ${info.PiecesDirtiedGood()}\n"
                            msg += "PiecesDirtiedBad: ${info.PiecesDirtiedBad()}\n"
                            msg += "TotalPeers: ${info.TotalPeers()}\n"
                            msg += "PendingPeers: ${info.PendingPeers()}\n"
                            msg += "ActivePeers: ${info.ActivePeers()}\n"
                            msg += "ConnectedSeeders: ${info.ConnectedSeeders()}\n"
                            msg += "HalfOpenPeers: ${info.HalfOpenPeers()}\n"
                            msg += "\n"
                        }
                    }
                    runOnUiThread { textViewInfo.text = msg }
                    Thread.sleep(1000)
                }
            }
        }
    }
}