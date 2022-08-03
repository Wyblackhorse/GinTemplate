//
global_requestAddress = "";
//登录接口
global_requestAddressJs_login = global_requestAddress+ "/v2/login";

global_requestAddressJs_getPayInformation = global_requestAddress+ "/v2/getPayInformation";

global_requestAddressJs_getPrepaidPhoneOrders = global_requestAddress+ "/v2/getPrepaidPhoneOrders";

global_requestAddressJs_getAllMoney = global_requestAddress+ "/v2/getAllMoney";

global_requestAddressJs_getDailyStatistics = global_requestAddress+ "/v2/getDailyStatistics";

global_requestAddressJs_getReceiveAddress = global_requestAddress+ "/v2/getReceiveAddress";

global_requestAddressJs_collection = global_requestAddress+ "/v2/collection";



var getRootPath_webStr = getRootPath_web();

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


		return projectName;
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


function get_unix_time(dateStr) {
	var newstr = dateStr.replace(/-/g, '/');
	var date = new Date(newstr);
	var time_str = date.getTime().toString();
	return time_str.substr(0, 10);
}
