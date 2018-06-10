package ru.yourok.torrserve.activitys

import android.content.Intent
import android.net.Uri
import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import java.net.URLDecoder

class AddDSActivity : AppCompatActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        var torrentLink = ""

        if (intent.action != null && intent.action.equals(Intent.ACTION_VIEW)) {
            intent.data?.let {
                torrentLink = URLDecoder.decode(it.toString(), "UTF-8")
            }
        }

        ///Intent send
        if (intent.action != null && intent.action.equals(Intent.ACTION_SEND)) {
            if (intent.getStringExtra(Intent.EXTRA_TEXT) != null)
                torrentLink = intent.getStringExtra(Intent.EXTRA_TEXT)
            if (intent.extras.get(Intent.EXTRA_STREAM) != null)
                torrentLink = intent.extras.get(Intent.EXTRA_STREAM).toString()
        }

        if (torrentLink.isEmpty()) {
            finish()
            return
        }

        val vintent = Intent(this, ViewActivity::class.java)
        vintent.setData(Uri.parse(torrentLink))
        vintent.action = Intent.ACTION_VIEW
        vintent.putExtra("DontSave", true)
        startActivity(vintent)

        finish()
    }
}