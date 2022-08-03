//
global_requestAddress = "";
//登录接口
js_global_requestAddress_login="/login";

//获取充值、昨日充提、过滤重复
js_global_requestAddress_getRecharge= "/getRecharge"

//获取玩家
js_global_requestAddress_getUserApp = "/getUserApp"

//获取钱包记录
js_global_requestAddress_getWalletRecord= "/getWalletRecord"

//获取提现
js_global_requestAddress_GetWithdraw = "/GetWithdraw"

//获取下注记录
js_global_requestAddress_GetBettingRecord= "/GetBettingRecord"

//获取用户玩家登录日志记录
js_global_requestAddress_GetAppUserLoginLog= "/GetAppUserLoginLog"

//日统计接口
js_global_requestAddress_SetStatistics= "/SetStatistics"

// //昨日充提、过滤重复
// js_global_requestAddress_GetRecharge = "/GetRecharge"




//上传接口
js_global_requestAddress_uploadFiles = "/uploadFiles"

//获取目录路径方法
function getRootPath_web() {

		//获取当前网址，如： http://localhost:8888/eeeeeee/aaaa/vvvv.html
		var curWwwPath = window.document.location.href;
		//获取主机地址之后的目录，如： uimcardprj/share/meun.jsp
		var pathName = window.document.location.pathname;
		var pos = curWwwPath.indexOf(pathName);
		//获取主机地址，如： http://localhost:8888
		var localhostPaht = curWwwPath.substring(0, pos);
		//获取带"/"的项目名，如：/abcd
		var projectName = pathName.substring(0, pathName.substr(1).indexOf('/') + 1);

		// return (localhostPaht + projectName);


		// console.log("当前网址:"+curWwwPath);
		// console.log("主机地址后的目录:"+pos+"----"+pathName);
		// console.log("主机地址:"+localhostPaht);
		// console.log("项目名:"+projectName);


		return localhostPaht;
}



//时间戳转日期时间型工具类
function formatDateTime(inputTime) {
	var date = new Date(inputTime);
	var y = date.getFullYear();
	var m = date.getMonth() + 1;
	m = m < 10 ? ('0' + m) : m;
	var d = date.getDate();
	d = d < 10 ? ('0' + d) : d;
	var h = date.getHours();
	h = h < 10 ? ('0' + h) : h;
	var minute = date.getMinutes();
	var second = date.getSeconds();
	minute = minute < 10 ? ('0' + minute) : minute;
	second = second < 10 ? ('0' + second) : second;
	return y + '-' + m + '-' + d+' '+h+':'+minute+':'+second;
}


function toDecimal2(x) {//金额处理两位小数点
	var f = parseFloat(x);
	if (isNaN(f)) {
		return false;
	}
	var f = Math.round(x*100)/100;
	var s = f.toString();
	var rs = s.indexOf('.');
	if (rs < 0) {
		rs = s.length;
		s += '.';
	}
	while (s.length <= rs + 2) {
		s += '0';
	}
	return s;
}


/**
 * 数字转整数 如 100000 转为10万
 * @param {需要转化的数} num
 * @param {需要保留的小数位数} point
 */
function tranNumber(num, point) {



	let numStr = num.toString()

	// console.log(numStr.length);
	// 一万以内直接返回
	if (numStr.length <=4) {
		return numStr;
	}
	//大于6位数是十万 (以10W分割 10W以下全部显示)
	else if (numStr.length > 4) {
		let decimal = numStr.substring(numStr.length - 4, numStr.length - 4 + point)
		// return parseFloat(parseInt(num / 10000) + ‘.’ + decimal) + ‘万’;
		return parseFloat(parseInt(num / 10000) + '.' + decimal) + '万';
	}
}




//验证是否为数字
function isNumber(value) { //验证是否为数字

	var patrn = /^(-)?\d+(\.\d+)?$/;

	if (patrn.exec(value) == null || value == "") {
		return false

	} else {
		return true

	}

}

// 小数转化为百分比
function toPercent(point){
	var str=Number(point*100).toFixed(1);
	str+="%";
	return str;
}

// 百分数转化为小数
function toPoint(percent){
	var str=percent.replace("%","");
	str= str/100;
	return str;
}


//校验日期格式
function checkDate(dateStr){
	var a = /^(\d{4})-(\d{2})-(\d{2})$/
	if (!a.test(dateStr)) {
		return false
	}else{
		return true;
	}
}
