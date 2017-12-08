//
//  ViewController.m
//  CentrifugoIOS
//
//  Created by nzlov on 2017/12/8.
//  Copyright © 2017年 nzlov. All rights reserved.
//

#import "ViewController.h"





@interface ViewController ()
@end

@implementation ViewController
@synthesize label;

- (void)viewDidLoad {
    [super viewDidLoad];
    // Do any additional setup after loading the view, typically from a nib.
    [self.label  setText: @"Connecting..."];
    dispatch_queue_t queue = dispatch_get_global_queue(DISPATCH_QUEUE_PRIORITY_DEFAULT, 0);
    dispatch_async(queue, ^{
        //创建令牌
        CentrifugeCredentials* creds = CentrifugeNewCredentials(@"42", @"1488055494", @"", @"24d0aa4d7c679e45e151d268044723d07211c6a9465d0e35ee35303d13c5eeff");
        //绑定连接事件
        CentrifugeEventHandler* eventHandler = CentrifugeNewEventHandler();
        [eventHandler onConnect:self];
        [eventHandler onDisconnect:self];
        
        //创建客户端连接
        CentrifugeClient *client = CentrifugeNew(@"ws://localhost:8000/connection/websocket", creds, eventHandler, CentrifugeDefaultConfig());
        //连接服务器
        [client connect:NULL];
        //绑定消息事件
        CentrifugeSubEventHandler*  subEventHandler = CentrifugeNewSubEventHandler();
        [subEventHandler onMessage:self];
        //订阅通道
        [client subscribe:@"public:chat" events:subEventHandler error:NULL];
    });
}


- (void)didReceiveMemoryWarning {
    [super didReceiveMemoryWarning];
    // Dispose of any resources that can be recreated.
}


- (void)onConnect:(CentrifugeClient *)p0 p1:(CentrifugeConnectContext *)p1 {
    dispatch_async( dispatch_get_main_queue(), ^{
        NSLog(@"Connected");
        [self.label setText:@"Connected"];
    });
}

- (void)onDisconnect:(CentrifugeClient *)p0 p1:(CentrifugeDisconnectContext *)p1 {
    
    dispatch_async( dispatch_get_main_queue(), ^{
        NSLog(@"Disconnect");
        [self.label setText:@"Disconnect"];
    });
}

- (void)onMessage:(CentrifugeSub *)p0 p1:(CentrifugeMessage *)p1 {
    NSLog(@"New Message:%@",p1.data);
    //获得消息后通知服务器已读
    BOOL a = NO;
    [p0 readMessage:p1.uid ret0_:&a error:NULL];
    
    dispatch_async( dispatch_get_main_queue(), ^{
        [self.label setText:p1.data];
    });
}


@end
