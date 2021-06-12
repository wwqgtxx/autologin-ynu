package main

import (
	"crypto"
	"crypto/hmac"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/dop251/goja"
)

var vm *goja.Runtime
var _xEncode goja.Callable
var _str2Bytes goja.Callable
var _ALPHA = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
var _StdEncoding = base64.NewEncoding(_ALPHA).WithPadding('=')

func pwd(passWord string, challenge string) string {
	hmacMD5 := hmac.New(crypto.MD5.New, []byte(challenge))
	hmacMD5.Write([]byte(passWord))
	return hex.EncodeToString(hmacMD5.Sum(nil))
}

func chksum(chkstr string) string {
	sha1 := crypto.SHA1.New()
	sha1.Write([]byte(chkstr))
	return hex.EncodeToString(sha1.Sum(nil))
}

func init() {
	vm = goja.New()
	_, err := vm.RunString(`
"use strict";
	function xEncode(str, key) {
		if (str == "") {
            return "";
        }
        var v = s(str, true),
            k = s(key, false);
        if (k.length < 4) {
            k.length = 4;
        }
        var n = v.length - 1,
            z = v[n],
            y = v[0],
            c = 0x86014019 | 0x183639A0,
            m,
            e,
            p,
            q = Math.floor(6 + 52 / (n + 1)),
            d = 0;
        while (0 < q--) {
            d = d + c & (0x8CE0D9BF | 0x731F2640);
            e = d >>> 2 & 3;
            for (p = 0; p < n; p++) {
                y = v[p + 1];
                m = z >>> 5 ^ y << 2;
                m += (y >>> 3 ^ z << 4) ^ (d ^ y);
                m += k[(p & 3) ^ e] ^ z;
                z = v[p] = v[p] + m & (0xEFB8D130 | 0x10472ECF);
            }
            y = v[0];
            m = z >>> 5 ^ y << 2;
            m += (y >>> 3 ^ z << 4) ^ (d ^ y);
            m += k[(p & 3) ^ e] ^ z;
            z = v[n] = v[n] + m & (0xBB390742 | 0x44C6F8BD);
        }
        return l(v, false);
	}
	function s(a, b) {
        var c = a.length, v = [];
        for (var i = 0; i < c; i += 4) {
            v[i >> 2] = a.charCodeAt(i) | a.charCodeAt(i + 1) << 8 | a.charCodeAt(i + 2) << 16 | a.charCodeAt(i + 3) << 24;
        }
        if (b) {
            v[v.length] = c;
        }
        return v;
    }
    function l(a, b) {
        var d = a.length, c = (d - 1) << 2;
        if (b) {
            var m = a[d - 1];
            if ((m < c - 3) || (m > c))
                return null;
            c = m;
        }
        for (var i = 0; i < d; i++) {
            a[i] = String.fromCharCode(a[i] & 0xff, a[i] >>> 8 & 0xff, a[i] >>> 16 & 0xff, a[i] >>> 24 & 0xff);
        }
        if (b) {
            return a.join('').substring(0, c);
        } else {
            return a.join('');
        }
    }

	function str2bytes(s){
		s = String(s);
		var i,t = [];
		for (i = 0; i < s.length; i++) {
			var x = s.charCodeAt(i);
			if (x > 255) {
				throw "INVALID_CHARACTER_ERR: DOM Exception 5"
			}
            t.push(x);
        }
		return t;
	}
`)
	if err != nil {
		Logger.Fatalln(err)
	}

	var ok bool
	_xEncode, ok = goja.AssertFunction(vm.Get("xEncode"))
	if !ok {
		Logger.Fatalln("can't get xEncode function")
	}
	_str2Bytes, ok = goja.AssertFunction(vm.Get("str2bytes"))
	if !ok {
		Logger.Fatalln("can't get str2Bytes function")
	}
}

func xEncode(d string, k string) (data []byte, err error) {
	v1, err := _xEncode(goja.Undefined(), vm.ToValue(d), vm.ToValue(k))
	if err != nil {
		return
	}
	v, err := _str2Bytes(goja.Undefined(), v1)
	if err != nil {
		return
	}
	err = vm.ExportTo(v, &data)
	return
}

func info(userName string, passWord string, acID string, userIP string, challenge string) string {
	dict := map[string]string{
		"username": userName,
		"password": passWord,
		"ip":       userIP,
		"acid":     acID,
		"enc_ver":  "srun_bx1",
	}
	jsonByte, err := json.Marshal(dict)
	if err != nil {
		Logger.Println(err)
	}

	v, err := xEncode(string(jsonByte), challenge)
	if err != nil {
		Logger.Println(err)
	}
	//Logger.Println(v)
	return "{SRBX1}" + _StdEncoding.EncodeToString(v)
}
