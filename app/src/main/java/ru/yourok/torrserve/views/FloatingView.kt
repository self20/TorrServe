package ru.yourok.torrserve.views

import android.app.Service
import android.content.Context
import android.content.res.Resources
import android.graphics.PixelFormat
import android.os.Build
import android.view.*
import android.widget.ImageView
import ru.yourok.torrserve.App
import ru.yourok.torrserve.R


class FloatingView {
    private var windowManager: WindowManager? = null
    private var view: View? = null

    fun create(): View? {
        view = (App.getContext().getSystemService(Context.LAYOUT_INFLATER_SERVICE) as LayoutInflater).inflate(R.layout.activity_floating, null)
        if (view == null)
            return null

        view!!.visibility = View.GONE
        windowManager = App.getContext().getSystemService(Service.WINDOW_SERVICE) as WindowManager

        val metrics = Resources.getSystem().getDisplayMetrics()
        val width = (210 * (metrics.densityDpi / 160f)).toInt()

        val myParams = if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O)
            WindowManager.LayoutParams(
                    width,
                    WindowManager.LayoutParams.WRAP_CONTENT,
                    WindowManager.LayoutParams.TYPE_PHONE,
                    WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE or
                            WindowManager.LayoutParams.FLAG_LAYOUT_IN_SCREEN,
                    PixelFormat.TRANSLUCENT)
        else
            WindowManager.LayoutParams(
                    width,
                    WindowManager.LayoutParams.WRAP_CONTENT,
                    WindowManager.LayoutParams.TYPE_APPLICATION_OVERLAY,
                    WindowManager.LayoutParams.FLAG_NOT_FOCUSABLE or
                            WindowManager.LayoutParams.FLAG_LAYOUT_IN_SCREEN,
                    PixelFormat.TRANSLUCENT)

        myParams.gravity = Gravity.BOTTOM or Gravity.RIGHT
        myParams.x = 0
        myParams.y = 100

        windowManager!!.addView(view, myParams)
        view!!.alpha = 0.25F

        view!!.setOnHoverListener(object : View.OnHoverListener {
            override fun onHover(v: View, event: MotionEvent): Boolean {
                when (event.action) {
                    MotionEvent.ACTION_HOVER_ENTER -> {
                        view!!.alpha = 1.0F
                    }
                    MotionEvent.ACTION_HOVER_EXIT -> {
                        view!!.alpha = 0.25F
                    }
                }
                return false
            }
        })

        view!!.setOnTouchListener(object : View.OnTouchListener {
            private var initialX: Int = 0
            private var initialY: Int = 0
            private var initialTouchX: Float = 0.0F
            private var initialTouchY: Float = 0.0F

            override fun onTouch(v: View, event: MotionEvent): Boolean {
                when (event.action) {
                    MotionEvent.ACTION_DOWN -> {
                        initialX = myParams.x
                        initialY = myParams.y
                        initialTouchX = event.rawX
                        initialTouchY = event.rawY
                        view?.alpha = 1.0F
                    }
                    MotionEvent.ACTION_UP -> {
                        view?.alpha = 0.25F
                    }
                    MotionEvent.ACTION_MOVE -> {
                        myParams.x = initialX - (event.rawX - initialTouchX).toInt()
                        myParams.y = initialY - (event.rawY - initialTouchY).toInt()
                        windowManager!!.updateViewLayout(v, myParams)
                        view?.alpha = 1.0F
                    }
                }
                return false
            }
        })
        return view
    }

    fun getView(): View? {
        return view
    }

    fun onCancel(l: View.OnClickListener) {
        val btn = view?.findViewById<ImageView>(R.id.buttonCancel)
        btn?.setOnClickListener(l)
    }

    fun remove() {
        view?.let {
            windowManager?.removeView(it)
        }
        view = null
    }

}