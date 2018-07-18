package ru.yourok.torrserve.activitys

import android.content.Intent
import android.os.Bundle
import android.os.Handler
import android.os.Looper
import android.support.v7.app.AppCompatActivity
import android.view.View
import android.widget.AdapterView
import android.widget.AdapterView.OnItemSelectedListener
import kotlinx.android.synthetic.main.activity_search.*
import org.json.JSONArray
import org.json.JSONObject
import ru.yourok.torrserve.R
import ru.yourok.torrserve.adapters.SearchTorrentAdapter
import ru.yourok.torrserve.adapters.SearchVideoAdapter
import ru.yourok.torrserve.serverhelper.ServerRequests
import ru.yourok.torrserve.utils.Dialog
import java.util.*
import kotlin.concurrent.thread


class SearchActivity : AppCompatActivity() {
    companion object {
        var config: JSONObject? = null
        var result: JSONArray? = null
    }

    val searchInfo: SerachInfo = SerachInfo()
    var currPage: Int = 1

    lateinit var searchVideoAdapter: SearchVideoAdapter
    lateinit var searchTorrentAdapter: SearchTorrentAdapter

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_search)
        searchVideoAdapter = SearchVideoAdapter(this)
        searchTorrentAdapter = SearchTorrentAdapter(this)
        result?.let {
            searchVideoAdapter.setResult(it)
            searchTorrentAdapter.setResult(it)
        }

        gridViewSearchVideo.adapter = searchVideoAdapter
        listViewSearchTorrent.adapter = searchTorrentAdapter
        thread {
            try {
                if (config == null)
                    config = ServerRequests.searchConfig()
                updateUI()
            } catch (e: Exception) {
                e.printStackTrace()
                finish()
            }
        }
        if (result == null)
            search()
        spinnerSearchType.setSelection(1)
        spinnerSearchType.onItemSelectedListener = object : OnItemSelectedListener {
            override fun onNothingSelected(p0: AdapterView<*>?) {
                searchInfo.SearchBy = 0
                updateUI()
            }

            override fun onItemSelected(p0: AdapterView<*>?, p1: View?, p2: Int, p3: Long) {
                searchInfo.SearchBy = p2
                updateUI()
            }
        }

        buttonMovies.isEnabled = false
        buttonMovies.setOnClickListener {
            searchInfo.SerachType = 0
            buttonMovies.isEnabled = false
            buttonShows.isEnabled = true
            buttonTorrent.isEnabled = true
            gridViewSearchVideo.visibility = View.VISIBLE
            listViewSearchTorrent.visibility = View.GONE
            searchInfo.Sort = ""
            searchInfo.Genres = ""
            currPage = 1
            updateUI()
        }

        buttonShows.setOnClickListener {
            searchInfo.SerachType = 1
            buttonMovies.isEnabled = true
            buttonShows.isEnabled = false
            buttonTorrent.isEnabled = true
            gridViewSearchVideo.visibility = View.VISIBLE
            listViewSearchTorrent.visibility = View.GONE
            searchInfo.Sort = ""
            searchInfo.Genres = ""
            currPage = 1
            updateUI()
        }

        buttonTorrent.setOnClickListener {
            searchInfo.SerachType = 2
            buttonMovies.isEnabled = true
            buttonShows.isEnabled = true
            buttonTorrent.isEnabled = false
            gridViewSearchVideo.visibility = View.GONE
            listViewSearchTorrent.visibility = View.VISIBLE
            updateUI()
        }

        gridViewSearchVideo.setOnItemClickListener { adapterView, view, i, l ->
            val intent = Intent(this, VideoReviewActivity::class.java)
            intent.action = Intent.ACTION_VIEW
            intent.flags = Intent.FLAG_ACTIVITY_NEW_TASK
            val js = (adapterView.adapter.getItem(i) as JSONObject?)
            js?.let {
                intent.putExtra("info", it.toString())
                startActivity(intent)
            }
        }

        //Year select
        val yearList = mutableListOf<String>()
        yearList.add("")
        for (i in Calendar.getInstance().get(Calendar.YEAR) downTo 1900)
            yearList.add(i.toString())
        buttonFTYear.setOnClickListener {
            Dialog.showListDialog(this, getString(R.string.year), yearList, false) { selStr, selInt ->
                searchInfo.Year = selStr[0]
            }
        }

        //Sort select
        buttonFTSort.setOnClickListener {
            val sortList = if (searchInfo.SerachType == 0)
                listOf("", "popularity.asc", "popularity.desc", "release_date.asc", "release_date.desc", "revenue.asc", "revenue.desc", "primary_release_date.asc", "primary_release_date.desc", "original_title.asc", "original_title.desc", "vote_average.asc", "vote_average.desc", "vote_count.asc", "vote_count.desc")
            else
                listOf("", "vote_average.desc", "vote_average.asc", "first_air_date.desc", "first_air_date.asc", "popularity.desc", "popularity.asc")
            Dialog.showListDialog(this, getString(R.string.sort), sortList, false) { selStr, selInt ->
                searchInfo.Sort = selStr[0]
            }
        }

        //Genres select
        buttonFTGenres.setOnClickListener {
            try {
                if (config == null)
                    return@setOnClickListener

                val genres = if (searchInfo.SerachType == 0)
                    config!!.getJSONObject("Genres").getJSONArray("MovieGenres")
                else
                    config!!.getJSONObject("Genres").getJSONArray("ShowGenres")

                val names = mutableListOf<String>()
                for (i in 0 until genres.length())
                    names.add(genres.getJSONObject(i).getString("name"))

                Dialog.showListDialog(this, getString(R.string.sort), names, true) { selStr, selInt ->
                    val genIds = mutableListOf<String>()
                    selInt.forEach {
                        val sel = genres.getJSONObject(it).getInt("id")
                        genIds.add(sel.toString())
                    }
                    searchInfo.Genres = genIds.joinToString(",")
                }
            } catch (e: Exception) {
                e.printStackTrace()
            }
        }

        buttonSearch.setOnClickListener { search() }

        buttonNext.setOnClickListener {
            currPage++
            search()
        }
        buttonPrev.setOnClickListener {
            if (currPage > 0) {
                currPage--
                search()
            }
        }

    }

    private fun search() {
        buttonSearch.isEnabled = false
        buttonNext.isEnabled = false
        buttonPrev.isEnabled = false
        progressBarLoading.visibility = View.VISIBLE
        thread {
            if (searchInfo.SerachType == 0) {//Movie
                if (searchInfo.SearchBy == 0) {
                    result = ServerRequests.searchMovies(listOf("type=search", "page=$currPage", "query=${editTextName.text}"))
                    Handler(Looper.getMainLooper()).post {
                        result?.let {
                            searchVideoAdapter.setResult(it)
                        }
                    }
                } else {
                    val params = mutableListOf("type=discover")

                    params.add("page=$currPage")

                    if (searchInfo.Year.isNotEmpty())
                        params.add("primary_release_year=${searchInfo.Year}")

                    if (searchInfo.Sort.isNotEmpty())
                        params.add("sort_by=${searchInfo.Sort}")

                    if (searchInfo.Genres.isNotEmpty())
                        params.add("with_genres=" + searchInfo.Genres)

                    result = ServerRequests.searchMovies(params)
                    Handler(Looper.getMainLooper()).post {
                        result?.let {
                            searchVideoAdapter.setResult(it)
                        }
                    }
                }
            } else if (searchInfo.SerachType == 1) {//Show
                if (searchInfo.SearchBy == 0) {
                    result = ServerRequests.searchShows(listOf("type=search", "page=$currPage", "query=${editTextName.text}"))
                    Handler(Looper.getMainLooper()).post {
                        result?.let {
                            searchVideoAdapter.setResult(it)
                        }
                    }
                } else {
                    val params = mutableListOf("type=discover")

                    params.add("primary_release_year=${searchInfo.Year}")
                    params.add("page=$currPage")

                    if (searchInfo.Sort.isNotEmpty())
                        params.add("sort_by=${searchInfo.Sort}")

                    if (searchInfo.Genres.isNotEmpty())
                        params.add("with_genres=" + searchInfo.Genres)

                    result = ServerRequests.searchShows(params)
                    Handler(Looper.getMainLooper()).post {
                        result?.let {
                            searchVideoAdapter.setResult(it)
                        }
                    }
                }
            } else if (searchInfo.SerachType == 2) {//Torrent
                val fts = editTextFilter.text.split(";")
                result = ServerRequests.searchTorrent(editTextName.text.toString(), fts.map { it.trim() })
                Handler(Looper.getMainLooper()).post {
                    result?.let {
                        searchTorrentAdapter.setResult(it)
                    }
                }
            }
            Handler(Looper.getMainLooper()).post {
                progressBarLoading.visibility = View.GONE
                buttonSearch.isEnabled = true
                buttonNext.isEnabled = true
                buttonPrev.isEnabled = true
                searchScroll.scrollTo(0, 0)
            }
        }
    }

    private fun updateUI() {
        Handler(Looper.getMainLooper()).post {
            if (searchInfo.SerachType == 2) {
                fTypeByName.visibility = View.VISIBLE
                fTypeTorrent.visibility = View.VISIBLE
                fTypeDiscover.visibility = View.GONE
            } else {
                fTypeTorrent.visibility = View.GONE
                if (searchInfo.SearchBy == 0) {
                    fTypeByName.visibility = View.VISIBLE
                    fTypeDiscover.visibility = View.GONE
                } else {
                    fTypeByName.visibility = View.GONE
                    fTypeDiscover.visibility = View.VISIBLE
                }
            }
        }
    }
}


class SerachInfo {
    var SerachType: Int = 0
    var SearchBy: Int = 1
    var Name: String = ""
    var Year: String = Calendar.getInstance().get(Calendar.YEAR).toString()
    var Sort: String = ""
    var Genres: String = ""
}