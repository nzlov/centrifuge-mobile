//
//  ViewController.h
//  CentrifugoIOS
//
//  Created by nzlov on 2017/12/8.
//  Copyright © 2017年 nzlov. All rights reserved.
//

#import <UIKit/UIKit.h>
@import Centrifuge;

@interface ViewController : UIViewController <CentrifugeConnectHandler,CentrifugeDisconnectHandler,CentrifugeMessageHandler>{
    UILabel *label;
}

@property(nonatomic,retain) IBOutlet UILabel * label;

@end

