//
//  ViewController.m
//  CentrifugoIOS
//
//  Created by nzlov on 2017/12/8.
//  Copyright © 2017年 nzlov. All rights reserved.
//

#import "ViewController.h"
#import <CommonCrypto/CommonDigest.h>
#import <CommonCrypto/CommonHMAC.h>




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
        CentrifugeCredentials* creds = CentrifugeNewCredentials(@"42", @"1488055494", @"", [self hmac:@"421488055494" withKey:@"109AF84FWF45AS4S5W8F"]);
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

- (NSString *)hmac:(NSString *)plaintext withKey:(NSString *)key
{
    const char *cKey  = [key cStringUsingEncoding:NSUTF8StringEncoding];
    const char *cData = [plaintext cStringUsingEncoding:NSUTF8StringEncoding];
    unsigned char cHMAC[CC_SHA256_DIGEST_LENGTH];
    CCHmac(kCCHmacAlgSHA256, cKey, strlen(cKey), cData, strlen(cData), cHMAC);
    NSData *HMACData = [NSData dataWithBytes:cHMAC length:sizeof(cHMAC)];
    const unsigned char *buffer = (const unsigned char *)[HMACData bytes];
    NSMutableString *HMAC = [NSMutableString stringWithCapacity:HMACData.length * 2];
    for (int i = 0; i < HMACData.length; ++i){
        [HMAC appendFormat:@"%02x", buffer[i]];
    }
    
    return HMAC;
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
