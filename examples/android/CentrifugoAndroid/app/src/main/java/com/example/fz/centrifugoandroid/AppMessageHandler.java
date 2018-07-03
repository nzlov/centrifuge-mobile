package com.example.fz.centrifugoandroid;

import android.app.Activity;
import android.content.Context;
import android.content.SharedPreferences;
import android.util.Log;
import android.widget.TextView;

import java.util.ArrayList;

import centrifuge.Message;
import centrifuge.MessageHandler;
import centrifuge.ReadHandler;
import centrifuge.Sub;
import centrifuge.SubscribeSuccessContext;
import centrifuge.SubscribeSuccessHandler;

import static android.content.Context.MODE_WORLD_WRITEABLE;

public class AppMessageHandler implements MessageHandler,ReadHandler {
    protected MainActivity context;

    public AppMessageHandler(Context context) {
        this.context = (MainActivity) context;
    }

    @Override
    public void onMessage(Sub sub, final Message message) {
        context.runOnUiThread(new Runnable() {
            @Override
            public void run() {
                TextView tv = (TextView) context.findViewById(R.id.text);
                tv.setText(message.getData());
            }
        });
        Log.i("centrifugo", message.toString());
        SharedPreferences.Editor editor = this.context.getSharedPreferences("centrifugo", Context.MODE_PRIVATE).edit();
        editor.putString("lastMsgid", message.getUID());
        editor.commit();
            try {
                //获得消息后通知服务器已读
                sub.readMessage(message.getUID());
            } catch (Exception e) {
                e.printStackTrace();
            }
    }

    @Override
    public void onRead(Sub sub, final String msgid) {
        Log.i("Read",sub.channel()+"-"+msgid);
    }

}
