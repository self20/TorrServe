package ru.yourok.torrserve.activitys

import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.webkit.WebView
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.Preferences
import ru.yourok.torrserve.serverhelper.ServerRequest

class SearchActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_search)
        val addr = ServerRequest.joinUrl(Preferences.getServerAddress(), "/search")
        val webview = findViewById<WebView>(R.id.webViewSearch)
        webview.getSettings().setJavaScriptEnabled(true);
        webview.loadUrl(addr)
    }
}
