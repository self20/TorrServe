package ru.yourok.torrserve.activitys

import android.content.Context
import android.graphics.Bitmap
import android.graphics.Point
import android.graphics.drawable.BitmapDrawable
import android.graphics.drawable.Drawable
import android.os.Build
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.support.v7.app.AppCompatActivity
import android.view.View
import com.squareup.picasso.Picasso
import kotlinx.android.synthetic.main.activity_video_review.*
import org.json.JSONArray
import org.json.JSONObject
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.SearchTorrentAdapter
import ru.yourok.torrserve.serverhelper.ServerRequests
import ru.yourok.torrserve.utils.Dialog
import java.lang.Exception
import kotlin.concurrent.thread

class VideoReviewActivity : AppCompatActivity() {
    companion object {
        var result: JSONArray? = null
    }

    lateinit var searchTorrentAdapter: SearchTorrentAdapter
    var info = JSONObject()
    lateinit var target: bgTarget

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_video_review)

        searchTorrentAdapter = SearchTorrentAdapter(this)
        listViewSearchTorrent.adapter = searchTorrentAdapter
        result?.let {
            searchTorrentAdapter.setResult(it)
        }

        target = bgTarget(this, backgroundContainer)

        if (intent.hasExtra("info")) {
            info = JSONObject(intent.getStringExtra("info"))

            if (info.has("poster_path")) {
                val poster = info.getString("poster_path")
                if (poster.isNotEmpty())
                    Picasso.get().load(poster).placeholder(R.color.lighter_gray).fit().centerCrop().into(imageViewPoster)
                else
                    imageViewPoster.setBackgroundResource(R.color.lighter_gray)
            }

            if (info.has("backdrop_path")) {
                val backdrop = info.getString("backdrop_path")
                if (backdrop.isNotEmpty()) {
                    val size = Point()
                    windowManager.defaultDisplay.getSize(size)
                    Picasso.get().load(backdrop).resize(size.x, size.y).centerCrop().into(target)
                }
            }

            var name = ""
            when {
                info.has("title") -> name = info.getString("title")
                info.has("original_title") -> name = info.getString("original_title")
                info.has("name") -> name = info.getString("name")
                info.has("original_name") -> name = info.getString("original_name")
            }

            var date = ""

            if (info.has("release_date")) {
                if (info.getString("release_date").isNotEmpty())
                    date = " (" + info.getString("release_date").substring(0, 4) + ")"
            }
            if (info.has("first_air_date")) {
                if (info.getString("first_air_date").isNotEmpty())
                    date = " (" + info.getString("first_air_date").substring(0, 4) + ")"
            }
            textViewVideoName.text = (name + date)

            if (info.has("overview"))
                textViewOverview.text = info.getString("overview")

            if (info.has("number_of_seasons")) {
                val seasons = info.getInt("number_of_seasons")
                if (seasons > 0) {
                    buttonSeasons.visibility = View.VISIBLE
                    buttonSeasons.setOnClickListener {
                        val sesList = mutableListOf<String>()
                        for (i in 1..seasons)
                            sesList.add(i.toString())

                        Dialog.showListDialog(this, getString(R.string.seasons), sesList, false) { selStr, selInt ->
                            search(name, selInt[0] + 1)
                        }
                    }
                }
            }
            if (result == null)
                search(name, 0)

        } else {
            finish()
            return
        }
    }

    private fun search(name: String, season: Int) {
        Handler(Looper.getMainLooper()).post {
            progressBarLoading.visibility = View.VISIBLE
            buttonSeasons.isEnabled = false
        }
        thread {
            try {
                var sesStr = ""
                if (info.has("number_of_seasons")) {
                    val seasons = info.getInt("number_of_seasons")
                    if (season > 0 && season < seasons) {
                        val sf = String.format("%02d", season)
                        sesStr = "S$sf|${sf}x"
                    }
                }

                result = ServerRequests.searchTorrent(name, listOf(sesStr))
                Handler(Looper.getMainLooper()).post {
                    result?.let {
                        searchTorrentAdapter.setResult(it)
                    }
                }
            } catch (e: Exception) {
            }
            Handler(Looper.getMainLooper()).post {
                progressBarLoading.visibility = View.GONE
                buttonSeasons.isEnabled = true
            }
        }
    }
}

class bgTarget(val context: Context, val view: View) : com.squareup.picasso.Target {
    override fun onPrepareLoad(placeHolderDrawable: Drawable?) {
    }

    override fun onBitmapFailed(e: Exception?, errorDrawable: Drawable?) {
        e?.printStackTrace()
    }

    override fun onBitmapLoaded(bitmap: Bitmap?, from: Picasso.LoadedFrom?) {
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.JELLY_BEAN) {
            view.background = BitmapDrawable(context.resources, bitmap)
        } else {
            view.setBackgroundDrawable(BitmapDrawable(context.resources, bitmap))
        }
    }
}