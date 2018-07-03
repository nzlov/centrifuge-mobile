package com.example.fz.centrifugoandroid;

import android.app.Activity;
import android.content.Context;
import android.widget.TextView;

import centrifuge.Client;
import centrifuge.ConnectContext;
import centrifuge.ConnectHandler;

public class AppConnectHandler implements ConnectHandler {
    protected MainActivity context;

    public AppConnectHandler(Context context) {
        this.context = (MainActivity) context;
    }

    @Override
    public void onConnect(Client client, ConnectContext connectContext) {
        context.runOnUiThread(new Runnable() {
            @Override
            public void run() {
                TextView tv = (TextView) context.findViewById(R.id.text);
                tv.setText("Connected");
            }
        });
    }
}
