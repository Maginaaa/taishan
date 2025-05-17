import sensors from 'sa-sdk-javascript'
// 一些动态的参数需要我们和服务端商议自行填写
sensors.init({
  name: 'sensors',
  server_url: 'https://track.cdfsunrise.com/receiver/api/gp?project=taishan&token=taishantoken',
  show_log: false, // 不输出log到控制台
  cross_subdomain: false, // 不在根域下设置cookie
  is_track_single_page: true, // 表示是否开启单页面自动采集 $pageview 功能，SDK 会在 url 改变之后自动采集web页面浏览事件 $pageview。
  heatmap: false,
  preset_properties: {
    //是否采集 $latest_utm 最近一次广告系列相关参数，默认值 true。
    latest_utm: false,
    //是否采集 $latest_traffic_source_type 最近一次流量来源类型，默认值 true。
    latest_traffic_source_type: false,
    //是否采集 $latest_search_keyword 最近一次搜索引擎关键字，默认值 true。
    latest_search_keyword: false,
    //是否采集 $latest_referrer 最近一次前向地址，默认值 true。
    latest_referrer: false,
    //是否采集 $latest_referrer_host 最近一次前向地址，1.14.8 以下版本默认是true，1.14.8 及以上版本默认是 false，需要手动设置为 true 开启。
    latest_referrer_host: false,
    //是否采集 $latest_landing_page 最近一次落地页地址，默认值 false。
    latest_landing_page: false,
    //是否采集 $url 页面地址作为公共属性，1.16.5 以下版本默认是 false，1.16.5 及以上版本默认是 true。
    url: false,
    //是否采集 $title 页面标题作为公共属性，1.16.5 以下版本默认是 false，1.16.5 及以上版本默认是 true。
    title: true
  }
  // heatmap: {
  //   // 是否开启点击图，default 表示开启，自动采集 $WebClick 事件，可以设置 'not_collect' 表示关闭。
  //   clickmap: 'not_collect',
  //   // 是否开启触达注意力图，not_collect 表示关闭，不会自动采集 $WebStay 事件，可以设置 'default' 表示开启。
  //   scroll_notice_map: 'not_collect',
  // },
})
//公共属性
sensors.registerPage({
  app_name: '泰山'
  // referrer: document.referrer,
  // path: location.href.split('?')[0],
})
//用于采集 $pageview 事件。
// sensors.quick('autoTrack', {
//   platform: 'h5',
// });
export default sensors
