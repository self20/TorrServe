package ru.yourok.torrserve.adapters

import android.app.Activity
import android.content.Context
import android.os.Handler
import android.os.Looper
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import android.widget.BaseAdapter
import android.widget.ImageView
import android.widget.TextView
import com.squareup.picasso.Picasso
import org.json.JSONArray
import ru.yourok.torrserve.R
import kotlin.concurrent.thread

class SearchVideoAdapter(val activity: Activity) : BaseAdapter() {
    var list: JSONArray = JSONArray()

    fun setResult(list: JSONArray) {
        this.list = list
        notifyDataSetChanged()
    }

    override fun getView(index: Int, convertView: View?, viewGroup: ViewGroup): View {
        val vi: View = convertView ?: (activity.getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater).inflate(R.layout.adapter_search_video_item, null)
        val info = list.getJSONObject(index)
        val image = vi.findViewById<ImageView>(R.id.imageViewPoster)
        if (info.has("poster_path")) {
            val poster = info.getString("poster_path")
            if (poster.isNotEmpty())
                thread {
                    val picass = Picasso.get().load(poster).placeholder(R.color.lighter_gray).fit().centerCrop()
                    Handler(Looper.getMainLooper()).post {
                        picass.into(image)
                    }
                }
            else
                image.setBackgroundResource(R.color.lighter_gray)
        }
        if (info.has("title"))
            vi.findViewById<TextView>(R.id.textViewMovieName).text = info.getString("title")
        else if (info.has("original_title"))
            vi.findViewById<TextView>(R.id.textViewMovieName).text = info.getString("original_title")
        else if (info.has("name"))
            vi.findViewById<TextView>(R.id.textViewMovieName).text = info.getString("name")
        else if (info.has("original_name"))
            vi.findViewById<TextView>(R.id.textViewMovieName).text = info.getString("original_name")

        var middle = ""
        if (info.has("release_date")) {
            if (info.getString("release_date").isNotEmpty())
                middle = info.getString("release_date").substring(0, 4)
        }
        if (info.has("first_air_date")) {
            if (info.getString("first_air_date").isNotEmpty())
                middle = info.getString("first_air_date").substring(0, 4)
        }

        if (info.has("number_of_seasons"))
            middle += " S${info.getInt("number_of_seasons")}"

        vi.findViewById<TextView>(R.id.textViewMovieDate).text = middle

        if (info.has("genres")) {
            val genres = mutableListOf<String>()
            for (i in 0 until info.getJSONArray("genres").length())
                genres.add(info.getJSONArray("genres").getJSONObject(i).getString("name"))
            vi.findViewById<TextView>(R.id.textViewMovieGenres).text = genres.joinToString(", ")
        }
        return vi
    }

    override fun getItem(p0: Int): Any? {
        try {
            return list.getJSONObject(p0)
        } catch (e: Exception) {
            return null
        }
    }

    override fun getItemId(p0: Int): Long {
        return p0.toLong()
    }

    override fun getCount(): Int {
        if (list.length() > 0) {
            if (list.getJSONObject(0).has("id"))
                return list.length()
        }
        return 0
    }
}