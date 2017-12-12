package com.example.fz.centrifugoandroid;

import android.support.v7.app.AppCompatActivity;
import android.os.Bundle;
import android.widget.TextView;

import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

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


    public String token(String secret,String user,String timestamp ,String info){
        try
        {
            // SHA 加密开始
            // 创建加密对象 并傳入加密類型
            MessageDigest messageDigest = MessageDigest.getInstance("SHA-256");
            // 传入要加密的字符串
            messageDigest.update(user.getBytes());
            messageDigest.update(timestamp.getBytes());
            messageDigest.update(info.getBytes());
            // 得到 byte 類型结果
            byte byteBuffer[] = messageDigest.digest();

            // 將 byte 轉換爲 string
            StringBuffer strHexString = new StringBuffer();
            // 遍歷 byte buffer
            for (int i = 0; i < byteBuffer.length; i++)
            {
                String hex = Integer.toHexString(0xff & byteBuffer[i]);
                if (hex.length() == 1)
                {
                    strHexString.append('0');
                }
                strHexString.append(hex);
            }
            // 得到返回結果
           return strHexString.toString();
        }
        catch (NoSuchAlgorithmException e)
        {
            e.printStackTrace();
        }
        return "";
    }

    @Override
    protected void onCreate(Bundle savedInstanceState) {
        super.onCreate(savedInstanceState);
        setContentView(R.layout.activity_main);
        TextView tv = (TextView) findViewById(R.id.text);

        //创建令牌
        Credentials creds = Centrifuge.newCredentials(
                "42", "1488055494", "",
                token("109AF84FWF45AS4S5W8F","42","1488055494","")
        );

        //绑定连接事件
        EventHandler events = Centrifuge.newEventHandler();
        ConnectHandler connectHandler = new AppConnectHandler(this);
        DisconnectHandler disconnectHandler = new AppDisconnectHandler(this);

        events.onConnect(connectHandler);
        events.onDisconnect(disconnectHandler);

        //创建客户端连接
        Client client = Centrifuge.new_(
                "ws://192.168.1.9:8000/connection/websocket",
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
            Sub sub = client.subscribe("test", subEvents);
        } catch (Exception e) {
            e.printStackTrace();
        }
    }
}
