var reqform = new function(){
    var that = this;
    that.load = function(clas){
	var j_main;
	var t_prepro = $('#prepro-templ').html();
	var prolist = null;

	if(clas == 0){
	    j_main = $('#reqalg');
	}else{
	    j_main = $('#reqlang');
	}

	$.post('/spt/d/req/getpre',{
	    'data':clas
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
			j_main.find('div.verify').hide();
			j_main.find('div.data').show();
		    }else{
			alert('驗證碼錯誤');
			j_btn.show();
			j_btn.siblings('span.msg').hide();
		    }
		});

		j_btn.hide();
		j_btn.siblings('span.msg').show();
	    });
	    j_main.find('div.data > button.send').on('click',function(e){
		var i;
		var j_textarea = j_main.find('div.data > textarea');
		var j_btn = $(this);
		var data = [];

		if(clas == 1){
		    var j_check;
		    var from = '';

		    j_check = j_main.find('div.data > input[name=from]:checked');
		    for(i = 0;i < j_check.length;i++){
			from += $(j_check[i]).val() + ',';
		    }
		    from += j_main.find('div.data > input[type=textbox][name=from]').val();
		    data.push(from);
		}

		for(i = 0;i < j_textarea.length;i++){
		    data.push($(j_textarea[i]).val());
		}

		$.post('/spt/d/req/data',{
		    'data':JSON.stringify(data)
		},function(res){
		    if(res.status == 'SUCCES'){
			j_main.find('div.data').hide();
			j_main.find('div.done').show();
		    }
		});
	    });
	});
    };
};
