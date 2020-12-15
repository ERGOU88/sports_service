package order

import (
  "errors"
  "fmt"
  nsqsvc "github.com/nsqio/go-nsq"
  "saisai/app_service/app/config"
  "saisai/app_service/app/log"
  "saisai/app_service/global/consts"
  "saisai/app_service/global/dao"
  "saisai/app_service/models"
  "saisai/app_service/models/mconsultant"
  "saisai/app_service/models/mnotify"
  "saisai/app_service/models/morder"
  "saisai/app_service/models/muser"
  "saisai/app_service/models/umeng"
  "saisai/app_service/nsqlx/protocol"
  "saisai/app_service/util"
  "time"
)

func connectConsumer(channel string) (*nsqsvc.Consumer, error) {
  nsqConfig := nsqsvc.NewConfig()
  fmt.Printf("初始化 topic: %v,channel:%v", consts.ORDER_EVENT_TOPIC, channel)
  consumer, err := nsqsvc.NewConsumer(consts.ORDER_EVENT_TOPIC, channel, nsqConfig)
  if err != nil {
    log.Log.Errorf("new consumer err:%v", err)
    return consumer, err
  }

  return consumer, err
}

func OrderConsumer(channel string) (consumer *nsqsvc.Consumer) {
  consumer, err := connectConsumer(channel)
  if err != nil {
    panic(fmt.Sprintf("consumer conn err:%s", err))
  }

  consumer.AddHandler(nsqsvc.HandlerFunc(NsqHandler))

  err = consumer.ConnectToNSQD(config.Global.NsqAddr)
  if err != nil {
    log.Log.Errorf("ConnectToNSQD err:%s", err)
  }

  return
}

func NsqHandler(msg *nsqsvc.Message) error {
  event := new(protocol.Event)
  if err := util.JsonFast.Unmarshal(msg.Body, event); err != nil {
    log.Log.Errorf("appointment_event: proto unmarshal event err:%s", err)
    return err
  }

  if err := handleOrderEvent(event); err != nil {
    msg.RequeueWithoutBackoff(time.Second * 3)
    log.Log.Errorf("handleOrderEvent err:%s", err)
    return err
  }

  return nil
}

func handleOrderEvent(event *protocol.Event) error {
  info := &protocol.OrderEvent{}
  if err := util.JsonFast.Unmarshal(event.Data, info); err != nil {
    log.Log.Errorf("order_event: proto unmarshal data err:%s", err)
    return err
  }

  switch event.EventType {
  // 订单超时
  case consts.ORDER_TIME_OUT:
    if err := OrderTimeOut(info); err != nil {
      log.Log.Errorf("order time out err:%s", err)
      return err
    }

  // 订单付款提示 超时前15分钟
  case consts.ORDER_PAYMENT_MSG:
    if err := OrderPaymentNotify(event.EventType, info); err != nil {
      log.Log.Errorf("order payment notify err:%s", err)
      return err
    }

  // 待咨询订单通知(开始前1天 用户端及咨询师端)
  case consts.ORDER_WAIT_CONSULT_ADVANCE_ONE_DAY:
    if err := OrderWaitConsultNotify(event.EventType, info); err != nil {
      log.Log.Errorf("wait consult notify err:%s", err)
      return err
    }

  // 待咨询订单通知(开始前1小时 用户端及咨询师端)
  case consts.ORDER_WAIT_CONSULT_ADVANCE_ONE_HOUR:
    if err := OrderWaitConsultNotify(event.EventType, info); err != nil {
      log.Log.Errorf("wait consult notify err:%s", err)
      return err
    }

  // 咨询师写评估报告通知（结束后1小时 咨询师端）
  case consts.ORDER_WRITE_REPORT_END_AN_HOUR:
    if err := OrderWriteReportNotify(event.EventType, info); err != nil {
      log.Log.Errorf("write report notify err:%s", err)
      return err
    }

  // 咨询师写评估报告通知（结束后24小时 咨询师端）
  case consts.ORDER_WRITE_REPORT_END_ONE_DAY:
    if err := OrderWriteReportNotify(event.EventType, info); err != nil {
      log.Log.Errorf("write report notify err:%s", err)
      return err
    }
  }

  return nil
}

// 咨询师写评估报告通知（结束后1小时/24小时 咨询师端）
func OrderWriteReportNotify(eventType int32, event *protocol.OrderEvent) error {
  session := dao.Engine.NewSession()
  defer session.Close()
  rmodel := morder.NewOrderModel(session)
  order := rmodel.GetOrder(event.OrderId)
  // 订单不存在
  if order == nil {
    log.Log.Errorf("order_event: order not found, orderId:%s", event.OrderId)
    return errors.New("order not exists, orderId:" + event.OrderId)
  }

  // 订单状态 != 3（咨询已结束 待咨询师填写评估报告） 则 不对订单做任何操作
  if order.Status != consts.COMPLETED_PAY_TYPE {
    log.Log.Errorf("order_event: don't need to change")
    return nil
  }

  info := &mconsultant.UserAppointmentForm{}
  if err := util.JsonFast.Unmarshal([]byte(order.Extra), event); err != nil {
    log.Log.Errorf("order_event: unmarshal err:%s", err)
    return nil
  }

  var msgType int32
  var content string
  if eventType == consts.ORDER_WAIT_CONSULT_ADVANCE_ONE_DAY {
    content = fmt.Sprintf("您好，您与用户%s的咨询已结结束24小时了，请及时填写评估报告", util.HideMobileNum(fmt.Sprint(info.MobileNum)))
    msgType = int32(consts.MSG_TYPE_WRITE_REPORT_END_ONE_DAY)

  } else {
    content = fmt.Sprintf("您好，您与用户%s的咨询已结结束1小时了，请及时填写评估报告", util.HideMobileNum(fmt.Sprint(info.MobileNum)))
    msgType = int32(consts.MSG_TYPE_WRITE_REPORT_END_AN_HOUR)
  }

  nmodel := mnotify.NewNotifyModel(session)
  now := int(time.Now().Unix())
  // 咨询师端通知
  msg := &models.SystemMessage{
    ReceiveId: order.ProductId,
    SendId: consts.SENDER,
    SendTime: now,
    SendType: int(eventType),
    SystemTopic: "写评估订单",
    SystemContent: content,
  }

  affected, err := nmodel.AddSystemNotify(msg)
  if err != nil || affected != 1 {
    log.Log.Errorf("order_trace: add write report msg err:%s, affected:%d, orderId:%s", err, affected, order.PayOrderId)
    return errors.New("send write report notify fail")
  }

  umodel := muser.NewUserModel(session)
  // 获取咨询师的信息（device_token 及 咨询师设备所属平台）
  user := umodel.FindUserByUserId(order.ProductId)
  if user == nil {
    log.Log.Errorf("order_event: user not found, consultant userId:%s", order.UserId)
    return nil
  }

  unreadNum := nmodel.GetUnreadSystemMsgNum(user.UserId)
  // 填写评估报告推送通知（咨询师端）
  PushNotify(user,"填写评估报告", content, msgType, unreadNum)

  return nil
}

// 待咨询订单通知(开始前1天/前1小时 通知用户端及咨询师端)
func OrderWaitConsultNotify(eventType int32, event *protocol.OrderEvent) error {
  session := dao.Engine.NewSession()
  defer session.Close()
  rmodel := morder.NewOrderModel(session)
  order := rmodel.GetOrder(event.OrderId)
  // 订单不存在
  if order == nil {
    log.Log.Errorf("order_event: order not found, orderId:%s", event.OrderId)
    return errors.New("order not exists, orderId:" + event.OrderId)
  }

  // 订单状态 != 2（已支付/待咨询） 则 不对订单做任何操作
  if order.Status != consts.PAID_PAY_TYPE {
    log.Log.Errorf("order_event: don't need to change")
    return nil
  }

  info := &mconsultant.UserAppointmentForm{}
  if err := util.JsonFast.Unmarshal([]byte(order.Extra), event); err != nil {
    log.Log.Errorf("order_event: unmarshal err:%s", err)
    return nil
  }

  var (
    userContent, consultantContent string
    msgType int32
  )
  if eventType == consts.ORDER_WAIT_CONSULT_ADVANCE_ONE_DAY {
    userContent = fmt.Sprintf("您好，您成功预约了%s 与%s咨询师的咨询服务，还有24小时开始，请在预约时间接听咨询师的沟通哦～",
      info.TimeNode, info.ConsultantName)

    consultantContent = fmt.Sprintf("您好，用户%s已成功预约了%s与您的咨询服务，还有24小时开始，请在预约时间主动联系用户哦~",
      util.HideMobileNum(fmt.Sprint(info.MobileNum)), info.TimeNode)

    msgType = int32(consts.MSG_TYPE_WAIT_CONSULT_ADVANCE_ONE_DAY)

  } else {
    userContent = fmt.Sprintf("您好，您成功预约了%s 与%s咨询师的咨询服务，还有1小时开始，请在预约时间接听咨询师的沟通哦～",
      info.TimeNode, info.ConsultantName)

    consultantContent = fmt.Sprintf("您好，用户%s已成功预约了%s与您的咨询服务，还有1小时开始，请在预约时间主动联系用户哦~",
      util.HideMobileNum(fmt.Sprint(info.MobileNum)), info.TimeNode)

    msgType = int32(consts.MSG_TYPE_WAIT_CONSULT_ADVANCE_AN_HOUR)
  }

  nmodel := mnotify.NewNotifyModel(session)
  now := int(time.Now().Unix())
  // 用户端通知
  sendUserMsg := &models.SystemMessage{
    ReceiveId: order.UserId,
    SendId: consts.SENDER,
    SendTime: now,
    SendType: int(eventType),
    SystemTopic: "待咨询订单",
    SystemContent: userContent,
  }

  // 咨询师端通知
  sendConsultantMsg := &models.SystemMessage{
    ReceiveId: order.ProductId,
    SendId: consts.SENDER,
    SendTime: now,
    SendType: int(eventType),
    SystemTopic: "待咨询订单",
    SystemContent: consultantContent,
  }

  msg := make([]*models.SystemMessage, 2)
  msg[0] = sendUserMsg
  msg[1] = sendConsultantMsg

  affected, err := nmodel.AddMultiSystemNotify(msg)
  if err != nil || affected != 2 {
    log.Log.Errorf("order_trace: add wait cnsult msg err:%s, affected:%d, orderId:%s", err, affected, order.PayOrderId)
    return errors.New("send wait consult msg fail")
  }

  umodel := muser.NewUserModel(session)
  // 用户信息
  user := umodel.FindUserByUserId(order.UserId)
  if user == nil {
    log.Log.Errorf("order_event: user not found, userId:%s", order.UserId)
    return nil
  }

  // 咨询师的用户信息
  consultant := umodel.FindUserByUserId(order.ProductId)
  if consultant == nil {
    log.Log.Errorf("order_event: user not found, consultant userId:%s", order.ProductId)
    return nil
  }

  unreadNum := nmodel.GetUnreadSystemMsgNum(user.UserId)
  // 待咨询订单推送（用户端）
  PushNotify(user, "待咨询订单", userContent, msgType, unreadNum)

  unreadNum = nmodel.GetUnreadSystemMsgNum(consultant.UserId)
  // 待咨询订单推送（咨询师端）
  PushNotify(consultant, "待咨询订单", consultantContent, msgType, unreadNum)
  return nil
}

// 推送通知
func PushNotify(user *models.User, title, content string, msgType int32, unreadNum int64) {
  extra := make(map[string]interface{}, 0)
  extra["unread_num"] = unreadNum
  // android推送
  if user.DeviceType == int(consts.ANDROID_PLATFORM) && user.DeviceToken != "" {
    client := umeng.New(umeng.SAI_ANDROID)
    if err := client.PushUnicastNotify(msgType, umeng.SAI_ANDROID, user.DeviceToken, title, content, extra); err != nil {
      log.Log.Errorf("order_event: push wait consult order notify by user err:%s", err)
    }
  }

  // iOS推送
  if user.DeviceType == int(consts.IOS_PLATFORM) && user.DeviceToken != "" {
    client := umeng.New(umeng.SAI_IOS)
    if err := client.PushUnicastNotify(msgType, umeng.SAI_IOS, user.DeviceToken, title, content, extra); err != nil {
      log.Log.Errorf("order_event: push wait consult order notify by user err:%s", err)
    }
  }
}

// 订单付款通知(超时前15分钟)
func OrderPaymentNotify(eventType int32, event *protocol.OrderEvent) error {
  session := dao.Engine.NewSession()
  defer session.Close()
  rmodel := morder.NewOrderModel(session)
  order := rmodel.GetOrder(event.OrderId)
  // 订单不存在
  if order == nil {
    log.Log.Errorf("order_event: order not found, orderId:%s", event.OrderId)
    return errors.New("order not exists, orderId:" + event.OrderId)
  }

  // 订单状态 != 0 （待支付） 表示 订单 已设为超时/已支付/已完成咨询  则 不对订单做任何操作
  if order.Status != consts.WAIT_PAY_TYPE {
    log.Log.Errorf("order_event: don't need to change")
    return nil
  }

  // 发送订单付款通知
  nmodel := mnotify.NewNotifyModel(session)

  content := fmt.Sprintf("您好，您于%s 预约的咨询师订单即将过期，请及时支付哦~",
    time.Unix(int64(order.CreateAt), 0).Format(consts.FORMAT_INFO))

  now := int(time.Now().Unix())
  // 用户端通知
  msg := &models.SystemMessage{
    ReceiveId: order.UserId,
    SendId: consts.SENDER,
    SendTime: now,
    SendType: int(eventType),
    SystemTopic: "待支付订单",
    SystemContent: content,
  }

  affected, err := nmodel.AddSystemNotify(msg)
  if err != nil || affected != 1 {
    log.Log.Errorf("order_trace: send payment notify err:%s, affected:%d, orderId:%s", err, affected, order.PayOrderId)
    return errors.New("send payment notify fail")
  }

  umodel := muser.NewUserModel(session)
  user := umodel.FindUserByUserId(order.UserId)
  if user == nil {
    log.Log.Errorf("order_event: user not found, userId:%s", order.UserId)
    return nil
  }

  unreadNum := nmodel.GetUnreadSystemMsgNum(order.UserId)
  PushNotify(user, "待支付订单", content, int32(consts.MSG_TYPE_PAYMENT_NOTIFY), unreadNum)
  return nil
}

// 订单超时（30分钟）
func OrderTimeOut(event *protocol.OrderEvent) error {
  // 先查询订单是否存在
  session := dao.Engine.NewSession()
  defer session.Close()
  rmodel := morder.NewOrderModel(session)
  order := rmodel.GetOrder(event.OrderId)
  // 订单不存在
  if order == nil {
    log.Log.Errorf("order_event: order not found, orderId:%s", event.OrderId)
    return errors.New("order not exists, orderId:" + event.OrderId)
  }

  // 订单状态 != 0 (待支付) 表示 订单 已设为超时/已支付/已完成咨询  则 不对订单做任何操作
  if order.Status != consts.WAIT_PAY_TYPE {
    log.Log.Errorf("order_event: don't need to change")
    return nil
  }

  // 将订单置为 超时/未支付
  order.Status = consts.UNPAID_PAY_TYPE
  order.UpdateAt = int(time.Now().Unix())
  if err := rmodel.UpdateOrderStatus(event.OrderId); err != nil {
    log.Log.Errorf("order_event: update order status fail, orderId:%s", event.OrderId)
    return errors.New("order_event: update order status fail")
  }

  return nil
}
