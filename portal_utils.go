package main

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"github.com/robertkrimen/otto"
)

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

func info(userName string, passWord string, acID string, userIP string, challenge string) string {
	vm := otto.New()
	_, _ = vm.Run(`
//jquery-base64 v1.0
"use strict";
base64 = (function($) {
    var _PADCHAR = "="
      , _ALPHA = "LVoJPiCN2R8G90yg+hmFHuacZ1OWMnrsSTXkYpUq/3dlbfKwv6xztjI7DeBE45QA"
      , _VERSION = "1.0";
    function _getbyte64(s, i) {
        var idx = _ALPHA.indexOf(s.charAt(i));
        if (idx === -1) {
            throw "Cannot decode base64"
        }
        return idx
    }
    function _setAlpha(s) {
        _ALPHA = s;
    }
    function _decode(s) {
        var pads = 0, i, b10, imax = s.length, x = [];
        s = String(s);
        if (imax === 0) {
            return s
        }
        if (imax % 4 !== 0) {
            throw "Cannot decode base64"
        }
        if (s.charAt(imax - 1) === _PADCHAR) {
            pads = 1;
            if (s.charAt(imax - 2) === _PADCHAR) {
                pads = 2
            }
            imax -= 4
        }
        for (i = 0; i < imax; i += 4) {
            b10 = (_getbyte64(s, i) << 18) | (_getbyte64(s, i + 1) << 12) | (_getbyte64(s, i + 2) << 6) | _getbyte64(s, i + 3);
            x.push(String.fromCharCode(b10 >> 16, (b10 >> 8) & 255, b10 & 255))
        }
        switch (pads) {
        case 1:
            b10 = (_getbyte64(s, i) << 18) | (_getbyte64(s, i + 1) << 12) | (_getbyte64(s, i + 2) << 6);
            x.push(String.fromCharCode(b10 >> 16, (b10 >> 8) & 255));
            break;
        case 2:
            b10 = (_getbyte64(s, i) << 18) | (_getbyte64(s, i + 1) << 12);
            x.push(String.fromCharCode(b10 >> 16));
            break
        }
        return x.join("")
    }
    function _getbyte(s, i) {
        var x = s.charCodeAt(i);
        if (x > 255) {
            throw "INVALID_CHARACTER_ERR: DOM Exception 5"
        }
        return x
    }
    function _encode(s) {
        if (arguments.length !== 1) {
            throw "SyntaxError: exactly one argument required"
        }
        s = String(s);
        var i, b10, x = [], imax = s.length - s.length % 3;
        if (s.length === 0) {
            return s
        }
        for (i = 0; i < imax; i += 3) {
            b10 = (_getbyte(s, i) << 16) | (_getbyte(s, i + 1) << 8) | _getbyte(s, i + 2);
            x.push(_ALPHA.charAt(b10 >> 18));
            x.push(_ALPHA.charAt((b10 >> 12) & 63));
            x.push(_ALPHA.charAt((b10 >> 6) & 63));
            x.push(_ALPHA.charAt(b10 & 63))
        }
        switch (s.length - imax) {
        case 1:
            b10 = _getbyte(s, i) << 16;
            x.push(_ALPHA.charAt(b10 >> 18) + _ALPHA.charAt((b10 >> 12) & 63) + _PADCHAR + _PADCHAR);
            break;
        case 2:
            b10 = (_getbyte(s, i) << 16) | (_getbyte(s, i + 1) << 8);
            x.push(_ALPHA.charAt(b10 >> 18) + _ALPHA.charAt((b10 >> 12) & 63) + _ALPHA.charAt((b10 >> 6) & 63) + _PADCHAR);
            break
        }
        return x.join("")
    }
    return {
        decode: _decode,
        encode: _encode,
        setAlpha: _setAlpha,
        VERSION: _VERSION
    }
}());

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
    function json(d) {
    	return JSON.stringify(d);
    }

    function info(d, k) {
        return "{SRBX1}" + base64.encode(xEncode(json(d), k));
    }

	function info2(username,password,client_ip,ac_id,token){
		var enc = "s" + "run" + "_bx1";
		var i = info({
        			username: username,
        			password: password,
        			ip: client_ip,
        			acid: ac_id,
        			enc_ver: "srun_bx1"
        		}, token)
		return i;
	}
`)
	info2, err := vm.Get("info2")
	if err != nil {
		Logger.Println(err)
	}
	v, err := info2.Call(otto.Value{}, userName, passWord, userIP, acID, challenge)
	if err != nil {
		Logger.Println(err)
	}
	//Logger.Println(v)
	return v.String()
}
