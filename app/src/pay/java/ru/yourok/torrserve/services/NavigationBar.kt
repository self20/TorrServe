package ru.yourok.torrserve.services

import android.content.Intent
import android.support.v7.app.AppCompatActivity
import android.widget.ListView
import co.zsmb.materialdrawerkt.builders.drawer
import co.zsmb.materialdrawerkt.builders.footer
import co.zsmb.materialdrawerkt.draweritems.badgeable.primaryItem
import co.zsmb.materialdrawerkt.draweritems.badgeable.secondaryItem
import co.zsmb.materialdrawerkt.draweritems.divider
import com.mikepenz.materialdrawer.Drawer
import ru.yourok.torrserve.R
import ru.yourok.torrserve.activitys.AddActivity
import ru.yourok.torrserve.activitys.SettingsActivity
import ru.yourok.torrserve.adapters.TorrentListAdapter
import ru.yourok.torrserve.serverhelper.ServerApi
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
                footer {
                    secondaryItem(R.string.clear_cache) {
                        icon = R.drawable.ic_clean_cache_black_24dp
                        selectable = false
                        onClick { _ ->
                            ServerApi.cleanCache("")
                            false
                        }
                    }
                    secondaryItem(R.string.exit) {
                        icon = R.drawable.ic_cancel_black_24dp
                        selectable = false
                        onClick { _ ->
                            TorrService.stopAndExit()
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

                }
                headerPadding = true
                stickyFooterDivider = true
                stickyFooterShadow = false
                selectedItem = -1
            }
        }
    }
}
