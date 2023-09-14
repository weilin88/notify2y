window.appPrefix = "/"; 
window.getAPIUrl = function(api){
	if(appPrefix === "/"){
		return api;
	}
	return appPrefix+api;
}
window.newRequest = function(argv){
    let req = {}
    if(argv) {
        req.Argvs = argv;
    } else {
        req.Argvs = [];
    }
    return JSON.stringify(req);
}                        
window.getFileName = function(url){
    return url.substring(url.lastIndexOf('/')+1);
} 
// (new Date()).Format("yyyy-MM-dd hh:mm:ss.S") ==> 2006-07-02 08:09:04.423
// (new Date()).Format("yyyy-M-d h:m:s.S")   ==> 2006-7-2 8:9:4.18
Date.prototype.Format = function (fmt) { //author: meizz
  var o = {
    "M+": this.getMonth() + 1, 
    "d+": this.getDate(), 
    "h+": this.getHours(), 
    "m+": this.getMinutes(),
    "s+": this.getSeconds(), 
    "q+": Math.floor((this.getMonth() + 3) / 3), 
    "S": this.getMilliseconds() 
  };
  if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (this.getFullYear() + "").substr(4 - RegExp.$1.length));
  for (var k in o)
  if (new RegExp("(" + k + ")").test(fmt)) fmt = fmt.replace(RegExp.$1, (RegExp.$1.length == 1) ? (o[k]) : (("00" + o[k]).substr(("" + o[k]).length)));
  return fmt;
}
function cylxFixed(num, digit) {
  if(Object.is(parseFloat(num), NaN)) {
    return NaN;
  }
  num = parseFloat(num);
  let ret = (Math.round((num + Number.EPSILON) * Math.pow(10, digit)) / Math.pow(10, digit)).toFixed(digit);
  return Number(ret);
}