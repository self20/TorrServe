package ru.yourok.torrserve.navigationBar

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.widget.ListView
import co.zsmb.materialdrawerkt.builders.drawer
import co.zsmb.materialdrawerkt.builders.footer
import co.zsmb.materialdrawerkt.draweritems.badgeable.primaryItem
import co.zsmb.materialdrawerkt.draweritems.badgeable.secondaryItem
import co.zsmb.materialdrawerkt.draweritems.divider
import com.mikepenz.materialdrawer.Drawer
import ru.yourok.torrserve.Donate
import ru.yourok.torrserve.R
import ru.yourok.torrserve.activitys.AddActivity
import ru.yourok.torrserve.activitys.SearchActivity
import ru.yourok.torrserve.activitys.SettingsActivity
import ru.yourok.torrserve.adapters.TorrentListAdapter
import ru.yourok.torrserve.serverhelper.ServerApi
import ru.yourok.torrserve.services.TorrService
import kotlin.concurrent.thread

object NavigationBar {
    fun setup(activity: AppCompatActivity): Drawer {
        with(activity) {
            return drawer {
                headerViewRes = R.layout.header

                primaryItem(R.string.add) {
                    icon = R.drawable.ic_add_black_24dp
                    selectable = false
                    onClick { _ ->
                        startActivity(Intent(activity, AddActivity::class.java))
                        false
                    }
                }

                divider {}
                primaryItem(R.string.remove_all) {
                    icon = R.drawable.ic_trash_black_24dp
                    selectable = false
                    onClick { _ ->
                        thread {
                            val torrList = ServerApi.list()
                            torrList.forEach {
                                ServerApi.rem(it.Hash)
                            }
                            runOnUiThread {
                                val ada = activity.findViewById<ListView>(R.id.listViewTorrent).adapter
                                (ada as TorrentListAdapter).updateList()
                            }
                        }
                        false
                    }
                }

                divider {}
                primaryItem(R.string.playlist) {
                    icon = R.drawable.ic_list_black_24dp
                    selectable = false
                    onClick { _ ->
                        thread {
                            if (ServerApi.list().isNotEmpty())
                                ServerApi.openPlayList()
                        }
                        false
                    }
                }

                divider {}
                primaryItem(R.string.donate) {
                    icon = R.drawable.ic_donate_black
                    selectable = false
                    onClick { _ ->
                        Donate.donateDialog(activity)
                        false
                    }
                }

                divider {}
                footer {
                    secondaryItem(R.string.search) {
                        icon = R.drawable.ic_search_black_24dp
                        selectable = false
                        onClick { _ ->
                            startActivity(Intent(activity, SearchActivity::class.java))
                            false
                        }
                    }
                    secondaryItem(R.string.clear_cache) {
                        icon = R.drawable.ic_clean_cache_black_24dp
                        selectable = false
                        onClick { _ ->
                            ServerApi.cleanCache("")
                            false
                        }
                    }
                    secondaryItem(R.string.settings) {
                        icon = R.drawable.ic_settings_black_24dp
                        selectable = false
                        onClick { _ ->
                            startActivity(Intent(activity, SettingsActivity::class.java))
                            false
                        }
                    }
                    secondaryItem(R.string.exit) {
                        icon = R.drawable.ic_cancel_black_24dp
                        selectable = false
                        onClick { _ ->
                            TorrService.exit()
                            false
                        }
                    }
                }
                headerPadding = true
                stickyFooterDivider = true
                stickyFooterShadow = false
                selectedItem = -1
            }
        }
    }
}
