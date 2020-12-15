package order
//
//import (
//  "sports_service/server/nsqlx/protocol"
//  "sports_service/server/tools/nsq"
//  "sports_service/server/util"
//  "time"
//  "sports_service/server/global/app/log"
//  "sports_service/server/global/consts"
//)
//
//// 事件延时消息
//func DeferredOEventMsg(userId, orderId string, eventType int32, delayTm time.Duration) {
//  log.Log.Debugf("event_trace: 延时处理，orderId = %s， eventType:%d, delayTm:%d", orderId, eventType, delayTm)
//  OrderEventNSQPub(userId, orderId, eventType, delayTm)
//}
//
//// 用户订单相关事件
//func OrderEventNSQPub(userId, orderId string, eventType int32, delayTm time.Duration) {
//  body := newOrderEvent(userId, orderId, eventType)
//  if err := nsq.NsqProducer.DeferredPublish(consts.ORDER_EVENT_TOPIC, delayTm, body); err != nil {
//    log.Log.Errorf("appointment_event: publish appointment event err:%s, uid:%s", err, userId)
//  }
//}
//
//func newOrderEvent(userId, orderId string, eventType int32) []byte {
//  event := new(protocol.Event)
//  event.Uid = userId
//  event.EventType = eventType
//  event.Ts = time.Now().Unix()
//
//  appointment := new(protocol.OrderEvent)
//  appointment.OrderId = orderId
//  data , _ := util.JsonFast.Marshal(appointment)
//
//  event.Data = data
//  b, err := util.JsonFast.Marshal(event)
//  if err != nil {
//    log.Log.Errorf("order_event: marshal err:%s", err)
//  }
//
//  return b
//}

