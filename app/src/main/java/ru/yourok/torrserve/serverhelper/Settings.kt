package ru.yourok.torrserve.serverhelper

import android.preference.PreferenceManager
import ru.yourok.torrserve.App

/**
 * Created by yourok on 20.02.18.
 */


object Preferences {

    fun getPlayer(): String {
        return get("Player", "0") as String
    }

    fun setPlayer(pkg: String) {
        set("Player", pkg)
    }

    fun isShowPreloadWnd(): Boolean {
        return get("ShowPreload", true) as Boolean
    }

    fun setShowPreloadWnd(v: Boolean) {
        set("ShowPreload", v)
    }

    fun isAutoStart(): Boolean {
        return get("AutoStart", false) as Boolean
    }

    fun setAutoStart(v: Boolean) {
        set("AutoStart", v)
    }

    fun getServerAddress(): String {
        return get("ServerAddress", "http://localhost:8090") as String
    }

    fun setServerAddress(addr: String) {
        set("ServerAddress", addr)
    }

    fun getLastViewDonate(): Long {
        return get("LastViewDonate", 0L) as Long
    }

    fun setLastViewDonate(l: Long) {
        set("LastViewDonate", l)
    }

    fun getSaveHosts(): List<String> {
        val prefs = PreferenceManager.getDefaultSharedPreferences(App.getContext())
        val ret = prefs.getStringSet("AutoCompleteHost", mutableSetOf())
        if (ret.isEmpty())
            ret.add("http://localhost:8090")
        return ret.toList()
    }

    fun addSaveHost(v: String) {
        val list = getSaveHosts()
        val set = list.toMutableSet()
        set.add(v)
        val prefs = PreferenceManager.getDefaultSharedPreferences(App.getContext())
        prefs.edit().putStringSet("AutoCompleteHost", set).apply()
    }

    private fun get(name: String, def: Any): Any? {
        val prefs = PreferenceManager.getDefaultSharedPreferences(App.getContext())
        if (prefs.all.containsKey(name))
            return prefs.all[name]
        return def
    }

    private fun set(name: String, value: Any?) {
        val prefs = PreferenceManager.getDefaultSharedPreferences(App.getContext())
        when (value) {
            is String -> prefs.edit().putString(name, value).apply()
            is Boolean -> prefs.edit().putBoolean(name, value).apply()
            is Float -> prefs.edit().putFloat(name, value).apply()
            is Int -> prefs.edit().putInt(name, value).apply()
            is Long -> prefs.edit().putLong(name, value).apply()
            is MutableSet<*>? -> prefs.edit().putStringSet(name, value as MutableSet<String>?).apply()
            else -> prefs.edit().putString(name, value.toString()).apply()
        }
    }
}