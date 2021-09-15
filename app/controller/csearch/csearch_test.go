package csearch

//func init() {
//  dao.AppEngine = dao.InitXorm("root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4", []string{"root:a3202381@tcp(127.0.0.1:3306)/fpv2?charset=utf8mb4"})
//  dao.InitRedis("127.0.0.1:6379", "")
//}

// 综合搜索
//func TestColligateSearch(t *testing.T) {
//  c, _ := gin.CreateTestContext(httptest.NewRecorder())
//  svc := New(c)
//  vlist, ulist, rlist := svc.ColligateSearch("202009101933004667", "yoozoo")
//  t.Logf("vlist:%v, ulist:%v, recommend:%+v\n", vlist, ulist, rlist)
//}
//
//func BenchmarkColligateSearch(b *testing.B) {
//  for i := 0; i < b.N; i++ {
//    c, _ := gin.CreateTestContext(httptest.NewRecorder())
//    svc := New(c)
//    vlist, ulist, rlist := svc.ColligateSearch("202009101933004667", "狗")
//    b.Logf("vlist:%v, ulist:%v, rlist:%v\n", vlist, ulist, rlist)
//  }
//}
