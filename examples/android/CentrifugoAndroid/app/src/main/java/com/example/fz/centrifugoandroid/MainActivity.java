package com.example.fz.centrifugoandroid;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.widget.TextView;

import centrifuge.Centrifuge;
import centrifuge.Client;
import centrifuge.Credentials;
import centrifuge.EventHandler;
import centrifuge.DisconnectHandler;
import centrifuge.ConnectHandler;
import centrifuge.MessageHandler;
import centrifuge.Sub;
import centrifuge.SubEventHandler;

public class MainActivity extends AppCompatActivity {

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        TextView tv = (TextView) findViewById(R.id.text);

        //创建令牌
        Credentials creds = Centrifuge.newCredentials(
                "42", "1488055494", "",
                "24d0aa4d7c679e45e151d268044723d07211c6a9465d0e35ee35303d13c5eeff"
        );

        //绑定连接事件
        EventHandler events = Centrifuge.newEventHandler();
        ConnectHandler connectHandler = new AppConnectHandler(this);
        DisconnectHandler disconnectHandler = new AppDisconnectHandler(this);

        events.onConnect(connectHandler);
        events.onDisconnect(disconnectHandler);

        //创建客户端连接
        Client client = Centrifuge.new_(
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
        tv.setText("Connected");

        //绑定消息事件
        SubEventHandler subEvents = Centrifuge.newSubEventHandler();
        MessageHandler messageHandler = new AppMessageHandler(this);
        subEvents.onMessage(messageHandler);

        try {
            //订阅通道
            Sub sub = client.subscribe("public:chat", subEvents);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
