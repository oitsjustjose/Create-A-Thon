package com.brogrammers.shopkinmanager;

/**
 * Many thanks to Jamie Avins for the Open Source Android App, FastBook
 * Repository Link: https://github.com/extjs/fastbook
 * <p>
 * This is a custom implementation of the same application, and functions
 * very similarly. Few things had to be changed, but due to API differences,
 * several refactors and port changes were required on our own part
 * <p>
 * --Jose Stovall / oitsjustjose [GitHub]
 */

import android.app.Activity;
import android.os.Bundle;
import android.view.KeyEvent;
import android.view.View;
import android.webkit.GeolocationPermissions;
import android.webkit.WebChromeClient;
import android.webkit.WebSettings;
import android.webkit.WebView;
import android.webkit.WebViewClient;

public class MainActivity extends Activity
{
    private WebView webView;

    class NimbleKitClient extends WebViewClient
    {
        @Override
        public boolean shouldOverrideUrlLoading(WebView view, String url)
        {
            if (url.startsWith("http:") || url.startsWith("https:"))
                view.loadUrl(url);
            return true;
        }
    }


    @Override
    protected void onCreate(Bundle savedInstanceState)
    {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        this.webView = (WebView) findViewById(R.id.webView1);
        webView.getSettings().setJavaScriptEnabled(true);
        webView.getSettings().setPluginState(WebSettings.PluginState.ON);
        webView.getSettings().setMixedContentMode(0);
        webView.setWebViewClient(new NimbleKitClient());

        webView.setWebChromeClient(new WebChromeClient()
        {
            public void onGeolocationPermissionsShowPrompt(String origin, GeolocationPermissions.Callback callback)
            {
                callback.invoke(origin, true, false);
            }
        });
        webView.getSettings().setDomStorageEnabled(true);
        webView.getSettings().setDatabaseEnabled(true);
        webView.setScrollBarStyle(View.SCROLLBARS_INSIDE_OVERLAY);
        webView.requestFocus(View.FOCUS_DOWN);
        if (savedInstanceState != null)
        {
            ((WebView)findViewById(R.id.webView1)).restoreState(savedInstanceState);
        }
        else
        {
            webView.loadUrl("http://oitsjustjose.github.io/Brogrammers/");
        }

    }

    @Override
    protected void onSaveInstanceState(Bundle outState )
    {
        super.onSaveInstanceState(outState);
        webView.saveState(outState);
    }

    @Override
    protected void onRestoreInstanceState(Bundle savedInstanceState)
    {
        super.onRestoreInstanceState(savedInstanceState);
        webView.restoreState(savedInstanceState);
    }

    @Override
    public boolean onKeyDown(int keyCode, KeyEvent event)
    {
        if (event.getAction() == KeyEvent.ACTION_DOWN)
        {
            switch (keyCode)
            {
                case KeyEvent.KEYCODE_BACK:
                    if (webView.canGoBack())
                        webView.goBack();
                    else
                        finish();
                    return true;
            }

        }
        return super.onKeyDown(keyCode, event);
    }
}