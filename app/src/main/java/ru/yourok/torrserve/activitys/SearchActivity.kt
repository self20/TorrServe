package ru.yourok.torrserve.activitys

import android.os.Bundle
import android.support.v7.app.AppCompatActivity
import android.webkit.WebResourceRequest
import android.webkit.WebView
import android.webkit.WebViewClient
import ru.yourok.torrserve.R
import ru.yourok.torrserve.serverhelper.ServerRequests


class SearchActivity : AppCompatActivity() {

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_search)
        val addr = ServerRequests.getHostUrl("/search")
        val webview = findViewById<WebView>(R.id.webViewSearch)
        webview.getSettings().setJavaScriptEnabled(true)
        webview.webViewClient = object : WebViewClient() {
            override fun shouldOverrideUrlLoading(view: WebView, req: WebResourceRequest): Boolean {
                view.loadUrl(req.toString())
                return false
            }
        }
        webview.loadUrl(addr)
    }
}
