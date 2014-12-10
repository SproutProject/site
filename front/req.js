var reqlang = new function(){
    var that = this;
    that.load = function(){
	var j_main = $('#reqlang');
    };
};
var reqalg = new function(){
    var that = this;
    that.load = function(){
	var j_main = $('#reqalg');
	var t_prepro = $('#prepro-templ').html();
	var prolist = null;

	$.post('/spt/d/req/getpre',{
	    'data':JSON.stringify(0)
	},function(res){
	    var i,j;
	    
	    prolist = res.data;
	    for(i = 0;i < prolist.length;i++){
		pro = prolist[i];
		pro.index = i + 1;
		for(j = 0;j < pro.Option.length;j++){
		    pro.Option[j] = {
			'desc':pro.Option[j],
			'value':j
		    };
		}
	    }

	    j_main.find('div.prepro').html(Mustache.render(t_prepro,prolist));

	    j_main.find('div.prepro > button').on('click',function(e){
		ans = [];
		for(i = 0;i < prolist.length;i++){
		    val = j_main.find('div.prepro input[name=' + (i + 1) + ']:checked').val();
		    ans[i] = parseInt(val);
		}

		$.post('/spt/d/req/checkpre',{
		    'data':JSON.stringify(ans)
		},function(res){
		    if(res.status == 'SUCCES'){
			j_main.find('div.prepro').hide();
			j_main.find('div.checkmail').show();
		    }else{
			alert('答案錯誤，請重新填寫');
			location.reload();
		    }
		});
	    });

	    j_main.find('div.checkmail > button.send').on('click',function(e){
		var j_btn = $(this);
		var mail = j_main.find('div.checkmail > input.mail').val();
		var repeat = j_main.find('div.checkmail > input.mail-repeat').val();

		if(mail != repeat){
		    alert('信箱輸入有錯誤');
		}else{
		    $.post('/spt/d/req/checkmail',{
			'data':mail
		    },function(res){
			if(res.status == 'SUCCES'){
			    j_main.find('div.checkmail').hide();
			    j_main.find('div.verify').show();
			}else{
			    alert('信箱尚未符合申請資格');
			    j_btn.show();
			    j_btn.siblings('span.msg').hide();
			}
		    });

		    j_btn.hide();
		    j_btn.siblings('span.msg').show();
		}
	    });
	    j_main.find('div.verify > button.send').on('click',function(e){
		var j_btn = $(this);
		verify = j_main.find('div.verify > input.verify').val();

		$.post('/spt/d/req/verify',{
		    'data':verify
		},function(res){
		    if(res.status == 'SUCCES'){
			j_btn.siblings('span.msg').hide();
		    }else{
			alert('驗證碼錯誤');
			j_btn.show();
			j_btn.siblings('span.msg').hide();
		    }
		});

		j_btn.hide();
		j_btn.siblings('span.msg').show();
	    });
	});
    };
};
