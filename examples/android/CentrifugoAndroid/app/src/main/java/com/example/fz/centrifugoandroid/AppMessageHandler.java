package com.example.fz.centrifugoandroid;

import android.app.Activity;
import android.content.Context;
import android.util.Log;
import android.widget.TextView;

import centrifuge.Message;
import centrifuge.MessageHandler;
import centrifuge.ReadHandler;
import centrifuge.Sub;

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
                TextView tv = (TextView) ((Activity) context).findViewById(R.id.text);
                tv.setText(message.getData());
            }
        });
        try {
            //获得消息后通知服务器已读
            sub.readMessage(message.getUID());
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void onRead(Sub sub, final String ch,final String msgid) {
        Log.i("Read",ch+"-"+msgid);
    }
}
