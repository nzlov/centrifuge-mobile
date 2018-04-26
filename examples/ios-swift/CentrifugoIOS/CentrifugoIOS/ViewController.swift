//
//  ViewController.swift
//  CentrifugoIOS
//
//  Created by Alexander Emelin on 25/02/2017.
//  Copyright © 2017 Alexander Emelin. All rights reserved.
//

import UIKit
import Centrifuge

class ConnectHandler : NSObject, CentrifugeConnectHandlerProtocol {
    var l: UILabel!
    
    func setLabel(l: UILabel!) {
        self.l = l
    }
    
    func onConnect(_ p0: CentrifugeClient!, p1: CentrifugeConnectContext!) {
        DispatchQueue.main.async{
            self.l.text = "Connected";
        }
    }
}

class DisconnectHandler : NSObject, CentrifugeDisconnectHandlerProtocol {
    var l: UILabel!
    
    func setLabel(l: UILabel!) {
        self.l = l
    }
    
    func onDisconnect(_ p0: CentrifugeClient!, p1: CentrifugeDisconnectContext!) {
        DispatchQueue.main.async{
            self.l.text = "Disconnected";
        }
    }
}

class MessageHandler : NSObject, CentrifugeMessageHandlerProtocol,CentrifugeReadHandlerProtocol {
    var l: UILabel!
    
    func setLabel(l: UILabel!) {
        self.l = l
    }
    
    func onMessage(_ p0: CentrifugeSub!, p1: CentrifugeMessage!) {
        DispatchQueue.main.async{
            self.l.text = p1.data()
        }
        do {
            var ok = UnsafeMutablePointer<ObjCBool>.allocate(capacity: 1)
            try  p0.readMessage(p1.uid(), ret0_: ok)
        } catch {
            return
        }
    }
    func onRead(_ p0: CentrifugeSub!, p1: String!) {
        print("Read"+p1)
    }
}

class ViewController: UIViewController {
    
    @IBOutlet weak var label: UILabel!
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        label.text = "Connecting..."
        
        DispatchQueue.main.async{
            let creds = CentrifugeNewCredentials(
                "42","ios_merchant", "1488055494", "", "17445e62d61dfd1fd9e81d0aede358bdeb490e3e9c6cd92f3fd661b72c95b37b"
            )
            
            let eventHandler = CentrifugeNewEventHandler()
            let connectHandler = ConnectHandler()
            connectHandler.setLabel(l: self.label)
            let disconnectHandler = DisconnectHandler()
            disconnectHandler.setLabel(l: self.label)
            
            eventHandler?.onConnect(connectHandler)
            eventHandler?.onDisconnect(disconnectHandler)
            
            
            let url = "ws://192.168.1.200:8000/connection/websocket"
            let client = CentrifugeNew(url, creds, eventHandler, CentrifugeDefaultConfig())
            
            do {
                try client?.connect()
            } catch {
                self.label.text = "Error on connect..."
                return
            }
            
            self.label.text = "Connected"
            
            let subEventHandler = CentrifugeNewSubEventHandler()
            let messageHandler = MessageHandler()
            messageHandler.setLabel(l: self.label)
            subEventHandler?.onMessage(messageHandler)
            subEventHandler?.onRead(messageHandler)
            
            var sub: CentrifugeSub!
            do {
                sub = try client?.subscribe("public:chat", events: subEventHandler)
            } catch {
                DispatchQueue.main.async{
                    self.label.text = "Subscribe error"
                }
                return
            }
            
        }
    }
    
    override func didReceiveMemoryWarning() {
        super.didReceiveMemoryWarning()
        // Dispose of any resources that can be recreated.
    }
    
    
}
