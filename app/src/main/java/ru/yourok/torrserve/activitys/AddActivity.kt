package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.Toast
import kotlinx.android.synthetic.main.activity_add.*
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.services.TorrService
import kotlin.concurrent.thread

class AddActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_add)
        wait(getString(R.string.starting_server)) {
            val run = TorrService.waitServer()
            runOnUiThread {
                if (run)
                    textViewStatus.setText("")
                else
                    textViewStatus.setText(R.string.error_server_start)
            }
        }
        ///Intent open
        val magnet = intent.data?.toString() ?: ""
        editTextTorrLink.setText(magnet)

        ///Intent send
        if (intent.action != null && intent.action.equals(Intent.ACTION_SEND)) {
            if (intent.getStringExtra(Intent.EXTRA_TEXT) != null)
                editTextTorrLink.setText(intent.getStringExtra(Intent.EXTRA_TEXT))
            if (intent.extras.get(Intent.EXTRA_STREAM) != null)
                editTextTorrLink.setText(intent.extras.get(Intent.EXTRA_STREAM).toString())
        }

        buttonAdd.setOnClickListener {
            wait("") {
                try {
                    ServerApi.add(editTextTorrLink.text.toString(), true)
                    finish()
                } catch (e: Exception) {
                    val msg = e.message ?: getString(R.string.error_add_torrent)
                    runOnUiThread {
                        Toast.makeText(this, msg, Toast.LENGTH_SHORT).show()
                    }
                }
            }
        }

        buttonCancel.setOnClickListener {
            finish()
        }
    }

    fun wait(message: String, fn: () -> Unit) {
        progressBar.visibility = View.VISIBLE
        textViewStatus.setText(message)
        thread {
            fn()
            runOnUiThread {
                progressBar.visibility = View.GONE
            }
        }
    }
}
