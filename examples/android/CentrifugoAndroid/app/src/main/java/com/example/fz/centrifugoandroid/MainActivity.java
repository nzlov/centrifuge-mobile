package com.example.fz.centrifugoandroid;

import android.content.SharedPreferences;
import android.os.Bundle;
import android.support.v7.app.AppCompatActivity;
import android.util.Log;
import android.widget.TextView;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;

import centrifuge.Centrifuge;
import centrifuge.Client;
import centrifuge.ConnectHandler;
import centrifuge.Credentials;
import centrifuge.DisconnectHandler;
import centrifuge.EventHandler;
import centrifuge.MicroResponseBody;
import centrifuge.Sub;
import centrifuge.SubEventHandler;

public class MainActivity extends AppCompatActivity {

    private Sub sub;
    private Client client;
    private SubEventHandler subEvents;

    /**
     * 将加密后的字节数组转换成字符串
     *
     * @param b 字节数组
     * @return 字符串
     */
    private static String byteArrayToHexString(byte[] b) {
        StringBuilder hs = new StringBuilder();
        String stmp;
        for (int n = 0; b != null && n < b.length; n++) {
            stmp = Integer.toHexString(b[n] & 0XFF);
            if (stmp.length() == 1)
                hs.append('0');
            hs.append(stmp);
        }
        return hs.toString().toLowerCase();
    }

    private static String token(String secret, String appkey, String user, String timestamp, String info) {
        String hash = "";
        try {
            Mac sha256_HMAC = Mac.getInstance("HmacSHA256");
            SecretKeySpec secret_key = new SecretKeySpec(secret.getBytes(), "HmacSHA256");
            sha256_HMAC.init(secret_key);
            byte[] bytes = sha256_HMAC.doFinal((user + appkey + timestamp + info).getBytes());
            hash = byteArrayToHexString(bytes);
        } catch (Exception e) {
            System.out.println("Error HmacSHA256 ===========" + e.getMessage());
        }
        return hash;
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        TextView tv = (TextView) findViewById(R.id.text);

        //创建令牌
        Credentials creds = Centrifuge.newCredentials(
                "42", "android_merchant", "1488055494", "",
                token("109AF84FWF45AS4S5W8F", "57a883afde","42", "1488055494", "")
        );

        //绑定连接事件
        EventHandler events = Centrifuge.newEventHandler();
        ConnectHandler connectHandler = new AppConnectHandler(this);
        DisconnectHandler disconnectHandler = new AppDisconnectHandler(this);

        events.onConnect(connectHandler);
        events.onDisconnect(disconnectHandler);

        //创建客户端连接
        client = Centrifuge.new_(
                "ws://192.168.1.200:8000/connection/websocket",
                creds,
                events,
                Centrifuge.defaultConfig()
        );

        try {
            //连接服务器
            client.connect();
        } catch (Exception e) {
            e.printStackTrace();
            tv.setText(e.getMessage());
            return;
        }
        try {
            //调用微服务
            MicroResponseBody body = client.micro("Activity.Call","{\"a\":1}");
            Log.d("micro",body.getData());
        } catch (Exception e) {
            e.printStackTrace();
            tv.setText(e.getMessage());
            return;
        }


        tv.setText("Connected");

        //绑定消息事件
        subEvents = Centrifuge.newSubEventHandler();
        AppMessageHandler messageHandler = new AppMessageHandler(this);
        subEvents.onMessage(messageHandler);
        subEvents.onRead(messageHandler);

        SharedPreferences sharedPreferences = this.getSharedPreferences("centrifugo", MODE_PRIVATE);

        try {
            //订阅通道
            //传入保存的最后MSGID 消息服务器自动返回MSGID之后的消息到Message Handler
            sub = client.subscribeWithLastMsgID("users:wfhtqp", sharedPreferences.getString("lastMsgid", "1"), subEvents);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void onResume() {
        super.onResume();
        SharedPreferences sharedPreferences = this.getSharedPreferences("centrifugo", MODE_PRIVATE);
        Log.i("c", "onResume");
        try {
            //订阅通道
            //传入保存的最后MSGID 消息服务器自动返回MSGID之后的消息到Message Handler
            sub = client.subscribeWithLastMsgID("users:wfhtqp", sharedPreferences.getString("lastMsgid", "1"), subEvents);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }

    @Override
    public void onPause() {
        super.onPause();
        Log.i("c", "onPause");
        try {
            sub.unsubscribe();
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
